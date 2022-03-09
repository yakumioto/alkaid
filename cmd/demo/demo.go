/*
 * Copyright (c) 2022. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/hkdf"
	"golang.org/x/crypto/pbkdf2"
)

func main() {
	fmt.Printf("Email: %v\nMaster password: %v\n", "yaku.mioto@gmail.com", "s9qn#UhhaDir5V2B")
	dk := pbkdf2.Key([]byte("s9qn#UhhaDir5V2B"), []byte("yaku.mioto@gmail.com"), 100000, 32, sha256.New)
	fmt.Println("Master Key:", base64.StdEncoding.EncodeToString(dk))
	enc := make([]byte, 32)
	hkdf.Expand(sha256.New, dk, []byte("enc")).Read(enc)
	mac := make([]byte, 32)
	hkdf.Expand(sha256.New, dk, []byte("mac")).Read(mac)
	fmt.Println("Stretched Master Key:", base64.StdEncoding.EncodeToString(append(enc, mac...)))
	fmt.Println("Encryption Key:", base64.StdEncoding.EncodeToString(enc))
	fmt.Println("MAC Key:", base64.StdEncoding.EncodeToString(mac))

	paddedText := pkcs7Padding([]byte("This is a secret."))
	symmetricKey, _ := base64.StdEncoding.DecodeString("p5Dq/6t/m3gRAPln7BSv5QGBVvwfZ3tGHYHdhkw/m9FHyR1TKEjK0A2lCLWWP0fdix0wWB/HRENJstO3ABu4MQ==")
	iv, err := base64.StdEncoding.DecodeString("c4LZ3sCdyOl7U7mzKuAUtg==")
	block, err := aes.NewCipher(symmetricKey[:32])
	if err != nil {
		panic(err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	ciphertext := make([]byte, len(paddedText))
	mode.CryptBlocks(ciphertext, paddedText)
	fmt.Println("ciphertext:", base64.StdEncoding.EncodeToString(ciphertext))

	hc := hmac.New(sha256.New, symmetricKey[32:])
	var text []byte
	text = append(text, iv...)
	// text = append(text, []byte("|")...)
	text = append(text, ciphertext...)
	hc.Write(text)
	fmt.Println("sign text:", base64.StdEncoding.EncodeToString(hc.Sum(nil)))

	// paddedText = pkcs7Padding(ciphertext)
	// symmetricKey, _ = base64.StdEncoding.DecodeString("p5Dq/6t/m3gRAPln7BSv5QGBVvwfZ3tGHYHdhkw/m9FHyR1TKEjK0A2lCLWWP0fdix0wWB/HRENJstO3ABu4MQ==")
	// iv, err = base64.StdEncoding.DecodeString("c4LZ3sCdyOl7U7mzKuAUtg==")
	// block, err = aes.NewCipher(symmetricKey[32:])
	// if err != nil {
	// 	panic(err)
	// }
	// mode = cipher.NewCBCEncrypter(block, iv)
	// signText := make([]byte, len(paddedText))
	// mode.CryptBlocks(signText, paddedText)
	// fmt.Println("sign text:", base64.StdEncoding.EncodeToString(signText))
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
