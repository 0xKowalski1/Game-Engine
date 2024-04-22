package engine

import (
	"0xKowalski/game/graphics"
	"0xKowalski/game/window"
	"log"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Engine struct {
	Window       *window.Window
	RenderSystem *graphics.RenderSystem
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

	err = graphics.InitOpenGL(win)
	if err != nil {
		log.Printf("Error initializing renderer: %v", err)
		win.Cleanup()
		return nil, err
	}

	rs, err := graphics.NewRenderSystem()
	if err != nil {
		return nil, err
	}

	engine := &Engine{
		Window:       win,
		RenderSystem: rs,
	}

	return engine, nil
}

func (e *Engine) Run() {
	entity := graphics.Entity(1)

	// Create a MeshComponent and assign it to the entity
	vertices := []float32{
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0,
	}
	indices := []uint32{0, 1, 2}

	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// Assuming positions are given as XYZ
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	gl.BindVertexArray(0) // Unbind the VAO

	comp := &graphics.RenderComponent{
		VAO:      vao,
		VBO:      vbo,
		EBO:      ebo,
		Vertices: vertices,
		Indices:  indices,
	}

	var renderComponents = make(map[graphics.Entity]*graphics.RenderComponent)
	renderComponents[entity] = comp
	entities := []graphics.Entity{entity}

	for !e.Window.GlfwWindow.ShouldClose() {
		e.RenderSystem.Render(entities, renderComponents)

		e.Window.GlfwWindow.SwapBuffers() // Swap buffers to display the frame
		glfw.PollEvents()
	}

	e.Cleanup()
}

func (e *Engine) Cleanup() {
	e.Window.Cleanup()
}
