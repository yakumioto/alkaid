/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yakumioto/alkaid/internal/common/jwt"
	"github.com/yakumioto/alkaid/internal/common/log"
	"github.com/yakumioto/alkaid/internal/errors"
	"github.com/yakumioto/alkaid/internal/restful"
)

var (
	logger = log.GetPackageLogger("middlewares.auth")
)

type Auth struct {
}

func (a *Auth) Name() string {
	return "Auth"
}

func (a *Auth) Sequence() int {
	return 4
}

func (a *Auth) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := restful.NewContext(c)
		authorization := c.GetHeader("Authorization")
		contents := strings.SplitN(authorization, " ", 2)
		if len(contents) != 2 {
			logger.Errorf("Authorization header incorrect format")
			ctx.Render(errors.NewError(http.StatusUnauthorized, errors.ErrUnauthorized,
				"Authorization header incorrect format"))
			c.Abort()
			return
		}

		typ := contents[0]
		credentials := contents[1]

		switch typ {
		case "Bearer":
			userCtx, err := jwt.VerifyTokenWithUser(credentials)
			if err != nil {
				logger.Errorf("JWT verify error: %v", err)
				ctx.Render(errors.NewError(http.StatusForbidden, errors.ErrUnauthorized,
					"jwt verification failed"))
				c.Abort()
				return
			}
			c.Set("UserContext", userCtx)
		}

		c.Next()
	}
}
