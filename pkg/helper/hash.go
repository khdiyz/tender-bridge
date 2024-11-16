package helper

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"tender-bridge/config"
)

func GenerateHash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password cannot empty")
	}

	salt := config.GetConfig().HashKey
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt))), nil
}
