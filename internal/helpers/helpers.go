// ChatGPT code so w?
package helpers

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
	"unicode"

	http "github.com/bogdanfinn/fhttp"
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

func WriteCookiesToFile(email, password string, cookies []*http.Cookie) error {
	outputDir := "data/output"
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Open the file in append mode, or create it if it doesn't exist
	file, err := os.OpenFile(fmt.Sprintf("%s/created.json", outputDir), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Structure for saving data
	data := struct {
		Email    string         `json:"email"`
		Password string         `json:"password"`
		Cookies  []*http.Cookie `json:"cookies"`
	}{
		Email:    email,
		Password: password,
		Cookies:  cookies,
	}

	// Write the data in JSON format
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("error encoding data to JSON: %w", err)
	}

	return nil
}

func GetCt0(length int) (string, error) {
	// Calculate how many bytes are needed to achieve the desired length in hex
	numBytes := length / 2
	if length%2 != 0 {
		numBytes++
	}

	// Generate random bytes
	bytes := make([]byte, numBytes)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Return the hexadecimal representation of the bytes
	return hex.EncodeToString(bytes)[:length], nil
}
