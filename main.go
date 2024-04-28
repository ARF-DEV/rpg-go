//go:build cgo
// +build cgo

package main

import (
	_ "image/jpeg"
	_ "image/png"
	"runtime"

	"log"

	"github.com/ARF-DEV/rpg-go/engine"
	"github.com/ARF-DEV/rpg-go/game"
	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func glfwInit() {
	runtime.LockOSThread()
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
}

// const WIDTH = 800
// const HEIGHT = 600

//	func glClearError() {
//		for gl.GetError() != gl.NO_ERROR {
//		}
//	}
//
//	func glCheckError() {
//		log.Printf("ERROR: %v", gl.GetError())
//	}
func main() {
	glfwInit()
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(game.WIDTH, game.HEIGHT, "Hello World", nil, nil)
	if err != nil {
		log.Println("failed to create window")
		panic(err)
	}
	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		log.Println("failed to init gl")
		panic(err)
	}
	gl.Viewport(0, 0, game.WIDTH, game.HEIGHT)

	sr := engine.CreateSpriteRenderer()
	shader, err := engine.CreateShader("assets/vertex.glsl", "assets/fragment.glsl")
	if err != nil {
		panic(err)
	}

	shader.Use()
	ortho := mgl32.Ortho(0, game.WIDTH, game.HEIGHT, 0, -1, 1)
	shader.SetMatrix4f("projection", ortho)
	shader.SetInt("texture0", 0)
	gl.UseProgram(0)

	engine.InputEvent.Init()
	gameInstance := game.Game{}
	gameInstance.Start(&engine.InputEvent)
	for !window.ShouldClose() {
		gameInstance.Update(window, &engine.InputEvent)
		gameInstance.Draw(window, &sr, &shader)
	}

	// window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {})
}
