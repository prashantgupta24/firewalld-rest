package db

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
)

var lock sync.Mutex

//FileType is the main struct for file database
type FileType struct {
	Path string
}

// Register interface with gob
func (fileType *FileType) Register(v interface{}) {
	gob.Register(v)
}

// Save saves a representation of v to the file at path.
func (fileType *FileType) Save(v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(fileType.Path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := marshal(v)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

// Load loads the file at path into v.
func (fileType *FileType) Load(v interface{}) error {
	fullPath, err := filepath.Abs(fileType.Path)
	if err != nil {
		return fmt.Errorf("could not locate absolute path : %v", err)
	}
	if fileExists(fileType.Path) {
		lock.Lock()
		defer lock.Unlock()
		f, err := os.Open(fileType.Path)
		if err != nil {
			return err
		}
		defer f.Close()
		log.Printf("Db file found: %v\n", fullPath)
		return unmarshal(f, v)
	}
	log.Printf("Db file not found, will be created here: %v\n", fullPath)
	return nil
}

// marshal is a function that marshals the object into an
// io.Reader.
var marshal = func(v interface{}) (io.Reader, error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)
	err := e.Encode(v)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf.Bytes()), nil
}

// unmarshal is a function that unmarshals the data from the
// reader into the specified value.
var unmarshal = func(r io.Reader, v interface{}) error {
	d := gob.NewDecoder(r)
	err := d.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
