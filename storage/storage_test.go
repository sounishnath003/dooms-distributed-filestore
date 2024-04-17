package storage

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathTransformFunc(t *testing.T) {
	key := "momsbestphoto"
	log.Println("testing out the function code key")
	pathName := CASPathTransformFunc(key)
	assert.Equal(t, "511be51e0eb/0bacc466548/5584399005f", pathName)
}

func TestStorage(t *testing.T) {
	opts := &StorageOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	log.Printf("storage options: %+v\n", opts)

	// Store initalization
	s := NewStorage(*opts)
	data := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("temp", data); err != nil {
		t.Error(err)
	}
}
