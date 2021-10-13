/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package middlewares

import (
	"github.com/gin-gonic/gin"
)

type Recovery struct{}

func (r *Recovery) Name() string {
	return "Recovery"
}

func (r *Recovery) Sequence() int {
	return 2
}

func (r *Recovery) HandlerFunc() gin.HandlerFunc {
	return gin.Recovery()
}
