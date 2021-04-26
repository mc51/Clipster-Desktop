package goey

import (
	"clipster/goey/base"
	"testing"
)

func TestHRMount(t *testing.T) {
	testingMountWidgets(t,
		&HR{},
		&HR{},
		&HR{},
	)
}

func TestHRClose(t *testing.T) {
	testingCloseWidgets(t,
		&HR{},
		&HR{},
	)
}

func TestHRUpdate(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&HR{},
		&HR{},
	}, []base.Widget{
		&HR{},
		&HR{},
	})
}
