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
		cacheMetrics CacheMetrics
		baseDir      string
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

	_ = os.MkdirAll(baseDir, os.ModePerm)

	fb := &FileBackend{
		baseDir:      baseDir,
		cacheMetrics: NewCacheMetrics("file","test"),
	}

	return fb
}

// Get reads a cache entry
func (fb *FileBackend) Get(key string) (entry *Entry, found bool) {
	key = escape.ReplaceAllString(key, ".")

	b, err := ioutil.ReadFile(filepath.Join(fb.baseDir, key))
	if err != nil {
		fb.cacheMetrics.countMiss()
		return nil, false
	}

	bb := bytes.NewBuffer(b)
	d := gob.NewDecoder(bb)
	entry = new(Entry)
	err = d.Decode(&entry)
	if err != nil {
		fb.cacheMetrics.countError("DecodeFailed")
		return nil, false
	}

	fb.cacheMetrics.countHit()
	return entry, true
}

// Set writes a cache entry
func (fb *FileBackend) Set(key string, entry *Entry) error {
	key = escape.ReplaceAllString(key, ".")

	gob.Register(entry)
	gob.Register(entry.Data)

	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(entry)
	if err != nil {
		fb.cacheMetrics.countError("EncodeFailed")
		return err
	}

	ioutil.WriteFile(filepath.Join(fb.baseDir, key), b.Bytes(), os.ModePerm)

	return nil
}

// Purge deletes a cache entry
func (fb *FileBackend) Purge(key string) error {
	key = escape.ReplaceAllString(key, ".")
	os.Remove(filepath.Join(fb.baseDir, key))

	return nil
}

func (fb *FileBackend) deleteFolderContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// Flush is not supported by FileBackend and does nothing
func (fb *FileBackend) Flush() error {
	return fb.deleteFolderContents(fb.baseDir)
}
