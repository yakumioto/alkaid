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

type Logger struct {
}

func (l *Logger) Name() string {
	return "Logger"
}

func (l *Logger) Sequence() int {
	return 1
}

func (l *Logger) HandlerFunc() gin.HandlerFunc {
	return gin.Logger()
}
