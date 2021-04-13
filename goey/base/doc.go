// Package base provides interfaces for the description, creation, and updating
// of GUI widgets.  Any users interested in creating new widgets will need to
// implement the interfaces Widget and Element described herein.
//
// GUI state is managed using three groups of types.  There are 'widgets', which
// are data-only representations of a desired GUI.  These widgets can be mounted
// to create 'elements', which manage actual, visible GUI resources.  The
// elements manage child elements, some number of platform-specific resources,
// called 'controls', or both.
//
// Additionally, this package contains geometric types.  These types support
// the automatic layout of widgets in a platform independent manner.  The
// layout algorithm is roughly based on Flutter.
package base
