/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package jwt

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/yakumioto/alkaid/internal/common/log"
	"github.com/yakumioto/alkaid/internal/services/users"
)

var (
	logger = log.GetPackageLogger("common.jwt")
	once   sync.Once
	t      *JWT
)

func Initialize(secret string, expires time.Duration) {
	logger.Infof("jwt woken expires is %v", expires)
	once.Do(func() {
		if t == nil {
			t = NewJWT(secret, expires)
		}
	})
}

func NewTokenWithUser(user *users.User, now int64) (string, error) {
	return t.NewTokenWithUser(user, now)
}

func VerifyTokenWithUser(tokenString string) (*users.UserContext, error) {
	return t.VerifyTokenWithUser(tokenString)
}

type JWT struct {
	secret  []byte
	expires time.Duration
}

func NewJWT(secret string, expires time.Duration) *JWT {
	return &JWT{
		secret:  []byte(secret),
		expires: expires,
	}
}

func (t *JWT) NewTokenWithUser(user *users.User, now int64) (string, error) {
	ctx := users.NewUserContext(user)
	ctx.SetExpiresAt(now + int64(t.expires.Seconds()))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ctx)
	return token.SignedString(t.secret)
}

func (t *JWT) VerifyTokenWithUser(tokenString string) (*users.UserContext, error) {
	token, err := jwt.ParseWithClaims(tokenString, &users.UserContext{},
		func(token *jwt.Token) (interface{}, error) {
			return t.secret, nil
		})

	return token.Claims.(*users.UserContext), err
}
