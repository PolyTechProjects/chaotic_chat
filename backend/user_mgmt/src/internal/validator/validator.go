package validator

import (
	"fmt"
	"strings"
)

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

func ValidateUrlTag(urlTag string) error {
	if strings.TrimSpace(urlTag) == "" && len(urlTag) != 0 {
		return fmt.Errorf("urlTag %s is blank", urlTag)
	}
	if len(urlTag) > 52 {
		return fmt.Errorf("urlTag %s is too short or too long", urlTag)
	}
	for _, c := range urlTag {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_') {
			return fmt.Errorf("urlTag %s contains forbidden characters", urlTag)
		}
	}
	return nil
}

func ValidateDescription(description string) error {
	if strings.TrimSpace(description) == "" && len(description) != 0 {
		return fmt.Errorf("description %s is blank", description)
	}
	if len(description) > 500 {
		return fmt.Errorf("description %s is too short or too long", description)
	}
	for _, c := range description {
		if c <= 0x1F {
			return fmt.Errorf("description %s contains forbidden characters", description)
		}
	}
	return nil
}
