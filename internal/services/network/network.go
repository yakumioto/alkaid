/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package network

import (
	"github.com/gin-gonic/gin"
	"github.com/yakumioto/glog"

	"github.com/yakumioto/alkaid/internal/services"
)

var (
	logger = glog.MustGetLogger("services.network")
)

type Service struct{}

func NewService() services.Service {
	return &Service{}
}

func (o *Service) RootPath() string {
	return "/network"
}

func (o *Service) DetailID() string {
	return "networkID"
}

func (o *Service) Post(ctx *gin.Context) {

}

func (o *Service) Patch(ctx *gin.Context) {

}

func (o *Service) Get(ctx *gin.Context) {

}

func (o *Service) GetList(ctx *gin.Context) {

}

func (o *Service) Delete(ctx *gin.Context) {

}
