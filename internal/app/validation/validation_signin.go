package validation

import (
	"regexp"
	"unicode"

	"forum/internal/model"
)

func ValidationFormSignIn(email, psw string) error {
	if email == "" || psw == "" {
		return model.ErrMessageInvalid
	}
	if !checkPatternEmailForSignIn(email) {
		return model.ErrMessageInvalid
	}
	if !isValidPasswordForSignIn(psw) {
		return model.ErrMessageInvalid
	}
	return nil
}

func checkPatternEmailForSignIn(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(email) {
		return true
	}
	return false
}

func isValidPasswordForSignIn(psw string) bool {
	var (
		hasMinLen bool
		hasUpper  bool
		hasLower  bool
		hasNumber bool
	)
	if len(psw) >= 7 {
		hasMinLen = true
	}
	for _, char := range psw {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber
}
