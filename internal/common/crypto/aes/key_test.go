package aes

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yakumioto/alkaid/internal/common/crypto"
)

func TestAesCBCPrivateKey_Encrypt(t *testing.T) {
	ack := &aesCBCPrivateKey{privateKey: []byte("test key"), algorithm: crypto.AES256}
	ciphertext, err := ack.Encrypt([]byte("hello world"))
	assert.NoError(t, err)
	t.Log(base64.StdEncoding.EncodeToString(ciphertext))
}

func TestAesCBCPrivateKey_Decrypt(t *testing.T) {
	ciphertext, _ := base64.StdEncoding.DecodeString("AgAAAAAAAAAAABV1g1cM9Ojzph0623F+LcDSVd0dZhpn8lXeI4jr5NZT")
	ack := &aesCBCPrivateKey{privateKey: []byte("test key"), algorithm: crypto.AES256}
	text, err := ack.Decrypt(ciphertext)
	assert.NoError(t, err)
	t.Log(string(text))
}
