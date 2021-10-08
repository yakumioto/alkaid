/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package restful

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakumioto/alkaid/internal/services/users"
)

type userController struct{}

func (u *userController) Name() string {
	return "create_user_url"
}

func (u *userController) Path() string {
	return "/users"
}

func (u *userController) Method() string {
	return http.MethodPost
}

func (u *userController) HandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		version, property := getVersionAndProperty(ctx)

		req := new(users.CreateRequest)
		if err := ctx.ShouldBindJSON(req); err != nil {
			render(ctx, property, err).Abort()
			return
		}

		switch version {
		default:
			user, err := users.CreateUser(req)
			if err != nil {
				render(ctx, property, err).Abort()
				return
			}

			render(ctx, property, user)
		}
	}
}
