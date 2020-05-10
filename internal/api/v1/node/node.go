/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package node

import (
	"github.com/gin-gonic/gin"
	"github.com/yakumioto/glog"
)

const (
	nodeID     = "node_id"
	nodeDetail = "/:" + nodeID
)

var (
	logger *glog.Logger
)

type Service struct{}

func (s *Service) Init(log *glog.Logger, rg *gin.RouterGroup) {
	logger = log.MustGetLogger("node")

	r := rg.Group("/node")
	r.POST("")
	r.GET("")
	r.GET(nodeDetail)
	r.PATCH(nodeDetail)
	r.DELETE(nodeDetail)

	logger.Infof("Service initialization success.")
}
