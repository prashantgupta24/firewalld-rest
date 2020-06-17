package db

import (
	"bytes"
	"encoding/gob"
	"io"
	"log"
	"os"
	"sync"
)

var lock sync.Mutex

// Register interface with gob
func Register(v interface{}) {
	gob.Register(v)
}

// Save saves a representation of v to the file at path.
func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(path)
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
func Load(path string, v interface{}) error {
	if fileExists(path) {
		lock.Lock()
		defer lock.Unlock()
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		return unmarshal(f, v)
	}
	log.Printf("File %s not found, will be created\n", path)
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
