// Package nopanic provides a utility to wrap function to prevent any panics
// from escaping.  If a panic occurs in the callback, the value of the panic
// will be converted to an error, and returned normally.
//
// This packages has the opinion that you should not panic with arbitrary
// values.  The argument to panic should either be a string, or a value that
// implements the error interface.
package nopanic
