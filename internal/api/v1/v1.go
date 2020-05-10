/*
 *  Copyright 2020. The Alkaid Authors. All rights reserved.
 *  Use of this source code is governed by a MIT-style
 *  license that can be found in the LICENSE file.
 *  Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/yakumioto/glog"
)

var (
	logger *glog.Logger
)

type Servicer interface {
	Init(*glog.Logger, *gin.RouterGroup)
}

func Init(e *gin.Engine, svcs ...Servicer) {
	logger = glog.MustGetLogger("api.v1")

	rg := e.Group("/v1")
	for _, svc := range svcs {
		svc.Init(logger, rg)
	}
}
