package gtkloop

// #cgo pkg-config: gtk+-3.0
// #include <gtk/gtk.h>
// #include "thunks.h"
import "C"

var (
	invokeFunction = make(chan func(), 1)
)

func Init() {
	C.gtk_init(nil, nil)
}

func Run() {
	C.gtk_main()
}

func Stop() {
	C.gtk_main_quit()
}

// MainContextInvoke is a wrapper around g_main_context_invoke.
func MainContextInvoke(function func()) {
	invokeFunction <- function
	C.loopMainContextInvoke()
}

// IdleAdd is a wrapper around g_idle_add.
func IdleAdd(function func()) {
	invokeFunction <- function
	C.loopIdleAdd()
}

//export mainContextInvokeCallback
func mainContextInvokeCallback() {
	fn := <-invokeFunction
	fn()
}
