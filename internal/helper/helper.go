// Package helper provides helper function to simplyfy repeating or predictable tasks
package helper

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/sweet-go/stdlib/helper"
)

// GenerateUploadedFilename generates a unique filename for uploaded file
func GenerateUploadedFilename(name string) string {
	return base64.StdEncoding.EncodeToString([]byte(helper.GenerateUniqueName() + strings.ReplaceAll(name, ".", "")))
}

// GenerateEncryptedFilename generates a unique filename for encrypted file
func GenerateEncryptedFilename(name string) string {
	return name + "_enc"
}

// GenerateDecryptedFilename generates a unique filename for decrypted file
func GenerateDecryptedFilename(name string) string {
	return name + "_dec"
}

// IsFileExists checks if a file exists
func IsFileExists(filename string) bool {
	c, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !c.IsDir()
}
