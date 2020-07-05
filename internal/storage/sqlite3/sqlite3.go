/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package sqlite3

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3" // database driver
	"github.com/yakumioto/glog"
	"xorm.io/xorm"

	"github.com/yakumioto/alkaid/internal/storage"
)

var logger = glog.MustGetLogger("storage.sqlite3")

type DB struct {
	path string
	mode string
	x    *xorm.Engine
}

func NewBD(path, mode string, schemas ...interface{}) (*DB, error) {
	dataSourceName := fmt.Sprintf("file:%s?mode=%s&cache=shared", path, mode)
	logger.Debuf("Data source name is: %v", dataSourceName)

	x, err := xorm.NewEngine("sqlite3", dataSourceName)
	if err != nil {
		logger.Errof("New database error: %v", err)
		return nil, err
	}

	if len(schemas) > 0 {
		err = x.Sync2(schemas...)
		if err != nil {
			logger.Errof("Synchronize database schema error: %v", err)
			return nil, err
		}
	}

	return &DB{
		path: path,
		mode: mode,
		x:    x,
	}, nil
}

func (db *DB) Create(bean interface{}) error {
	has, err := db.Exist(bean)
	if err != nil {
		return err
	}

	if has {
		return storage.ErrAlreadyExist
	}

	_, err = db.x.InsertOne(bean)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Exist(bean interface{}) (bool, error) {
	return db.x.Exist(bean)
}

func (db *DB) Update(candiBean, bean interface{}) error {
	has, err := db.Exist(candiBean)
	if err != nil {
		return err
	}

	if !has {
		return storage.ErrNotExist
	}

	_, err = db.x.Update(bean, candiBean)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) Get(_ string, bean interface{}) error {
	has, err := db.x.Get(bean)
	if err != nil {
		return err
	}

	if !has {
		return storage.ErrNotExist
	}

	return nil
}

func (db *DB) Query(beans interface{}, sqlAndArgs ...interface{}) error {
	return db.x.Find(beans, sqlAndArgs...)
}

func (db *DB) Delete(bean interface{}) error {
	_, err := db.x.Delete(bean)
	return err
}

func (db *DB) Close() {
	_ = db.x.Close()
}
