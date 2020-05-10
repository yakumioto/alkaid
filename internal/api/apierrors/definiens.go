/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package apierrors

type ErrorCode int64

const (
	BadRequest          = 400
	DataAlreadyExists   = 4001
	DataNotExists       = 4002
	InternalServerError = 500
)

var (
	errorsDefines = map[ErrorCode]*error{
		BadRequest: {
			Code:    BadRequest,
			Message: "Bad request",
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
