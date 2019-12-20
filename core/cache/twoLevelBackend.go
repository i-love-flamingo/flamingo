package cache

import (
	"errors"
	"fmt"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// twoLevelBackend instance representation
	twoLevelBackend struct {
		firstBackend  Backend
		secondBackend Backend
		logger        flamingo.Logger
	}

	TwoLevelBackendConfig struct {
		FirstLevel  Backend
		SecondLevel Backend
	}

	TwoLevelBackendFactory struct {
		logger flamingo.Logger
		config TwoLevelBackendConfig
	}
)

// Inject TwoLevelBackendFactory dependencies
func (f *TwoLevelBackendFactory) Inject(logger flamingo.Logger) *TwoLevelBackendFactory {
	f.logger = logger
	return f
}

// Inject TwoLevelBackendFactory dependencies
func (f *TwoLevelBackendFactory) SetConfig(config TwoLevelBackendConfig) *TwoLevelBackendFactory {
	f.config = config
	return f
}

// Inject TwoLevelBackendFactory dependencies
func (f *TwoLevelBackendFactory) Build() (Backend, error) {
	return &twoLevelBackend{
		firstBackend:  f.config.FirstLevel,
		secondBackend: f.config.SecondLevel,
		logger:        f.logger,
	}, nil
}

// Get entry by key
func (mb *twoLevelBackend) Get(key string) (entry *Entry, found bool) {
	entry, found = mb.firstBackend.Get(key)
	if found {
		return entry, found
	}

	entry, found = mb.secondBackend.Get(key)
	if found {
		go func() {
			_ = mb.firstBackend.Set(key, entry)
		}()
		return entry, found
	}

	return nil, false
}

// Set entry for key
func (mb *twoLevelBackend) Set(key string, entry *Entry) (err error) {
	errorList := []error{}

	err = mb.firstBackend.Set(key, entry)
	if err != nil {
		errorList = append(errorList, err)
		mb.logger.WithField("category", "twoLevelBackend").Error(fmt.Sprintf("Failed to set key %v with error %v", key, err))
	}

	err = mb.secondBackend.Set(key, entry)
	if err != nil {
		errorList = append(errorList, err)
		mb.logger.WithField("category", "twoLevelBackend").Error(fmt.Sprintf("Failed to set key %v with error %v", key, err))
	}

	if 2 == len(errorList) {
		return errors.New("all backends failed")
	}

	return nil
}

// Purge entry by key
func (mb *twoLevelBackend) Purge(key string) (err error) {
	errorList := []error{}

	err = mb.firstBackend.Purge(key)
	if err != nil {
		errorList = append(errorList, err)
		mb.logger.WithField("category", "twoLevelBackend").Error(fmt.Sprintf("Failed Purge with error %v", err))
	}

	err = mb.secondBackend.Purge(key)
	if err != nil {
		errorList = append(errorList, err)
		mb.logger.WithField("category", "twoLevelBackend").Error(fmt.Sprintf("Failed Purge with error %v", err))
	}

	if 0 != len(errorList) {
		return fmt.Errorf("Not all backends succeeded to Purge key %v, Errors: %v", key, errorList)
	}

	return nil
}

// Flush the whole cache
func (mb *twoLevelBackend) Flush() (err error) {
	errorList := []error{}

	err = mb.firstBackend.Flush()
	if err != nil {
		errorList = append(errorList, err)
		mb.logger.WithField("category", "twoLevelBackend").Error(fmt.Sprintf("Failed Flush error %v", err))
	}

	err = mb.secondBackend.Flush()
	if err != nil {
		errorList = append(errorList, err)
		mb.logger.WithField("category", "twoLevelBackend").Error(fmt.Sprintf("Failed Flush error %v", err))
	}

	if 0 != len(errorList) {
		return fmt.Errorf("Not all backends succeeded to Flush. Errors: %v", errorList)
	}

	return nil
}
