// Package goey provides a declarative, cross-platform GUI.
// The range of controls, their supported properties and events, should roughly
// match what is available in HTML.  However, properties and events may be
// limited to support portability.  Additionally, styling of the controls
// will be limited, with the look of controls matching the native platform.
//
// The minimal GUI example application is clipster/goey/example/onebutton,
// and additional example applications are in the example folder.
//
// Screenshots
//
// There are screenshots of some of the example applications linked in the
// README, located at https://clipster/goey/src/default/README.md.
//
// Windows
//
// To get properly themed controls, a manifest is required.  Please
// look at the source for the example applications for an example.  This file
// needs to be compiled with github.com/akavel/rsrc to create a .syso that will
// be recognize by the go build program.  Additionally, you could use build flags
// (-ldflags="-H windowsgui") to change the type of application built.
//
// CGO is not required.
//
// Linux
//
// Although this package does not use CGO, some of its dependencies
// do.  The build machine also requires that GTK+ 3 is installed.
// This should be installed before issuing `go get` or you will have error
// messages during the building of some of the dependencies.
//
// On Ubuntu:
//
//      sudo apt-get install libgtk-3-dev
//
// Darwin (MacOS)
//
// A port to darwin using the Cocoa API is in the repository, but is only
// available under the macos branch.  Development has been done on linux using
// GNUstep, so the build tags will need to be updated to build on a darwin
// system.
package goey
