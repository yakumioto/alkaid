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

type createUserController struct{}

func (c *createUserController) Name() string {
	return "create_user_url"
}

func (c *createUserController) Path() string {
	return "/users"
}

func (c *createUserController) Method() string {
	return http.MethodPost
}

func (c *createUserController) HandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		version, property := getVersionAndProperty(ctx)

		req := new(users.CreateRequest)
		if err := ctx.ShouldBindJSON(req); err != nil {
			render(ctx, property, err).Abort()
			return
		}

		switch version {
		case "v1":
			fallthrough
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
