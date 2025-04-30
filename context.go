/*
 * Copyright (c) 2025 yakumioto <yaku.mioto@gmail.com>
 * All rights reserved.
 */

package alkaid

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Context struct {
	*gin.Context
}

func (c *Context) RenderAny(obj any) *gin.Context {
	statusCode := http.StatusOK

	switch v := obj.(type) {
	case nil:
		c.Status(statusCode)
		return c.Context
	case int:
		if v < 0 || v > 999 {
			v = http.StatusInternalServerError
		}
		c.Status(v)
		return c.Context
	case error:
		if v != nil {
			statusCode = http.StatusInternalServerError
		}
	}

	switch format := c.GetHeader("Accept"); format {
	case "application/xml":
		c.XML(statusCode, obj)
	case "application/yaml":
		c.YAML(statusCode, obj)
	case "application/javascript":
		c.JSONP(statusCode, obj)
	case "text/plain":
		if data, ok := obj.(string); ok {
			c.String(statusCode, data)
			break
		}
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/svg+xml":
		if data, ok := obj.([]byte); ok {
			c.Data(statusCode, format, data)
			break
		}
	case "audio/mpeg", "audio/mp3", "audio/mpga", "audio/mpg":
		if data, ok := obj.([]byte); ok {
			c.Data(statusCode, format, data)
			break
		}
	case "application/octet-stream":
		if data, ok := obj.([]byte); ok {
			c.Data(statusCode, format, data)
			break
		}
	default:
		c.JSON(statusCode, obj)
	}

	return c.Context
}
