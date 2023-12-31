package jls

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAESGCM(t *testing.T) {
	// check len
	cipherTextAndMac, err := Encrypt([]byte("key"), []byte("key"), []byte("key"))
	assert.NotNil(t, err)
	assert.Empty(t, cipherTextAndMac)

	key := sha256.New()
	key.Write([]byte("key"))
	plainText := []byte("cipherText")

	nonce := sha256.New()
	nonce.Write([]byte("nonce"))

	cipherTextAndMac, err = Encrypt(nonce.Sum(nil), plainText, key.Sum(nil))
	assert.NotEmpty(t, cipherTextAndMac)
	assert.Nil(t, err)

	cipherTextAndMac2, err := Encrypt(nonce.Sum(nil), []byte("cipherText"), key.Sum(nil))
	assert.Equal(t, cipherTextAndMac, cipherTextAndMac2)

	res, err := Decrypt(nonce.Sum(nil), cipherTextAndMac, key.Sum(nil))
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
	assert.Equal(t, plainText, res)

	nonce.Write([]byte("key"))
	res, err = Decrypt(nonce.Sum(nil), cipherTextAndMac, key.Sum(nil))
	assert.NotNil(t, err)
	assert.Empty(t, res)
}

func TestSeq(t *testing.T) {
	key := sha256.New()
	key.Write([]byte("1"))
	plainText := []byte("3")

	nonce := sha256.New()
	nonce.Write([]byte("2"))

	cipherTextAndMac, _ := Encrypt(nonce.Sum(nil), plainText, key.Sum(nil))
	temp := make([]byte, 1)
	temp[0] = 36
	assert.Equal(t, cipherTextAndMac[:1], temp)
}
