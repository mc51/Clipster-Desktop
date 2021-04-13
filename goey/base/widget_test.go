package base

import (
	"fmt"
)

func ExampleKind_String() {
	kind := NewKind("guitest/goey/base.Example")

	fmt.Println("Kind is", kind.String())

	// Output:
	// Kind is guitest/goey/base.Example
}
