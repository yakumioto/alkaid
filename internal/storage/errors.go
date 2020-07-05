/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package storage

import "errors"

var (
	ErrAlreadyExist = errors.New("already exist")
	ErrNotExist     = errors.New("not exist")
)
