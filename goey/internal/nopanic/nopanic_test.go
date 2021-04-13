package nopanic_test

import (
	"guitest/goey/internal/nopanic"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func ExampleWrap() {
	myerr := errors.New("No luck!")

	// Wrap the callback to prevent the escape of any panics.
	err := nopanic.Wrap(func() error {
		return myerr
	})
	// Check for errors.
	if err != nil {
		// Print the
		fmt.Println(err)
	}

	// Output:
	// No luck!
}

func ExampleWrap_2() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered...")
			fmt.Println(r)
		}
	}()

	// Wrap the callback to prevent the escape of any panics.
	err := nopanic.Wrap(func() error {
		// Bad things sometimes happen to good people...
		panic("No luck!")
	})
	// Check for errors.
	if err != nil {
		// If this is a PanicError, we rethrow.
		if value, ok := err.(nopanic.PanicError); ok {
			panic(value.Value())
		}
		// Otherwise, continue with normal error handling.
		fmt.Println("Normal error...")
		fmt.Println(err)
	}

	// Output:
	// Recovered...
	// No luck!
}

func ExampleUnwrap() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered...")
			fmt.Println(r)
		}
	}()

	// Wrap the callback to prevent the escape of any panics.
	err := nopanic.Unwrap(nopanic.Wrap(func() error {
		// Bad things sometimes happen to good people...
		panic("No luck!")
	}))
	// Check for errors.
	if err != nil {
		// Otherwise, continue with normal error handling.
		fmt.Println("Normal error...")
		fmt.Println(err)
	}

	// Output (not reproducible because of stack trace):
	// Recovered...
	// No luck!
	//
	// goroutine 1 [running]:
	// runtime/debug.Stack(0x0, 0x0, 0xc04202b9c8)
	//         C:/Go/src/runtime/debug/stack.go:24 +0x80
	// guitest/goey/internal/nopanic.New(0x5190c0, 0xc04200cf30, 0xc04200cf30, 0xc04202ba78, 0x40ef6c, 0xc04200cf30)
	//         guitest/goey/internal/nopanic/_test/_obj_test/nopanic.go:35 +0x33
	// guitest/goey/internal/nopanic.Wrap.func1(0xc04202bb58)
	//         guitest/goey/internal/nopanic/_test/_obj_test/nopanic.go:45 +0x71
	// panic(0x5190c0, 0xc04200cf30)
	//         C:/Go/src/runtime/panic.go:489 +0x2dd
	//     ... several frame omitted ...
	// main.main()
	//         guitest/goey/internal/nopanic/_test/_testmain.go:98 +0x201
}

func TestWrap(t *testing.T) {
	err1 := errors.New("No luck!")
	cases := []struct {
		in  func() error
		out interface{}
	}{
		{func() error { return nil }, nil},
		{func() error { return err1 }, err1},
		{func() error { panic("No luck!") }, "No luck!"},
		{func() error { panic(err1) }, err1},
	}

	for i, v := range cases {
		out := nopanic.Wrap(v.in)
		if err, ok := out.(nopanic.PanicError); ok {
			if !reflect.DeepEqual(err.Value(), v.out) {
				t.Errorf("Case %d:  got %v, want %v", i, err.Value(), v.out)
			}
			if err.Stack() == "" {
				t.Errorf("Case %d:  missing stack trace", i)
			}
			if err.Error() == "" {
				t.Errorf("Case %d:  missing stack trace", i)
			}
			if !strings.HasSuffix(err.Error(), err.Stack()) {
				t.Errorf("Case %d: stack trace is not a suffix of the error message", i)
			}
		} else if !reflect.DeepEqual(out, v.out) {
			t.Errorf("Case %d:  got %v, want %v", i, out, v.out)
		}
	}
}

func TestUnwrap(t *testing.T) {
	err1 := errors.New("No luck!")

	cases := []struct {
		in  error
		out error
	}{
		{nil, nil},
		{err1, err1},
	}

	for i, v := range cases {
		out := nopanic.Unwrap(v.in)
		if !reflect.DeepEqual(out, v.out) {
			t.Errorf("Case %d:  got %v, want %v", i, out, v.out)
		}
	}
}

func TestUnwrap2(t *testing.T) {
	err1 := errors.New("No luck!")

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Missing panic")
		}
	}()

	nopanic.Unwrap(nopanic.Wrap(func() error {
		panic(err1)
	}))
}
