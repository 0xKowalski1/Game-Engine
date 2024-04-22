package graphics

import (
	"0xKowalski/game/window"
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Renderer struct {
	VAO           uint32
	VBO           uint32
	ShaderProgram *ShaderProgram
}

func InitRenderer(win *window.Window) (*Renderer, error) {
	win.GlfwWindow.MakeContextCurrent()

	if err := gl.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize OpenGL: %w", err)
	}

	errorCode := gl.GetError()
	if errorCode != gl.NO_ERROR {
		log.Fatalf("OpenGL error after initialization: %d", errorCode)
		return nil, fmt.Errorf("OpenGL error: %d", errorCode)
	}

	// Set the initial viewport
	width, height := win.GetWidthAndHeight()
	gl.Viewport(0, 0, int32(width), int32(height))
	// Set callback so viewport changes when window size changes
	win.GlfwWindow.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	})

	renderer := &Renderer{}

	// Set up shaders
	shaderProgram, err := InitShaderProgram("assets/shaders/vertex.glsl", "assets/shaders/fragment.glsl")
	if err != nil {
		return nil, err
	}
	renderer.ShaderProgram = shaderProgram

	// Other
	gl.GenVertexArrays(1, &renderer.VAO)
	gl.GenBuffers(1, &renderer.VBO)

	vertices := []float32{
		-0.5, -0.5, 0.0, // Vertex 1 (X, Y, Z)
		0.5, -0.5, 0.0, // Vertex 2 (X, Y, Z)
		0.0, 0.5, 0.0, // Vertex 3 (X, Y, Z)
	}

	gl.BindVertexArray(renderer.VAO)
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.VBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// Position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.Ptr(nil))
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	return renderer, nil
}

func (r *Renderer) Render() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear the color and depth buffers
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)                   // Set the clear color to a dark blue

	r.ShaderProgram.Use()
	gl.BindVertexArray(r.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.BindVertexArray(0)
}

func (r *Renderer) Cleanup() {
	gl.DeleteBuffers(1, &r.VBO)
	gl.DeleteVertexArrays(1, &r.VAO)
}
