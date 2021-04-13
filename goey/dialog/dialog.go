package dialog

// Dialog contains some functionality used by all dialog types.
type Dialog struct {
	dialogImpl
	err error
}

// Err returns the first error that was encountered while building the dialog.
func (d *Dialog) Err() error {
	return d.err
}
