package errors_test

import (
	stderrors "errors"
	"testing"

	"github.com/phayes/errors"
)

var (
	Strger  = stringer{}
	ErrFoo  = errors.New("Fooey")
	ErrBar  = errors.New("Barf")
	ErrStd  = stderrors.New("This is a stanard error from the standard library")
	ErrStd2 = stderrors.New("Another standard error from the standard library")
	ErrFmt  = errors.Newf("%s", Strger)
)

type stringer struct{}

func (s stringer) String() string {
	return "stringer out"
}

func TestFooWrappingStdError(t *testing.T) {
	err := FooWrappingStdError()
	if !errors.Equal(err, ErrFoo) || !errors.Equal(ErrFoo, err) {
		t.Error("Foo not determined to be equal to an errors.Error based on itself")
		return
	}
	if !errors.IsA(err, ErrStd) {
		t.Error("Error that wraps standard library not determined to contain the standard")
		return
	}
	if err.Error() != "Fooey. This is a stanard error from the standard library" {
		t.Error("String genertation not correct for FooWrappingStdError")
		return
	}
}

func TestFooWrappingBar(t *testing.T) {
	err := FooWrappingBar()
	if !errors.Equal(err, ErrFoo) || !errors.Equal(ErrFoo, err) {
		t.Error("Foo not determined to be equal to an errors.Error based on itself")
		return
	}
	if !errors.IsA(err, ErrFoo) {
		t.Error("Error that wraps Bar not determined to contain it")
		return
	}
	if err.Error() != "Fooey. Barf" {
		t.Error("String genertation not correct for FooWrappingBar")
		return
	}
}

func TestStringWrappingFoo(t *testing.T) {
	err := StringWrappingFoo()
	if err.Error() != "String. Fooey" {
		t.Error("String error for StringWrappingFoo")
		return
	}
	if err.(errors.Error).Inner() != ErrFoo {
		t.Error("Foo not inner to the string error")
		return
	}
}

func TestFooWrappingFmt(t *testing.T) {
	err := FmtWrappingFoo()
	if err.Error() != "stringer out. Fooey" {
		t.Error("String error for FooWrappingFmt")
		return
	}
}

func TestAppend(t *testing.T) {
	err := errors.Append(ErrFoo, errors.New("inner string"))
	if err.Error() != "Fooey. inner string" {
		t.Error("String error for Appends")
		return
	}
	err = errors.Appends(ErrFoo, "inner string")
	if err.Error() != "Fooey. inner string" {
		t.Error("String error for Appends")
		return
	}
	err = errors.Appendf(ErrFoo, "%s", Strger)
	if err.Error() != "Fooey. stringer out" {
		t.Error("String error for Appendf")
		return
	}

	// Test appending to a stderr
	err = errors.Append(ErrStd, errors.New("inner string"))
	if err.Error() != "This is a stanard error from the standard library. inner string" {
		t.Error("String error for Appends")
		return
	}
	err = errors.Append(ErrStd, ErrStd2)
	if err.Error() != "This is a stanard error from the standard library. Another standard error from the standard library" {
		t.Error("String error for Appends")
		return
	}
	err = errors.Appends(ErrStd, "inner string")
	if err.Error() != "This is a stanard error from the standard library. inner string" {
		t.Error("String error for Appends")
		return
	}
	err = errors.Appendf(ErrStd, "%s", Strger)
	if err.Error() != "This is a stanard error from the standard library. stringer out" {
		t.Error("String error for Appendf")
		return
	}

}

func TestEquality(t *testing.T) {
	if ErrFmt.Error() != "stringer out" {
		t.Error("wrong output for FmtErr")
		return
	}
	if !errors.Equal(ErrStd, ErrStd) {
		t.Error("ErrStd equality error")
		return
	}
	if !errors.Equal(ErrFoo, ErrFoo) {
		t.Error("ErrFoo equality error")
		return
	}
	if errors.Equal(ErrStd, ErrStd2) {
		t.Error("ErrStd and ErrStd2 found to be equal")
		return
	}
	if errors.Equal(ErrStd, ErrFmt) {
		t.Error("ErrStd and ErrFmt found to be equal")
		return
	}
	if errors.IsA(ErrStd, ErrFmt) {
		t.Error("ErrStd and ErrFmt returned true for IsA")
		return
	}
}
func TestCause(t *testing.T) {
	if errors.Cause(ErrStd) != ErrStd {
		t.Error("Cause of standard error should be itself")
	}
	if errors.Cause(FooWrappingStdError()) != ErrStd {
		t.Error("Cause of FooWrappingStdError should be ErrStd")
	}
	if !errors.Equal(errors.Cause(StdErrorWrappingFoo()), ErrFoo) {
		t.Error("Cause of StdErrorWrappingFoo should be Foo")
	}

}

func StdErrorWrappingFoo() error {
	return errors.Wrap(ErrFoo, ErrStd)
}

func FooWrappingStdError() error {
	return errors.Wrap(ErrStd, ErrFoo)
}

func FooWrappingBar() error {
	return errors.Wrap(ErrBar, ErrFoo)
}

func StringWrappingFoo() error {
	return errors.Wraps(ErrFoo, "String")
}

func FmtWrappingFoo() error {
	return errors.Wrapf(ErrFoo, "%s", Strger)
}
