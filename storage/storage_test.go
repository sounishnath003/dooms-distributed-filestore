package storage

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestphoto"
	log.Println("testing out the function code key")
	expectedOriginalKey := "511be51e0eb0bacc4665485584399005ff19c5c3"
	expectedPathname := "511be/51e0e/b0bac/c4665/48558/43990/05ff1/9c5c3"

	pathKey := CASPathTransformFunc(key)
	assert.Equal(t, expectedPathname, pathKey.Pathname)
	assert.Equal(t, expectedOriginalKey, pathKey.Filename)
}

func TestStorage(t *testing.T) {
	opts := &StorageOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	log.Printf("storage options: %+v\n", opts)

	// Store initalization
	s := NewStorage(*opts)
	key := "momspecials"
	data := []byte("some jpg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	r, err := s.Read(key)
	assert.Equal(t, nil, err)

	b, _ := io.ReadAll(r)
	assert.Equal(t, data, b)

}
