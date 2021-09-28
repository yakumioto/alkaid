package crypto

const (
	ECDSAP256 Algorithm = "ECDSAP256"
	ECDSAP384 Algorithm = "ECDSAP384"

	AES128 Algorithm = "AES128"
	AES192 Algorithm = "AES192"
	AES256 Algorithm = "AES256"
)

type Algorithm string

type ECDSAP256KeyGenOpts struct{}

func (opts *ECDSAP256KeyGenOpts) Algorithm() Algorithm {
	return ECDSAP256
}

type ECDSAP384KeyGenOpts struct{}

func (opts *ECDSAP384KeyGenOpts) Algorithm() Algorithm {
	return ECDSAP384
}

type AES128KeyImportOpts struct{}

func (opts *AES128KeyImportOpts) Algorithm() Algorithm {
	return AES128
}

type AES192KeyImportOpts struct{}

func (opts *AES192KeyImportOpts) Algorithm() Algorithm {
	return AES192
}

type AES256KeyImportOpts struct{}

func (opts *AES256KeyImportOpts) Algorithm() Algorithm {
	return AES256
}
