// This package provides an example application built using the goey package
// that rebuilds the classic Todos tutorial application.
package main

import (
	"fmt"
	"strconv"

	"guitest/goey"
	"guitest/goey/base"
	"guitest/goey/loop"
)

var (
	mainWindow *goey.Window
)

func main() {
	err := loop.Run(createWindow)
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func createWindow() error {
	mw, err := goey.NewWindow("Example", render())
	if err != nil {
		return err
	}
	mw.SetScroll(false, false)
	mainWindow = mw
	return nil
}

func update() {
	err := mainWindow.SetChild(render())
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func render() base.Widget {
	widgets := []base.Widget{
		&goey.Label{Text: "What needs to be done:"},
		&goey.TextInput{Placeholder: "Enter todo description.", OnEnterKey: onNewTodoItem},
	}
	count1, count2 := getItemCounts()
	if count2 > 0 {
		widgets = append(widgets, &goey.HR{})
		widgets = append(widgets, &goey.Label{Text: "There are " + strconv.Itoa(count2) + " waiting item(s)."})
		for i, v := range Model {
			if !v.Completed {
				index := i
				widgets = append(widgets, &goey.Checkbox{Text: v.Text, Value: v.Completed,
					OnChange: func(newValue bool) {
						Model[index].Completed = newValue
						update()
					}})
			}
		}
	}
	if count1 > 0 {
		widgets = append(widgets, &goey.HR{})
		widgets = append(widgets, &goey.Label{Text: "There are " + strconv.Itoa(count1) + " completed item(s)."})
		for i, v := range Model {
			if v.Completed {
				index := i
				widgets = append(widgets, &goey.Checkbox{Text: v.Text, Value: v.Completed,
					OnChange: func(newValue bool) {
						Model[index].Completed = newValue
						update()
					}})
			}
		}
	}

	return &goey.Padding{
		Insets: goey.DefaultInsets(),
		Child:  &goey.VBox{Children: widgets},
	}
}

func onNewTodoItem(value string) {
	Model = append(Model, TodoItem{Text: value})
	update()
}
