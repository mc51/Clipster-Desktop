package goey

import (
	"testing"

	"clipster/goey/base"
)

func TestEmptyMount(t *testing.T) {
	testingMountWidgets(t,
		&Empty{},
		&Empty{},
		&Empty{},
	)
}

func TestEmptyClose(t *testing.T) {
	testingCloseWidgets(t,
		&Empty{},
		&Empty{},
	)
}

func TestEmptyUpdate(t *testing.T) {
	testingUpdateWidgets(t, []base.Widget{
		&Empty{},
		&Empty{},
	}, []base.Widget{
		&Empty{},
		&Empty{},
	})
}
