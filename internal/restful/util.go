/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package restful

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yakumioto/alkaid/internal/errors"
)

func render(ctx *gin.Context, property string, obj interface{}) *gin.Context {
	statusCode := http.StatusOK

	if obj, ok := obj.(errors.Error); ok {
		statusCode = obj.StatusCode
	}

	switch property {
	case "xml":
		ctx.XML(statusCode, obj)
	case "yaml":
		ctx.YAML(statusCode, obj)
	case "json":
		ctx.JSON(statusCode, obj)
	case "javascript":
		ctx.JSONP(statusCode, obj)
	default:
		ctx.JSON(statusCode, obj)
	}

	return ctx
}

func getVersionAndProperty(ctx *gin.Context) (string, string) {
	accept := ctx.GetHeader("Accept")
	data := strings.SplitN(strings.SplitN(accept, "/", 2)[1], ".", 2)

	return data[0], data[1]
}
