/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package sqlite3

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBD(t *testing.T) {
	db, err := NewBD("test.sqlite3", "memory")
	assert.NoError(t, err, "new database error")
	db.Close()
}

func TestDB_Create(t *testing.T) {
	type table struct {
		ID   string `xorm:"'id' PRIMARY KEY UNIQUE NOT NULL"`
		Name string `xorm:"name"`
	}

	db, _ := NewBD("test.sqlite3", "memory", new(table))
	defer db.Close()

	tcs := []struct {
		input    *table
		expected error
	}{
		{
			input:    &table{ID: "id1", Name: "name1"},
			expected: nil,
		},
		{
			input:    &table{ID: "id1", Name: "name1"},
			expected: errors.New("already exist"),
		},
	}

	for _, tc := range tcs {
		err := db.Create(tc.input)
		assert.Equal(t, tc.expected, err)
	}
}

func TestDB_Get(t *testing.T) {
	type table struct {
		ID   string `xorm:"'id' PRIMARY KEY UNIQUE NOT NULL"`
		Name string `xorm:"name"`
	}

	db, _ := NewBD("test.sqlite3", "memory", new(table))
	defer db.Close()

	data := &table{ID: "id1"}
	err := db.Get("", data)
	assert.Equal(t, errors.New("not exist"), err)

	_ = db.Create(&table{ID: "id1", Name: "name1"})
	err = db.Get("", data)
	assert.Equal(t, nil, err)

	assert.Equal(t, &table{ID: "id1", Name: "name1"}, data)
}

func TestDB_Update(t *testing.T) {
	type table struct {
		ID   string `xorm:"'id' PRIMARY KEY UNIQUE NOT NULL"`
		Name string `xorm:"name"`
	}

	db, _ := NewBD("test.sqlite3", "memory", new(table))
	defer db.Close()

	data := &table{ID: "id1"}

	err := db.Update(&table{ID: "id1"}, &table{ID: "id1", Name: "name1"})
	assert.Equal(t, errors.New("not exist"), err)

	_ = db.Create(&table{ID: "id1", Name: "name1"})
	err = db.Update(&table{ID: "id1"}, &table{ID: "id1", Name: "name2"})
	assert.Equal(t, nil, err)
	_ = db.Get("", data)
	assert.Equal(t, &table{ID: "id1", Name: "name2"}, data)
}

func TestDB_Query(t *testing.T) {
	type table struct {
		ID   string `xorm:"'id' PRIMARY KEY UNIQUE NOT NULL"`
		Name string `xorm:"name"`
	}

	db, _ := NewBD("test.sqlite3", "memory", new(table))
	defer db.Close()

	datas := make([]*table, 0)
	_ = db.Create(&table{ID: "id1", Name: "name1"})
	err := db.Query(&datas, &table{Name: "name1"})
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(datas))
}

func TestDB_Delete(t *testing.T) {
	type table struct {
		ID   string `xorm:"'id' PRIMARY KEY UNIQUE NOT NULL"`
		Name string `xorm:"name"`
	}

	db, _ := NewBD("test.sqlite3", "memory", new(table))
	defer db.Close()

	_ = db.Create(&table{ID: "id1", Name: "name1"})
	err := db.Delete(&table{ID: "id1"})
	assert.Equal(t, nil, err)
	err = db.Get("", &table{ID: "id1"})
	assert.Equal(t, errors.New("not exist"), err)
}
