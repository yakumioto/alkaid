/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package main

import (
	"fmt"

	"github.com/yakumioto/alkaid/internal/services/users"
)

func main() {
	fmt.Println(users.RoleRoot.String())
}
