/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package sqlite3

import (
	"github.com/yakumioto/alkaid/internal/common/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const Driver = "sqlite3"

type sqlite3 struct {
	db *gorm.DB
}

func NewDB(path string) (*sqlite3, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &sqlite3{
		db: db,
	}, nil
}

func (s *sqlite3) AutoMigrate(dst ...interface{}) error {
	return s.db.AutoMigrate(dst...)
}

func (s *sqlite3) Create(value interface{}) error {
	if tx := s.db.Create(value); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *sqlite3) Update(values interface{}, options *storage.UpdateOptions) error {
	if options == nil {
		return storage.ErrNeedUpdateOptions
	}

	if tx := s.db.Model(values).Where(options.Query, options.Args).Updates(values); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *sqlite3) FindByID(dest interface{}, conditions ...interface{}) error {
	if tx := s.db.First(dest, conditions); tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return storage.ErrNotFound
		}
		return tx.Error
	}

	return nil
}

func (s *sqlite3) FindByQuery(dest interface{}, options *storage.QueryOptions) error {
	if options == nil {
		options = storage.NewQueryOptions()
	}

	tx := s.db.Order(options.GetOrder()).Limit(options.GetLimit()).Offset(options.GetOffset())

	if where := options.GetWhere(); where != nil {
		tx.Where(where.Query, where.Args)
	}

	if ors := options.GetOrs(); ors != nil {
		for _, or := range ors {
			tx.Or(or.Query, or.Args)
		}
	}

	if not := options.GetNot(); not != nil {
		tx.Not(not.Query, not.Args)
	}

	tx.Find(dest)

	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *sqlite3) Delete(value interface{}, conditions ...interface{}) error {
	if tx := s.db.Delete(value, conditions); tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (s *sqlite3) Begin() storage.Storage {
	return &sqlite3{db: s.db.Begin()}
}

func (s *sqlite3) Commit() error {
	if tx := s.db.Commit(); tx.Error != nil {
		return tx.Error
	}

	return nil
}
