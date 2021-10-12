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

type Mode string

const (
	ProductMode Mode = "product"
	DevelopMode Mode = "develop"
)

type Service interface {
	RegisterMiddlewares(...Middleware)
	RegisterHandlerMiddlewares(...Middleware)
	RegisterControllers(...Controller)
	Run(addr string) error
}

type Middleware interface {
	Name() string
	Priority() int
	HandlerFunc() gin.HandlerFunc
}

type Controller interface {
	Name() string
	Path() string
	Method() string
	HandlerFuncChain() []gin.HandlerFunc
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
	mode           Mode
	requestTimeout time.Duration
}

type OptionFunc func(opt *options)

func WithMode(mode Mode) OptionFunc {
	return func(opt *options) {
		opt.mode = mode
	}
}

func WithRequestTimeout(duration time.Duration) OptionFunc {
	return func(opt *options) {
		opt.requestTimeout = duration
	}
}

var (
	defaultOptions = options{
		mode:           ProductMode,
		requestTimeout: 10 * time.Second,
	}
)

type service struct {
	opts               *options
	engine             *gin.Engine
	middlewares        middlewares
	handlerMiddlewares middlewares
}

func NewService(optsFunc ...OptionFunc) Service {
	opts := defaultOptions
	for _, f := range optsFunc {
		f(&opts)
	}

	return &service{
		opts:   &opts,
		engine: gin.New(),
	}
}

func (s *service) RegisterMiddlewares(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		middleware := middleware
		s.middlewares = append(s.middlewares, middleware)
	}

	sort.Sort(s.middlewares)
}

func (s *service) RegisterHandlerMiddlewares(middlewares ...Middleware) {
	for _, middleware := range middlewares {
		middleware := middleware
		s.handlerMiddlewares = append(s.handlerMiddlewares, middleware)
	}

	sort.Sort(s.handlerMiddlewares)
}

func (s *service) RegisterControllers(controllers ...Controller) {
	hMiddlewaresChain := make([]gin.HandlerFunc, 0, len(s.handlerMiddlewares))

	for _, hMiddleware := range s.handlerMiddlewares {
		hMiddlewaresChain = append(hMiddlewaresChain, hMiddleware.HandlerFunc())
	}

	for _, controller := range controllers {
		handlerChain := make([]gin.HandlerFunc, 0)
		handlerChain = append(handlerChain, hMiddlewaresChain...)
		handlerChain = append(handlerChain, controller.HandlerFuncChain()...)
		s.engine.Handle(controller.Method(), controller.Path(), handlerChain...)
	}
}

func (s *service) Run(addr string) error {
	switch s.opts.mode {
	case ProductMode:
		gin.SetMode(gin.ReleaseMode)
	case DevelopMode:
		gin.SetMode(gin.DebugMode)
	}

	return s.engine.Run(addr)
}
