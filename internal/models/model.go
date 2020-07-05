/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package models

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/yakumioto/alkaid/internal/storage"
	"github.com/yakumioto/alkaid/internal/storage/sqlite3"
)

var (
	db storage.DBProvider
)

func InitModels() error {
	var err error
	scheams := []interface{}{new(Organization)}

	driver := viper.GetString("core.database.driver")
	switch driver {
	case "sqlite3":
		db, err = sqlite3.NewBD(viper.GetString("core.database.address"), "rwc", scheams...)
	default:
		return fmt.Errorf("unknown driver: %v", driver)
	}

	return err
}
