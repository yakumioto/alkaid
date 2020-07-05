/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package services

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service interface {
	RootPath() string
	DetailID() string
	Post(ctx *gin.Context)
	Patch(ctx *gin.Context)
	Get(ctx *gin.Context)
	GetList(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

func NewServices(r *gin.Engine, servs ...Service) {
	for _, srv := range servs {
		relativePath := fmt.Sprintf("/:%s", srv.DetailID())

		rg := r.Group(srv.RootPath())
		rg.Handle(http.MethodPost, "/", srv.Post)
		rg.Handle(http.MethodGet, "/", srv.GetList)
		rg.Handle(http.MethodPatch, relativePath, srv.Patch)
		rg.Handle(http.MethodGet, relativePath, srv.Get)
		rg.Handle(http.MethodDelete, relativePath, srv.Delete)
	}
}
