package dingo

import (
	"errors"
	"fmt"
	"testing"
)

func TryModule(module Module) (resultingError error) {
	defer func() {
		if err := recover(); err != nil {
			if err, ok := err.(error); ok {
				resultingError = err
				return
			}
			resultingError = errors.New(fmt.Sprint(err))
		}
	}()

	injector := NewInjector()
	injector.InitModules(module)
	return nil
}

type (
	tryModuleOk   struct{}
	tryModuleFail struct{}
)

func (t *tryModuleOk) Configure(injector *Injector) {
	injector.Bind(new(string)).ToInstance("test")
}

func (t *tryModuleFail) Configure(injector *Injector) {
	injector.Bind(new(int)).ToInstance("test")
}

func TestTryModule(t *testing.T) {
	err := TryModule(new(tryModuleOk))
	if err != nil {
		t.Errorf("tryModuleOk{} failed during module load, error: %q", err)
	}

	err = TryModule(new(tryModuleFail))
	if err == nil {
		t.Errorf("tryModuleFail{} did not fail during module load, error: %q", err)
	}
}
