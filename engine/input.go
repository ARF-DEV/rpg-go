package engine

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type InputEventListener interface {
	UpdateOnInput(in *Input)
}

type Input struct {
	Keys     [1024]bool
	PrevKeys [1024]bool
	// TODO mouse event
	singleEventSubs []InputEventListener
}

var InputEvent Input = Input{}

func (e *Input) keycallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mode glfw.ModifierKey) {
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}

	// fmt.Println("owkaodka")
	e.PrevKeys = e.Keys
	if key >= 0 && key < 1024 {
		if action == glfw.Press {
			e.Keys[key] = true
		}
		if action == glfw.Release {
			e.Keys[key] = false
		}
	}

	for _, sub := range e.singleEventSubs {
		sub.UpdateOnInput(e)
	}
}

func (e *Input) Init() {
	glfw.GetCurrentContext().SetKeyCallback(glfw.KeyCallback(e.keycallback))
}

func (s *Input) AddSubcsriber(sub InputEventListener) {
	s.singleEventSubs = append(s.singleEventSubs, sub)
}

// func (e *Input) HandleEvent() {
// 	glfw.PollEvents()
// }
