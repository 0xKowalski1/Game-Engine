package engine

import (
	"0xKowalski/game/graphics"
	"0xKowalski/game/window"
	"log"
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

	renderer, err := graphics.InitRenderer()
	if err != nil {
		log.Printf("Error initializing renderer: %v", err)
		win.Terminate()
		return nil, err
	}

	return &Engine{Renderer: renderer, Window: win}, nil
}

func (e *Engine) Cleanup() {
	e.Window.Terminate()
}
