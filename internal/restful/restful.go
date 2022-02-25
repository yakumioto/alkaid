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
	"github.com/yakumioto/alkaid/internal/common/util"
)

const (
	ReleaseMode = "release"
	DebugMode   = "debug"
)

type Service interface {
	RegisterMiddlewares(...Middleware)
	RegisterHandlerMiddlewares(...Middleware)
	RegisterControllers(...Controller)
	Run(addr string) error
}

type Middleware interface {
	Name() string
	Sequence() int
	HandlerFunc() gin.HandlerFunc
}

type Controller interface {
	Name() string
	Path() string
	Method() string
	HandlerFuncChain() []gin.HandlerFunc
}

type middlewares []Middleware

type controllers []Controller

func (m middlewares) Len() int {
	return len(m)
}

func (m middlewares) Less(i, j int) bool {
	return m[i].Sequence() < m[j].Sequence()
}

func (m middlewares) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

type options struct {
	mode           string
	requestTimeout time.Duration
}

type OptionFunc func(opt *options)

func WithMode(mode string) OptionFunc {
	return func(opt *options) {
		switch mode {
		case ReleaseMode:
			opt.mode = ReleaseMode
		case DebugMode:
			opt.mode = DebugMode
		default:
			opt.mode = DebugMode
		}
	}
}

func WithRequestTimeout(duration time.Duration) OptionFunc {
	return func(opt *options) {
		opt.requestTimeout = duration
	}
}

var (
	defaultOptions = options{
		mode:           ReleaseMode,
		requestTimeout: 10 * time.Second,
	}
)

type service struct {
	opts               *options
	engine             *gin.Engine
	middlewares        middlewares
	handlerMiddlewares middlewares
	controllers        controllers
}

func NewService(optsFunc ...OptionFunc) Service {
	opts := defaultOptions
	for _, f := range optsFunc {
		f(&opts)
	}

	gin.SetMode(opts.mode)

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
	for _, controller := range controllers {
		controller := controller
		s.controllers = append(s.controllers, controller)
	}
}

func (s *service) Run(addr string) error {
	for _, middleware := range s.middlewares {
		s.engine.Use(middleware.HandlerFunc())
	}

	hMiddlewaresChain := make([]gin.HandlerFunc, 0, len(s.handlerMiddlewares))

	for _, hMiddleware := range s.handlerMiddlewares {
		hMiddlewaresChain = append(hMiddlewaresChain, hMiddleware.HandlerFunc())
	}

	for _, controller := range s.controllers {
		handlerChain := make([]gin.HandlerFunc, 0)
		handlerChain = append(handlerChain, hMiddlewaresChain...)
		handlerChain = append(handlerChain, controller.HandlerFuncChain()...)
		s.engine.Handle(controller.Method(), controller.Path(), handlerChain...)
	}

	return s.engine.Run(addr)
}

type Base struct{}

func (c *Base) RenderFormat(ctx *gin.Context) string {
	return ctx.GetString("AcceptFormat")
}

func (c *Base) MatchVersion(ctx *gin.Context, version string) bool {
	if ctx.GetString("AcceptVersion") != version {
		ctx.Next()
		return false
	}

	return true
}

func (c *Base) Render(ctx *gin.Context, obj interface{}) *gin.Context {
	return util.Render(ctx, c.RenderFormat(ctx), obj)
}
