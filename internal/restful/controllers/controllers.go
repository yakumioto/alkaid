/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package controllers

import "github.com/yakumioto/alkaid/internal/common/log"

var (
	logger = log.GetPackageLogger("restful.controllers")
)

// type Controllers struct{}
//
// func (c *Controllers) RenderFormat(ctx *gin.Context) string {
// 	return ctx.GetString("AcceptFormat")
// }
//
// func (c *Controllers) MatchVersion(ctx *gin.Context, version string) bool {
// 	if ctx.GetString("AcceptVersion") != version {
// 		ctx.Next()
// 		return false
// 	}
//
// 	return true
// }
//
// func (c *Controllers) Render(ctx *gin.Context, obj interface{}) *gin.Context {
// 	return util.Render(ctx, c.RenderFormat(ctx), obj)
// }
