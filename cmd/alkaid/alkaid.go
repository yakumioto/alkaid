/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yakumioto/glog"

	"github.com/yakumioto/alkaid/internal/api/routers"
)

func main() {
	r := gin.Default()
	routers.AddRouters(r)

	if err := r.Run(":8080"); err != nil {
		glog.Fatalf("service startup failed: %s", err)
	}
}
