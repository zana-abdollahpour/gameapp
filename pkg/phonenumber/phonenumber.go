package phonenumber

import "regexp"

func IsValid(phoneNumber string) bool {
	pattern := `^(0|\+98)?9\d{9}$`
	matched, _ := regexp.MatchString(pattern, phoneNumber)
	return matched
}
