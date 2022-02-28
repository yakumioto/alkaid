/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
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
	"github.com/yakumioto/alkaid/internal/services/systems"
	"github.com/yakumioto/alkaid/internal/versions"
)

type InitializeSystem struct {
}

func (i *InitializeSystem) Name() string {
	return "initialize_system"
}

func (i *InitializeSystem) Path() string {
	return "/initialize"
}

func (i *InitializeSystem) Method() string {
	return http.MethodPost
}

func (i *InitializeSystem) HandlerFuncChain() []gin.HandlerFunc {
	handler := func(ctx *restful.Context) {
		req := new(systems.InitRequest)
		if err := ctx.ShouldBindJSON(req); err != nil {
			ctx.Render(err).Abort()
			return
		}

		sys, err := systems.SystemInit(req)
		if err != nil {
			ctx.Render(err).Abort()
			return
		}

		ctx.Render(sys)
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
