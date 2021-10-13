/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package middlewares

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ResolveVersion struct {
}

func (r *ResolveVersion) Name() string {
	return "ResolveVersion"
}

func (r *ResolveVersion) Sequence() int {
	return 3
}

func (r *ResolveVersion) HandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		mimeType := strings.SplitN(ctx.GetHeader("Accept"), "/", 2)[1]
		data := strings.SplitN(mimeType, "+", 2)

		version := strings.SplitN(data[0], ".", 3)[2]
		property := data[1]

		ctx.Set("Version", version)
		ctx.Set("Property", property)

		logrus.Debugf("resolve version middleware: version is [%v], property is [%v]", version, property)
		ctx.Next()
	}
}
