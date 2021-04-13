package base

// CloseElements closes all of the elements contained in a slice.
// This is a utility function to help containers close all of their children
// when required.
func CloseElements(children []Element) {
	for _, v := range children {
		v.Close()
	}
}

// DiffChild adds and removes controls in a GUI to reconcile differences
// between the desired and current GUI state.  Depending on the kind for both
// lhs and rhs, the current element may either be updated or replaced.
//
// The element lhs should be considered as 'sunk' by this function, and will be
// closed if necessary.  On the other hand, the caller will be responsible
// for the returned element.  Note that the returned element may be non-nil
// even in the presence of an error.
//
// If the rhs is nil, DiffChild will still return a non-nil element.  See
// the function Method for more details.
func DiffChild(parent Control, lhs Element, rhs Widget) (Element, error) {
	// If the rhs is empty, then make sure we delete the lhs if necessary
	if rhs == nil {
		if lhs != nil {
			lhs.Close()
		}
		return (*nilElement)(nil), nil
	}

	// if the lhs is empty, then create a new element
	if lhs == nil {
		return rhs.Mount(parent)
	}

	// Can we propagate properties rather than mounting a new element?
	if kind1, kind2 := lhs.Kind(), rhs.Kind(); kind1 == kind2 {
		err := lhs.UpdateProps(rhs)
		return lhs, err
	}

	// Need to replace the element.
	newChild, err := rhs.Mount(parent)
	if err != nil {
		return lhs, err
	}
	lhs.Close()
	return newChild, nil
}

// DiffChildren will update the list of elements to match the desired state.
// As necessary, widgets will be mounted to create new elements, and existing
// elements will be updated or closed to reconcile differences between the
// desired and current GUI state.
//
// The elements contained in lhs should be considered as 'sunk' by this function,
// and will be closed if necessary.  On the other hand, the caller will be
// responsible for all of the returned elements.  Note that the returned slice
// of elements may be non-nil even in the presence of an error.  Use
// CloseElements as necessary to avoid leaking controls.
//
// DiffChildren will try to reuse the underlying array from lhs for the
// returned slice.
func DiffChildren(parent Control, lhs []Element, rhs []Widget) ([]Element, error) {
	// If the new tree does not contain any children, then we can trivially
	// match the tree by deleting the actual widgets.
	if len(rhs) == 0 {
		for _, v := range lhs {
			v.Close()
		}
		return nil, nil
	}

	// If the old tree does not contain any children, then we can trivially
	// match the tree by mounting all of the widgets.
	if len(lhs) == 0 {
		c := make([]Element, 0, len(rhs))

		for _, v := range rhs {
			mountedChild, err := v.Mount(parent)
			if err != nil {
				return nil, err
			}
			c = append(c, mountedChild)
		}

		return c, nil
	}

	// Delete excessive children
	if len(lhs) > len(rhs) {
		for _, v := range lhs[len(rhs):] {
			v.Close()
		}
		lhs = lhs[:len(rhs)]
	}

	// Update kind (if necessary) and properties for all of the currently
	// existing children.
	for i := range lhs {
		if kind1, kind2 := lhs[i].Kind(), rhs[i].Kind(); kind1 == kind2 {
			err := lhs[i].UpdateProps(rhs[i])
			if err != nil {
				return lhs, err
			}
		} else {
			mountedWidget, err := rhs[i].Mount(parent)
			if err != nil {
				return lhs, err
			}
			lhs[i].Close()
			lhs[i] = mountedWidget
		}
	}

	// Mount any remaining children.
	for _, v := range rhs[len(lhs):] {
		mountedWidget, err := v.Mount(parent)
		if err != nil {
			return lhs, err
		}
		lhs = append(lhs, mountedWidget)
	}

	return lhs, nil
}
