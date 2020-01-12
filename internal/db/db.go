/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package db

import (
	"errors"
	"os"
	"path/filepath"

	// sqlite3
	_ "github.com/mattn/go-sqlite3"
	"github.com/yakumioto/glog"
	"xorm.io/xorm"
)

var (
	logger *glog.Logger
	x      *xorm.Engine
)

func Init(dbPath, options string) error {
	if dbPath == "" {
		return errors.New("no db path")
	}

	logger = glog.MustGetLogger("db")

	dbDir := filepath.Dir(dbPath)
	logger.Debuf("DB storage directory: %s", dbDir)

	err := os.MkdirAll(dbDir, 0755)
	if err != nil {
		return err
	}

	if options != "" {
		x, err = xorm.NewEngine("sqlite3", dbPath+"?"+options)
	} else {
		x, err = xorm.NewEngine("sqlite3", dbPath)
	}
	if err != nil {
		return err
	}

	err = x.StoreEngine("InnoDB").Sync2(new(CA), new(MSP), new(Network), new(Organization))
	if err != nil {
		return err
	}

	return nil
}
