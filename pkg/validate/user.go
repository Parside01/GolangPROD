package validatate

import (
	"errors"
	"regexp"
	"solution/models"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

var (
	ErrNotValidLogin     = errors.New("invalid login")
	ErrNotValidEmail     = errors.New("invalid email")
	ErrNotValidPhone     = errors.New("invalid phone")
	ErrNotValidPass      = errors.New("invalid password")
	ErrNotValidCode      = errors.New("invalid country code")
	ErrNotValidImageLink = errors.New("invalid image link")
)

func IsValidUser(user *models.User) error {
	if ok := IsValidLogin(user.Login); !ok {
		return ErrNotValidLogin
	}
	if ok := IsValidEmail(user.Email); !ok {
		return ErrNotValidEmail
	}
	//TODO: доделать норм валидацию номера телефона
	if ok := IsValidPhone(user.Phone); !ok {
		return ErrNotValidPhone
	}
	if ok := IsValidImageLink(user.Image); !ok {
		return ErrNotValidImageLink
	}
	if ok := IsValidAlpha2(user.CountryCode); !ok {
		return ErrNotValidCode
	}
	return nil
}
func IsValidImageLink(link string) bool {
	return (validation.Validate(link, validation.Required, validation.Length(-1, 200)) == nil) || (link == "")
}

func IsValidLogin(login string) bool {
	return validation.Validate(login, validation.Required, validation.Length(4, 30), validation.Match(regexp.MustCompile("[a-zA-Z0-9]+"))) == nil
}

// TODO: Add more validation
func IsValidPassword(pass string) bool {
	return validation.Validate(pass, validation.Required, validation.Length(6, 32), validation.Match(regexp.MustCompile("[a-zA-Z0-9]+"))) == nil
}

func IsValidEmail(email string) bool {
	return validation.Validate(email, validation.Required, is.Email) == nil
}

func IsValidPhone(phone string) bool {
	return validation.Validate(phone, validation.Required, validation.Match(regexp.MustCompile(`[0-9]+`))) == nil
}
