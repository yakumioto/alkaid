package crypto

type KeyType int

const (
	ECDSAType KeyType = iota
	AESCBCType
)

type Key interface {
	Bytes() ([]byte, error)
	SKI() []byte
	Symmetric() bool
	Private() bool
	PublicKey() (Key, error)
	Sign(digest []byte) ([]byte, error)
	Verify(hash, sig []byte) bool
	Encrypt(src []byte) ([]byte, error)
	Decrypt(src []byte) ([]byte, error)
}

type KeyGenerator interface {
	KeyGen(opts KeyGenOpts) (Key, error)
}

type KeyGenOpts interface {
	Algorithm() string
}

type KeyImporter interface {
	KeyImport(raw interface{}, opts KeyImportOpts) (Key, error)
}

type KeyImportOpts interface {
	Algorithm() string
}

//func KeyGen(opts KeyGenOpts) (Key, error) {
//	switch opts.Algorithm() {
//	case ECDSAP256, ECDSAP384:
//		kg := &ecdsa.KeyGenerator{}
//		return kg.KeyGen(opts)
//	}
//
//	return nil, errors.New("not found key generator")
//}
//
//func KeyImport(raw interface{}, opts KeyImportOpts) (Key, error) {
//	switch opts.Algorithm() {
//	case AES128, AES192, AES256:
//		kg := &aes.KeyImport{}
//		return kg.KeyImport(raw, opts)
//	}
//
//	return nil, errors.New("not found key importer")
//}
