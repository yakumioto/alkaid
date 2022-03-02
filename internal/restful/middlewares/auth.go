/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/yakumioto/alkaid/internal/common/jwt"
	"github.com/yakumioto/alkaid/internal/common/log"
	"github.com/yakumioto/alkaid/internal/errors"
	"github.com/yakumioto/alkaid/internal/restful"
	"github.com/yakumioto/alkaid/internal/services/users"
)

var (
	logger = log.GetPackageLogger("middlewares.auth")
)

type Auth struct {
	enforcer *casbin.Enforcer
}

func NewAuth(model, policy string) *Auth {
	enforcer, err := casbin.NewEnforcer(model, policy)
	if err != nil {
		logger.Panicf("new enforcer error: %v", err)
	}

	return &Auth{
		enforcer: enforcer,
	}
}

func (a *Auth) Name() string {
	return "Auth"
}

func (a *Auth) Sequence() int {
	return 4
}

func (a *Auth) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			noAuthorization bool
			userCtx         *users.UserContext
			err             error
		)

		ctx := restful.NewContext(c)
		authorization := c.GetHeader("Authorization")
		contents := strings.SplitN(authorization, " ", 2)
		if len(contents) != 2 {
			noAuthorization = true
		}

		if noAuthorization {
			if ok, _ := a.enforcer.Enforce("*", "*", ctx.FullPath(), ctx.Request.Method); ok {
				ctx.Next()
				return
			}

			ctx.Render(errors.NewError(http.StatusUnauthorized, errors.ErrUnauthorized,
				"no access"))
			ctx.Abort()
			return
		}

		typ := contents[0]
		credentials := contents[1]

		switch typ {
		case "Bearer":
			userCtx, err = jwt.VerifyTokenWithUser(credentials)
			if err != nil {
				logger.Errorf("JWT verify error: %v", err)
				ctx.Render(errors.NewError(http.StatusForbidden, errors.ErrUnauthorized,
					"jwt verification failed"))
				c.Abort()
				return
			}
			c.Set("UserContext", userCtx)
		}

		if ok, _ := a.enforcer.Enforce(fmt.Sprintf(
			"%v::role", userCtx.Role.String()), userCtx.ResourceID, ctx.FullPath(), ctx.Request.Method); !ok {
			ctx.Render(errors.NewError(http.StatusUnauthorized, errors.ErrUnauthorized,
				"no access"))
			ctx.Abort()
			return
		}

		c.Next()
	}
}
