// 		DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
				// Version 2, December 2004

// Copyright 2023 Jimmy Huang

// Everyone is permitted to copy and distribute verbatim or modified
// copies of this license document, and changing it is allowed as long
// as the name is changed.

		// DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
// TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

// 0. You just DO WHAT THE FUCK YOU WANT TO.

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

	return &FakeRandom{PWD: pwd.Sum(nil), IV: iv.Sum(nil)}
}

func (f *FakeRandom) Build() error {
	n := make([]byte, 16)
	random := make([]byte, 32)
	io.ReadFull(rand.Reader, n)

	random, err := Encrypt(f.IV, n, f.PWD)
	if err != nil {
		return err
	}
	f.Random = random
	f.N = n
	return nil
}

func (f *FakeRandom) Check(random []byte) (bool, error) {
	n, err := Decrypt(f.IV, random, f.PWD)
	if err != nil {
		return false, err
	}
	f.Random = random
	f.N = n
	return true, nil
}
