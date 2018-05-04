package cache

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"os"
	"regexp"
)

type (
	FileBackend struct{}
)

var (
	_ Backend = &FileBackend{}

	escape = regexp.MustCompile(`[^a-zA-Z0-9.]`)
)

func (*FileBackend) Get(key string) (entry *Entry, found bool) {
	key = escape.ReplaceAllString(key, ".")

	b, err := ioutil.ReadFile("/tmp/cache/" + key)
	if err != nil {
		return nil, false
	}

	bb := bytes.NewBuffer(b)
	d := gob.NewDecoder(bb)
	entry = new(Entry)
	d.Decode(&entry)
	return entry, true
}

func (*FileBackend) Set(key string, entry *Entry) {
	key = escape.ReplaceAllString(key, ".")

	gob.Register(entry)

	b := new(bytes.Buffer)
	gob.NewEncoder(b).Encode(entry)

	ioutil.WriteFile("/tmp/cache/"+key, b.Bytes(), os.ModePerm)
}

func (*FileBackend) Purge(key string) {
	key = escape.ReplaceAllString(key, ".")
	os.Remove("/tmp/cache/" + key)
}

func (*FileBackend) PurgeTags(tags []string) {}
func (*FileBackend) Flush()                  {}
