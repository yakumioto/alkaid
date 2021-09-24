package crypto

const (
	ECDSAP256 = "ECDSAP256"
	ECDSAP384 = "ECDSAP384"

	AES128 = "AES128"
	AES192 = "AES192"
	AES256 = "AES256"
)

type ECDSAP256KeyGenOpts struct{}

func (opts *ECDSAP256KeyGenOpts) Algorithm() string {
	return ECDSAP256
}

type ECDSAP384KeyGenOpts struct{}

func (opts *ECDSAP384KeyGenOpts) Algorithm() string {
	return ECDSAP384
}

type AES128KeyImportOpts struct{}

func (opts *AES128KeyImportOpts) Algorithm() string {
	return AES128
}

type AES192KeyImportOpts struct{}

func (opts *AES192KeyImportOpts) Algorithm() string {
	return AES192
}

type AES256KeyImportOpts struct{}

func (opts *AES256KeyImportOpts) Algorithm() string {
	return AES256
}
