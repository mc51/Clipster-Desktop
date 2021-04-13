// This package provides an example application built using the goey package
// that demonstrates most of the controls that are available.
package main

import (
	"fmt"
	"time"

	"guitest/goey"
	"guitest/goey/base"
	"guitest/goey/loop"
)

var (
	currentCD  string
	s1, s2, s3 bool
	showLorem  bool
	window     *goey.Window
)

func main() {
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	w, err := goey.NewWindow("Controls", renderWindow())
	if err != nil {
		return err
	}
	w.SetScroll(false, true)
	window = w
	return nil
}

func updateWindow() {
	window.SetChild(renderWindow())
}

func renderWindow() base.Widget {
	widget := &goey.Tabs{
		Insets: goey.DefaultInsets(),
		Children: []goey.TabItem{
			{
				Caption: "Input",
				Child: &goey.VBox{
					Children: []base.Widget{
						&goey.Label{Text: "Text input:"},
						&goey.TextInput{Value: "Some input...", Placeholder: "Type some text here.  And some more.  And something really long.",
							OnChange: func(v string) { println("text input ", v) }, OnEnterKey: func(v string) { println("t1* ", v) }},
						&goey.Label{Text: "Password input:"},
						&goey.TextInput{Value: "", Placeholder: "Don't share", Password: true,
							OnChange: func(v string) { println("password input ", v) }},
						&goey.Label{Text: "Integer input:"},
						&goey.IntInput{Value: 3, Placeholder: "Please enter a number",
							Min: -100, Max: 100,
							OnChange: func(v int64) { println("int input ", v) }},
						&goey.Label{Text: "Date input:"},
						&goey.DateInput{Value: time.Now().Add(24 * time.Hour),
							OnChange: func(v time.Time) { println("date input: ", v.String()) }},
						&goey.Label{Text: "Select input (combobox):"},
						&goey.SelectInput{Items: []string{"Choice 1", "Choice 2", "Choice 3"},
							OnChange: func(v int) { println("select input: ", v) }},
						&goey.Label{Text: "Number input:"},
						&goey.Slider{Value: 25, Min: 0, Max: 100,
							OnChange: func(v float64) { println("slider input: ", v) }},
						&goey.HR{},
						&goey.Expand{Child: &goey.TextArea{Value: "", Placeholder: "Room to write",
							OnChange: func(v string) { println("text area: ", v) },
						}},
					},
				},
			},
			{
				Caption: "Buttons",
				Child: &goey.VBox{
					Children: []base.Widget{
						&goey.HBox{Children: []base.Widget{
							&goey.Button{Text: "Left 1", Default: true},
							&goey.Button{Text: "Left 2"},
						}},
						&goey.HBox{
							Children: []base.Widget{
								&goey.Button{Text: "Center"},
							},
							AlignMain: goey.MainCenter,
						},
						&goey.HBox{
							Children: []base.Widget{
								&goey.Button{Text: "D1"},
								&goey.Button{Text: "D2", Disabled: true},
								&goey.Button{Text: "D3"},
							},
							AlignMain: goey.MainEnd,
						},
						&goey.HR{},
						&goey.Label{Text: "Check boxes:"},
						&goey.Checkbox{Value: true, Text: "Please click on the checkbox A",
							OnChange: func(v bool) { println("check box input: ", v) }},
						&goey.Checkbox{Text: "Please click on the checkbox B",
							OnChange: func(v bool) { println("check box input: ", v) }},
					},
				},
			},
			{
				Caption: "Lorem",
				Child: &goey.VBox{
					Children: []base.Widget{
						&goey.P{Text: lorem, Align: goey.JustifyFull},
						&goey.P{Text: "This is a paragraph, but without much text.", Align: goey.JustifyLeft},
						&goey.P{Text: "This is a paragraph, but without much text.", Align: goey.JustifyCenter},
						&goey.P{Text: "This is a paragraph, but without much text.", Align: goey.JustifyRight},
					},
					AlignMain: goey.MainCenter,
				},
			},
		},
	}
	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child:  widget,
	}
}
