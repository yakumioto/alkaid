package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"

	"github.com/yakumioto/alkaid/internal/common/crypto"
	"github.com/yakumioto/alkaid/internal/common/util"
)

type aesCBCPrivateKey struct {
	privateKey []byte
	algorithm  string
}

func (a *aesCBCPrivateKey) Bytes() ([]byte, error) {
	return nil, errors.New("not supported")
}

func (a *aesCBCPrivateKey) SKI() []byte {
	hash := sha256.New()
	hash.Write([]byte{0x01})
	hash.Write(a.privateKey)
	return hash.Sum(nil)
}

func (a *aesCBCPrivateKey) Symmetric() bool {
	return true
}

func (a *aesCBCPrivateKey) Private() bool {
	return true
}

func (a *aesCBCPrivateKey) PublicKey() (crypto.Key, error) {
	return nil, errors.New("cannot call this method on a symmetric key")
}

func (a *aesCBCPrivateKey) Sign(_ []byte) ([]byte, error) {
	return nil, errors.New("cannot call this method on a symmetric key")
}

func (a *aesCBCPrivateKey) Verify(_, _ []byte) bool {
	return false
}

func (a *aesCBCPrivateKey) Encrypt(text []byte) ([]byte, error) {
	paddedText := pkcs7Padding(text)

	iv, err := a.randomIV(aes.BlockSize)
	if err != nil {
		return nil, fmt.Errorf("random iv error: %v", err)
	}
	stretchKey := a.stretchKey(iv)

	block, err := aes.NewCipher(stretchKey)
	if err != nil {
		return nil, fmt.Errorf("new chipher error: %v", err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	dst := make([]byte, len(paddedText))
	mode.CryptBlocks(dst, paddedText)

	ciphertext := make([]byte, 0, 10+len(iv)+len(dst))
	ciphertext = append(ciphertext, util.GetVarintBytesWithMaxVarintLen64(int64(crypto.AESCBCType))...)
	ciphertext = append(ciphertext, iv...)
	ciphertext = append(ciphertext, dst...)

	return ciphertext, nil
}

func (a *aesCBCPrivateKey) Decrypt(ciphertext []byte) ([]byte, error) {
	typ := crypto.KeyType(util.MustVarint(ciphertext[:10]))
	iv := ciphertext[10 : 10+16]
	src := ciphertext[10+16:]

	stretchKey := a.stretchKey(iv)

	if typ != crypto.AESCBCType {
		return nil, fmt.Errorf("type does not match: %v", typ)
	}

	block, err := aes.NewCipher(stretchKey)
	if err != nil {
		return nil, fmt.Errorf("new chipher error: %v", err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	paddedText := make([]byte, len(src))
	mode.CryptBlocks(paddedText, src)

	return pkcs7UnPadding(paddedText), nil
}

func (a *aesCBCPrivateKey) randomIV(len int) ([]byte, error) {
	iv := make([]byte, len)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	return iv, nil
}

func (a *aesCBCPrivateKey) stretchKey(iv []byte) []byte {
	keyLen := 0

	switch a.algorithm {
	case crypto.AES128:
		keyLen = 128 / 8
	case crypto.AES192:
		keyLen = 192 / 8
	case crypto.AES256:
		keyLen = 256 / 8
	}

	return util.PBKDF2WithSha256(a.privateKey, iv, keyLen)
}

func pkcs7Padding(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize

	var paddingText []byte
	if padding == 0 {
		paddingText = bytes.Repeat([]byte{byte(aes.BlockSize)}, aes.BlockSize)
	} else {
		paddingText = bytes.Repeat([]byte{byte(padding)}, padding)
	}
	return append(src, paddingText...)
}

func pkcs7UnPadding(src []byte) []byte {
	unPadding := int(src[len(src)-1])
	return src[:(len(src) - unPadding)]
}
