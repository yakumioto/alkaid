/*
 * Copyright (c) 2025 yakumioto <yaku.mioto@gmail.com>
 * All rights reserved.
 */

package main

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakumioto/alkaid"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	alk := alkaid.NewAlkaid(
		alkaid.WithEngineMode(gin.DebugMode),
		alkaid.WithEnableUrlVersionPrefix(),
		alkaid.WithVersionHeader("X-Version"),
		alkaid.WithAuthenticator(alkaid.NewNoneNoneAuthenticator()),
	)

	alk.RegisterService(&alkaid.Service{
		Name: "users",
		Controllers: []alkaid.Controller{
			&UserController{},
		},
	})

	_ = alk.Run(":8080")
}

type UserController struct {
}

func (h *UserController) NoAuth() bool {
	return false
}

func (h *UserController) Paths() []*alkaid.PathInfo {
	return []*alkaid.PathInfo{
		{
			Role: "*",
			Path: "/users",
		},
		{
			Role: "*",
			Path: "/organizations/:organization_id/users",
		},
	}
}

func (h *UserController) Method() string {
	return http.MethodGet
}

func (h *UserController) DefaultHandlerFunc() alkaid.HandlerFunc {
	return func(c *alkaid.Context) {
		c.RenderAny(gin.H{
			"message": "OK",
		})
	}
}

func (h *UserController) VersionHandlerFuncs() []alkaid.VersionHandlerFunc {
	return []alkaid.VersionHandlerFunc{
		func() (*alkaid.HandlerVersion, alkaid.HandlerFunc) {
			return alkaid.NewHandlerVersion(1), func(c *alkaid.Context) {
				c.RenderAny(gin.H{
					"message": "v1.0.0",
				})
			}
		},
		func() (*alkaid.HandlerVersion, alkaid.HandlerFunc) {
			return alkaid.NewHandlerVersion(1, 1, 1), func(c *alkaid.Context) {
				c.RenderAny(gin.H{
					"message": "v1.1.1",
				})
			}
		},
		func() (*alkaid.HandlerVersion, alkaid.HandlerFunc) {
			return alkaid.NewHandlerVersion(1, 2, 1), func(c *alkaid.Context) {
				c.RenderAny(gin.H{
					"message": "v1.2.1",
				})
			}
		},
		func() (*alkaid.HandlerVersion, alkaid.HandlerFunc) {
			return alkaid.NewHandlerVersion(2), func(c *alkaid.Context) {
				c.RenderAny(gin.H{
					"message": "v2.0.0",
				})
			}
		},
		func() (*alkaid.HandlerVersion, alkaid.HandlerFunc) {
			return alkaid.NewHandlerVersion(3), func(c *alkaid.Context) {
				c.RenderAny(gin.H{
					"message": "OK",
				})
			}
		},
	}
}
