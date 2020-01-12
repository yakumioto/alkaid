/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package errors

import "encoding/json"

type Errors struct {
	Errors []*Error `json:"errors,omitempty"`
}

func (e *Errors) Error() string {
	data, _ := json.Marshal(e)
	return string(data)
}

type Error struct {
	Code    int64  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func NewErrors(apicode ...APICode) *Errors {
	err := new(Errors)
	errs := make([]*Error, 0)

	for _, code := range apicode {
		if errorsDefines[code] != nil {
			errs = append(errs, errorsDefines[code])
		}
	}

	if len(errs) != 0 {
		err.Errors = errs
	}

	return err
}
