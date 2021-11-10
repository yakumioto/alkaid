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
	"github.com/spf13/viper"
	"github.com/yakumioto/alkaid/internal/services/users"
)

var (
	once    sync.Once
	secret  []byte
	expires time.Duration
)

func Initialize() {
	once.Do(func() {
		secret = []byte(viper.GetString("auth.jwt.secret"))
		expires = viper.GetDuration("auth.jwt.expires")
	})
}

func NewTokenWithUser(user *users.User, now int64) (string, error) {
	ctx := users.NewUserContext(user)
	ctx.SetExpiresAt(now + int64(expires.Seconds()))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, ctx)
	return token.SignedString(secret)
}

func VerifyTokenWithUser(tokenString string) (*users.UserContext, error) {
	token, err := jwt.ParseWithClaims(tokenString, &users.UserContext{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

	return token.Claims.(*users.UserContext), err
}
