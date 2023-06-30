package jls

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

const ivLen = 64 // bytes

func Encrypt(nonce []byte, plaintext []byte, key []byte) ([]byte, error) {
	if len(nonce) != ivLen {
		return nil, errors.New("wrong nonce len")
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithNonceSize(c, ivLen)

	if err != nil {
		return nil, err
	}

	cipertextAndMac := gcm.Seal(nonce, nonce, plaintext, nil)
	return cipertextAndMac[ivLen:], nil
}

func Decrypt(nonce []byte, cipherTextAndMac []byte, key []byte) ([]byte, error) {
	if len(nonce) != ivLen {
		return nil, errors.New("wrong nonce len")
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithNonceSize(c, ivLen)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(nonce) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonceAndPlainText, err := gcm.Open(nonce, nonce, cipherTextAndMac, nil)

	if err != nil {
		return nonceAndPlainText, err
	}

	if len(nonceAndPlainText) <= ivLen {
		return nonceAndPlainText, errors.New("wrong res len")
	}

	return nonceAndPlainText[ivLen:], err
}
