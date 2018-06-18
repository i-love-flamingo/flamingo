package router

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.8
The original interface "Router" can be found in go.aoe.com/flamingo/framework/event
*/
import (
	"context"
	"sync/atomic"
	"time"

	"github.com/gojuno/minimock"
	event "flamingo.me/flamingo/framework/event"

	testify_assert "github.com/stretchr/testify/assert"
)

//RouterMock implements go.aoe.com/flamingo/framework/event.Router
type RouterMock struct {
	t minimock.Tester

	DispatchFunc       func(ctx context.Context, p event.Event)
	DispatchCounter    uint64
	DispatchPreCounter uint64
	DispatchMock       mRouterMockDispatch
}

//NewRouterMock returns a mock for go.aoe.com/flamingo/framework/event.Router
func NewRouterMock(t minimock.Tester) *RouterMock {
	m := &RouterMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DispatchMock = mRouterMockDispatch{mock: m}

	return m
}

type mRouterMockDispatch struct {
	mock             *RouterMock
	mockExpectations *RouterMockDispatchParams
}

//RouterMockDispatchParams represents input parameters of the Router.Dispatch
type RouterMockDispatchParams struct {
	p event.Event
}

//Expect sets up expected params for the Router.Dispatch
func (m *mRouterMockDispatch) Expect(p event.Event) *mRouterMockDispatch {
	m.mockExpectations = &RouterMockDispatchParams{p}
	return m
}

//Return sets up a mock for Router.Dispatch to return Return's arguments
func (m *mRouterMockDispatch) Return() *RouterMock {
	m.mock.DispatchFunc = func(ctx context.Context, p event.Event) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Router.Dispatch method
func (m *mRouterMockDispatch) Set(f func(ctx context.Context, p event.Event)) *RouterMock {
	m.mock.DispatchFunc = f
	return m.mock
}

//Dispatch implements go.aoe.com/flamingo/framework/event.Router interface
func (m *RouterMock) Dispatch(ctx context.Context, p event.Event) {
	atomic.AddUint64(&m.DispatchPreCounter, 1)
	defer atomic.AddUint64(&m.DispatchCounter, 1)

	if m.DispatchMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.DispatchMock.mockExpectations, RouterMockDispatchParams{p},
			"Router.Dispatch got unexpected parameters")

		if m.DispatchFunc == nil {

			m.t.Fatal("No results are set for the RouterMock.Dispatch")

			return
		}
	}

	if m.DispatchFunc == nil {
		m.t.Fatal("Unexpected call to RouterMock.Dispatch")
		return
	}

	m.DispatchFunc(ctx, p)
}

//DispatchMinimockCounter returns a count of RouterMock.DispatchFunc invocations
func (m *RouterMock) DispatchMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DispatchCounter)
}

//DispatchMinimockPreCounter returns the value of RouterMock.Dispatch invocations
func (m *RouterMock) DispatchMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DispatchPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RouterMock) ValidateCallCounters() {

	if m.DispatchFunc != nil && atomic.LoadUint64(&m.DispatchCounter) == 0 {
		m.t.Fatal("Expected call to RouterMock.Dispatch")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *RouterMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *RouterMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *RouterMock) MinimockFinish() {

	if m.DispatchFunc != nil && atomic.LoadUint64(&m.DispatchCounter) == 0 {
		m.t.Fatal("Expected call to RouterMock.Dispatch")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *RouterMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *RouterMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.DispatchFunc == nil || atomic.LoadUint64(&m.DispatchCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.DispatchFunc != nil && atomic.LoadUint64(&m.DispatchCounter) == 0 {
				m.t.Error("Expected call to RouterMock.Dispatch")
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
func (m *RouterMock) AllMocksCalled() bool {

	if m.DispatchFunc != nil && atomic.LoadUint64(&m.DispatchCounter) == 0 {
		return false
	}

	return true
}
