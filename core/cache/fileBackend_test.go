package cache_test

import (
	"bytes"
	"encoding/gob"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"flamingo.me/flamingo/core/cache"
)

type (
	testStruct struct {
		S string
		B bool
		I int
	}
)

var (
	// Assert the interface is matched
	_ cache.Backend = &cache.FileBackend{}

	update = flag.Bool("update", false, "update .golden files")
)

func TestFileBackendGet(t *testing.T) {
	gob.Register(testStruct{})
	type args struct {
		key string
	}
	tests := []struct {
		name      string
		args      args
		wantEntry *cache.Entry
		wantFound bool
	}{
		{
			name: "string",
			args: args{
				key: "get.string",
			},
			wantEntry: &cache.Entry{
				Meta: cache.Meta{},
				Data: "bar",
			},
			wantFound: true,
		},
		{
			name: "struct",
			args: args{
				key: "get.struct",
			},
			wantEntry: &cache.Entry{
				Meta: cache.Meta{},
				Data: testStruct{
					S: "string",
					B: true,
					I: -17,
				},
			},
			wantFound: true,
		},
		{
			name: "not-found",
			args: args{
				key: "get.not.found",
			},
			wantEntry: nil,
			wantFound: false,
		},
		{
			name: "invalid file content",
			args: args{
				key: "get.invalid",
			},
			wantEntry: nil,
			wantFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cache.NewFileBackend(filepath.Join("testdata", "file_backend"))

			gotEntry, gotFound := f.Get(tt.args.key)
			if !reflect.DeepEqual(gotEntry, tt.wantEntry) {
				t.Errorf("FileBackend.Get() gotEntry = %v, want %v", gotEntry, tt.wantEntry)
			}
			if gotFound != tt.wantFound {
				t.Errorf("FileBackend.Get() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func TestFileBackendSet(t *testing.T) {
	type args struct {
		key   string
		entry *cache.Entry
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "string",
			args: args{
				key: "set.string",
				entry: &cache.Entry{
					Meta: cache.Meta{},
					Data: "bar",
				},
			},
		},
		{
			name: "struct",
			args: args{
				key: "set.struct",
				entry: &cache.Entry{
					Meta: cache.Meta{},
					Data: testStruct{
						S: "string",
						B: true,
						I: -17,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedCacheFileName := filepath.Join("testdata", "file_backend", tt.args.key)
			defer func() { os.Remove(expectedCacheFileName) }()

			f := cache.NewFileBackend(filepath.Join("testdata", "file_backend"))
			f.Set(tt.args.key, tt.args.entry)

			written, err := ioutil.ReadFile(expectedCacheFileName)
			if err != nil {
				t.Fatal("cache entry not written")
			}

			if *update {
				t.Log("update golden file")
				f.Set(tt.args.key+".golden", tt.args.entry)
			}

			golden, err := ioutil.ReadFile(expectedCacheFileName + ".golden")
			if err != nil {
				t.Fatalf("failed reading .golden: %s", err)
			}
			if !bytes.Equal(written, golden) {
				t.Errorf("saved entry does not match .golden file")
			}
		})
	}
}

func TestFileBackend_Purge(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "purge",
			args: args{
				key: "purge.string",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cache.NewFileBackend(filepath.Join("testdata", "file_backend"))
			f.Set(tt.args.key, &cache.Entry{
				Meta: cache.Meta{},
				Data: "bar",
			})
			f.Purge(tt.args.key)

			if _, err := os.Stat(filepath.Join("testdata", "file_backend", tt.args.key)); err == nil {
				t.Error("cache entry was not deleted")
			}
		})
	}
}
