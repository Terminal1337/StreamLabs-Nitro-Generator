// ChatGPT code so w?
package helpers

import (
	"math/rand"
	"time"
	"unicode"
)

const (
	lowerChars   = "abcdefghijklmnopqrstuvwxyz"
	upperChars   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numberChars  = "0123456789"
	specialChars = "!@#$%^&*()-_=+[]{}|;:,.<>?/`~"
)

// GeneratePassword generates a random password of given length
func GeneratePassword(length int) string {
	// Combine all possible characters
	allChars := lowerChars + upperChars + numberChars + specialChars

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	var password string
	for i := 0; i < length; i++ {
		password += string(allChars[rand.Intn(len(allChars))])
	}

	// Ensure password contains at least one uppercase, one lowercase, one number, and one special character
	if !containsRequiredChars(password) {
		return GeneratePassword(length) // Retry if the password doesn't meet the requirements
	}

	return password
}

// containsRequiredChars checks if the password contains at least one lowercase, one uppercase, one number, and one special character
func containsRequiredChars(password string) bool {
	hasLower, hasUpper, hasNumber, hasSpecial := false, false, false, false

	for _, c := range password {
		switch {
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsDigit(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasNumber && hasSpecial
}
