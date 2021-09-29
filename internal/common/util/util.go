package util

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/lithammer/shortuuid"
	"golang.org/x/crypto/pbkdf2"
)

func HashPassword(password, salt string, iter int) string {
	passwordHash := pbkdf2.Key([]byte(password), []byte(salt), iter, 32, sha256.New)
	passwordHashBase64 := base64.StdEncoding.EncodeToString(passwordHash)
	return fmt.Sprintf("%d.%s", iter, passwordHashBase64)
}

// PBKDF2WithSha256 用于生成扩展密钥
func PBKDF2WithSha256(password, salt []byte, keyLen int) []byte {
	return pbkdf2.Key(password, salt, 1, keyLen, sha256.New)
}

func ValidatePassword(password, salt, passwordHash string) bool {
	iter, err := strconv.Atoi(strings.SplitN(passwordHash, ".", 2)[0])
	if err != nil {
		return false
	}

	if passwordHash != HashPassword(password, salt, iter) {
		return false
	}

	return true
}

func GenResourceID(namespace string) string {
	return fmt.Sprintf("%s-%s", namespace, shortuuid.New())
}

func GetVarintBytesWithMaxVarintLen64(x int64) []byte {
	buf := make([]byte, 10)
	binary.PutVarint(buf, x)
	return buf
}

func MustVarint(buf []byte) int64 {
	num, _ := binary.Varint(buf)
	return num
}
