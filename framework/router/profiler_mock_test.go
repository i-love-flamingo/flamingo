package router

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.8
The original interface "Profiler" can be found in go.aoe.com/flamingo/framework/profiler
*/
import (
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	profiler "go.aoe.com/flamingo/framework/profiler"

	testify_assert "github.com/stretchr/testify/assert"
)

//ProfilerMock implements go.aoe.com/flamingo/framework/profiler.Profiler
type ProfilerMock struct {
	t minimock.Tester

	ProfileFunc       func(p string, p1 string) (r profiler.ProfileFinishFunc)
	ProfileCounter    uint64
	ProfilePreCounter uint64
	ProfileMock       mProfilerMockProfile
}

//NewProfilerMock returns a mock for go.aoe.com/flamingo/framework/profiler.Profiler
func NewProfilerMock(t minimock.Tester) *ProfilerMock {
	m := &ProfilerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.ProfileMock = mProfilerMockProfile{mock: m}

	return m
}

type mProfilerMockProfile struct {
	mock             *ProfilerMock
	mockExpectations *ProfilerMockProfileParams
}

//ProfilerMockProfileParams represents input parameters of the Profiler.Profile
type ProfilerMockProfileParams struct {
	p  string
	p1 string
}

//Expect sets up expected params for the Profiler.Profile
func (m *mProfilerMockProfile) Expect(p string, p1 string) *mProfilerMockProfile {
	m.mockExpectations = &ProfilerMockProfileParams{p, p1}
	return m
}

//Return sets up a mock for Profiler.Profile to return Return's arguments
func (m *mProfilerMockProfile) Return(r profiler.ProfileFinishFunc) *ProfilerMock {
	m.mock.ProfileFunc = func(p string, p1 string) profiler.ProfileFinishFunc {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Profiler.Profile method
func (m *mProfilerMockProfile) Set(f func(p string, p1 string) (r profiler.ProfileFinishFunc)) *ProfilerMock {
	m.mock.ProfileFunc = f
	return m.mock
}

//Profile implements go.aoe.com/flamingo/framework/profiler.Profiler interface
func (m *ProfilerMock) Profile(p string, p1 string) (r profiler.ProfileFinishFunc) {
	atomic.AddUint64(&m.ProfilePreCounter, 1)
	defer atomic.AddUint64(&m.ProfileCounter, 1)

	if m.ProfileMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ProfileMock.mockExpectations, ProfilerMockProfileParams{p, p1},
			"Profiler.Profile got unexpected parameters")

		if m.ProfileFunc == nil {

			m.t.Fatal("No results are set for the ProfilerMock.Profile")

			return
		}
	}

	if m.ProfileFunc == nil {
		m.t.Fatal("Unexpected call to ProfilerMock.Profile")
		return
	}

	return m.ProfileFunc(p, p1)
}

//ProfileMinimockCounter returns a count of ProfilerMock.ProfileFunc invocations
func (m *ProfilerMock) ProfileMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ProfileCounter)
}

//ProfileMinimockPreCounter returns the value of ProfilerMock.Profile invocations
func (m *ProfilerMock) ProfileMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ProfilePreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProfilerMock) ValidateCallCounters() {

	if m.ProfileFunc != nil && atomic.LoadUint64(&m.ProfileCounter) == 0 {
		m.t.Fatal("Expected call to ProfilerMock.Profile")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ProfilerMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ProfilerMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ProfilerMock) MinimockFinish() {

	if m.ProfileFunc != nil && atomic.LoadUint64(&m.ProfileCounter) == 0 {
		m.t.Fatal("Expected call to ProfilerMock.Profile")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ProfilerMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ProfilerMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.ProfileFunc == nil || atomic.LoadUint64(&m.ProfileCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.ProfileFunc != nil && atomic.LoadUint64(&m.ProfileCounter) == 0 {
				m.t.Error("Expected call to ProfilerMock.Profile")
			}

			m.t.Fatalf("Some mocks were not called on time: %s", timeout)
			return
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

//AllMocksCalled returns true if all mocked methods were called before the execution of AllMocksCalled,
//it can be used with assert/require, i.e. assert.True(mock.AllMocksCalled())
func (m *ProfilerMock) AllMocksCalled() bool {

	if m.ProfileFunc != nil && atomic.LoadUint64(&m.ProfileCounter) == 0 {
		return false
	}

	return true
}
