/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakumioto/alkaid/internal/common/util"
	"github.com/yakumioto/alkaid/internal/services/users"
	"github.com/yakumioto/alkaid/internal/versions"
)

type CreateUser struct{}

func (c *CreateUser) Name() string {
	return "create_user"
}

func (c *CreateUser) Path() string {
	return "/users"
}

func (c *CreateUser) Method() string {
	return http.MethodPost
}

func (c *CreateUser) HandlerFuncChain() []gin.HandlerFunc {
	handler := func(ctx *gin.Context) {
		property := ctx.GetString("Property")

		req := new(users.CreateRequest)
		if err := ctx.ShouldBindJSON(req); err != nil {
			util.Render(ctx, property, err).Abort()
			return
		}

		user, err := users.Create(req)
		if err != nil {
			util.Render(ctx, property, err).Abort()
			return
		}

		util.Render(ctx, property, user)
	}

	return []gin.HandlerFunc{
		func(ctx *gin.Context) {
			if ctx.GetString("Version") != versions.V1 {
				ctx.Next()
			}

			handler(ctx)
		},
		func(ctx *gin.Context) {
			handler(ctx)
		},
	}
}
