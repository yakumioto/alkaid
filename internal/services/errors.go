/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package services

import (
	"encoding/json"
	"fmt"
)

type err struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Errors struct {
	Errors []*err `json:"errors,omitempty"`
}

func (e *Errors) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

func NewError(code int, v ...interface{}) *Errors {
	return new(Errors).Add(code, v...)
}

func (e *Errors) Add(code int, v ...interface{}) *Errors {
	if len(v) == 0 {
		if errorsDefines[code] != nil {
			e.Errors = append(e.Errors, errorsDefines[code])
			return e
		}
	}

	e.Errors = append(e.Errors, &err{
		Code:    int64(code),
		Message: fmt.Sprint(v...),
	})

	return e
}
