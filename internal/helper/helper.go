package helper

import (
	"encoding/base64"
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
