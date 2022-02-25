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
	"github.com/yakumioto/alkaid/internal/services/users"
	"github.com/yakumioto/alkaid/internal/versions"
)

type CreateUser struct {
	Controllers
}

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
		req := new(users.CreateRequest)
		if err := ctx.ShouldBindJSON(req); err != nil {
			c.Render(ctx, err).Abort()
			return
		}

		user := new(users.User)

		userCtx, _ := ctx.Get("UserContext")
		if err := user.Create(req, userCtx.(*users.UserContext)); err != nil {
			c.Render(ctx, err).Abort()
			return
		}

		c.Render(ctx, user)
	}

	return []gin.HandlerFunc{
		func(ctx *gin.Context) {
			if !c.MatchVersion(ctx, versions.V1) {
				return
			}

			handler(ctx)
			ctx.Abort()
		},
		func(ctx *gin.Context) {
			handler(ctx)
		},
	}
}

type GetUserDetailByID struct {
	Controllers
}

func (c *GetUserDetailByID) Name() string {
	return "find_user_by_id"
}

func (c *GetUserDetailByID) Path() string {
	return "/users/:id"
}

func (c *GetUserDetailByID) Method() string {
	return http.MethodGet
}

func (c *GetUserDetailByID) HandlerFuncChain() []gin.HandlerFunc {
	handler := func(ctx *gin.Context) {
		id := ctx.Param("id")

		user := new(users.User)
		if err := user.GetDetailByID(id); err != nil {
			c.Render(ctx, err).Abort()
			return
		}

		c.Render(ctx, user)
	}

	return []gin.HandlerFunc{
		func(ctx *gin.Context) {
			if !c.MatchVersion(ctx, versions.V1) {
				return
			}

			handler(ctx)
			ctx.Abort()
		},
		func(ctx *gin.Context) {
			handler(ctx)
		},
	}
}
