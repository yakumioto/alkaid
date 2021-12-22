/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yakumioto/alkaid/internal/common/log"
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
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		stop := time.Now()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Debugf("Status Code: %d | Time Consuming: %v | Client IP: %s | Request Method: %s | Request Path: %s",
			c.Writer.Status(),
			stop.Sub(start),
			c.ClientIP(),
			c.Request.Method,
			path)

		// todo: c.Error 是否需要使用
	}
}
