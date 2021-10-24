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
	"github.com/yakumioto/alkaid/internal/common/log"
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
		// XXX: 是否需要基于特定的供应商实现 Accept 的解析
		// eg：application/vnd.[alkaid][.version][+json]
		// 感觉不是太有必要

		mimeType := strings.SplitN(ctx.GetHeader("Accept"), "/", 2)[1]
		data := strings.SplitN(mimeType, "+", 2)

		version := strings.SplitN(data[0], ".", 3)[2]
		property := data[1]

		ctx.Set("Version", version)
		ctx.Set("Property", property)

		log.Debugf("resolve version middleware: version is [%v], property is [%v]", version, property)
		ctx.Next()
	}
}
