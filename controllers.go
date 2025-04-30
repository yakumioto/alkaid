/*
 * Copyright (c) 2025 yakumioto <yaku.mioto@gmail.com>
 * All rights reserved.
 */

package alkaid

import (
	"fmt"
)

type HandlerVersion struct {
	Major int
	Minor int
	Patch int
}

func NewHandlerVersion(version ...int) *HandlerVersion {
	handlerVersion := &HandlerVersion{}

	switch len(version) {
	case 3:
		handlerVersion.Patch = version[2]
		fallthrough
	case 2:
		handlerVersion.Minor = version[1]
		fallthrough
	case 1:
		handlerVersion.Major = version[0]
	}

	return handlerVersion
}

func (hv *HandlerVersion) GreaterThan(other *HandlerVersion) bool {
	if hv.Major > other.Major {
		return true
	}

	if hv.Major == other.Major && hv.Minor > other.Minor {
		return true
	}

	if hv.Major == other.Major && hv.Minor == other.Minor && hv.Patch > other.Patch {
		return true
	}

	return false
}

func (hv *HandlerVersion) MajorString() string {
	return fmt.Sprintf("%d", hv.Major)
}

func (hv *HandlerVersion) String() string {
	return fmt.Sprintf("%d.%d.%d", hv.Major, hv.Minor, hv.Patch)
}

type HandlerFunc func(*Context)
type VersionHandlerFunc func() (version *HandlerVersion, handlerFunc HandlerFunc)

type PathInfo struct {
	Role string
	Path string
}

type Controller interface {
	NoAuth() bool
	Paths() []*PathInfo
	Method() string
	DefaultHandlerFunc() HandlerFunc
	VersionHandlerFuncs() []VersionHandlerFunc
}
