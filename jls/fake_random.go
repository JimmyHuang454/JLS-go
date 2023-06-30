package jls

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"io"
)

type FakeRandom struct {
	N      []byte
	Random []byte

	PWD []byte
	IV  []byte
}

func NewFakeRandom(PWD []byte, IV []byte) *FakeRandom {
	pwd := sha256.New()
	pwd.Write(PWD)
	iv := sha512.New()
	iv.Write(IV)

	return &FakeRandom{PWD: pwd.Sum(nil), IV: iv.Sum(nil), Random: make([]byte, 32), N: make([]byte, 16)}
}

func (f *FakeRandom) Build() error {
	n := make([]byte, 16)
	random := make([]byte, 32)
	io.ReadFull(rand.Reader, n)

	random, err := Encrypt(f.IV, n, f.PWD)
	if err != nil {
		return err
	}
	copy(f.Random, random)
	copy(f.N, n)
	return nil
}

func (f *FakeRandom) Check(random []byte) (bool, error) {
	n, err := Decrypt(f.IV, random, f.PWD)
	if err != nil {
		return false, err
	}
	copy(f.Random, random)
	copy(f.N, n)
	return true, nil
}
