package engine

import (
	"0xKowalski/game/graphics"
	"0xKowalski/game/window"
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Engine struct {
	Renderer *graphics.Renderer
	Window   *window.Window
}

func InitEngine() (*Engine, error) {
	winConfig := window.WindowConfig{
		Title:  "Game Window",
		Width:  800,
		Height: 600,
	}

	win, err := window.InitWindow(winConfig)
	if err != nil {
		log.Printf("Failed to create window: %v", err)
		return nil, err
	}

	renderer, err := graphics.InitRenderer(win)
	if err != nil {
		log.Printf("Error initializing renderer: %v", err)
		win.Terminate()
		return nil, err
	}

	engine := &Engine{
		Renderer: renderer,
		Window:   win,
	}

	return engine, nil
}

func (e *Engine) Run() {
	for !e.Window.GlfwWindow.ShouldClose() {
		e.Renderer.Render() // Call the renderer to draw the frame

		e.Window.GlfwWindow.SwapBuffers() // Swap buffers to display the frame
		glfw.PollEvents()
	}

	e.Cleanup()
}

func (e *Engine) Cleanup() {
	e.Window.Terminate()
}
