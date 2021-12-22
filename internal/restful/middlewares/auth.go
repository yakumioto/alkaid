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

type Auth struct {
}

func (a *Auth) Name() string {
	return "Auth"
}

func (a *Auth) Sequence() int {
	return 4
}

func (a *Auth) HandlerFunc() gin.HandlerFunc {
	// TODO implement me
	panic("implement me")
}
