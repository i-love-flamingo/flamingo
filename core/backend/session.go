package backend

/*
// SessionBackendCreater psc
type SessionBackendCreater func(string) SessionBackender

var sessionBackends map[string]SessionBackendCreater

func init() {
	sessionBackends = make(map[string]SessionBackendCreater)
}

// RegisterSessionBackend rsb
func RegisterSessionBackend(name string, sbc SessionBackendCreater) {
	sessionBackends[name] = sbc
}

// CreateSessionBackend csb
func CreateSessionBackend(dsn string) SessionBackender {
	cfg, err := url.Parse(dsn)
	if err != nil {
		panic(err)
	}
	if sessionBackends[cfg.Scheme] == nil {
		panic("no such session backend " + cfg.Scheme)
	}
	return sessionBackends[cfg.Scheme](dsn)
}

// SessionBackender sb
type SessionBackender interface {
	Init(echo.Context) Sessioner
}

// Sessioner s
type Sessioner interface {
	Set(string, interface{})
	Get(string) interface{}
	Persist() bool
	ID() string
}
*/
