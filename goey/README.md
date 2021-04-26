# Goey

Package goey provides a declarative, cross-platform GUI for the
[Go](https://golang.org/) language. The range of controls, their supported
properties and events, should roughly match what is available in HTML. However,
properties and events may be limited to support portability. Additionally,
styling of the controls will be limited, with the look of controls matching the
native platform.

[![Documentation](https://godoc.org/clipster/goey?status.svg)](http://godoc.org/clipster/goey)
[![Go Report Card](https://goreportcard.com/badge/clipster/goey)](https://goreportcard.com/report/clipster/goey) 
[![Windows Build Status](https://ci.appveyor.com/api/projects/status/3n6qnl555b5sho70?svg=true)](https://ci.appveyor.com/project/rj/goey) 

## Install

The package can be installed from the command line using the
[go](https://golang.org/cmd/go/) tool.  However, depending on your OS, please
check for special instructions below.

    go get clipster/goey

### Windows

No special instructions are required to build this package on windows.  CGO is not used.

### Linux

This package requires the use of CGO to access GTK, which must be installed.  The GTK libraries should be installed before issuing `go get` or you will have error messages during the building of some of the internal packages.

On Ubuntu:

    sudo apt-get install libgtk-3-dev

#### Linux with GNUstep

This package can be built to target Cocoa using GNUstep, which must be installed.  Most users are unlikely to want to use this option, but it can be useful for development.  The libraries for GNUstep must be installed before issuing `go get` or you will have error message during the building of some of the internal packages.

On Ubuntu:

    sudo apt-get install gnustep-devel

To force the use of GNUstep, build using the build tag `cocoa`.

### BSD

This package requires the use of CGO to access GTK, which must be installed.  The GTK libraries should be installed before issuing `go get` or you will have error messages during the building of some of the internal packages.

### MacOS

There is a in-progress port for Cocoa.  It is currently being developed using GNUstep on Linux, but has been developed based on documentation from Apple.  All controls, except for the date control (which is not available in GNUstep), are implemented.  However, additional testing, especially on Darwin, is still required.

If you can either test on Macs, or provide build systems, please contact us.

## Getting Started

Package documentation and examples are on [godoc](https://godoc.org/clipster/goey).

The minimal GUI example application is [onebutton](https://godoc.org/clipster/goey/example/onebutton), and additional example applications are in the example folder.  Some of the example show the options available for the widgets, for example [align](https://godoc.org/clipster/goey/example/align) and [paragraph](https://godoc.org/clipster/goey/example/paragraph).

New layout widgets can be developed entirely in Go.  For testing, a mock widget is provided in the [`mock` package](https://godoc.org/clipster/goey/mock).

### Windows

To get properly themed controls, a manifest is required. Please look at the
source code for the example applications for an example. The manifest needs to
be compiled with `github.com/akavel/rsrc` to create a .syso that will be
recognize by the go build program. Additionally, you could use build flags
(`-ldflags="-H windowsgui"`) to change the type of application built.

## Screenshots

| Windows    | Linux (GTK) | MacOS (Cocoa) |
|:----------:|:-----------:|:-------------:|
|![Screenshot](https://clipster/goey/raw/master/example/onebutton/onebutton_windows.png)|![Screenshot](https://clipster/goey/raw/master/example/onebutton/onebutton_gtk.png)|![Screenshot](https://clipster/goey/raw/master/example/onebutton/onebutton_cocoa.png)|
|![Screenshot](https://clipster/goey/raw/master/example/twofields/twofields_windows.png)|![Screenshot](https://clipster/goey/raw/master/example/twofields/twofields_gtk.png)|![Screenshot](https://clipster/goey/raw/master/example/twofields/twofields_cocoa.png)|
|![Screenshot](https://clipster/goey/raw/master/example/decoration/decoration_windows.png)|![Screenshot](https://clipster/goey/raw/master/example/decoration/decoration_gtk.png)|![Screenshot](https://clipster/goey/raw/master/example/decoration/decoration_cocoa.png)|
|![Screenshot](https://clipster/goey/raw/master/example/colour/colour_windows.png)|![Screenshot](https://clipster/goey/raw/master/example/colour/colour_gtk.png)|![Screenshot](https://clipster/goey/raw/master/example/colour/colour_cocoa.png)|
|![Screenshot](https://clipster/goey/raw/master/example/feettometer/feettometer_windows.png)|![Screenshot](https://clipster/goey/raw/master/example/feettometer/feettometer_gtk.png)|![Screenshot](https://clipster/goey/raw/master/example/feettometer/feettometer_cocoa.png)|
|![Screenshot](https://clipster/goey/raw/master/example/controls/controls1_windows.png)|![Screenshot](https://clipster/goey/raw/master/example/controls/controls1_gtk.png)|![Screenshot](https://clipster/goey/raw/master/example/controls/controls1_cocoa.png)|
|![Screenshot](https://clipster/goey/raw/master/example/controls/controls2_windows.png)|![Screenshot](https://clipster/goey/raw/master/example/controls/controls2_gtk.png)|![Screenshot](https://clipster/goey/raw/master/example/controls/controls2_cocoa.png)|
|![Screenshot](https://clipster/goey/raw/master/example/controls/controls3_windows.png)|![Screenshot](https://clipster/goey/raw/master/example/controls/controls3_gtk.png)|![Screenshot](https://clipster/goey/raw/master/example/controls/controls3_cocoa.png)|

## Contribute

Feedback and PRs welcome.

In particular, if anyone has the expertise to provide a port for MacOS, that would provide support for all major desktop operating systems.

[![Go Report Card](https://goreportcard.com/badge/clipster/goey)](https://goreportcard.com/report/clipster/goey)


## License

BSD (c) Robert Johnstone
