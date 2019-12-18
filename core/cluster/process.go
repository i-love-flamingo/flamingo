package cluster

import "time"

type (

	SingleProcess interface {
		//Do runs the given function if there is no other function with that key running
		Do(key string, fn func()) (bool, error)
	}

	DistributedLock interface {
		//GetLock - call to requests a distributed lock - if not possible. returns "false" in case Lock is not available without blocking
		GetLock(key string, max *time.Duration) (bool, error)
		//Free the distributed Lock
		Unlock(key string)
	}

	SingleProcessImpl struct {
		lock DistributedLock
	}
)

//Do runs the func if Lock can be obtained
func (s *SingleProcessImpl) Do(key string, fn func()) (bool, error) {
	haslock, err := s.lock.GetLock(key, nil)
	if err != nil {
		return false, err
	}
	if haslock{
		defer s.lock.Unlock(key)
	}
	if haslock{
		fn()
		return true, nil
	}
	return false, nil
}