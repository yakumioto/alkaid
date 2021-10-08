/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package restful

import (
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	Name() string
	Priority() int
	HandlerFunc() gin.HandlerFunc
}

type Controller interface {
	Name() string
	Path() string
	Method() string
	HandlerFunc() gin.HandlerFunc
}

type middlewares []Middleware

func (m middlewares) Len() int {
	return len(m)
}

func (m middlewares) Less(i, j int) bool {
	return m[i].Priority() < m[j].Priority()
}

func (m middlewares) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type options struct {
	mode           string
	requestTimeout time.Duration
}

type Service interface {
}

type service struct {
	middlewares middlewares
	engine      *gin.Engine
}

func (s *service) RegisterMiddlewares(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		middleware := middleware
		s.middlewares = append(s.middlewares, middleware)
	}

	sort.Sort(s.middlewares)
}

func (s *service) RegisterControllers(controllers ...Controller) {
}
