package cache

import (
	"errors"
	"fmt"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// TwoLevelBackend instance representation
	TwoLevelBackend struct {
		backends []Backend
		logger   flamingo.Logger
	}
	// TwoLevelBackendOptions representation
	TwoLevelBackendOptions struct {
		Backends []Backend
	}
)

// NewTwoLevelBackend creates a TwoLevelBackend isntance
func NewTwoLevelBackend(options TwoLevelBackendOptions, logger flamingo.Logger) *TwoLevelBackend {
	return &TwoLevelBackend{
		backends: options.Backends,
		logger:   logger,
	}
}

// Inject TwoLevelBackend dependencies
func (mb *TwoLevelBackend) Inject(logger flamingo.Logger) {
	mb.logger = logger
}

// Get entry by key
func (mb *TwoLevelBackend) Get(key string) (entry *Entry, found bool) {
	for _, backend := range mb.backends {
		entry, found := backend.Get(key)
		if found {
			return entry, found
		}
	}

	return nil, false
}

// Set entry for key
func (mb *TwoLevelBackend) Set(key string, entry *Entry) error {
	errorList := []error{}
	for _, backend := range mb.backends {
		err := backend.Set(key, entry)
		if err != nil {
			errorList = append(errorList, err)
			mb.logger.WithField("category", "twoLevelBackend").Error(fmt.Sprintf("Failed to set key %v with error %v", key, err))
		}
	}

	if len(mb.backends) == len(errorList) {
		return errors.New("all backends failed")
	}

	return nil
}

// Purge entry by key
func (mb *TwoLevelBackend) Purge(key string) error {
	errorList := []error{}
	for _, backend := range mb.backends {
		err := backend.Purge(key)
		if err != nil {
			errorList = append(errorList, err)
			mb.logger.WithField("category", "twoLevelBackend").Error(fmt.Sprintf("Failed Purge with error %v", err))
		}
	}

	if 0 != len(errorList) {
		return fmt.Errorf("Not all backends succeeded to Purge key %v, Errors: %v", key, errorList)
	}

	return nil
}

// Flush the whole cache
func (mb *TwoLevelBackend) Flush() error {
	errorList := []error{}
	for _, backend := range mb.backends {
		err := backend.Flush()
		if err != nil {
			errorList = append(errorList, err)
			mb.logger.WithField("category", "twoLevelBackend").Error(fmt.Sprintf("Failed Flush error %v", err))
		}
	}

	if 0 != len(errorList) {
		return fmt.Errorf("Not all backends succeeded to Flush. Errors: %v", errorList)
	}

	return nil
}
