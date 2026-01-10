package servers

import (
	"errors"
	"fme_backend/internal/config"
	"fme_backend/internal/models"
	"fme_backend/internal/responses"
	"fme_backend/internal/schemas"
	"fme_backend/internal/utilities"
	"time"

	"gorm.io/gorm"
)

type UserServer struct{}

const (
	OTP_EXPIRES = 10
)

func (u *UserServer) CreateFmeUser(db *gorm.DB, data schemas.UserCreateSchema) error {
	err, _ := createUserWithRole(db, data, 1)
	return err
}

func (u *UserServer) SuspendUser(db *gorm.DB, userId int, instanceId, role string) error {
	var instance models.User
	if err := db.First(&instance, instanceId).Error; err != nil {
		return err
	}
	if err := validateOwnership(db, userId, role, instance); err != nil {
		return errors.New(responses.USER_SUSPEND_UNAUTHORIZED)
	}
	return suspendUser(db, instance)
}

func (u *UserServer) ActivateUser(db *gorm.DB, userId int, instanceId, role string) error {
	var instance models.User
	if err := db.First(&instance, instanceId).Error; err != nil {
		return err
	}
	if err := validateOwnership(db, userId, role, instance); err != nil {
		return errors.New(responses.USER_ACTIVATE_UNAUTHORIZED)
	}
	return activateUser(db, instance)

}

func (u *UserServer) Login(db *gorm.DB, data schemas.LoginSchema) (schemas.LoginResp, error) {
	var user models.User
	if err := db.Where("email = ?", data.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.LoginResp{}, errors.New(responses.INVALID_LOGIN)
		}
		return schemas.LoginResp{}, err
	}
	if !user.IsActive {
		return schemas.LoginResp{}, errors.New(responses.USER_SUSPENDED)
	}
	if err := utilities.CompareHashAndPassword(data.Password, user.Password); err != nil {
		return schemas.LoginResp{}, nil
	}
	token_resp, err := utilities.GeneratJWT(user.Email, user.ID, uint(user.Role))
	if err != nil {
		return schemas.LoginResp{}, err
	}
	return token_resp, nil
}

func (u *UserServer) RequestOtp(db *gorm.DB, data schemas.RequestOtpSchema) error {
	var user models.User
	if err := db.Where("email = ?", data.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(responses.NO_USER)
		}
		return err
	}
	if !user.IsActive {
		return errors.New(responses.USER_SUSPENDED)
	}
	user.OTP = utilities.GenerateOtp()
	user.OTPExpiresAt = time.Now().Add(time.Minute * OTP_EXPIRES)
	if err := config.DB.Save(&user).Error; err != nil {
		return err
	}
	go utilities.SendOtpEmail(user.Email, user.OTP)
	return nil
}

func (u *UserServer) VerifyOtp(db *gorm.DB, data schemas.VerifyOtpSchema) error {
	var user models.User
	if err := db.Where("email = ?", data.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(responses.NO_USER)
		}
		return err
	}
	if user.OTP != data.Otp || time.Now().After(user.OTPExpiresAt) {
		return errors.New(responses.INVALID_OTP)
	}
	updates := map[string]interface{}{
		"otp_verified":   true,
		"otp":            "",
		"otp_expires_at": time.Now().Add(time.Minute * 4),
	}
	if err := db.Model(&user).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

func (u *UserServer) ChangePassword(db *gorm.DB, data schemas.ChangePasswordSchema) error {
	var user models.User
	if err := db.Where("email = ?", data.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New(responses.NO_USER)
		}
		return err
	}
	if !user.OTPVerified || time.Now().After(user.OTPExpiresAt) {
		return errors.New(responses.PASSWORD_CHANGE_FAIL)
	}
	password_hash, err := utilities.HashPassword(data.Password)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"password":       string(password_hash),
		"otp_verified":   false,
		"otp_expires_at": time.Now(),
	}
	if err = db.Model(&user).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

// --internals--
func validateOwnership(db *gorm.DB, adminUserId int, adminRole string, targetInstance models.User) error {
	switch adminRole {
	case utilities.USER_ROLE_FME:
		if targetInstance.Role == 2 || targetInstance.Role == 3 || targetInstance.Role == 4 || targetInstance.Role == 6 {
			return nil
		}

	case utilities.USER_ROLE_MDA:
		mda_id, err := getMdaIdFromUserId(db, adminUserId)
		if err != nil {
			return err
		}

		switch targetInstance.Role {
		case 3:
			var instanceMdaId uint
			if err := db.Table("stcs").Select("mda_id").Where("user_id = ?", targetInstance.ID).Scan(&instanceMdaId).Error; err != nil {
				return err
			}
			if instanceMdaId == mda_id {
				return nil
			}

		case 4, 6:
			var instanceData struct {
				MdaId uint
				StcId uint
			}
			if err := db.Table("students").Select("mda_id, stc_id").Where("user_id = ?", targetInstance.ID).Scan(&instanceData).Error; err != nil {
				return err
			}
			switch {
			case instanceData.MdaId != 0:
				if mda_id == instanceData.MdaId {
					return nil
				}
			case instanceData.StcId != 0:
				parentMdaId, err := getMdaIdFromStcId(db, instanceData.StcId)
				if err != nil {
					return err
				}
				if mda_id == parentMdaId {
					return nil
				}
			}
		}

	case utilities.USER_ROLE_STC:
		stc_id, err := getStcIdFromUserId(db, adminUserId)
		if err != nil {
			return err
		}
		if targetInstance.Role == 4 || targetInstance.Role == 6 {
			var studentStcId uint
			if err := db.Table("students").Select("stc_id").Where("user_id = ?", targetInstance.ID).Scan(&studentStcId).Error; err != nil {
				return err
			}
			if studentStcId == stc_id {
				return nil
			}
		}
	}
	return errors.New(responses.UNAUTHORIZED_ACCESS)
}

func createMdaUser(tx *gorm.DB, data schemas.UserCreateSchema) (error, uint) {
	return createUserWithRole(tx, data, 2)
}

func createStcUser(tx *gorm.DB, data schemas.UserCreateSchema) (error, uint) {
	return createUserWithRole(tx, data, 3)
}

func createStudentUser(tx *gorm.DB, data schemas.UserCreateSchema) (error, uint) {
	return createUserWithRole(tx, data, 4)
}

func createEmployerUser(tx *gorm.DB, data schemas.UserCreateSchema) (error, uint) {
	return createUserWithRole(tx, data, 5)
}

func createUserWithRole(tx *gorm.DB, data schemas.UserCreateSchema, role int) (error, uint) {
	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", data.Email).Count(&count)
	if count > 0 {
		return errors.New(responses.USER_EXISTS), 0
	}

	hashed_password, err := utilities.HashPassword(data.Password)
	if err != nil {
		return errors.New(responses.UNABLE_TO_HASH_PASSWORD), 0
	}
	user := models.User{
		Email:        data.Email,
		Password:     hashed_password,
		OTPExpiresAt: time.Now(),
		Role:         role,
		IsActive:     true,
	}

	result := tx.Create(&user)
	if result.Error != nil {
		return errors.New(responses.USER_CREATE_FAILED), 0
	}

	return nil, user.ID
}

func getMdaIdFromUserId(db *gorm.DB, user_id int) (uint, error) {
	var mda_id uint
	err := db.Table("mdas").Where("user_id = ?", user_id).Pluck("id", &mda_id).Error
	if err != nil {
		return 0, err
	}
	return mda_id, nil
}

func getMdaIdFromStcId(db *gorm.DB, stc_id uint) (uint, error) {
	var mda_id uint
	err := db.Table("stcs").Where("id = ?", stc_id).Select("mda_id").Scan(&mda_id).Error
	if err != nil {
		return 0, err
	}
	return mda_id, nil
}

func getStcIdFromUserId(db *gorm.DB, user_id int) (uint, error) {
	var stc_id uint
	if err := db.Table("stcs").Select("id").Where("user_id = ?", user_id).Scan(&stc_id).Error; err != nil {
		return 0, err
	}
	return stc_id, nil
}

func suspendUser(db *gorm.DB, instance models.User) error {
	instance.IsActive = false
	result := db.Save(&instance)
	return result.Error
}

func activateUser(db *gorm.DB, instance models.User) error {
	instance.IsActive = true
	result := db.Save(&instance)
	return result.Error
}
