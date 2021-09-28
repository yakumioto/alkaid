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
