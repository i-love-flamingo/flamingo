package cache

import (
	"errors"
	"fmt"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type (
	// TwoLevelBackend instance representation
	TwoLevelBackend struct {
		firstBackend  Backend
		secondBackend Backend
		logger        flamingo.Logger
	}
)

// NewTwoLevelBackend creates a TwoLevelBackend isntance
func NewTwoLevelBackend(firstBackend Backend, secondBackend Backend) *TwoLevelBackend {
	return &TwoLevelBackend{
		firstBackend:  firstBackend,
		secondBackend: secondBackend,
		logger:        flamingo.NullLogger{},
	}
}

// Inject TwoLevelBackend dependencies
func (mb *TwoLevelBackend) Inject(logger flamingo.Logger) {
	mb.logger = logger
}

// Get entry by key
func (mb *TwoLevelBackend) Get(key string) (entry *Entry, found bool) {
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
func (mb *TwoLevelBackend) Set(key string, entry *Entry) (err error) {
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
func (mb *TwoLevelBackend) Purge(key string) (err error) {
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
func (mb *TwoLevelBackend) Flush() (err error) {
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
