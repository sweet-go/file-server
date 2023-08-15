package helper

import (
	"encoding/base64"
	"os"
	"strings"

	"github.com/sweet-go/stdlib/helper"
)

func GenerateUploadedFilename(name string) string {
	return base64.StdEncoding.EncodeToString([]byte(helper.GenerateUniqueName() + strings.ReplaceAll(name, ".", "")))
}

func GenerateEncryptedFilename(name string) string {
	return name + "_enc"
}

func GenerateDecryptedFilename(name string) string {
	return name + "_dec"
}

func IsFileExists(filename string) bool {
	c, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !c.IsDir()
}
