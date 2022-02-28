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
	"github.com/yakumioto/alkaid/internal/versions"
)

type Health struct {
}

func (c *Health) Name() string {
	return "check_health"
}

func (c *Health) Path() string {
	return "/health"
}

func (c *Health) Method() string {
	return http.MethodGet
}

func (c *Health) HandlerFuncChain() []gin.HandlerFunc {
	handler := func(ctx *restful.Context) {
		ctx.Render(gin.H{
			"status":  "ok",
			"version": versions.V1,
		})
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
