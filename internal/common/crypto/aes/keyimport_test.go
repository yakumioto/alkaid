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

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/yakumioto/alkaid/internal/common/crypto"
)

func TestKeyImport(t *testing.T) {
	tcs := []struct {
		qPassword interface{}
		aError    error
	}{
		{
			[]byte("test_key"),
			nil,
		},
		{
			struct {
			}{},
			errors.New("only supports []byte type of key"),
		},
	}

	for _, tc := range tcs {
		ki := &keyImporter{}
		_, err := ki.KeyImport(tc.qPassword, &crypto.AES128KeyImportOpts{})
		if tc.aError == nil {
			assert.NoError(t, err)
			continue
		}

		assert.Contains(t, err.Error(), tc.aError.Error())
	}

}
