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
	"github.com/yakumioto/alkaid/internal/versions"
)

type Health struct{}

func (h *Health) Name() string {
	return "check_health"
}

func (h *Health) Path() string {
	return "/health"
}

func (h *Health) Method() string {
	return http.MethodGet
}

func (h *Health) HandlerFuncChain() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		func(ctx *gin.Context) {
			if ctx.GetString("Version") != versions.V1 {
				ctx.Next()
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"version": versions.V1,
			})

			ctx.Abort()
		},
		func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{
				"status": "ok",
			})
		},
	}
}
