/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package errors

type Code int

const (
	ErrServerUnknownError   Code = 100001
	ErrUnauthorized         Code = 100002
	ErrForbidden            Code = 100003
	ErrBadRequestParameters Code = 100004

	ErrUserNotFount        Code = 200001
	ErrUserCreateVerifying Code = 200002
)
