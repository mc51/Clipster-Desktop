package icons

import (
	"fmt"
	"time"

	"guitest/goey"
	"guitest/goey/base"
	"guitest/goey/loop"
)

const Build = rune(0xe869)
const Brightness1 = rune(0xe3a6)
const Brightness7 = rune(0xe3ac)

func ExampleIcon() {
	render := func(r rune) base.Widget {
		return &goey.Padding{
			Insets: goey.DefaultInsets(),
			Child: &goey.Align{
				Child: Icon(r),
			},
		}
	}

	createWindow := func() error {
		// Add the controls
		window, err := goey.NewWindow("Icons", render(Build))
		if err != nil {
			return err
		}

		go func() {
			for i := Brightness1; i <= Brightness7; i++ {
				time.Sleep(1 * time.Second)
				err := loop.Do(func() error {
					return window.SetChild(render(i))
				})
				if err != nil {
					panic(err)
				}
			}
			time.Sleep(1 * time.Second)

			err := loop.Do(func() error {
				window.Close()
				return nil
			})
			if err != nil {
				panic(err)
			}
		}()

		return nil
	}

	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("OK")

	// Output:
	// OK
}
