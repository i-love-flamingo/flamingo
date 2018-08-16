package cache

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

type (
	// FileBackend is a cache backend which saves the data in files
	FileBackend struct {
		baseDir string
	}
)

const defaultBaseDir = "/tmp/cache"

var (
	escape = regexp.MustCompile(`[^a-zA-Z0-9.]`)
)

// NewFileBackend returns a FileBackend operating in the given baseDir
func NewFileBackend(baseDir string) *FileBackend {
	if baseDir == "" {
		baseDir = defaultBaseDir
	}

	return &FileBackend{
		baseDir: baseDir,
	}
}

// Get reads a cache entry
func (fb *FileBackend) Get(key string) (entry *Entry, found bool) {
	key = escape.ReplaceAllString(key, ".")

	b, err := ioutil.ReadFile(filepath.Join(fb.baseDir, key))
	if err != nil {
		return nil, false
	}

	bb := bytes.NewBuffer(b)
	d := gob.NewDecoder(bb)
	entry = new(Entry)
	err = d.Decode(&entry)
	if err != nil {
		return nil, false
	}

	return entry, true
}

// Set writes a cache entry
func (fb *FileBackend) Set(key string, entry *Entry) {
	key = escape.ReplaceAllString(key, ".")

	gob.Register(entry)
	gob.Register(entry.Data)

	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(entry)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(filepath.Join(fb.baseDir, key), b.Bytes(), os.ModePerm)
}

// Purge deletes a cache entry
func (fb *FileBackend) Purge(key string) {
	key = escape.ReplaceAllString(key, ".")
	os.Remove(filepath.Join(fb.baseDir, key))
}

// PurgeTags is not supported by FileBackend and does nothing
func (*FileBackend) PurgeTags(tags []string) {}

// Flush is not supported by FileBackend and does nothing
func (*FileBackend) Flush() {}
