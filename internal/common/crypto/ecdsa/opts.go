package ecdsa

var (
	P256 = "P256"
	P384 = "P384"
)

type P256KeyGenOpts struct{}

func (opts *P256KeyGenOpts) Algorithm() string {
	return P256
}

type P384KeyGenOpts struct{}

func (opts *P384KeyGenOpts) Algorithm() string {
	return P384
}
