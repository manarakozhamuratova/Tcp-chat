package validation

import (
	"regexp"
	"unicode"

	"forum/internal/model"
)

func ValidationFormSignUp(uname, umail, psw, psw2 string) error {
	if uname == "" || umail == "" || psw == "" || psw2 == "" {
		return model.ErrMessageInvalid
	}
	if psw != psw2 {
		return model.ErrMessageInvalid
	}
	if !checkPatternName(uname) {
		return model.ErrMessageInvalid
	}
	if !checkPatternEmail(umail) {
		return model.ErrMessageInvalid
	}
	if !isValidPassword(psw) {
		return model.ErrMessageInvalid
	}
	return nil
}

func checkPatternName(uname string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9]+$")
	if re.MatchString(uname) {
		return true
	}
	return false
}

func checkPatternEmail(umail string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(umail) {
		return true
	}
	return false
}

func isValidPassword(psw string) bool {
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
