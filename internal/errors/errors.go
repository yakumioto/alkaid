/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package errors

import "fmt"

type Error struct {
	StatusCode int    `json:"-"`
	Code       int    `json:"code,omitempty"`
	Message    string `json:"message,omitempty"`
}

func (r *Error) Error() string {
	return r.Message
}

func NewError(status, code int, message string) *Error {
	return &Error{
		StatusCode: status,
		Code:       code,
		Message:    message,
	}
}

func NewErrorf(status, code int, format string, a ...interface{}) *Error {
	return &Error{
		StatusCode: status,
		Code:       code,
		Message:    fmt.Sprintf(format, a...),
	}
}
