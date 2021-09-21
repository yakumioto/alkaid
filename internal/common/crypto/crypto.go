package crypto

type Key interface {
	Bytes() ([]byte, error)
	SKI() []byte
	Symmetric() bool
	Private() bool
	PublicKey() (Key, error)
}

type KeyGenerator interface {
	KeyGen(opts KeyGenOpts) (Key, error)
}

type KeyGenOpts interface {
	Algorithm() string
}
