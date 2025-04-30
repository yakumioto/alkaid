/*
 * Copyright (c) 2025 yakumioto <yaku.mioto@gmail.com>
 * All rights reserved.
 */

package alkaid

import (
	"fmt"
	"log/slog"
	"strings"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	stringadapter "github.com/casbin/casbin/v2/persist/string-adapter"
)

type Authenticator interface {
	Initialize() error
	RegisterPolicyWithController(Controller, bool) Authenticator
	RegisterPolicyWithString(params ...any) Authenticator
	HandlerFunc() HandlerFunc
}

type CasbinAuthenticator struct {
	once        sync.Once
	enforcer    *casbin.Enforcer
	model       string
	policy      []string
	handlerFunc HandlerFunc
}

func NewCasbinAuthenticator(model string, handlerFunc HandlerFunc) *CasbinAuthenticator {
	return &CasbinAuthenticator{
		model:       model,
		handlerFunc: handlerFunc,
	}
}

func (a *CasbinAuthenticator) Initialize() error {
	var err error

	a.once.Do(func() {
		var m model.Model
		m, err = model.NewModelFromString(a.model)
		if err != nil {
			return
		}

		a.enforcer, err = casbin.NewEnforcer(m, stringadapter.NewAdapter(strings.Join(a.policy, "\n")))
		if err != nil {
			return
		}
	})

	slog.Debug("CasbinAuthenticator initialized", "model", a.model, "policy", strings.Join(a.policy, "\n"))

	return err
}

func (a *CasbinAuthenticator) RegisterPolicyWithController(c Controller, enableUrlVersionPrefix bool) Authenticator {
	for _, pi := range c.Paths() {
		a.policy = append(a.policy, fmt.Sprintf("p, %s, %s, %s",
			pi.Role,
			pi.Path,
			c.Method(),
		))
	}

	if enableUrlVersionPrefix {
		for _, pi := range c.Paths() {
			versionSet := make(map[string]struct{})
			for _, vh := range c.VersionHandlerFuncs() {
				version, _ := vh()
				if _, exists := versionSet[version.MajorString()]; !exists {
					versionSet[version.MajorString()] = struct{}{}
					a.policy = append(a.policy, fmt.Sprintf("p, %s, %s, %s",
						pi.Role,
						"/v"+version.MajorString()+pi.Path,
						c.Method(),
					))
				}
			}
		}
	}

	return a
}

func (a *CasbinAuthenticator) RegisterPolicyWithString(params ...any) Authenticator {
	for _, param := range params {
		switch v := param.(type) {
		case string:
			a.policy = append(a.policy, v)
		case []string:
			a.policy = append(a.policy, v...)
		}
	}

	return a
}

func (a *CasbinAuthenticator) HandlerFunc() HandlerFunc {
	return a.handlerFunc
}

type NoneAuthenticator struct {
	handlerFunc HandlerFunc
}

func NewNoneNoneAuthenticator() Authenticator {
	return &NoneAuthenticator{}
}

func (n *NoneAuthenticator) Initialize() error {
	return nil
}

func (n *NoneAuthenticator) RegisterPolicyWithController(c Controller, enableUrlVersionPrefix bool) Authenticator {
	return n
}

func (n *NoneAuthenticator) RegisterPolicyWithString(params ...any) Authenticator {
	return n
}

func (n *NoneAuthenticator) HandlerFunc() HandlerFunc {
	return n.handlerFunc
}
