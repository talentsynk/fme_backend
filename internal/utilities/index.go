package utilities

import (
	"errors"
	"fme_backend/internal/config"
	"fme_backend/internal/responses"
	"fme_backend/internal/schemas"
	"fmt"
	"math/rand"
	"net/mail"
	"reflect"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	TOKEN_EXP_TIME = 30
)

func HashPassword(password string) (string, error) {
	raw_bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(raw_bytes), err
}

func CompareHashAndPassword(password, hash string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return errors.New(responses.INVALID_LOGIN)
	}
	return nil

}

func IsEmpty(v any) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return val.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return val.IsNil()
	case reflect.Bool:
		return !val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return val.Interface() == reflect.Zero(val.Type()).Interface()
	default:
		return reflect.DeepEqual(v, reflect.Zero(val.Type()).Interface())
	}
}

func IsEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func IsValidPassword(password string) bool {
	return true
}

func GeneratJWT(email string, user_id, role uint) (schemas.LoginResp, error) {
	role_string := GetRoleString(role)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":     email,
		"role":    role_string,
		"user_id": user_id,
		"exp":     time.Now().Add(time.Hour * TOKEN_EXP_TIME).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.GetHashSecret()))
	if err != nil {
		return schemas.LoginResp{}, err
	}
	return schemas.LoginResp{
		Token: tokenString,
		Role:  role_string,
	}, nil
}

func GetRoleString(role uint) string {
	var role_str string
	switch role {
	case 1:
		role_str = USER_ROLE_FME
	case 2:
		role_str = USER_ROLE_MDA
	case 3:
		role_str = USER_ROLE_STC
	case 4:
		role_str = USER_ROLE_STUDENT
	case 5:
		role_str = USER_ROLE_EMPLOYER
	case 6:
		role_str = USER_ROLE_ARTISAN
	}
	return role_str
}

func GenerateOtp() string {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	otp := rng.Intn(100000)
	otpString := fmt.Sprintf("%05d", otp)

	return otpString
}
