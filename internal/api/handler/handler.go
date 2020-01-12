/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yakumioto/glog"
)

var (
	logger *glog.Logger
)

func Init() {
	logger = glog.MustGetLogger("handler")
}

func returnInternalServerError(ctx *gin.Context, format string, v ...interface{}) {
	logger.Errof(format, v...)
	ctx.Status(http.StatusInternalServerError)
}
