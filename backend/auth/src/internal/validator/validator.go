package validator

import (
	"fmt"
	"strings"

	"github.com/nyaruka/phonenumbers"
)

func ValidateLogin(login string) error {
	phonenumber, err := phonenumbers.Parse(login, "RU")
	if err != nil {
		return err
	}
	if !phonenumbers.IsValidNumber(phonenumber) {
		return fmt.Errorf("invalid phone number %s", phonenumber)
	}
	return nil
}

func ValidateName(name string) error {
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name %s is blank", name)
	}
	if len(name) < 1 || len(name) > 52 {
		return fmt.Errorf("name %s is too short or too long", name)
	}
	for _, c := range name {
		if c <= 0x1F {
			return fmt.Errorf("name %s contains forbidden characters", name)
		}
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 6 || len(password) > 32 {
		return fmt.Errorf("password %s too short or too long", password)
	}
	lowerCase := false
	upperCase := false
	digit := false
	specSymbol := false
	for _, c := range password {
		if c >= 'a' && c <= 'z' {
			lowerCase = true
		} else if c >= 'A' && c <= 'Z' {
			upperCase = true
		} else if c >= '0' && c <= '9' {
			digit = true
		} else if isValidSpecSymbol(c) {
			specSymbol = true
		} else {
			return fmt.Errorf("password %s contains forbidden characters", password)
		}
	}
	if !(lowerCase && upperCase && digit && specSymbol) {
		return fmt.Errorf("invalid password %s", password)
	}
	return nil
}

func isValidSpecSymbol(c rune) bool {
	if c == '~' ||
		c == '`' ||
		c == '!' ||
		c == '@' ||
		c == '#' ||
		c == '$' ||
		c == '%' ||
		c == '^' ||
		c == '&' ||
		c == '*' ||
		c == '(' ||
		c == ')' ||
		c == '-' ||
		c == '_' ||
		c == '+' ||
		c == '=' ||
		c == '[' ||
		c == ']' ||
		c == '{' ||
		c == '}' ||
		c == '|' ||
		c == ',' ||
		c == '.' ||
		c == '?' ||
		c == 0x27 {
		return true
	}
	return false
}
