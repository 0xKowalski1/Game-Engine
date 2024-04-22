package graphics

import (
	"0xKowalski/game/window"
	"fmt"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Renderer struct {
}

func InitRenderer(win *window.Window) (*Renderer, error) {

	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize OpenGL: %w", err)
	}

	// Set the initial viewport
	width, height := win.GetWidthAndHeight()
	gl.Viewport(0, 0, int32(width), int32(height))
	// Set callback so viewport changes when window size changes
	win.GlfwWindow.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	return &Renderer{}, nil
}

func (r *Renderer) Render() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear the color and depth buffers
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)                   // Set the clear color to a dark blue

}
