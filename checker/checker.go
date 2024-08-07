package checker

import "regexp"

const (
	phoneRegex = `^[0][8]\d{8,14}$`
	emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

type AllowedDataType interface {
	int64 | string | float64
}

// Check wether the slice contains the input x
func ArrayInt64Contains(slice []int64, x int64) bool {
	for _, i := range slice {
		if i == x {
			return true
		}
	}
	return false
}

// Check wether the slice in int64, string, or float data type contains the input x
func ArrayContains[T AllowedDataType](slice []T, x T) bool {
	for _, i := range slice {
		if i == x {
			return true
		}
	}
	return false
}

// Remove duplicate values in int64, string, or float data type slice
func ArrayDeduplicate[T AllowedDataType](values []T) []T {
	var list []T
	var mapIsDuplicate = make(map[T]bool)
	for _, value := range values {
		if !mapIsDuplicate[value] {
			mapIsDuplicate[value] = true
			list = append(list, value)
		}
	}
	return list
}

func IsPhoneNumber(phone string) bool {
	phoneRegex := regexp.MustCompile(phoneRegex)
	return phoneRegex.MatchString(phone)
}

func IsEmail(email string) bool {
	emailRegex := regexp.MustCompile(emailRegex)
	return emailRegex.MatchString(email)
}
