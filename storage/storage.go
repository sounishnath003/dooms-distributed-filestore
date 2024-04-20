package storage

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "dooms-distrubuted-network-filestore"

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLength := len(hashStr) / blockSize

	paths := make([]string, sliceLength)

	for i := 0; i < sliceLength; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

type PathKey struct {
	Pathname string
	Filename string
}

type PathTransformFunc func(string) PathKey

func (p PathKey) FirstPathname() string {
	paths := strings.Split(p.Pathname, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}

type StorageOpts struct {
	Root string
	PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	// complexity increament ....
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

type Storage struct {
	StorageOpts
}

func NewStorage(opts StorageOpts) *Storage {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Storage{
		StorageOpts: opts,
	}
}

func (s *Storage) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	fi, err := os.Stat(pathKey.FullPath())
	if err != nil {
		return false
	}
	log.Printf("info the target file: %s\n", fi)
	return true
}

func (s *Storage) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("deleted [%s] from the disk", s.Root + "/" + pathKey.Pathname)
	}()
	return os.RemoveAll(s.Root + "/" + pathKey.FirstPathname())
}

func (s *Storage) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		log.Printf("an error has been occurred: %s\n", err)
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	if err != nil {
		log.Printf("an error has been occurred: %s\n", err)
		return nil, err
	}

	return buf, nil
}

func (s *Storage) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	return os.Open(s.Root + "/" + pathKey.FullPath())
}

func (s *Storage) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)

	if err := os.MkdirAll(s.Root+"/"+pathKey.Pathname, os.ModePerm); err != nil {
		log.Printf("make directory creation error: %s\n", err)
		return err
	}

	fullPath := s.Root + "/" + pathKey.FullPath()
	log.Print("printing the fullPath = ", fullPath)

	f, err := os.Create(fullPath)
	if err != nil {
		log.Printf("creating filename errror: %s\n", err)
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s\n", n, fullPath)

	return nil
}
