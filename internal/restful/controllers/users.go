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
	"github.com/yakumioto/alkaid/internal/restful"
	"github.com/yakumioto/alkaid/internal/services/users"
	"github.com/yakumioto/alkaid/internal/versions"
)

type CreateUser struct {
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
	handler := func(ctx *restful.Context) {
		req := new(users.CreateRequest)
		if err := ctx.ShouldBindJSON(req); err != nil {
			ctx.Render(err).Abort()
			return
		}

		userCtx, _ := ctx.Get("UserContext")
		user, err := users.Create(req, userCtx.(*users.UserContext))
		if err != nil {
			ctx.Render(err).Abort()
			return
		}

		ctx.Render(user)
	}

	return []gin.HandlerFunc{
		func(c *gin.Context) {
			ctx := restful.NewContext(c)
			if !ctx.MatchVersion(versions.V1) {
				return
			}

			handler(ctx)
			ctx.Abort()
		},
		func(c *gin.Context) {
			ctx := restful.NewContext(c)
			handler(ctx)
		},
	}
}

type GetUserDetailByID struct {
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
	handler := func(ctx *restful.Context) {
		id := ctx.Param("id")

		user, err := users.GetDetailByID(id)
		if err != nil {
			ctx.Render(err).Abort()
			return
		}

		ctx.Render(user)
	}

	return []gin.HandlerFunc{
		func(c *gin.Context) {
			ctx := restful.NewContext(c)
			if !ctx.MatchVersion(versions.V1) {
				return
			}

			handler(ctx)
			ctx.Abort()
		},
		func(c *gin.Context) {
			ctx := restful.NewContext(c)
			handler(ctx)
		},
	}
}
