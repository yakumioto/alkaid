/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package sqlite3

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func testMockSQL() (*sqlite3, sqlmock.Sqlmock) {
	stdDB, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(&sqlite.Dialector{Conn: stdDB})
	if err != nil {
		panic(err)
	}

	return &sqlite3{
		db: db,
	}, mock
}

func TestNewDB(t *testing.T) {
	_, err := NewDB("file::memory:?cache=shared")
	if err != nil {
		assert.NoErrorf(t, err, "new db error: %v", err)
	}
}

type TestDocuments struct {
	Message string
}

func TestSqlite3_Create(t *testing.T) {
	db, mock := testMockSQL()

	test := &TestDocuments{
		Message: "Hello World!",
	}

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `test_documents`").
		WithArgs(test.Message).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if err := db.Create(test); err != nil {
		assert.NoErrorf(t, err, "create data error: %v", err)
	}
}
