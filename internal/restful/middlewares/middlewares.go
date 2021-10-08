/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package middlewares

import "github.com/gin-gonic/gin"

type Middleware interface {
	Name() string
	Priority() int
	HandlerFunc() gin.HandlerFunc
}

type Middlewares []Middleware

func (m Middlewares) Len() int {
	return len(m)
}

func (m Middlewares) Less(i, j int) bool {
	return m[i].Priority() < m[j].Priority()
}

func (m Middlewares) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}
