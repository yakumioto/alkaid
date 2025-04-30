/*
 * Copyright (c) 2025 yakumioto <yaku.mioto@gmail.com>
 * All rights reserved.
 */

package alkaid

import (
	"github.com/gin-gonic/gin"
)

const (
	ReleaseMode = gin.ReleaseMode
	DebugMode   = gin.DebugMode
)

// Alkaid is a web framework for building RESTful API services based on Gin.
type Alkaid struct {
	// s is the list of registered services
	s []*Service
	// m is the list of registered middlewares
	m []Middleware
	// a is the authenticator
	a Authenticator
	// e is the Gin engine
	e *gin.Engine
	// rootPath is the root path of the API
	rootPath string
	// headerVersionName is the header name for the version
	headerVersionName string
	// enableUrlVersionPrefix is the flag to enable the version prefix in the URL
	enableUrlVersionPrefix bool
}

type AlkaidOption func(*Alkaid)

// NewAlkaid creates a new Alkaid instance.
func NewAlkaid(opts ...AlkaidOption) *Alkaid {
	a := &Alkaid{
		s: []*Service{},
		m: []Middleware{},
		a: &NoneAuthenticator{},
		e: gin.New(),
	}

	for _, opt := range opts {
		opt(a)
	}

	return a
}

// WithAuthenticator sets the authenticator for the Alkaid instance.
func WithAuthenticator(a Authenticator) AlkaidOption {
	return func(alk *Alkaid) {
		alk.a = a
	}
}

// WithRootPath sets the root path of the API.
func WithRootPath(rootPath string) AlkaidOption {
	return func(alk *Alkaid) {
		alk.rootPath = rootPath
	}
}

// WithEngineOption sets the options for the Gin engine.
func WithEngineOption(opts ...gin.OptionFunc) AlkaidOption {
	return func(alk *Alkaid) {
		alk.e = alk.e.With(opts...)
	}
}

func WithDefaultEngineOption(opts ...gin.OptionFunc) AlkaidOption {
	return func(alk *Alkaid) {
		alk.e = gin.Default()
		alk.e = alk.e.With(opts...)
	}
}

func WithEngineMode(mode string) AlkaidOption {
	return func(alk *Alkaid) {
		gin.SetMode(mode)
	}
}

func WithVersionHeader(header string) AlkaidOption {
	return func(alk *Alkaid) {
		alk.headerVersionName = header
	}
}

func WithEnableUrlVersionPrefix() AlkaidOption {
	return func(alk *Alkaid) {
		alk.enableUrlVersionPrefix = true
	}
}

// RegisterService registers one or more services to the Alkaid instance.
func (alk *Alkaid) RegisterService(service ...*Service) *Alkaid {
	alk.s = append(alk.s, service...)

	return alk
}

// RegisterMiddleware registers one or more middlewares to the Alkaid instance.
func (alk *Alkaid) RegisterMiddleware(middleware ...Middleware) *Alkaid {
	alk.m = append(alk.m, middleware...)

	return alk
}

// Run starts the Alkaid instance, listening and handling requests.
func (alk *Alkaid) Run(addr ...string) error {
	// Register middlewares
	for _, middleware := range alk.m {
		alk.e.Use(alk.wrapHandler(middleware.HandlerFunc()))
	}

	// Register services
	for _, s := range alk.s {
		alk.registerControllers(s.Controllers)
	}

	if err := alk.a.Initialize(); err != nil {
		return err
	}

	// Start the Gin engine, listening and handling requests
	return alk.e.Run(addr...)
}

// registerControllers registers controllers under the specified path prefix
func (alk *Alkaid) registerControllers(controllers []Controller) {
	for _, c := range controllers {
		// If the controller has not disabled authentication, register the authentication policy
		if !c.NoAuth() {
			alk.a.RegisterPolicyWithController(c, alk.enableUrlVersionPrefix)
		}

		for _, p := range c.Paths() {
			// Register the route handler functions of the controller
			alk.e.Handle(c.Method(),
				alk.rootPath+p.Path,
				alk.genHandlerFuncChain(c.DefaultHandlerFunc(), c.VersionHandlerFuncs())...)

			// Unique version handler functions
			if alk.enableUrlVersionPrefix {
				alk.registerUrlVersionPrefixHandlerFuncs(p, c)
			}
		}
	}
}

func (alk *Alkaid) registerUrlVersionPrefixHandlerFuncs(patchInfo *PathInfo, c Controller) {
	versionSet := make(map[string]VersionHandlerFunc)

	for _, vh := range c.VersionHandlerFuncs() {
		version, _ := vh()

		if oldVH, exists := versionSet[version.MajorString()]; !exists {
			versionSet[version.MajorString()] = vh
		} else {
			oldVersion, _ := oldVH()
			if version.GreaterThan(oldVersion) {
				versionSet[version.MajorString()] = vh
			}
		}
	}

	for _, vh := range versionSet {
		v, h := vh()
		fullPath := "v" + v.MajorString() + alk.rootPath + patchInfo.Path
		alk.e.Handle(c.Method(), fullPath, alk.genHandlerFuncChain(h, nil)...)
	}
}

func (alk *Alkaid) genHandlerFuncChain(defaultHandlerFunc HandlerFunc, versionHandlerFuncs []VersionHandlerFunc) []gin.HandlerFunc {
	if len(versionHandlerFuncs) == 0 {
		return []gin.HandlerFunc{
			alk.wrapHandler(alk.a.HandlerFunc()),
			alk.wrapHandler(defaultHandlerFunc),
		}
	}

	return []gin.HandlerFunc{
		alk.wrapHandler(alk.a.HandlerFunc()),
		alk.wrapVersionHandler(versionHandlerFuncs),
		alk.wrapHandler(defaultHandlerFunc),
	}
}

func (alk *Alkaid) wrapHandler(handler HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{Context: c}

		if handler != nil {
			handler(ctx)
		}
	}
}

func (alk *Alkaid) wrapVersionHandler(versionHandlerFuncs []VersionHandlerFunc) gin.HandlerFunc {
	versionHandlerMap := make(map[string]HandlerFunc)

	for _, vh := range versionHandlerFuncs {
		version, handlerFunc := vh()
		versionHandlerMap[version.String()] = handlerFunc
	}

	return func(c *gin.Context) {
		ctx := &Context{Context: c}

		handlerFunc, ok := versionHandlerMap[c.GetHeader(alk.headerVersionName)]
		if !ok {
			ctx.Next()
			return
		}

		if handlerFunc != nil {
			handlerFunc(ctx)
		}

		ctx.Abort()
	}
}
