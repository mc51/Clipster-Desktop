// Package mock provides a mock widget to be used for testing the layout
// algorithms of container widgets.
//
// Unlike most other widgets, the type for the element is also made public.
// Container widgets can therefore directly created instances, as if they
// were mounted, to test layout algorithms.
//
// If used in a real GUI, the mock object will not create any controls.
package mock
