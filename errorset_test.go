package errors_test

import (
	"testing"

	"github.com/phayes/errors"
)

var (
	ErrFoo2 = errors.New("Foo")
	ErrBar2 = errors.New("Bar")
)

func TestErrorSet(t *testing.T) {
	errset := errors.NewErrorSet()

	if errset.HasErrors() {
		t.Error("Empty set has errors")
		return
	}

	errset.Add("foo", ErrFoo2)
	errset.Add("bar", ErrBar2)

	if !errset.HasErrors() {
		t.Error("Non-empty set does not have errors")
		return
	}

	foo := errset.Get("foo")
	if foo != ErrFoo2 {
		t.Error("Unable to set and get foo")
		return
	}

	bar := errset.Get("bar")
	if bar != ErrBar2 {
		t.Error("Unable to set and get bar")
		return
	}

	all := errset.GetAll()

	if len(all) != 2 {
		t.Error("Bad len for GetAll()")
		return
	}

	if all["foo"] != foo {
		t.Error("Unable get foo from all")
		return
	}
	if all["bar"] != bar {
		t.Error("Unable get bar from all")
		return
	}

	if errset.Error() != "foo: Foo. bar: Bar" && errset.Error() != "bar: Bar. foo: Foo" {
		t.Error("Wrong errset.Error() output")
		return
	}
}
