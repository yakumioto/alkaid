/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package aes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumioto/alkaid/internal/common/crypto"
)

func TestAesKeyImport_KeyImport(t *testing.T) {
	ki := &keyImporter{}
	_, err := ki.KeyImport([]byte("test"), &crypto.AES128KeyImportOpts{})
	assert.NoError(t, err)
}
