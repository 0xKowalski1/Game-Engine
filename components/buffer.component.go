package components

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

type BufferComponent struct {
	VAO uint32
	VBO uint32
	EBO uint32
}

func NewBufferComponent(vertices []float32, indices []uint32) *BufferComponent {
	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.Ptr(nil)) // Position attribute
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.Ptr(uintptr(3*4))) // Texture coordinates
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

	return &BufferComponent{
		VAO: vao,
		VBO: vbo,
		EBO: ebo,
	}
}
