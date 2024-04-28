package engine

import (
	"0xKowalski/game/entities"
	"0xKowalski/game/input"
	"0xKowalski/game/systems"
	"0xKowalski/game/window"
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// GLFW event handling must run on the main OS thread.
	runtime.LockOSThread()
}

type Engine struct {
	LastFrame float64

	Window       *window.Window
	InputManager *input.InputManager

	// ECS
	EntityStore *entities.EntityStore

	// Systems
	RenderSystem *systems.RenderSystem
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

	entityStore := entities.NewEntityStore()

	rs, err := systems.NewRenderSystem(win, entityStore)
	if err != nil {
		return nil, err
	}

	inputManager := input.NewInputManager(win.GlfwWindow)

	engine := &Engine{
		Window:       win,
		InputManager: inputManager,

		// Ecs
		EntityStore: entityStore,

		//Systems
		RenderSystem: rs,
	}

	return engine, nil
}

func (e *Engine) Run(gameLoop func()) {
	for !e.Window.GlfwWindow.ShouldClose() {
		glfw.PollEvents()

		e.InputManager.Update()

		gameLoop()

		e.RenderSystem.Update()

		e.Window.GlfwWindow.SwapBuffers() // Swap buffers to display the frame
		e.LastFrame = glfw.GetTime()
	}

	e.Cleanup()
}

func (e *Engine) Cleanup() {
	e.Window.Cleanup()
}
