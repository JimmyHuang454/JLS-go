package jls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeRandom(t *testing.T) {
	fakeRandom1 := NewFakeRandom([]byte("abc"), []byte("abc"))
	fakeRandom2 := NewFakeRandom([]byte("abc"), []byte("abc"))
	fakeRandom3 := NewFakeRandom([]byte("abc"), []byte("ab"))
	assert.Equal(t, fakeRandom1.IV, fakeRandom2.IV)
	assert.Equal(t, fakeRandom1.PWD, fakeRandom2.PWD)
	assert.Equal(t, fakeRandom1.PWD, fakeRandom3.PWD)
	assert.NotEqual(t, fakeRandom1.IV, fakeRandom3.IV)

	err := fakeRandom1.Build()
	assert.Nil(t, err)

	isValid, err := fakeRandom2.Check(fakeRandom1.Random)
	assert.Nil(t, err)
	assert.Equal(t, isValid, true)

	assert.Equal(t, fakeRandom1.N, fakeRandom2.N)
	assert.Equal(t, fakeRandom1.Random, fakeRandom2.Random)

	isValid, err = fakeRandom3.Check(fakeRandom1.Random)
	assert.NotNil(t, err)
	assert.Equal(t, isValid, false)

	random1 := fakeRandom1.Random
	err = fakeRandom1.Build()
	assert.Nil(t, err)
	assert.NotEqual(t, random1, fakeRandom1.Random)
}
