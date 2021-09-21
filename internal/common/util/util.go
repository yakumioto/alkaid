package util

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/pbkdf2"
)

func PBKDF2L32Sha256WithBase64(password, salt string, iter int) string {
	passwordHash := pbkdf2.Key([]byte(password), []byte(salt), iter, 32, sha256.New)
	passwordHashBase64 := base64.StdEncoding.EncodeToString(passwordHash)
	return fmt.Sprintf("%d.%s", iter, passwordHashBase64)
}

func ValidatePassword(password, salt, passwordHash string) bool {
	iter, err := strconv.Atoi(strings.SplitN(passwordHash, ".", 2)[0])
	if err != nil {
		return false
	}

	if passwordHash != PBKDF2L32Sha256WithBase64(password, salt, iter) {
		return false
	}

	return true
}

func GenResourceID() string {
	return uuid.NewString()
}
