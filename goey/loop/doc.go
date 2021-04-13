// Package loop provides a GUI event loop.  The event loop will be locked to
// an OS thread and not return until all GUI elements have been destroyed.
// However, callbacks can be dispatched to the event loop, and will be executed
// on the same thread.
//
// Cocoa: Unlike the other platform, the run loop not only needs to be locked
// to an OS thread, it must be locked to the main thread.  This is the first
// thread in an application.  For this reason, the init function for this
// package will call runtime.LockOSThread() to ensure that the thread is
// locked as soon as possible.  Additionally, Run must be called from same
// goroutine (i.e., same OS thread) as used to run main.
package loop
