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
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	// empty path test
	err := Init("", "")
	assert.EqualError(t, err, "no db path")

	os.TempDir()
	dbPath := os.TempDir() + "/alkaid/test/db/test.sqlite3"
	_ = Init(dbPath, "")
	assert.Equal(t, true, checkPath(dbPath))

	err = Init("file:test.sqlite3", "mode=memory&chahe=shared")
	assert.NoError(t, err)
}

func checkPath(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}

	return true
}

func testInit() {
	_ = Init("file:test.sqlite3", "mode=memory&chahe=shared")
}
