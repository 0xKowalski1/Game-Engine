package components

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

type RenderComponent struct {
	VAO uint32
	VBO uint32
	EBO uint32

	Vertices []float32
	Indices  []uint32
}

func NewRenderComponent() *RenderComponent {
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

	return &RenderComponent{
		VAO:      vao,
		VBO:      vbo,
		EBO:      ebo,
		Vertices: vertices,
		Indices:  indices,
	}

}
