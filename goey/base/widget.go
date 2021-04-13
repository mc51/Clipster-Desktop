package base

// Kind identifies the different kinds of widgets.  Most widgets have two
// concrete types associated with their behavior.  First, there is a type with data
// to describe the widget when unmounted, which should implement the interface
// Widget.  Second, there is a type with a handle to the windowing system when
// mounted, which should implement the interface Element.  Automatic reconciliation
// of two widget trees relies on Kind to match the unmounted and mounted widgets.
//
// Note that comparison of kinds is done by address, and not done using the value of any fields.
// Any internal state is simply to help with debugging.
type Kind struct {
	name string
}

// NewKind creates a new kind.  The name should identify the type used for the widget,
// but is currently unused.
func NewKind(name string) Kind {
	return Kind{name}
}

// String returns the string with the name of the widget and element kind.
func (k Kind) String() string {
	return k.name
}

// Widget is an interface that wraps any type describing part of a GUI.
// A widget can be 'mounted' to create controls using the platform GUI.
type Widget interface {
	// Kind returns the concrete type's Kind.  The returned value should
	// be constant, and the same for all instances of a concrete type.
	// Users should not need to use this method directly.
	Kind() *Kind
	// Mount creates a widget or control in the GUI.  The newly created widget
	// will be a child of the widget specified by parent.  If non-nil, the returned
	// Element must have a matching kind.
	Mount(parent Control) (Element, error)
}

// Element is an interface that wraps any type representing a control, or group
// of controls, created using the platform GUI.  An element represents an
// instantiation of a Widget into visible parts of the GUI.
type Element interface {
	// NativeElement provides platform-dependent methods.  These should
	// not be used by client libraries, but exist for the internal implementation
	// of platform dependent code.
	NativeElement

	// Close removes the element from the GUI, and frees any associated resources.
	Close()
	// Kind returns the concrete type for the Element.
	// Users should not need to use this method directly.
	Kind() *Kind
	// Layout determines the best size for an element that satisfies the
	// constraints.
	Layout(Constraints) Size
	// MinIntrinsicHeight returns the minimum height that this element requires
	// to be correctly displayed.
	MinIntrinsicHeight(width Length) Length
	// MinIntrinsicWidth returns the minimum width that this element requires
	// to be correctly displayed.
	MinIntrinsicWidth(height Length) Length
	// SetBounds updates the position of the widget.
	SetBounds(bounds Rectangle)
	// UpdateProps will update the properties of the widget.  The Kind for
	// the parameter data must match the Kind for the interface.
	UpdateProps(data Widget) error
}
