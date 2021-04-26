package icons

import (
	"clipster/goey"
	"clipster/goey/base"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// Icon describes a widget that shows an icon as an image.
type Icon rune

var (
	kind   = base.NewKind("clipster/goey/icons.Icon")
	assets struct {
		font *truetype.Font
		face font.Face
	}
)

func init() {
	var err error
	assets.font, err = truetype.Parse(file0[:])
	if err != nil {
		panic("internal error: failed to parse embedded truetype file")
	}

	assets.face = truetype.NewFace(assets.font, &truetype.Options{Size: 32})
}

// New returns a new widget description an image showing the icon with the
// specified rune.
func New(r rune) Icon {
	return Icon(r)
}

// Kind returns the concrete type for use in the Widget interface.
// Users should not need to use this method directly.
func (i Icon) Kind() *base.Kind {
	return &kind
}

// Mount creates a control in the GUI to display the icon.
// The newly created widget will be a child of the widget specified by parent.
func (i Icon) Mount(parent base.Control) (base.Element, error) {
	img, err := DrawImage(rune(i))
	if err != nil {
		return nil, err
	}

	widget := goey.Img{Image: img}
	elem, err := widget.Mount(parent)
	if err != nil {
		return nil, err
	}

	return &iconElement{parent, elem, rune(i)}, nil
}
