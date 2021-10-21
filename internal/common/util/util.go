/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package util

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	stdErrors "errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lithammer/shortuuid"
	"github.com/yakumioto/alkaid/internal/errors"
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

func Render(ctx *gin.Context, property string, obj interface{}) *gin.Context {
	statusCode := http.StatusOK

	if obj, ok := obj.(error); ok {
		err := new(errors.Error)
		if stdErrors.As(obj, &err) {
			statusCode = err.StatusCode
		}
	}

	switch property {
	case "xml":
		ctx.XML(statusCode, obj)
	case "yaml":
		ctx.YAML(statusCode, obj)
	case "json":
		ctx.JSON(statusCode, obj)
	case "javascript":
		ctx.JSONP(statusCode, obj)
	default:
		ctx.JSON(statusCode, obj)
	}

	return ctx
}
