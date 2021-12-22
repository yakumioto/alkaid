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
		accept := ctx.GetHeader("Accept")
		r.handler(ctx, accept)
		ctx.Next()
	}
}

// handler 处理 accept 的 调用版本以及格式化方式
// case: No Accept
// case: application/json
// case: application/vnd.alkaid[.version][+json,xml,yaml,javascript]
func (r *ResolveVersion) handler(ctx *gin.Context, accept string) {
	var (
		version string
		format  string
	)
	if accept == "" {
		return
	}

	data := strings.SplitN(accept, "/", 2)
	if len(data) != 2 || data[0] != "application" {
		return
	}

	// 处理需要返回的格式
	data = strings.SplitN(data[1], "+", 2)
	if len(data) == 2 {
		format = data[1]
		ctx.Set("AcceptFormat", format)
	}

	data = strings.SplitN(data[0], ".", 3)
	if len(data) == 3 {
		version = data[2]
		ctx.Set("AcceptVersion", version)
	}

	log.Debugf("resolve version middleware: version is [%v], format is [%v]", version, format)
}
