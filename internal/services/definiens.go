/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package services

const (
	BadRequest          = 400
	NotFound            = 404
	DataAlreadyExists   = 4001
	DataNotExists       = 4002
	InternalServerError = 500
)

var (
	errorsDefines = map[int]*err{
		BadRequest: {
			Code:    BadRequest,
			Message: "Bad request",
		},
		NotFound: {
			Code:    NotFound,
			Message: "Not Found",
		},
		DataAlreadyExists: {
			Code:    DataAlreadyExists,
			Message: "Data already exists",
		},
		DataNotExists: {
			Code:    DataNotExists,
			Message: "Data not exists",
		},
		InternalServerError: {
			Code:    InternalServerError,
			Message: "Internal server error",
		},
	}
)
