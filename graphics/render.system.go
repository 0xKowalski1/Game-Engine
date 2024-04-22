package graphics

import (
	"github.com/go-gl/gl/v4.3-core/gl"
)

type RenderSystem struct {
	ShaderProgram *ShaderProgram
}

func NewRenderSystem() (*RenderSystem, error) {
	rs := new(RenderSystem)
	shaderProgram, err := InitShaderProgram("assets/shaders/vertex.glsl", "assets/shaders/fragment.glsl")
	if err != nil {
		return nil, err
	}

	rs.ShaderProgram = shaderProgram

	return rs, nil
}

type Entity uint32
type RenderComponent struct {
	VAO uint32
	VBO uint32
	EBO uint32

	Vertices []float32
	Indices  []uint32
}

func (rs *RenderSystem) Render(entities []Entity, renderComponents map[Entity]*RenderComponent) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear the color and depth buffers
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)                   // Set the clear color to a dark blue

	rs.ShaderProgram.Use()
	for _, entity := range entities {
		mesh := renderComponents[entity]
		gl.BindVertexArray(mesh.VAO)
		gl.DrawElements(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, nil)
		gl.BindVertexArray(0)
	}
}
