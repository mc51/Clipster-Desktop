package base

import (
	"fmt"
)

func ExampleKind_String() {
	kind := NewKind("clipster/goey/base.Example")

	fmt.Println("Kind is", kind.String())

	// Output:
	// Kind is clipster/goey/base.Example
}
