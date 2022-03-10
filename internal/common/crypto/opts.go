/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package crypto

const (
	EcdsaP256 Algorithm = "ECDSA_P256"
	EcdsaP384 Algorithm = "ECDSA_P384"

	AesCbc128 Algorithm = "AES_CBC_128"
	AesCbc192 Algorithm = "AES_CBC_192"
	AesCbc256 Algorithm = "AES_CBC_256"

	HmacSha256 Algorithm = "HMAC_SHA256"
	HmacSha512 Algorithm = "HMAC_SHA512"
)

type Algorithm string

type ECDSAP256KeyGenOpts struct{}

func (opts *ECDSAP256KeyGenOpts) Algorithm() Algorithm {
	return EcdsaP256
}

type ECDSAP384KeyGenOpts struct{}

func (opts *ECDSAP384KeyGenOpts) Algorithm() Algorithm {
	return EcdsaP384
}

type AES128KeyImportOpts struct{}

func (opts *AES128KeyImportOpts) Algorithm() Algorithm {
	return AesCbc128
}

type AES192KeyImportOpts struct{}

func (opts *AES192KeyImportOpts) Algorithm() Algorithm {
	return AesCbc192
}

type AES256KeyImportOpts struct{}

func (opts *AES256KeyImportOpts) Algorithm() Algorithm {
	return AesCbc256
}

type HMACSha256ImportOpts struct{}

func (opts *HMACSha256ImportOpts) Algorithm() Algorithm {
	return HmacSha256
}

type HMACSha512ImportOpts struct{}

func (opts *HMACSha512ImportOpts) Algorithm() Algorithm {
	return HmacSha512
}
