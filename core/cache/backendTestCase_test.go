package cache

import (
	"encoding/gob"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type backendTestEntry struct {
	Content string
}

type (
	// BackendTestCase representations
	BackendTestCase struct {
		t            *testing.T
		backend      Backend
		tagsInResult bool
	}
)

func init() {
	gob.Register(new(backendTestEntry))
}

func NewBackendTestCase(t *testing.T, backend Backend, tagsInResult bool) *BackendTestCase {
	return &BackendTestCase{
		t:            t,
		backend:      backend,
		tagsInResult: tagsInResult,
	}
}

func (tc *BackendTestCase) RunTests() {
	tc.testSetGetPurge()

	tc.testFlush()

	if _, ok := tc.backend.(TagSupportingBackend); ok {
		tc.testPurgeTags()
	}
}

func (tc *BackendTestCase) testSetGetPurge() {
	entry := tc.buildEntry("ASDF", []string{"eins", "zwei"})
	wantedEntry := tc.buildWanted(entry)

	tc.setAndCompareEntry("ONE_KEY", entry, wantedEntry)
	tc.setAndCompareEntry("ANOTHER_KEY", entry, wantedEntry)

	err := tc.backend.Purge("ONE_KEY")
	if err != nil {
		tc.t.Fatalf("Purge Key Failed: %v", err)
	}

	tc.shouldNotExists("ONE_KEY")

	tc.getAndCompareEntry("ANOTHER_KEY", wantedEntry)
}

func (tc *BackendTestCase) testFlush() {
	entry := tc.buildEntry("ASDF", []string{"eins", "zwei"})

	tc.setEntry("ONE_KEY", entry)
	tc.setEntry("ANOTHERKEY_KEY", entry)

	err := tc.backend.Flush()
	if err != nil {
		tc.t.Fatalf("Flush Failed: %v", err)
	}

	tc.shouldNotExists("ONE_KEY")
	tc.shouldNotExists("ANOTHERKEY_KEY")
}

func (tc *BackendTestCase) testPurgeTags() {
	entry := tc.buildEntry("ASDF", []string{"eins", "zwei"})
	entryWithoutTags := tc.buildEntry("ASDF", []string{})

	tc.setEntry("ONE_KEY", entry)
	tc.setEntry("ANOTHERKEY_KEY", entry)
	tc.setEntry("THIRD_KEY", entryWithoutTags)

	tagsToPurge := []string{"eins"}
	err := tc.backend.(TagSupportingBackend).PurgeTags(tagsToPurge)
	if err != nil {
		tc.t.Fatalf("Purge Tags Failed: %v", err)
	}

	tc.shouldNotExists("ONE_KEY")
	tc.shouldNotExists("ANOTHERKEY_KEY")
	tc.shouldExists("THIRD_KEY")
}

func (tc *BackendTestCase) setEntry(key string, entry *Entry) {
	err := tc.backend.Set(key, entry)
	if err != nil {
		tc.t.Fatalf("Failed to set Entry for key %v with error: %v", key, err)
	}
	tc.shouldExists(key)
}

func (tc *BackendTestCase) setAndCompareEntry(key string, entry *Entry, wanted *Entry) {
	tc.setEntry(key, entry)
	tc.getAndCompareEntry(key, wanted)
}

func (tc *BackendTestCase) getAndCompareEntry(key string, wanted *Entry) {
	entry := tc.shouldExists(key)
	tc.mustBeEqual(entry, wanted)
}

func (tc *BackendTestCase) mustBeEqual(entry *Entry, wanted *Entry) {
	if entry.Meta.Gracetime != wanted.Meta.Gracetime {
		tc.t.Fatalf("Entry gracetimes are not equal %v, want %v", entry.Meta.Gracetime, wanted.Meta.Gracetime)
	}

	if entry.Meta.Lifetime != wanted.Meta.Lifetime {
		tc.t.Fatalf("Entry lifetimes are not equal %v, want %v", entry.Meta.Lifetime, wanted.Meta.Lifetime)
	}

	if !assert.Equal(tc.t, entry.Meta.Tags, wanted.Meta.Tags) {
		tc.t.Fatalf("Entry Meta.Tags are not equal %v, want %v", entry.Meta.Tags, wanted.Meta.Tags)
	}

	if !reflect.DeepEqual(entry.Data, wanted.Data) {
		tc.t.Fatalf("Entry data are not equal %v, want %v", entry.Data, wanted.Data)
	}
}

func (tc *BackendTestCase) shouldExists(key string) *Entry {
	entry, found := tc.backend.Get(key)
	if !found {
		tc.t.Fatalf("Failed to get Entry with key: %v", key)
	}
	return entry
}

func (tc *BackendTestCase) shouldNotExists(key string) {
	entry, found := tc.backend.Get(key)
	if found {
		tc.t.Fatalf("Entry with key %v should not exists, but returns %v", key, entry)
	}
}

func (tc *BackendTestCase) buildEntry(content string, tags []string) *Entry {
	return &Entry{
		Data: &backendTestEntry{
			Content: content,
		},
		Meta: Meta{
			lifetime:  time.Now().Add(time.Minute * 3),
			gracetime: time.Now().Add(time.Minute * 30),
			Tags:      tags,
		},
	}
}

func (tc *BackendTestCase) buildWanted(orig *Entry) *Entry {
	var meta Meta
	if tc.tagsInResult {
		meta = Meta{
			lifetime:  time.Now().Add(time.Minute * 3),
			gracetime: time.Now().Add(time.Minute * 30),
			Tags:      orig.Meta.Tags,
		}
	} else {
		meta = Meta{
			lifetime:  orig.Meta.lifetime,
			gracetime: orig.Meta.gracetime,
		}
	}
	return &Entry{
		Data: orig.Data,
		Meta: meta,
	}
}
