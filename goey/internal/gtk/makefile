all: gtk.go
	go build

gtk.go: thunks.h thunks.i
	swig -go -cgo -intgosize 64 thunks.i

fmt:
	-$(RM) gtk.go
	clang-format-8 -i *.c *.h
	gofmt -w *.go
