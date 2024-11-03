package scripts

import (
	"fme_backend/internal/config"
	myuser "fme_backend/internal/user"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func verifyFmeDetails() (string,string,bool) {
	// verify that the email is present
	if config.GetFmeEmail() == "" || config.GetFmePassWord() == "" {
		return "","",false
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(config.GetFmePassWord()), 10)
	if err != nil {
		fmt.Println("Error hashing fme Password")
		return "","",false
	}
	return config.GetFmeEmail(),string(hash),true

}

func CreateFmeAtStart() {
	// check if the fme user has been created

	var fmeUsers []uint
	err := config.DB.Table("users").
			Select("users.id").
			Where("users.role = ?",1).Find(&fmeUsers).Error
	if err != nil {
		fmt.Println("Unable to confirm fme existence or inexistence")
		return
	}
	if len(fmeUsers) != 0 {
		fmt.Println("FmeUser already exists")
		return
	}

	// create the fme user
	fmeEmail,fmePassword,success :=verifyFmeDetails()
	if !success {
		fmt.Println("Error verifying Fme details")
		return
	}

	fmeUser := myuser.User{
		Email: fmeEmail,
		Password: fmePassword,
		Role: 1,
		IsActive: true,
		OTPExpiresAt: time.Now(),
	}
	
	// create the user in the database
	err = config.DB.Create(&fmeUser).Error

	if err != nil {
		fmt.Println("Error creating Fme object in db")
		return
	}	
}