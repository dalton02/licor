package validator

import (
	"net/mail"
	"reflect"
	"strconv"
	"time"
	"unicode"
)

func isEmail(value reflect.Value) string {
	val, test := value.Interface().(string)

	if !test {
		return "must be a valid email"
	}
	_, err := mail.ParseAddress(val)
	if err != nil {
		return "must be a valid email"
	}
	return ""
}
func isDate(value reflect.Value) string {
	val, test := value.Interface().(string)
	if !test {
		return "must be a valid date"
	}
	layout := "02-01-2006"
	_, err := time.Parse(layout, val)
	if err != nil {
		return "must be a valid date"
	}
	return ""
}
func isNumericString(value reflect.Value) string {
	val, test := value.Interface().(string)
	if !test {
		return "a numeric string"
	}
	_, err := strconv.Atoi(val)
	if err != nil {
		return "a numeric string"
	}
	return ""
}

func isBooleanString(value reflect.Value) string {
	val, test := value.Interface().(string)
	if !test {
		return "must be a true or false string"
	}
	_, err := strconv.Atoi(val)
	if err != nil {
		return "must be a true or false string"
	}
	return ""
}
func isStrongPassword(value reflect.Value) string {
	val, test := value.Interface().(string)
	if !test {
		return "must be a string"

	}

	hasNumber := false
	hasUpper := false
	hasSpecial := false
	hasLower := false

	for _, c := range val {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		case unicode.IsLetter(c) || c == ' ':
		default:
			//return false, false, false, false
		}
	}

	if !hasNumber || !hasSpecial || !hasUpper || !hasLower {
		return "must be a strong password"
	}

	return ""

}
