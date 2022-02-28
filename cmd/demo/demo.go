/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package main

import (
	"log"

	"github.com/yakumioto/alkaid/internal/common/storage"
	"github.com/yakumioto/alkaid/internal/common/storage/sqlite3"
	"github.com/yakumioto/alkaid/internal/services/users"
)

func main() {
	db, _ := sqlite3.NewDB("testData/alkaid.db")

	storage.Initialize(db)

	user := new(users.User)
	if err := user.FindByIDOrEmail("root"); err != nil {
		log.Println(err)
	}

	log.Println(user)
}
