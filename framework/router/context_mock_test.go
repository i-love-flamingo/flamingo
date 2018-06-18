package router

/*
DO NOT EDIT!
This code was generated automatically using github.com/gojuno/minimock v1.8
The original interface "Context" can be found in go.aoe.com/flamingo/framework/web
*/
import (
	http "net/http"
	"sync/atomic"
	time "time"

	"github.com/gojuno/minimock"
	sessions "github.com/gorilla/sessions"
	event "flamingo.me/flamingo/framework/event"
	profiler "flamingo.me/flamingo/framework/profiler"
	web "flamingo.me/flamingo/framework/web"

	testify_assert "github.com/stretchr/testify/assert"
)

//ContextMock implements go.aoe.com/flamingo/framework/web.Context
type ContextMock struct {
	t minimock.Tester

	DeadlineFunc       func() (r time.Time, r1 bool)
	DeadlineCounter    uint64
	DeadlinePreCounter uint64
	DeadlineMock       mContextMockDeadline

	DoneFunc func() (r <-chan struct {
	})
	DoneCounter    uint64
	DonePreCounter uint64
	DoneMock       mContextMockDone

	ErrFunc       func() (r error)
	ErrCounter    uint64
	ErrPreCounter uint64
	ErrMock       mContextMockErr

	EventRouterFunc       func() (r event.Router)
	EventRouterCounter    uint64
	EventRouterPreCounter uint64
	EventRouterMock       mContextMockEventRouter

	FormFunc       func(p string) (r []string, r1 error)
	FormCounter    uint64
	FormPreCounter uint64
	FormMock       mContextMockForm

	Form1Func       func(p string) (r string, r1 error)
	Form1Counter    uint64
	Form1PreCounter uint64
	Form1Mock       mContextMockForm1

	FormAllFunc       func() (r map[string][]string)
	FormAllCounter    uint64
	FormAllPreCounter uint64
	FormAllMock       mContextMockFormAll

	IDFunc       func() (r string)
	IDCounter    uint64
	IDPreCounter uint64
	IDMock       mContextMockID

	LoadParamsFunc       func(p map[string]string)
	LoadParamsCounter    uint64
	LoadParamsPreCounter uint64
	LoadParamsMock       mContextMockLoadParams

	MustFormFunc       func(p string) (r []string)
	MustFormCounter    uint64
	MustFormPreCounter uint64
	MustFormMock       mContextMockMustForm

	MustForm1Func       func(p string) (r string)
	MustForm1Counter    uint64
	MustForm1PreCounter uint64
	MustForm1Mock       mContextMockMustForm1

	MustParam1Func       func(p string) (r string)
	MustParam1Counter    uint64
	MustParam1PreCounter uint64
	MustParam1Mock       mContextMockMustParam1

	MustQueryFunc       func(p string) (r []string)
	MustQueryCounter    uint64
	MustQueryPreCounter uint64
	MustQueryMock       mContextMockMustQuery

	MustQuery1Func       func(p string) (r string)
	MustQuery1Counter    uint64
	MustQuery1PreCounter uint64
	MustQuery1Mock       mContextMockMustQuery1

	Param1Func       func(p string) (r string, r1 error)
	Param1Counter    uint64
	Param1PreCounter uint64
	Param1Mock       mContextMockParam1

	ParamAllFunc       func() (r map[string]string)
	ParamAllCounter    uint64
	ParamAllPreCounter uint64
	ParamAllMock       mContextMockParamAll

	ProfileFunc       func(p string, p1 string) (r profiler.ProfileFinishFunc)
	ProfileCounter    uint64
	ProfilePreCounter uint64
	ProfileMock       mContextMockProfile

	ProfilerFunc       func() (r profiler.Profiler)
	ProfilerCounter    uint64
	ProfilerPreCounter uint64
	ProfilerMock       mContextMockProfiler

	PushFunc       func(p string, p1 *http.PushOptions) (r error)
	PushCounter    uint64
	PushPreCounter uint64
	PushMock       mContextMockPush

	QueryFunc       func(p string) (r []string, r1 error)
	QueryCounter    uint64
	QueryPreCounter uint64
	QueryMock       mContextMockQuery

	Query1Func       func(p string) (r string, r1 error)
	Query1Counter    uint64
	Query1PreCounter uint64
	Query1Mock       mContextMockQuery1

	QueryAllFunc       func() (r map[string][]string)
	QueryAllCounter    uint64
	QueryAllPreCounter uint64
	QueryAllMock       mContextMockQueryAll

	RequestFunc       func() (r *http.Request)
	RequestCounter    uint64
	RequestPreCounter uint64
	RequestMock       mContextMockRequest

	SessionFunc       func() (r *sessions.Session)
	SessionCounter    uint64
	SessionPreCounter uint64
	SessionMock       mContextMockSession

	ValueFunc       func(p interface{}) (r interface{})
	ValueCounter    uint64
	ValuePreCounter uint64
	ValueMock       mContextMockValue

	WithValueFunc       func(p interface{}, p1 interface{}) (r web.Context)
	WithValueCounter    uint64
	WithValuePreCounter uint64
	WithValueMock       mContextMockWithValue

	WithVarsFunc       func(p map[string]string) (r web.Context)
	WithVarsCounter    uint64
	WithVarsPreCounter uint64
	WithVarsMock       mContextMockWithVars
}

//NewContextMock returns a mock for go.aoe.com/flamingo/framework/web.Context
func NewContextMock(t minimock.Tester) *ContextMock {
	m := &ContextMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.DeadlineMock = mContextMockDeadline{mock: m}
	m.DoneMock = mContextMockDone{mock: m}
	m.ErrMock = mContextMockErr{mock: m}
	m.EventRouterMock = mContextMockEventRouter{mock: m}
	m.FormMock = mContextMockForm{mock: m}
	m.Form1Mock = mContextMockForm1{mock: m}
	m.FormAllMock = mContextMockFormAll{mock: m}
	m.IDMock = mContextMockID{mock: m}
	m.LoadParamsMock = mContextMockLoadParams{mock: m}
	m.MustFormMock = mContextMockMustForm{mock: m}
	m.MustForm1Mock = mContextMockMustForm1{mock: m}
	m.MustParam1Mock = mContextMockMustParam1{mock: m}
	m.MustQueryMock = mContextMockMustQuery{mock: m}
	m.MustQuery1Mock = mContextMockMustQuery1{mock: m}
	m.Param1Mock = mContextMockParam1{mock: m}
	m.ParamAllMock = mContextMockParamAll{mock: m}
	m.ProfileMock = mContextMockProfile{mock: m}
	m.ProfilerMock = mContextMockProfiler{mock: m}
	m.PushMock = mContextMockPush{mock: m}
	m.QueryMock = mContextMockQuery{mock: m}
	m.Query1Mock = mContextMockQuery1{mock: m}
	m.QueryAllMock = mContextMockQueryAll{mock: m}
	m.RequestMock = mContextMockRequest{mock: m}
	m.SessionMock = mContextMockSession{mock: m}
	m.ValueMock = mContextMockValue{mock: m}
	m.WithValueMock = mContextMockWithValue{mock: m}
	m.WithVarsMock = mContextMockWithVars{mock: m}

	return m
}

type mContextMockDeadline struct {
	mock *ContextMock
}

//Return sets up a mock for Context.Deadline to return Return's arguments
func (m *mContextMockDeadline) Return(r time.Time, r1 bool) *ContextMock {
	m.mock.DeadlineFunc = func() (time.Time, bool) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Deadline method
func (m *mContextMockDeadline) Set(f func() (r time.Time, r1 bool)) *ContextMock {
	m.mock.DeadlineFunc = f
	return m.mock
}

//Deadline implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Deadline() (r time.Time, r1 bool) {
	atomic.AddUint64(&m.DeadlinePreCounter, 1)
	defer atomic.AddUint64(&m.DeadlineCounter, 1)

	if m.DeadlineFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.Deadline")
		return
	}

	return m.DeadlineFunc()
}

//DeadlineMinimockCounter returns a count of ContextMock.DeadlineFunc invocations
func (m *ContextMock) DeadlineMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DeadlineCounter)
}

//DeadlineMinimockPreCounter returns the value of ContextMock.Deadline invocations
func (m *ContextMock) DeadlineMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DeadlinePreCounter)
}

type mContextMockDone struct {
	mock *ContextMock
}

//Return sets up a mock for Context.Done to return Return's arguments
func (m *mContextMockDone) Return(r <-chan struct {
}) *ContextMock {
	m.mock.DoneFunc = func() <-chan struct {
	} {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Done method
func (m *mContextMockDone) Set(f func() (r <-chan struct {
})) *ContextMock {
	m.mock.DoneFunc = f
	return m.mock
}

//Done implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Done() (r <-chan struct {
}) {
	atomic.AddUint64(&m.DonePreCounter, 1)
	defer atomic.AddUint64(&m.DoneCounter, 1)

	if m.DoneFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.Done")
		return
	}

	return m.DoneFunc()
}

//DoneMinimockCounter returns a count of ContextMock.DoneFunc invocations
func (m *ContextMock) DoneMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.DoneCounter)
}

//DoneMinimockPreCounter returns the value of ContextMock.Done invocations
func (m *ContextMock) DoneMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.DonePreCounter)
}

type mContextMockErr struct {
	mock *ContextMock
}

//Return sets up a mock for Context.Err to return Return's arguments
func (m *mContextMockErr) Return(r error) *ContextMock {
	m.mock.ErrFunc = func() error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Err method
func (m *mContextMockErr) Set(f func() (r error)) *ContextMock {
	m.mock.ErrFunc = f
	return m.mock
}

//Err implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Err() (r error) {
	atomic.AddUint64(&m.ErrPreCounter, 1)
	defer atomic.AddUint64(&m.ErrCounter, 1)

	if m.ErrFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.Err")
		return
	}

	return m.ErrFunc()
}

//ErrMinimockCounter returns a count of ContextMock.ErrFunc invocations
func (m *ContextMock) ErrMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ErrCounter)
}

//ErrMinimockPreCounter returns the value of ContextMock.Err invocations
func (m *ContextMock) ErrMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ErrPreCounter)
}

type mContextMockEventRouter struct {
	mock *ContextMock
}

//Return sets up a mock for Context.EventRouter to return Return's arguments
func (m *mContextMockEventRouter) Return(r event.Router) *ContextMock {
	m.mock.EventRouterFunc = func() event.Router {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.EventRouter method
func (m *mContextMockEventRouter) Set(f func() (r event.Router)) *ContextMock {
	m.mock.EventRouterFunc = f
	return m.mock
}

//EventRouter implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) EventRouter() (r event.Router) {
	atomic.AddUint64(&m.EventRouterPreCounter, 1)
	defer atomic.AddUint64(&m.EventRouterCounter, 1)

	if m.EventRouterFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.EventRouter")
		return
	}

	return m.EventRouterFunc()
}

//EventRouterMinimockCounter returns a count of ContextMock.EventRouterFunc invocations
func (m *ContextMock) EventRouterMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.EventRouterCounter)
}

//EventRouterMinimockPreCounter returns the value of ContextMock.EventRouter invocations
func (m *ContextMock) EventRouterMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.EventRouterPreCounter)
}

type mContextMockForm struct {
	mock             *ContextMock
	mockExpectations *ContextMockFormParams
}

//ContextMockFormParams represents input parameters of the Context.Form
type ContextMockFormParams struct {
	p string
}

//Expect sets up expected params for the Context.Form
func (m *mContextMockForm) Expect(p string) *mContextMockForm {
	m.mockExpectations = &ContextMockFormParams{p}
	return m
}

//Return sets up a mock for Context.Form to return Return's arguments
func (m *mContextMockForm) Return(r []string, r1 error) *ContextMock {
	m.mock.FormFunc = func(p string) ([]string, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Form method
func (m *mContextMockForm) Set(f func(p string) (r []string, r1 error)) *ContextMock {
	m.mock.FormFunc = f
	return m.mock
}

//Form implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Form(p string) (r []string, r1 error) {
	atomic.AddUint64(&m.FormPreCounter, 1)
	defer atomic.AddUint64(&m.FormCounter, 1)

	if m.FormMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.FormMock.mockExpectations, ContextMockFormParams{p},
			"Context.Form got unexpected parameters")

		if m.FormFunc == nil {

			m.t.Fatal("No results are set for the ContextMock.Form")

			return
		}
	}

	if m.FormFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.Form")
		return
	}

	return m.FormFunc(p)
}

//FormMinimockCounter returns a count of ContextMock.FormFunc invocations
func (m *ContextMock) FormMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FormCounter)
}

//FormMinimockPreCounter returns the value of ContextMock.Form invocations
func (m *ContextMock) FormMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FormPreCounter)
}

type mContextMockForm1 struct {
	mock             *ContextMock
	mockExpectations *ContextMockForm1Params
}

//ContextMockForm1Params represents input parameters of the Context.Form1
type ContextMockForm1Params struct {
	p string
}

//Expect sets up expected params for the Context.Form1
func (m *mContextMockForm1) Expect(p string) *mContextMockForm1 {
	m.mockExpectations = &ContextMockForm1Params{p}
	return m
}

//Return sets up a mock for Context.Form1 to return Return's arguments
func (m *mContextMockForm1) Return(r string, r1 error) *ContextMock {
	m.mock.Form1Func = func(p string) (string, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Form1 method
func (m *mContextMockForm1) Set(f func(p string) (r string, r1 error)) *ContextMock {
	m.mock.Form1Func = f
	return m.mock
}

//Form1 implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Form1(p string) (r string, r1 error) {
	atomic.AddUint64(&m.Form1PreCounter, 1)
	defer atomic.AddUint64(&m.Form1Counter, 1)

	if m.Form1Mock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.Form1Mock.mockExpectations, ContextMockForm1Params{p},
			"Context.Form1 got unexpected parameters")

		if m.Form1Func == nil {

			m.t.Fatal("No results are set for the ContextMock.Form1")

			return
		}
	}

	if m.Form1Func == nil {
		m.t.Fatal("Unexpected call to ContextMock.Form1")
		return
	}

	return m.Form1Func(p)
}

//Form1MinimockCounter returns a count of ContextMock.Form1Func invocations
func (m *ContextMock) Form1MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.Form1Counter)
}

//Form1MinimockPreCounter returns the value of ContextMock.Form1 invocations
func (m *ContextMock) Form1MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.Form1PreCounter)
}

type mContextMockFormAll struct {
	mock *ContextMock
}

//Return sets up a mock for Context.FormAll to return Return's arguments
func (m *mContextMockFormAll) Return(r map[string][]string) *ContextMock {
	m.mock.FormAllFunc = func() map[string][]string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.FormAll method
func (m *mContextMockFormAll) Set(f func() (r map[string][]string)) *ContextMock {
	m.mock.FormAllFunc = f
	return m.mock
}

//FormAll implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) FormAll() (r map[string][]string) {
	atomic.AddUint64(&m.FormAllPreCounter, 1)
	defer atomic.AddUint64(&m.FormAllCounter, 1)

	if m.FormAllFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.FormAll")
		return
	}

	return m.FormAllFunc()
}

//FormAllMinimockCounter returns a count of ContextMock.FormAllFunc invocations
func (m *ContextMock) FormAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.FormAllCounter)
}

//FormAllMinimockPreCounter returns the value of ContextMock.FormAll invocations
func (m *ContextMock) FormAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.FormAllPreCounter)
}

type mContextMockID struct {
	mock *ContextMock
}

//Return sets up a mock for Context.ID to return Return's arguments
func (m *mContextMockID) Return(r string) *ContextMock {
	m.mock.IDFunc = func() string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.ID method
func (m *mContextMockID) Set(f func() (r string)) *ContextMock {
	m.mock.IDFunc = f
	return m.mock
}

//ID implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) ID() (r string) {
	atomic.AddUint64(&m.IDPreCounter, 1)
	defer atomic.AddUint64(&m.IDCounter, 1)

	if m.IDFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.ID")
		return
	}

	return m.IDFunc()
}

//IDMinimockCounter returns a count of ContextMock.IDFunc invocations
func (m *ContextMock) IDMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.IDCounter)
}

//IDMinimockPreCounter returns the value of ContextMock.ID invocations
func (m *ContextMock) IDMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.IDPreCounter)
}

type mContextMockLoadParams struct {
	mock             *ContextMock
	mockExpectations *ContextMockLoadParamsParams
}

//ContextMockLoadParamsParams represents input parameters of the Context.LoadParams
type ContextMockLoadParamsParams struct {
	p map[string]string
}

//Expect sets up expected params for the Context.LoadParams
func (m *mContextMockLoadParams) Expect(p map[string]string) *mContextMockLoadParams {
	m.mockExpectations = &ContextMockLoadParamsParams{p}
	return m
}

//Return sets up a mock for Context.LoadParams to return Return's arguments
func (m *mContextMockLoadParams) Return() *ContextMock {
	m.mock.LoadParamsFunc = func(p map[string]string) {
		return
	}
	return m.mock
}

//Set uses given function f as a mock of Context.LoadParams method
func (m *mContextMockLoadParams) Set(f func(p map[string]string)) *ContextMock {
	m.mock.LoadParamsFunc = f
	return m.mock
}

//LoadParams implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) LoadParams(p map[string]string) {
	atomic.AddUint64(&m.LoadParamsPreCounter, 1)
	defer atomic.AddUint64(&m.LoadParamsCounter, 1)

	if m.LoadParamsMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.LoadParamsMock.mockExpectations, ContextMockLoadParamsParams{p},
			"Context.LoadParams got unexpected parameters")

		if m.LoadParamsFunc == nil {

			m.t.Fatal("No results are set for the ContextMock.LoadParams")

			return
		}
	}

	if m.LoadParamsFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.LoadParams")
		return
	}

	m.LoadParamsFunc(p)
}

//LoadParamsMinimockCounter returns a count of ContextMock.LoadParamsFunc invocations
func (m *ContextMock) LoadParamsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.LoadParamsCounter)
}

//LoadParamsMinimockPreCounter returns the value of ContextMock.LoadParams invocations
func (m *ContextMock) LoadParamsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.LoadParamsPreCounter)
}

type mContextMockMustForm struct {
	mock             *ContextMock
	mockExpectations *ContextMockMustFormParams
}

//ContextMockMustFormParams represents input parameters of the Context.MustForm
type ContextMockMustFormParams struct {
	p string
}

//Expect sets up expected params for the Context.MustForm
func (m *mContextMockMustForm) Expect(p string) *mContextMockMustForm {
	m.mockExpectations = &ContextMockMustFormParams{p}
	return m
}

//Return sets up a mock for Context.MustForm to return Return's arguments
func (m *mContextMockMustForm) Return(r []string) *ContextMock {
	m.mock.MustFormFunc = func(p string) []string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.MustForm method
func (m *mContextMockMustForm) Set(f func(p string) (r []string)) *ContextMock {
	m.mock.MustFormFunc = f
	return m.mock
}

//MustForm implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) MustForm(p string) (r []string) {
	atomic.AddUint64(&m.MustFormPreCounter, 1)
	defer atomic.AddUint64(&m.MustFormCounter, 1)

	if m.MustFormMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.MustFormMock.mockExpectations, ContextMockMustFormParams{p},
			"Context.MustForm got unexpected parameters")

		if m.MustFormFunc == nil {

			m.t.Fatal("No results are set for the ContextMock.MustForm")

			return
		}
	}

	if m.MustFormFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.MustForm")
		return
	}

	return m.MustFormFunc(p)
}

//MustFormMinimockCounter returns a count of ContextMock.MustFormFunc invocations
func (m *ContextMock) MustFormMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MustFormCounter)
}

//MustFormMinimockPreCounter returns the value of ContextMock.MustForm invocations
func (m *ContextMock) MustFormMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MustFormPreCounter)
}

type mContextMockMustForm1 struct {
	mock             *ContextMock
	mockExpectations *ContextMockMustForm1Params
}

//ContextMockMustForm1Params represents input parameters of the Context.MustForm1
type ContextMockMustForm1Params struct {
	p string
}

//Expect sets up expected params for the Context.MustForm1
func (m *mContextMockMustForm1) Expect(p string) *mContextMockMustForm1 {
	m.mockExpectations = &ContextMockMustForm1Params{p}
	return m
}

//Return sets up a mock for Context.MustForm1 to return Return's arguments
func (m *mContextMockMustForm1) Return(r string) *ContextMock {
	m.mock.MustForm1Func = func(p string) string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.MustForm1 method
func (m *mContextMockMustForm1) Set(f func(p string) (r string)) *ContextMock {
	m.mock.MustForm1Func = f
	return m.mock
}

//MustForm1 implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) MustForm1(p string) (r string) {
	atomic.AddUint64(&m.MustForm1PreCounter, 1)
	defer atomic.AddUint64(&m.MustForm1Counter, 1)

	if m.MustForm1Mock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.MustForm1Mock.mockExpectations, ContextMockMustForm1Params{p},
			"Context.MustForm1 got unexpected parameters")

		if m.MustForm1Func == nil {

			m.t.Fatal("No results are set for the ContextMock.MustForm1")

			return
		}
	}

	if m.MustForm1Func == nil {
		m.t.Fatal("Unexpected call to ContextMock.MustForm1")
		return
	}

	return m.MustForm1Func(p)
}

//MustForm1MinimockCounter returns a count of ContextMock.MustForm1Func invocations
func (m *ContextMock) MustForm1MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MustForm1Counter)
}

//MustForm1MinimockPreCounter returns the value of ContextMock.MustForm1 invocations
func (m *ContextMock) MustForm1MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MustForm1PreCounter)
}

type mContextMockMustParam1 struct {
	mock             *ContextMock
	mockExpectations *ContextMockMustParam1Params
}

//ContextMockMustParam1Params represents input parameters of the Context.MustParam1
type ContextMockMustParam1Params struct {
	p string
}

//Expect sets up expected params for the Context.MustParam1
func (m *mContextMockMustParam1) Expect(p string) *mContextMockMustParam1 {
	m.mockExpectations = &ContextMockMustParam1Params{p}
	return m
}

//Return sets up a mock for Context.MustParam1 to return Return's arguments
func (m *mContextMockMustParam1) Return(r string) *ContextMock {
	m.mock.MustParam1Func = func(p string) string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.MustParam1 method
func (m *mContextMockMustParam1) Set(f func(p string) (r string)) *ContextMock {
	m.mock.MustParam1Func = f
	return m.mock
}

//MustParam1 implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) MustParam1(p string) (r string) {
	atomic.AddUint64(&m.MustParam1PreCounter, 1)
	defer atomic.AddUint64(&m.MustParam1Counter, 1)

	if m.MustParam1Mock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.MustParam1Mock.mockExpectations, ContextMockMustParam1Params{p},
			"Context.MustParam1 got unexpected parameters")

		if m.MustParam1Func == nil {

			m.t.Fatal("No results are set for the ContextMock.MustParam1")

			return
		}
	}

	if m.MustParam1Func == nil {
		m.t.Fatal("Unexpected call to ContextMock.MustParam1")
		return
	}

	return m.MustParam1Func(p)
}

//MustParam1MinimockCounter returns a count of ContextMock.MustParam1Func invocations
func (m *ContextMock) MustParam1MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MustParam1Counter)
}

//MustParam1MinimockPreCounter returns the value of ContextMock.MustParam1 invocations
func (m *ContextMock) MustParam1MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MustParam1PreCounter)
}

type mContextMockMustQuery struct {
	mock             *ContextMock
	mockExpectations *ContextMockMustQueryParams
}

//ContextMockMustQueryParams represents input parameters of the Context.MustQuery
type ContextMockMustQueryParams struct {
	p string
}

//Expect sets up expected params for the Context.MustQuery
func (m *mContextMockMustQuery) Expect(p string) *mContextMockMustQuery {
	m.mockExpectations = &ContextMockMustQueryParams{p}
	return m
}

//Return sets up a mock for Context.MustQuery to return Return's arguments
func (m *mContextMockMustQuery) Return(r []string) *ContextMock {
	m.mock.MustQueryFunc = func(p string) []string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.MustQuery method
func (m *mContextMockMustQuery) Set(f func(p string) (r []string)) *ContextMock {
	m.mock.MustQueryFunc = f
	return m.mock
}

//MustQuery implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) MustQuery(p string) (r []string) {
	atomic.AddUint64(&m.MustQueryPreCounter, 1)
	defer atomic.AddUint64(&m.MustQueryCounter, 1)

	if m.MustQueryMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.MustQueryMock.mockExpectations, ContextMockMustQueryParams{p},
			"Context.MustQuery got unexpected parameters")

		if m.MustQueryFunc == nil {

			m.t.Fatal("No results are set for the ContextMock.MustQuery")

			return
		}
	}

	if m.MustQueryFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.MustQuery")
		return
	}

	return m.MustQueryFunc(p)
}

//MustQueryMinimockCounter returns a count of ContextMock.MustQueryFunc invocations
func (m *ContextMock) MustQueryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MustQueryCounter)
}

//MustQueryMinimockPreCounter returns the value of ContextMock.MustQuery invocations
func (m *ContextMock) MustQueryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MustQueryPreCounter)
}

type mContextMockMustQuery1 struct {
	mock             *ContextMock
	mockExpectations *ContextMockMustQuery1Params
}

//ContextMockMustQuery1Params represents input parameters of the Context.MustQuery1
type ContextMockMustQuery1Params struct {
	p string
}

//Expect sets up expected params for the Context.MustQuery1
func (m *mContextMockMustQuery1) Expect(p string) *mContextMockMustQuery1 {
	m.mockExpectations = &ContextMockMustQuery1Params{p}
	return m
}

//Return sets up a mock for Context.MustQuery1 to return Return's arguments
func (m *mContextMockMustQuery1) Return(r string) *ContextMock {
	m.mock.MustQuery1Func = func(p string) string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.MustQuery1 method
func (m *mContextMockMustQuery1) Set(f func(p string) (r string)) *ContextMock {
	m.mock.MustQuery1Func = f
	return m.mock
}

//MustQuery1 implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) MustQuery1(p string) (r string) {
	atomic.AddUint64(&m.MustQuery1PreCounter, 1)
	defer atomic.AddUint64(&m.MustQuery1Counter, 1)

	if m.MustQuery1Mock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.MustQuery1Mock.mockExpectations, ContextMockMustQuery1Params{p},
			"Context.MustQuery1 got unexpected parameters")

		if m.MustQuery1Func == nil {

			m.t.Fatal("No results are set for the ContextMock.MustQuery1")

			return
		}
	}

	if m.MustQuery1Func == nil {
		m.t.Fatal("Unexpected call to ContextMock.MustQuery1")
		return
	}

	return m.MustQuery1Func(p)
}

//MustQuery1MinimockCounter returns a count of ContextMock.MustQuery1Func invocations
func (m *ContextMock) MustQuery1MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.MustQuery1Counter)
}

//MustQuery1MinimockPreCounter returns the value of ContextMock.MustQuery1 invocations
func (m *ContextMock) MustQuery1MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.MustQuery1PreCounter)
}

type mContextMockParam1 struct {
	mock             *ContextMock
	mockExpectations *ContextMockParam1Params
}

//ContextMockParam1Params represents input parameters of the Context.Param1
type ContextMockParam1Params struct {
	p string
}

//Expect sets up expected params for the Context.Param1
func (m *mContextMockParam1) Expect(p string) *mContextMockParam1 {
	m.mockExpectations = &ContextMockParam1Params{p}
	return m
}

//Return sets up a mock for Context.Param1 to return Return's arguments
func (m *mContextMockParam1) Return(r string, r1 error) *ContextMock {
	m.mock.Param1Func = func(p string) (string, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Param1 method
func (m *mContextMockParam1) Set(f func(p string) (r string, r1 error)) *ContextMock {
	m.mock.Param1Func = f
	return m.mock
}

//Param1 implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Param1(p string) (r string, r1 error) {
	atomic.AddUint64(&m.Param1PreCounter, 1)
	defer atomic.AddUint64(&m.Param1Counter, 1)

	if m.Param1Mock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.Param1Mock.mockExpectations, ContextMockParam1Params{p},
			"Context.Param1 got unexpected parameters")

		if m.Param1Func == nil {

			m.t.Fatal("No results are set for the ContextMock.Param1")

			return
		}
	}

	if m.Param1Func == nil {
		m.t.Fatal("Unexpected call to ContextMock.Param1")
		return
	}

	return m.Param1Func(p)
}

//Param1MinimockCounter returns a count of ContextMock.Param1Func invocations
func (m *ContextMock) Param1MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.Param1Counter)
}

//Param1MinimockPreCounter returns the value of ContextMock.Param1 invocations
func (m *ContextMock) Param1MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.Param1PreCounter)
}

type mContextMockParamAll struct {
	mock *ContextMock
}

//Return sets up a mock for Context.ParamAll to return Return's arguments
func (m *mContextMockParamAll) Return(r map[string]string) *ContextMock {
	m.mock.ParamAllFunc = func() map[string]string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.ParamAll method
func (m *mContextMockParamAll) Set(f func() (r map[string]string)) *ContextMock {
	m.mock.ParamAllFunc = f
	return m.mock
}

//ParamAll implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) ParamAll() (r map[string]string) {
	atomic.AddUint64(&m.ParamAllPreCounter, 1)
	defer atomic.AddUint64(&m.ParamAllCounter, 1)

	if m.ParamAllFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.ParamAll")
		return
	}

	return m.ParamAllFunc()
}

//ParamAllMinimockCounter returns a count of ContextMock.ParamAllFunc invocations
func (m *ContextMock) ParamAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ParamAllCounter)
}

//ParamAllMinimockPreCounter returns the value of ContextMock.ParamAll invocations
func (m *ContextMock) ParamAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ParamAllPreCounter)
}

type mContextMockProfile struct {
	mock             *ContextMock
	mockExpectations *ContextMockProfileParams
}

//ContextMockProfileParams represents input parameters of the Context.Profile
type ContextMockProfileParams struct {
	p  string
	p1 string
}

//Expect sets up expected params for the Context.Profile
func (m *mContextMockProfile) Expect(p string, p1 string) *mContextMockProfile {
	m.mockExpectations = &ContextMockProfileParams{p, p1}
	return m
}

//Return sets up a mock for Context.Profile to return Return's arguments
func (m *mContextMockProfile) Return(r profiler.ProfileFinishFunc) *ContextMock {
	m.mock.ProfileFunc = func(p string, p1 string) profiler.ProfileFinishFunc {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Profile method
func (m *mContextMockProfile) Set(f func(p string, p1 string) (r profiler.ProfileFinishFunc)) *ContextMock {
	m.mock.ProfileFunc = f
	return m.mock
}

//Profile implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Profile(p string, p1 string) (r profiler.ProfileFinishFunc) {
	atomic.AddUint64(&m.ProfilePreCounter, 1)
	defer atomic.AddUint64(&m.ProfileCounter, 1)

	if m.ProfileMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ProfileMock.mockExpectations, ContextMockProfileParams{p, p1},
			"Context.Profile got unexpected parameters")

		if m.ProfileFunc == nil {

			m.t.Fatal("No results are set for the ContextMock.Profile")

			return
		}
	}

	if m.ProfileFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.Profile")
		return
	}

	return m.ProfileFunc(p, p1)
}

//ProfileMinimockCounter returns a count of ContextMock.ProfileFunc invocations
func (m *ContextMock) ProfileMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ProfileCounter)
}

//ProfileMinimockPreCounter returns the value of ContextMock.Profile invocations
func (m *ContextMock) ProfileMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ProfilePreCounter)
}

type mContextMockProfiler struct {
	mock *ContextMock
}

//Return sets up a mock for Context.Profiler to return Return's arguments
func (m *mContextMockProfiler) Return(r profiler.Profiler) *ContextMock {
	m.mock.ProfilerFunc = func() profiler.Profiler {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Profiler method
func (m *mContextMockProfiler) Set(f func() (r profiler.Profiler)) *ContextMock {
	m.mock.ProfilerFunc = f
	return m.mock
}

//Profiler implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Profiler() (r profiler.Profiler) {
	atomic.AddUint64(&m.ProfilerPreCounter, 1)
	defer atomic.AddUint64(&m.ProfilerCounter, 1)

	if m.ProfilerFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.Profiler")
		return
	}

	return m.ProfilerFunc()
}

//ProfilerMinimockCounter returns a count of ContextMock.ProfilerFunc invocations
func (m *ContextMock) ProfilerMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ProfilerCounter)
}

//ProfilerMinimockPreCounter returns the value of ContextMock.Profiler invocations
func (m *ContextMock) ProfilerMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ProfilerPreCounter)
}

type mContextMockPush struct {
	mock             *ContextMock
	mockExpectations *ContextMockPushParams
}

//ContextMockPushParams represents input parameters of the Context.Push
type ContextMockPushParams struct {
	p  string
	p1 *http.PushOptions
}

//Expect sets up expected params for the Context.Push
func (m *mContextMockPush) Expect(p string, p1 *http.PushOptions) *mContextMockPush {
	m.mockExpectations = &ContextMockPushParams{p, p1}
	return m
}

//Return sets up a mock for Context.Push to return Return's arguments
func (m *mContextMockPush) Return(r error) *ContextMock {
	m.mock.PushFunc = func(p string, p1 *http.PushOptions) error {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Push method
func (m *mContextMockPush) Set(f func(p string, p1 *http.PushOptions) (r error)) *ContextMock {
	m.mock.PushFunc = f
	return m.mock
}

//Push implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Push(p string, p1 *http.PushOptions) (r error) {
	atomic.AddUint64(&m.PushPreCounter, 1)
	defer atomic.AddUint64(&m.PushCounter, 1)

	if m.PushMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.PushMock.mockExpectations, ContextMockPushParams{p, p1},
			"Context.Push got unexpected parameters")

		if m.PushFunc == nil {

			m.t.Fatal("No results are set for the ContextMock.Push")

			return
		}
	}

	if m.PushFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.Push")
		return
	}

	return m.PushFunc(p, p1)
}

//PushMinimockCounter returns a count of ContextMock.PushFunc invocations
func (m *ContextMock) PushMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.PushCounter)
}

//PushMinimockPreCounter returns the value of ContextMock.Push invocations
func (m *ContextMock) PushMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.PushPreCounter)
}

type mContextMockQuery struct {
	mock             *ContextMock
	mockExpectations *ContextMockQueryParams
}

//ContextMockQueryParams represents input parameters of the Context.Query
type ContextMockQueryParams struct {
	p string
}

//Expect sets up expected params for the Context.Query
func (m *mContextMockQuery) Expect(p string) *mContextMockQuery {
	m.mockExpectations = &ContextMockQueryParams{p}
	return m
}

//Return sets up a mock for Context.Query to return Return's arguments
func (m *mContextMockQuery) Return(r []string, r1 error) *ContextMock {
	m.mock.QueryFunc = func(p string) ([]string, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Query method
func (m *mContextMockQuery) Set(f func(p string) (r []string, r1 error)) *ContextMock {
	m.mock.QueryFunc = f
	return m.mock
}

//Query implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Query(p string) (r []string, r1 error) {
	atomic.AddUint64(&m.QueryPreCounter, 1)
	defer atomic.AddUint64(&m.QueryCounter, 1)

	if m.QueryMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.QueryMock.mockExpectations, ContextMockQueryParams{p},
			"Context.Query got unexpected parameters")

		if m.QueryFunc == nil {

			m.t.Fatal("No results are set for the ContextMock.Query")

			return
		}
	}

	if m.QueryFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.Query")
		return
	}

	return m.QueryFunc(p)
}

//QueryMinimockCounter returns a count of ContextMock.QueryFunc invocations
func (m *ContextMock) QueryMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.QueryCounter)
}

//QueryMinimockPreCounter returns the value of ContextMock.Query invocations
func (m *ContextMock) QueryMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.QueryPreCounter)
}

type mContextMockQuery1 struct {
	mock             *ContextMock
	mockExpectations *ContextMockQuery1Params
}

//ContextMockQuery1Params represents input parameters of the Context.Query1
type ContextMockQuery1Params struct {
	p string
}

//Expect sets up expected params for the Context.Query1
func (m *mContextMockQuery1) Expect(p string) *mContextMockQuery1 {
	m.mockExpectations = &ContextMockQuery1Params{p}
	return m
}

//Return sets up a mock for Context.Query1 to return Return's arguments
func (m *mContextMockQuery1) Return(r string, r1 error) *ContextMock {
	m.mock.Query1Func = func(p string) (string, error) {
		return r, r1
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Query1 method
func (m *mContextMockQuery1) Set(f func(p string) (r string, r1 error)) *ContextMock {
	m.mock.Query1Func = f
	return m.mock
}

//Query1 implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Query1(p string) (r string, r1 error) {
	atomic.AddUint64(&m.Query1PreCounter, 1)
	defer atomic.AddUint64(&m.Query1Counter, 1)

	if m.Query1Mock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.Query1Mock.mockExpectations, ContextMockQuery1Params{p},
			"Context.Query1 got unexpected parameters")

		if m.Query1Func == nil {

			m.t.Fatal("No results are set for the ContextMock.Query1")

			return
		}
	}

	if m.Query1Func == nil {
		m.t.Fatal("Unexpected call to ContextMock.Query1")
		return
	}

	return m.Query1Func(p)
}

//Query1MinimockCounter returns a count of ContextMock.Query1Func invocations
func (m *ContextMock) Query1MinimockCounter() uint64 {
	return atomic.LoadUint64(&m.Query1Counter)
}

//Query1MinimockPreCounter returns the value of ContextMock.Query1 invocations
func (m *ContextMock) Query1MinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.Query1PreCounter)
}

type mContextMockQueryAll struct {
	mock *ContextMock
}

//Return sets up a mock for Context.QueryAll to return Return's arguments
func (m *mContextMockQueryAll) Return(r map[string][]string) *ContextMock {
	m.mock.QueryAllFunc = func() map[string][]string {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.QueryAll method
func (m *mContextMockQueryAll) Set(f func() (r map[string][]string)) *ContextMock {
	m.mock.QueryAllFunc = f
	return m.mock
}

//QueryAll implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) QueryAll() (r map[string][]string) {
	atomic.AddUint64(&m.QueryAllPreCounter, 1)
	defer atomic.AddUint64(&m.QueryAllCounter, 1)

	if m.QueryAllFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.QueryAll")
		return
	}

	return m.QueryAllFunc()
}

//QueryAllMinimockCounter returns a count of ContextMock.QueryAllFunc invocations
func (m *ContextMock) QueryAllMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.QueryAllCounter)
}

//QueryAllMinimockPreCounter returns the value of ContextMock.QueryAll invocations
func (m *ContextMock) QueryAllMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.QueryAllPreCounter)
}

type mContextMockRequest struct {
	mock *ContextMock
}

//Return sets up a mock for Context.Request to return Return's arguments
func (m *mContextMockRequest) Return(r *http.Request) *ContextMock {
	m.mock.RequestFunc = func() *http.Request {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Request method
func (m *mContextMockRequest) Set(f func() (r *http.Request)) *ContextMock {
	m.mock.RequestFunc = f
	return m.mock
}

//Request implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Request() (r *http.Request) {
	atomic.AddUint64(&m.RequestPreCounter, 1)
	defer atomic.AddUint64(&m.RequestCounter, 1)

	if m.RequestFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.Request")
		return
	}

	return m.RequestFunc()
}

//RequestMinimockCounter returns a count of ContextMock.RequestFunc invocations
func (m *ContextMock) RequestMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.RequestCounter)
}

//RequestMinimockPreCounter returns the value of ContextMock.Request invocations
func (m *ContextMock) RequestMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.RequestPreCounter)
}

type mContextMockSession struct {
	mock *ContextMock
}

//Return sets up a mock for Context.Session to return Return's arguments
func (m *mContextMockSession) Return(r *sessions.Session) *ContextMock {
	m.mock.SessionFunc = func() *sessions.Session {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Session method
func (m *mContextMockSession) Set(f func() (r *sessions.Session)) *ContextMock {
	m.mock.SessionFunc = f
	return m.mock
}

//Session implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Session() (r *sessions.Session) {
	atomic.AddUint64(&m.SessionPreCounter, 1)
	defer atomic.AddUint64(&m.SessionCounter, 1)

	if m.SessionFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.Session")
		return
	}

	return m.SessionFunc()
}

//SessionMinimockCounter returns a count of ContextMock.SessionFunc invocations
func (m *ContextMock) SessionMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.SessionCounter)
}

//SessionMinimockPreCounter returns the value of ContextMock.Session invocations
func (m *ContextMock) SessionMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.SessionPreCounter)
}

type mContextMockValue struct {
	mock             *ContextMock
	mockExpectations *ContextMockValueParams
}

//ContextMockValueParams represents input parameters of the Context.Value
type ContextMockValueParams struct {
	p interface{}
}

//Expect sets up expected params for the Context.Value
func (m *mContextMockValue) Expect(p interface{}) *mContextMockValue {
	m.mockExpectations = &ContextMockValueParams{p}
	return m
}

//Return sets up a mock for Context.Value to return Return's arguments
func (m *mContextMockValue) Return(r interface{}) *ContextMock {
	m.mock.ValueFunc = func(p interface{}) interface{} {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.Value method
func (m *mContextMockValue) Set(f func(p interface{}) (r interface{})) *ContextMock {
	m.mock.ValueFunc = f
	return m.mock
}

//Value implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) Value(p interface{}) (r interface{}) {
	atomic.AddUint64(&m.ValuePreCounter, 1)
	defer atomic.AddUint64(&m.ValueCounter, 1)

	if m.ValueMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.ValueMock.mockExpectations, ContextMockValueParams{p},
			"Context.Value got unexpected parameters")

		if m.ValueFunc == nil {

			m.t.Fatal("No results are set for the ContextMock.Value")

			return
		}
	}

	if m.ValueFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.Value")
		return
	}

	return m.ValueFunc(p)
}

//ValueMinimockCounter returns a count of ContextMock.ValueFunc invocations
func (m *ContextMock) ValueMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.ValueCounter)
}

//ValueMinimockPreCounter returns the value of ContextMock.Value invocations
func (m *ContextMock) ValueMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.ValuePreCounter)
}

type mContextMockWithValue struct {
	mock             *ContextMock
	mockExpectations *ContextMockWithValueParams
}

//ContextMockWithValueParams represents input parameters of the Context.WithValue
type ContextMockWithValueParams struct {
	p  interface{}
	p1 interface{}
}

//Expect sets up expected params for the Context.WithValue
func (m *mContextMockWithValue) Expect(p interface{}, p1 interface{}) *mContextMockWithValue {
	m.mockExpectations = &ContextMockWithValueParams{p, p1}
	return m
}

//Return sets up a mock for Context.WithValue to return Return's arguments
func (m *mContextMockWithValue) Return(r web.Context) *ContextMock {
	m.mock.WithValueFunc = func(p interface{}, p1 interface{}) web.Context {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.WithValue method
func (m *mContextMockWithValue) Set(f func(p interface{}, p1 interface{}) (r web.Context)) *ContextMock {
	m.mock.WithValueFunc = f
	return m.mock
}

//WithValue implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) WithValue(p interface{}, p1 interface{}) (r web.Context) {
	atomic.AddUint64(&m.WithValuePreCounter, 1)
	defer atomic.AddUint64(&m.WithValueCounter, 1)

	if m.WithValueMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.WithValueMock.mockExpectations, ContextMockWithValueParams{p, p1},
			"Context.WithValue got unexpected parameters")

		if m.WithValueFunc == nil {

			m.t.Fatal("No results are set for the ContextMock.WithValue")

			return
		}
	}

	if m.WithValueFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.WithValue")
		return
	}

	return m.WithValueFunc(p, p1)
}

//WithValueMinimockCounter returns a count of ContextMock.WithValueFunc invocations
func (m *ContextMock) WithValueMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WithValueCounter)
}

//WithValueMinimockPreCounter returns the value of ContextMock.WithValue invocations
func (m *ContextMock) WithValueMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WithValuePreCounter)
}

type mContextMockWithVars struct {
	mock             *ContextMock
	mockExpectations *ContextMockWithVarsParams
}

//ContextMockWithVarsParams represents input parameters of the Context.WithVars
type ContextMockWithVarsParams struct {
	p map[string]string
}

//Expect sets up expected params for the Context.WithVars
func (m *mContextMockWithVars) Expect(p map[string]string) *mContextMockWithVars {
	m.mockExpectations = &ContextMockWithVarsParams{p}
	return m
}

//Return sets up a mock for Context.WithVars to return Return's arguments
func (m *mContextMockWithVars) Return(r web.Context) *ContextMock {
	m.mock.WithVarsFunc = func(p map[string]string) web.Context {
		return r
	}
	return m.mock
}

//Set uses given function f as a mock of Context.WithVars method
func (m *mContextMockWithVars) Set(f func(p map[string]string) (r web.Context)) *ContextMock {
	m.mock.WithVarsFunc = f
	return m.mock
}

//WithVars implements go.aoe.com/flamingo/framework/web.Context interface
func (m *ContextMock) WithVars(p map[string]string) (r web.Context) {
	atomic.AddUint64(&m.WithVarsPreCounter, 1)
	defer atomic.AddUint64(&m.WithVarsCounter, 1)

	if m.WithVarsMock.mockExpectations != nil {
		testify_assert.Equal(m.t, *m.WithVarsMock.mockExpectations, ContextMockWithVarsParams{p},
			"Context.WithVars got unexpected parameters")

		if m.WithVarsFunc == nil {

			m.t.Fatal("No results are set for the ContextMock.WithVars")

			return
		}
	}

	if m.WithVarsFunc == nil {
		m.t.Fatal("Unexpected call to ContextMock.WithVars")
		return
	}

	return m.WithVarsFunc(p)
}

//WithVarsMinimockCounter returns a count of ContextMock.WithVarsFunc invocations
func (m *ContextMock) WithVarsMinimockCounter() uint64 {
	return atomic.LoadUint64(&m.WithVarsCounter)
}

//WithVarsMinimockPreCounter returns the value of ContextMock.WithVars invocations
func (m *ContextMock) WithVarsMinimockPreCounter() uint64 {
	return atomic.LoadUint64(&m.WithVarsPreCounter)
}

//ValidateCallCounters checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ContextMock) ValidateCallCounters() {

	if m.DeadlineFunc != nil && atomic.LoadUint64(&m.DeadlineCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Deadline")
	}

	if m.DoneFunc != nil && atomic.LoadUint64(&m.DoneCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Done")
	}

	if m.ErrFunc != nil && atomic.LoadUint64(&m.ErrCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Err")
	}

	if m.EventRouterFunc != nil && atomic.LoadUint64(&m.EventRouterCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.EventRouter")
	}

	if m.FormFunc != nil && atomic.LoadUint64(&m.FormCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Form")
	}

	if m.Form1Func != nil && atomic.LoadUint64(&m.Form1Counter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Form1")
	}

	if m.FormAllFunc != nil && atomic.LoadUint64(&m.FormAllCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.FormAll")
	}

	if m.IDFunc != nil && atomic.LoadUint64(&m.IDCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.ID")
	}

	if m.LoadParamsFunc != nil && atomic.LoadUint64(&m.LoadParamsCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.LoadParams")
	}

	if m.MustFormFunc != nil && atomic.LoadUint64(&m.MustFormCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.MustForm")
	}

	if m.MustForm1Func != nil && atomic.LoadUint64(&m.MustForm1Counter) == 0 {
		m.t.Fatal("Expected call to ContextMock.MustForm1")
	}

	if m.MustParam1Func != nil && atomic.LoadUint64(&m.MustParam1Counter) == 0 {
		m.t.Fatal("Expected call to ContextMock.MustParam1")
	}

	if m.MustQueryFunc != nil && atomic.LoadUint64(&m.MustQueryCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.MustQuery")
	}

	if m.MustQuery1Func != nil && atomic.LoadUint64(&m.MustQuery1Counter) == 0 {
		m.t.Fatal("Expected call to ContextMock.MustQuery1")
	}

	if m.Param1Func != nil && atomic.LoadUint64(&m.Param1Counter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Param1")
	}

	if m.ParamAllFunc != nil && atomic.LoadUint64(&m.ParamAllCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.ParamAll")
	}

	if m.ProfileFunc != nil && atomic.LoadUint64(&m.ProfileCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Profile")
	}

	if m.ProfilerFunc != nil && atomic.LoadUint64(&m.ProfilerCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Profiler")
	}

	if m.PushFunc != nil && atomic.LoadUint64(&m.PushCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Push")
	}

	if m.QueryFunc != nil && atomic.LoadUint64(&m.QueryCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Query")
	}

	if m.Query1Func != nil && atomic.LoadUint64(&m.Query1Counter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Query1")
	}

	if m.QueryAllFunc != nil && atomic.LoadUint64(&m.QueryAllCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.QueryAll")
	}

	if m.RequestFunc != nil && atomic.LoadUint64(&m.RequestCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Request")
	}

	if m.SessionFunc != nil && atomic.LoadUint64(&m.SessionCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Session")
	}

	if m.ValueFunc != nil && atomic.LoadUint64(&m.ValueCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Value")
	}

	if m.WithValueFunc != nil && atomic.LoadUint64(&m.WithValueCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.WithValue")
	}

	if m.WithVarsFunc != nil && atomic.LoadUint64(&m.WithVarsCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.WithVars")
	}

}

//CheckMocksCalled checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish method or use Finish method of minimock.Controller
func (m *ContextMock) CheckMocksCalled() {
	m.Finish()
}

//Finish checks that all mocked methods of the interface have been called at least once
//Deprecated: please use MinimockFinish or use Finish method of minimock.Controller
func (m *ContextMock) Finish() {
	m.MinimockFinish()
}

//MinimockFinish checks that all mocked methods of the interface have been called at least once
func (m *ContextMock) MinimockFinish() {

	if m.DeadlineFunc != nil && atomic.LoadUint64(&m.DeadlineCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Deadline")
	}

	if m.DoneFunc != nil && atomic.LoadUint64(&m.DoneCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Done")
	}

	if m.ErrFunc != nil && atomic.LoadUint64(&m.ErrCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Err")
	}

	if m.EventRouterFunc != nil && atomic.LoadUint64(&m.EventRouterCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.EventRouter")
	}

	if m.FormFunc != nil && atomic.LoadUint64(&m.FormCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Form")
	}

	if m.Form1Func != nil && atomic.LoadUint64(&m.Form1Counter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Form1")
	}

	if m.FormAllFunc != nil && atomic.LoadUint64(&m.FormAllCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.FormAll")
	}

	if m.IDFunc != nil && atomic.LoadUint64(&m.IDCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.ID")
	}

	if m.LoadParamsFunc != nil && atomic.LoadUint64(&m.LoadParamsCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.LoadParams")
	}

	if m.MustFormFunc != nil && atomic.LoadUint64(&m.MustFormCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.MustForm")
	}

	if m.MustForm1Func != nil && atomic.LoadUint64(&m.MustForm1Counter) == 0 {
		m.t.Fatal("Expected call to ContextMock.MustForm1")
	}

	if m.MustParam1Func != nil && atomic.LoadUint64(&m.MustParam1Counter) == 0 {
		m.t.Fatal("Expected call to ContextMock.MustParam1")
	}

	if m.MustQueryFunc != nil && atomic.LoadUint64(&m.MustQueryCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.MustQuery")
	}

	if m.MustQuery1Func != nil && atomic.LoadUint64(&m.MustQuery1Counter) == 0 {
		m.t.Fatal("Expected call to ContextMock.MustQuery1")
	}

	if m.Param1Func != nil && atomic.LoadUint64(&m.Param1Counter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Param1")
	}

	if m.ParamAllFunc != nil && atomic.LoadUint64(&m.ParamAllCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.ParamAll")
	}

	if m.ProfileFunc != nil && atomic.LoadUint64(&m.ProfileCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Profile")
	}

	if m.ProfilerFunc != nil && atomic.LoadUint64(&m.ProfilerCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Profiler")
	}

	if m.PushFunc != nil && atomic.LoadUint64(&m.PushCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Push")
	}

	if m.QueryFunc != nil && atomic.LoadUint64(&m.QueryCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Query")
	}

	if m.Query1Func != nil && atomic.LoadUint64(&m.Query1Counter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Query1")
	}

	if m.QueryAllFunc != nil && atomic.LoadUint64(&m.QueryAllCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.QueryAll")
	}

	if m.RequestFunc != nil && atomic.LoadUint64(&m.RequestCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Request")
	}

	if m.SessionFunc != nil && atomic.LoadUint64(&m.SessionCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Session")
	}

	if m.ValueFunc != nil && atomic.LoadUint64(&m.ValueCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.Value")
	}

	if m.WithValueFunc != nil && atomic.LoadUint64(&m.WithValueCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.WithValue")
	}

	if m.WithVarsFunc != nil && atomic.LoadUint64(&m.WithVarsCounter) == 0 {
		m.t.Fatal("Expected call to ContextMock.WithVars")
	}

}

//Wait waits for all mocked methods to be called at least once
//Deprecated: please use MinimockWait or use Wait method of minimock.Controller
func (m *ContextMock) Wait(timeout time.Duration) {
	m.MinimockWait(timeout)
}

//MinimockWait waits for all mocked methods to be called at least once
//this method is called by minimock.Controller
func (m *ContextMock) MinimockWait(timeout time.Duration) {
	timeoutCh := time.After(timeout)
	for {
		ok := true
		ok = ok && (m.DeadlineFunc == nil || atomic.LoadUint64(&m.DeadlineCounter) > 0)
		ok = ok && (m.DoneFunc == nil || atomic.LoadUint64(&m.DoneCounter) > 0)
		ok = ok && (m.ErrFunc == nil || atomic.LoadUint64(&m.ErrCounter) > 0)
		ok = ok && (m.EventRouterFunc == nil || atomic.LoadUint64(&m.EventRouterCounter) > 0)
		ok = ok && (m.FormFunc == nil || atomic.LoadUint64(&m.FormCounter) > 0)
		ok = ok && (m.Form1Func == nil || atomic.LoadUint64(&m.Form1Counter) > 0)
		ok = ok && (m.FormAllFunc == nil || atomic.LoadUint64(&m.FormAllCounter) > 0)
		ok = ok && (m.IDFunc == nil || atomic.LoadUint64(&m.IDCounter) > 0)
		ok = ok && (m.LoadParamsFunc == nil || atomic.LoadUint64(&m.LoadParamsCounter) > 0)
		ok = ok && (m.MustFormFunc == nil || atomic.LoadUint64(&m.MustFormCounter) > 0)
		ok = ok && (m.MustForm1Func == nil || atomic.LoadUint64(&m.MustForm1Counter) > 0)
		ok = ok && (m.MustParam1Func == nil || atomic.LoadUint64(&m.MustParam1Counter) > 0)
		ok = ok && (m.MustQueryFunc == nil || atomic.LoadUint64(&m.MustQueryCounter) > 0)
		ok = ok && (m.MustQuery1Func == nil || atomic.LoadUint64(&m.MustQuery1Counter) > 0)
		ok = ok && (m.Param1Func == nil || atomic.LoadUint64(&m.Param1Counter) > 0)
		ok = ok && (m.ParamAllFunc == nil || atomic.LoadUint64(&m.ParamAllCounter) > 0)
		ok = ok && (m.ProfileFunc == nil || atomic.LoadUint64(&m.ProfileCounter) > 0)
		ok = ok && (m.ProfilerFunc == nil || atomic.LoadUint64(&m.ProfilerCounter) > 0)
		ok = ok && (m.PushFunc == nil || atomic.LoadUint64(&m.PushCounter) > 0)
		ok = ok && (m.QueryFunc == nil || atomic.LoadUint64(&m.QueryCounter) > 0)
		ok = ok && (m.Query1Func == nil || atomic.LoadUint64(&m.Query1Counter) > 0)
		ok = ok && (m.QueryAllFunc == nil || atomic.LoadUint64(&m.QueryAllCounter) > 0)
		ok = ok && (m.RequestFunc == nil || atomic.LoadUint64(&m.RequestCounter) > 0)
		ok = ok && (m.SessionFunc == nil || atomic.LoadUint64(&m.SessionCounter) > 0)
		ok = ok && (m.ValueFunc == nil || atomic.LoadUint64(&m.ValueCounter) > 0)
		ok = ok && (m.WithValueFunc == nil || atomic.LoadUint64(&m.WithValueCounter) > 0)
		ok = ok && (m.WithVarsFunc == nil || atomic.LoadUint64(&m.WithVarsCounter) > 0)

		if ok {
			return
		}

		select {
		case <-timeoutCh:

			if m.DeadlineFunc != nil && atomic.LoadUint64(&m.DeadlineCounter) == 0 {
				m.t.Error("Expected call to ContextMock.Deadline")
			}

			if m.DoneFunc != nil && atomic.LoadUint64(&m.DoneCounter) == 0 {
				m.t.Error("Expected call to ContextMock.Done")
			}

			if m.ErrFunc != nil && atomic.LoadUint64(&m.ErrCounter) == 0 {
				m.t.Error("Expected call to ContextMock.Err")
			}

			if m.EventRouterFunc != nil && atomic.LoadUint64(&m.EventRouterCounter) == 0 {
				m.t.Error("Expected call to ContextMock.EventRouter")
			}

			if m.FormFunc != nil && atomic.LoadUint64(&m.FormCounter) == 0 {
				m.t.Error("Expected call to ContextMock.Form")
			}

			if m.Form1Func != nil && atomic.LoadUint64(&m.Form1Counter) == 0 {
				m.t.Error("Expected call to ContextMock.Form1")
			}

			if m.FormAllFunc != nil && atomic.LoadUint64(&m.FormAllCounter) == 0 {
				m.t.Error("Expected call to ContextMock.FormAll")
			}

			if m.IDFunc != nil && atomic.LoadUint64(&m.IDCounter) == 0 {
				m.t.Error("Expected call to ContextMock.ID")
			}

			if m.LoadParamsFunc != nil && atomic.LoadUint64(&m.LoadParamsCounter) == 0 {
				m.t.Error("Expected call to ContextMock.LoadParams")
			}

			if m.MustFormFunc != nil && atomic.LoadUint64(&m.MustFormCounter) == 0 {
				m.t.Error("Expected call to ContextMock.MustForm")
			}

			if m.MustForm1Func != nil && atomic.LoadUint64(&m.MustForm1Counter) == 0 {
				m.t.Error("Expected call to ContextMock.MustForm1")
			}

			if m.MustParam1Func != nil && atomic.LoadUint64(&m.MustParam1Counter) == 0 {
				m.t.Error("Expected call to ContextMock.MustParam1")
			}

			if m.MustQueryFunc != nil && atomic.LoadUint64(&m.MustQueryCounter) == 0 {
				m.t.Error("Expected call to ContextMock.MustQuery")
			}

			if m.MustQuery1Func != nil && atomic.LoadUint64(&m.MustQuery1Counter) == 0 {
				m.t.Error("Expected call to ContextMock.MustQuery1")
			}

			if m.Param1Func != nil && atomic.LoadUint64(&m.Param1Counter) == 0 {
				m.t.Error("Expected call to ContextMock.Param1")
			}

			if m.ParamAllFunc != nil && atomic.LoadUint64(&m.ParamAllCounter) == 0 {
				m.t.Error("Expected call to ContextMock.ParamAll")
			}

			if m.ProfileFunc != nil && atomic.LoadUint64(&m.ProfileCounter) == 0 {
				m.t.Error("Expected call to ContextMock.Profile")
			}

			if m.ProfilerFunc != nil && atomic.LoadUint64(&m.ProfilerCounter) == 0 {
				m.t.Error("Expected call to ContextMock.Profiler")
			}

			if m.PushFunc != nil && atomic.LoadUint64(&m.PushCounter) == 0 {
				m.t.Error("Expected call to ContextMock.Push")
			}

			if m.QueryFunc != nil && atomic.LoadUint64(&m.QueryCounter) == 0 {
				m.t.Error("Expected call to ContextMock.Query")
			}

			if m.Query1Func != nil && atomic.LoadUint64(&m.Query1Counter) == 0 {
				m.t.Error("Expected call to ContextMock.Query1")
			}

			if m.QueryAllFunc != nil && atomic.LoadUint64(&m.QueryAllCounter) == 0 {
				m.t.Error("Expected call to ContextMock.QueryAll")
			}

			if m.RequestFunc != nil && atomic.LoadUint64(&m.RequestCounter) == 0 {
				m.t.Error("Expected call to ContextMock.Request")
			}

			if m.SessionFunc != nil && atomic.LoadUint64(&m.SessionCounter) == 0 {
				m.t.Error("Expected call to ContextMock.Session")
			}

			if m.ValueFunc != nil && atomic.LoadUint64(&m.ValueCounter) == 0 {
				m.t.Error("Expected call to ContextMock.Value")
			}

			if m.WithValueFunc != nil && atomic.LoadUint64(&m.WithValueCounter) == 0 {
				m.t.Error("Expected call to ContextMock.WithValue")
			}

			if m.WithVarsFunc != nil && atomic.LoadUint64(&m.WithVarsCounter) == 0 {
				m.t.Error("Expected call to ContextMock.WithVars")
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
func (m *ContextMock) AllMocksCalled() bool {

	if m.DeadlineFunc != nil && atomic.LoadUint64(&m.DeadlineCounter) == 0 {
		return false
	}

	if m.DoneFunc != nil && atomic.LoadUint64(&m.DoneCounter) == 0 {
		return false
	}

	if m.ErrFunc != nil && atomic.LoadUint64(&m.ErrCounter) == 0 {
		return false
	}

	if m.EventRouterFunc != nil && atomic.LoadUint64(&m.EventRouterCounter) == 0 {
		return false
	}

	if m.FormFunc != nil && atomic.LoadUint64(&m.FormCounter) == 0 {
		return false
	}

	if m.Form1Func != nil && atomic.LoadUint64(&m.Form1Counter) == 0 {
		return false
	}

	if m.FormAllFunc != nil && atomic.LoadUint64(&m.FormAllCounter) == 0 {
		return false
	}

	if m.IDFunc != nil && atomic.LoadUint64(&m.IDCounter) == 0 {
		return false
	}

	if m.LoadParamsFunc != nil && atomic.LoadUint64(&m.LoadParamsCounter) == 0 {
		return false
	}

	if m.MustFormFunc != nil && atomic.LoadUint64(&m.MustFormCounter) == 0 {
		return false
	}

	if m.MustForm1Func != nil && atomic.LoadUint64(&m.MustForm1Counter) == 0 {
		return false
	}

	if m.MustParam1Func != nil && atomic.LoadUint64(&m.MustParam1Counter) == 0 {
		return false
	}

	if m.MustQueryFunc != nil && atomic.LoadUint64(&m.MustQueryCounter) == 0 {
		return false
	}

	if m.MustQuery1Func != nil && atomic.LoadUint64(&m.MustQuery1Counter) == 0 {
		return false
	}

	if m.Param1Func != nil && atomic.LoadUint64(&m.Param1Counter) == 0 {
		return false
	}

	if m.ParamAllFunc != nil && atomic.LoadUint64(&m.ParamAllCounter) == 0 {
		return false
	}

	if m.ProfileFunc != nil && atomic.LoadUint64(&m.ProfileCounter) == 0 {
		return false
	}

	if m.ProfilerFunc != nil && atomic.LoadUint64(&m.ProfilerCounter) == 0 {
		return false
	}

	if m.PushFunc != nil && atomic.LoadUint64(&m.PushCounter) == 0 {
		return false
	}

	if m.QueryFunc != nil && atomic.LoadUint64(&m.QueryCounter) == 0 {
		return false
	}

	if m.Query1Func != nil && atomic.LoadUint64(&m.Query1Counter) == 0 {
		return false
	}

	if m.QueryAllFunc != nil && atomic.LoadUint64(&m.QueryAllCounter) == 0 {
		return false
	}

	if m.RequestFunc != nil && atomic.LoadUint64(&m.RequestCounter) == 0 {
		return false
	}

	if m.SessionFunc != nil && atomic.LoadUint64(&m.SessionCounter) == 0 {
		return false
	}

	if m.ValueFunc != nil && atomic.LoadUint64(&m.ValueCounter) == 0 {
		return false
	}

	if m.WithValueFunc != nil && atomic.LoadUint64(&m.WithValueCounter) == 0 {
		return false
	}

	if m.WithVarsFunc != nil && atomic.LoadUint64(&m.WithVarsCounter) == 0 {
		return false
	}

	return true
}
