package cache

import (
	"errors"
	"fmt"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// MultiLevelBackend instance representation
	MultiLevelBackend struct {
		backends []Backend
		logger   flamingo.Logger
	}
	// MultiLevelBackendOptions representation
	MultiLevelBackendOptions struct {
		Backends []Backend
	}
)

// NewMultiLevelBackend creates a MultiLevelBackend isntance
func NewMultiLevelBackend(options MultiLevelBackendOptions, logger flamingo.Logger) *MultiLevelBackend {
	return &MultiLevelBackend{
		backends: options.Backends,
		logger:   logger,
	}
}

// Inject MultiLevelBackend dependencies
func (mb *MultiLevelBackend) Inject(logger flamingo.Logger) {
	mb.logger = logger
}

// Get entry by key
func (mb *MultiLevelBackend) Get(key string) (entry *Entry, found bool) {
	for _, backend := range mb.backends {
		entry, found := backend.Get(key)
		if found {
			return entry, found
		}
	}

	return nil, false
}

// Set entry for key
func (mb *MultiLevelBackend) Set(key string, entry *Entry) error {
	errorList := []error{}
	for _, backend := range mb.backends {
		err := backend.Set(key, entry)
		if err != nil {
			errorList = append(errorList, err)
			mb.logger.WithField("category", "multiLevelBackend").Error(fmt.Sprintf("Failed to set key %v with error %v", key, err))
		}
	}

	if len(mb.backends) == len(errorList) {
		return errors.New("all backends failed")
	}

	return nil
}

// Purge entry by key
func (mb *MultiLevelBackend) Purge(key string) error {
	errorList := []error{}
	for _, backend := range mb.backends {
		err := backend.Purge(key)
		if err != nil {
			errorList = append(errorList, err)
			mb.logger.WithField("category", "multiLevelBackend").Error(fmt.Sprintf("Failed Purge with error %v", err))
		}
	}

	if 0 != len(errorList) {
		return errors.New("not all backends succeeded")
	}

	return nil
}

// PurgeTags purges erntries by tags
func (mb *MultiLevelBackend) PurgeTags(tags []string) error {
	if mb.FlushSupport() {
		return errors.New("Not supported by this backend")
	}

	errorList := []error{}
	for _, backend := range mb.backends {
		err := backend.PurgeTags(tags)
		if err != nil {
			errorList = append(errorList, err)
			mb.logger.WithField("category", "multiLevelBackend").Error(fmt.Sprintf("Failed PurgeTags %v with error %v", tags, err))
		}
	}

	if 0 != len(errorList) {
		return errors.New("not all backends succeeded")
	}

	return nil
}

// Flush the whole cache
func (mb *MultiLevelBackend) Flush() error {
	if !mb.FlushSupport() {
		return errors.New("not supported by this backend")
	}

	errorList := []error{}
	for _, backend := range mb.backends {
		err := backend.Flush()
		if err != nil {
			errorList = append(errorList, err)
			mb.logger.WithField("category", "multiLevelBackend").Error(fmt.Sprintf("Failed Flush error %v", err))
		}
	}

	if 0 != len(errorList) {
		return errors.New("ot all backends succeeded")
	}

	return nil
}

// FlushSupport returns true, if all registered backends supports the flush method
func (mb *MultiLevelBackend) FlushSupport() bool {
	for _, backend := range mb.backends {
		if !backend.FlushSupport() {
			return false
		}
	}

	return true
}

// TagSupport returns true, if all registered backends supports tags
func (mb *MultiLevelBackend) TagSupport() bool {
	for _, backend := range mb.backends {
		if !backend.TagSupport() {
			return false
		}
	}

	return true
}
