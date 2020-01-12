/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package errors

type APICode int64

const (
	BadAuthenticationData = 4001
	DataAlreadyExists     = 4002
	DataNotExists         = 4003
)

var (
	errorsDefines = map[APICode]*Error{
		BadAuthenticationData: {
			Code:    BadAuthenticationData,
			Message: "Bad authentication data",
		},
		DataAlreadyExists: {
			Code:    DataAlreadyExists,
			Message: "Data already exists",
		},
		DataNotExists: {
			Code:    DataNotExists,
			Message: "Data not exists",
		},
	}
)
