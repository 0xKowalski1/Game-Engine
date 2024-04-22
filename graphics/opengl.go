package graphics

import (
	"0xKowalski/game/window"
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func InitOpenGL(win *window.Window) error {
	win.GlfwWindow.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return fmt.Errorf("failed to initialize OpenGL: %w", err)
	}

	errorCode := gl.GetError()
	if errorCode != gl.NO_ERROR {
		log.Fatalf("OpenGL error after initialization: %d", errorCode)
		return fmt.Errorf("OpenGL error: %d", errorCode)
	}

	// Set the initial viewport
	width, height := win.GetWidthAndHeight()
	gl.Viewport(0, 0, int32(width), int32(height))
	// Set callback so viewport changes when window size changes
	win.GlfwWindow.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	return nil
}
