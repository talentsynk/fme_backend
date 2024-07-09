package myuser

import (
	//Inbuilt packages
	"fmt"
	"net/http"
	"time"

	//Project packages
	"fme_backend/internal/config"
	"fme_backend/internal/utilities"

	//External packages
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Create FME User
func CreateFmeUser(c *gin.Context) {
	//Read the request body and binds it to the schema variable
	if c.Bind(&UserCreateSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request body",
		})
		return
	}

	var userCheck User
	config.DB.Where("email= ?", UserCreateSchema.Email).First(&userCheck)
	if userCheck.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User already exists",
		})
		return
	}
	

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(UserCreateSchema.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// setup the user create instance
	user:= User{
		Email: UserCreateSchema.Email,
		Password: string(hash),
		OTPExpiresAt: time.Now(),
		Role:1,
		IsActive: true,
	}

	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	c.JSON(200, gin.H{"message": "User created successfully"})
}

//Suspend the user
func SuspendUser(c *gin.Context) {
		// Receive the part parameter.
		id:= c.Param("id")
		if id == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "path parameter not provided",
			})
			return
		}

		// get the instance  
		var instance User
		instance_result := config.DB.First(&instance, id)
		if instance_result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "instance does not exist",
			})
			return
		}

		// get the user and confirm permission
		userId ,userexists := c.Get("userID")
		if !userexists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Problem with the authorization token",
			})
		}

		var user User
		user_result := config.DB.First(&user, userId)
		if user_result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Problem with authorization token",
			})
			return
		}

		switch user.Role {
		case 1:
			if (instance.Role == 2 || instance.Role ==3 || instance.Role ==4){
				instance.IsActive = false
				result:= config.DB.Save(&instance)
				if result.Error !=nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "Unable to update the user record",
					})
					return
				}
				c.JSON(200, gin.H{"message": "User suspended successfully"})
				return
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid role for suspending user"})
				return
			}
			
			
		case 2:
			// Get user MDA ID
			var userMdaId uint 
			err := config.DB.Table("mdas").
			Where("user_id = ?", user.ID).
			Pluck("id",&userMdaId).Error
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "MdaAccount has issues",
				})
				return
			}
			//if the instance is an stc
			if instance.Role == 3 {
				//get the related mda id
				var insanceMdaId uint 
				err := config.DB.Table("stcs").Pluck("mda_id",&insanceMdaId).Error
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{
						"message": "MdaAccount has issues",
					})
					return
				}
				// compare if the instance is an stc under the user
				if (insanceMdaId == userMdaId) {
					instance.IsActive = false
					result:= config.DB.Save(&instance)
					if result.Error !=nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"message": "Cannot suspend this user",
						})
						return
					}
					c.JSON(200, gin.H{"message": "User suspended successfully"})
					return
				} else {	//if not return unauthorized
					c.JSON(http.StatusUnauthorized, gin.H{
						"message": "Cannot suspend this user",
					})
					return
				}
			} else if (instance.Role ==4){	// if the instance is a student
				// get the related data to determine if the student is an mda student or stc student
				var instanceData struct{
					MdaId uint
					StcId uint
				}
				err := config.DB.Table("students").
				Select("mda_id, stc_id").
				Where("user_id = ?", instance.ID).
				Scan(&instanceData).Error
				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{
						"message": "cannot get the instance data",
					})
					return
				}
				fmt.Println(instanceData)
				if instanceData.MdaId != 0 {	// if this is the student is under an mda directly - Mda student
					if userMdaId != instanceData.MdaId { // check if this student is under the authenticated mda
							c.JSON(http.StatusUnauthorized, gin.H{
								"message": "cannot get the instance data",
							})
							return
						}
					instance.IsActive = false	// deactivate the user
					result:= config.DB.Save(&instance)
					if result.Error !=nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"message": "Unable to update the user record",
						})
						return
					}
					c.JSON(200, gin.H{"message": "User suspended successfully"})
					return
				} else if instanceData.StcId != 0 {	// if the student is an stc student
					var insatnceStcMdaId uint	// get the related mda id 

					err := config.DB.Table("stcs").
					Where("id = ?", instanceData.StcId).
					Select("mda_id").
					Scan(&insatnceStcMdaId).Error
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"message": "Unable to update the user record",
						})
						return
					}

					if insatnceStcMdaId != userMdaId {	// check if the student is under the authenticated mda
						c.JSON(http.StatusUnauthorized, gin.H{
							"message": "Unable to suspend this user",
						})
						return
					}
					instance.IsActive = false	//deactivate the student
					result:= config.DB.Save(&instance)
					if result.Error !=nil {
						c.JSON(http.StatusInternalServerError, gin.H{
							"message": "Unable to update the user record",
						})
						return
					}
					c.JSON(200, gin.H{"message": "User suspended successfully"})
					return
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{
						"message": "Cannot suspend this user",
					})
					return
				}

			}

		case 3:
			//get user stcid
			var userStcId uint
			err := config.DB.
					Table("stcs").
					Select("id").
					Where("user_id = ?", user.ID).
					Scan(&userStcId).Error
			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error trying to suspend this user",
				})
				return
			}

			if instance.Role != 4{
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Cannot suspend this user",
				})
				return
			}

			// get student stc id
			var studentStcId uint
			nerr := config.DB.
					Table("students").
					Select("stc_id").
					Where("user_id = ?", instance.ID).
					Scan(&studentStcId).Error

			if nerr != nil{
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Error trying to suspend this user",
				})
				return
			}

			//check if student is under this stc
			if (studentStcId != userStcId) {
				c.JSON(http.StatusUnauthorized,gin.H{
					"message": "Unable to suspend this user",
				})
				return
			}

			instance.IsActive = false	//deactivate the student
			result:= config.DB.Save(&instance)
			if result.Error !=nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Unable to update the user record",
				})
				return
			}
			c.JSON(200, gin.H{"message": "User suspended successfully"})
			return
			
	

		default:
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Unable to update the user record",
			})
			return
		}
	}

// Activate User
func ActivateUser(c *gin.Context) {
	// Receive the part parameter.
	id:= c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "path parameter not provided",
		})
		return
	}

	// get the instance  
	var instance User
	instance_result := config.DB.First(&instance, id)
	if instance_result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "instance does not exist",
		})
		return
	}

	// get the user and confirm permission
	userId ,userexists := c.Get("userID")
	if !userexists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with the authorization token",
		})
		return
	}

	var user User
	user_result := config.DB.First(&user, userId)
	if user_result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Problem with authorization token",
		})
		return
	}

	switch user.Role {
	case 1:
		if instance.Role == 2 || instance.Role ==3 || instance.Role == 4{
			instance.IsActive = false
			result:= config.DB.Save(&instance)
			if result.Error !=nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "Unable to update the user record",
				})
				return
			}
			c.JSON(200, gin.H{"message": "User activated successfully"})
			return

		}

	case 2:
		// Get user MDA ID
		var userMdaId uint 
		err := config.DB.Table("mdas").
		Where("user_id = ?", user.ID).
		Pluck("id",&userMdaId).Error
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "MdaAccount has issues",
			})
			return
		}
		//if the instance is an stc
		if instance.Role == 3 {
			//get the related mda id
			var insanceMdaId uint 
			err := config.DB.Table("stcs").Pluck("mda_id",&insanceMdaId).Error
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "MdaAccount has issues",
				})
				return
			}
			// compare if the instance is an stc under the user
			if (insanceMdaId == userMdaId) {
				instance.IsActive = true
				result:= config.DB.Save(&instance)
				if result.Error !=nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"message": "Cannot suspend this user",
					})
					return
				}
				c.JSON(200, gin.H{"message": "User reactivated successfully"})
				return
			} else {	//if not return unauthorized
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Cannot suspend this user",
				})
				return
			}
		} else if (instance.Role ==4){	// if the instance is a student
			// get the related data to determine if the student is an mda student or stc student
			var instanceData struct{
				MdaId uint
				StcId uint
			}
			err := config.DB.Table("students").
			Select("mda_id, stc_id").
			Where("user_id = ?", instance.ID).
			Scan(&instanceData).Error
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "cannot get the instance data",
				})
				return
			}
			fmt.Println(instanceData)
			if instanceData.MdaId != 0 {	// if this is the student is under an mda directly - Mda student
				if userMdaId != instanceData.MdaId { // check if this student is under the authenticated mda
						c.JSON(http.StatusUnauthorized, gin.H{
							"message": "cannot get the instance data",
						})
						return
					}
				instance.IsActive = true	// deactivate the user
				result:= config.DB.Save(&instance)
				if result.Error !=nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"message": "Unable to update the user record",
					})
					return
				}
				c.JSON(200, gin.H{"message": "User reactivated successfully"})
				return
			} else if instanceData.StcId != 0 {	// if the student is an stc student
				var insatnceStcMdaId uint	// get the related mda id 

				err := config.DB.Table("stcs").
				Where("id = ?", instanceData.StcId).
				Select("mda_id").
				Scan(&insatnceStcMdaId).Error
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"message": "Unable to update the user record",
					})
					return
				}

				if insatnceStcMdaId != userMdaId {	// check if the student is under the authenticated mda
					c.JSON(http.StatusUnauthorized, gin.H{
						"message": "Unable to suspend this user",
					})
					return
				}
				instance.IsActive = true	//deactivate the student
				result:= config.DB.Save(&instance)
				if result.Error !=nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "Unable to update the user record",
					})
					return
				}
				c.JSON(200, gin.H{"message": "User reactivated successfully"})
				return
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Cannot suspend this user",
				})
				return
			}

		}

	case 3:
		//get user stcid
		var userStcId uint
		err := config.DB.
				Table("stcs").
				Select("id").
				Where("user_id = ?", user.ID).
				Scan(&userStcId).Error
		if err != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error trying to suspend this user",
			})
			return
		}

		if instance.Role != 4{
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Cannot suspend this user",
			})
			return
		}

		// get student stc id
		var studentStcId uint
		nerr := config.DB.
				Table("students").
				Select("stc_id").
				Where("user_id = ?", instance.ID).
				Scan(&studentStcId).Error

		if nerr != nil{
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Error trying to suspend this user",
			})
			return
		}

		//check if student is under this stc
		if (studentStcId != userStcId) {
			c.JSON(http.StatusUnauthorized,gin.H{
				"message": "Unable to suspend this user",
			})
			return
		}

		instance.IsActive = true	//deactivate the student
		result:= config.DB.Save(&instance)
		if result.Error !=nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to update the user record",
			})
			return
		}
		c.JSON(200, gin.H{"message": "User reactivated successfully"})
		return
	
	default:
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unable to update the user record",
		})
		return

	}

	}
	
// Login 
func Login(c *gin.Context) {
	var user User
	// receive the request body
	if c.Bind(&LoginSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect request schema",
		})
		return
	}

	// get the user record
	config.DB.Where("email= ?", LoginSchema.Email).First(&user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect email details",
		})
		return
	}

	// check if the user is active
	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Suspended User",
		})
		return
	}

	// compare the passwords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(LoginSchema.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect login details",
		})
		return
	}

	// generate the jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Email,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(config.GetHashSecret()))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to authenticate this user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jwt":        tokenString,
		"token_type": "Bearer",
		"expires_in": 86400, //time in seconds
		"message": "succesful login",
		"role":user.Role,
	})

	fmt.Println("Stored Password: ", user.Password)
    fmt.Println("Input Password: ", LoginSchema.Password)
}

// Request Otp 
func RequestOtp(c *gin.Context) {
	// bind request schema 
	if c.Bind(&RequestOtpSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect request schema",
		})
		return
	}


	// check if email is registered and fetch user data
	var user User
	config.DB.First(&user, "email = ?", RequestOtpSchema.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email is not in the record",
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "This account has been suspended",
		})
		return
	}

	// generate otp and expiry time 
	user.OTP = utilities.GenerateOtp()
	user.OTPExpiresAt = time.Now().Add(time.Minute*5)
	result:= config.DB.Save(&user)
	if result.Error !=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update the user record",
		})
		return
	}


	// send otp -remember this must be changed to mail
	c.JSON(http.StatusOK, gin.H{
		"message": "Otp generated succesfully",
		"otp": user.OTP,
	})
}

// change password
func VerifyOtp(c *gin.Context) {
	// bind the schema
	if c.Bind(&VerifyOtpSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect request schema",
		})
		return
	}

	// get the user and verify otp
	var user User
	config.DB.First(&user, "email = ?", VerifyOtpSchema.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email is not in the record",
		})
		return
	}

	if user.OTP != VerifyOtpSchema.Otp {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect Otp",
		})
		return
	}

	user.OTPVerified = true
	result:= config.DB.Save(&user)
	if result.Error !=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update the user record",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Otp verified succesfully",
	})

}


func ChangePassword(c *gin.Context) {
	// bind the schema
	if c.Bind(&ChangePasswordSchema) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect request schema",
		})
		return
	}

	// get the user and verify otp
	var user User
	config.DB.First(&user, "email = ?", ChangePasswordSchema.Email)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Email is not in the record",
		})
		return
	}

	//check otp verification
	if !user.OTPVerified || time.Now().After(user.OTPExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad otp",
		})
		return
	}


	// change and hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(ChangePasswordSchema.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user.Password = string(hash)
	user.OTPVerified = false
	user.OTPExpiresAt = time.Now()

	result:= config.DB.Save(&user)
	if result.Error !=nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to update the user record",
		})
		return
	}

	// return success
	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed succesfully",
	})
}

func CreateMdaUser(tx *gorm.DB, email, password string) (bool,string,uint){
	var userCheck User
	tx.Where("email= ?", email).First(&userCheck) // Use tx for checking email
	if userCheck.ID != 0 {
	  return false, "user already exists", 0
	}
  
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
	  return false, "Unable to hash password", 0
	}
  
	// Setup the user create instance
	user := User{
	  Email: email,
	  Password: string(hash),
	  OTPExpiresAt: time.Now(),
	  Role: 2,
	  IsActive: true,
	}
  
	result := tx.Create(&user) // Use tx for user creation
	if result.Error != nil {
	  return false, "failed to create user", 0
	}
  
	return true, "user created succesfully", user.ID
}

func CreateStcUser(tx *gorm.DB, email, password string) (bool,string,uint){
	var userCheck User
	tx.Where("email= ?", email).First(&userCheck) // Use tx for checking email
	if userCheck.ID != 0 {
	  return false, "user already exists", 0
	}
  
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
	  return false, "Unable to hash password", 0
	}
  
	// Setup the user create instance
	user := User{
	  Email: email,
	  Password: string(hash),
	  OTPExpiresAt: time.Now(),
	  Role: 3,
	  IsActive: true,
	}
  
	result := tx.Create(&user) // Use tx for user creation
	if result.Error != nil {
	  return false, "failed to create user", 0
	}
  
	return true, "user created succesfully", user.ID
}

func CreateStudentUser(tx *gorm.DB, email, password string) (bool, string, uint) {
	var userCheck User
	tx.Where("email= ?", email).First(&userCheck) // Use tx for checking email
	if userCheck.ID != 0 {
	  return false, "user already exists", 0
	}
  
	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
	  return false, "Unable to hash password", 0
	}

  
	// Setup the user create instance
	user := User{
	  Email: email,
	  Password: string(hash),
	  OTPExpiresAt: time.Now(),
	  Role: 4,
	  IsActive: true,
	}
  
	result := tx.Create(&user) // Use tx for user creation
	if result.Error != nil {
	  return false, "failed to create user", 0
	}
  
	return true, "user created succesfully", user.ID
  }


//   func  CreateEmployerUser(tx *gorm.DB, email, password string )(bool, string, uint){
// 	var userCheck User
// 	tx.Where("email=?", email).First(&userCheck)
// 	if userCheck.ID != 0 {
// 		return false, "user already exist", 0
// 	}

// 	// hash the password 
// 	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
// 	if err != nil {
// 		return false, "Unable to hash password", 0
// 	}


// 	  user := User{
// 		Email: email,
// 		Password: string(hash),
// 		OTPExpiresAt: time.Now(),
// 		Role: 5,
// 		IsActive: true,
// 	  }


// 	  result := tx.Create(&user)
// 	  if result.Error != nil {
// 		return false, "failed to create user", 0
// 	  }

// 	  return true, "user created successfully", user.ID

//   }




func CreateEmployerUser(tx *gorm.DB, email, password string) (bool, string, uint) {
    var userCheck User
    tx.Where("email = ?", email).First(&userCheck)
    if userCheck.ID != 0 {
        return false, "user already exists", 0
    }

    // Hash the password
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
    if err != nil {
        return false, "Unable to hash password", 0
    }

    user := User{
        Email:        email,
        Password:     string(hash),
        OTPExpiresAt: time.Now(),
        Role:         5,
        IsActive:     true,
    }

    result := tx.Create(&user)
    if result.Error != nil {
        return false, "failed to create user", 0
    }

    return true, "user created successfully", user.ID
}