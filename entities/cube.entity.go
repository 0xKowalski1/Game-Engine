package entities

import (
	"0xKowalski/game/components"
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

var defaultVertices = []float32{
	// Front face
	// Position         // Tex coords  // Normals
	-0.5, -0.5, -0.5, 0.0, 0.0, 0.0, 0.0, -1.0,
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, 0.0, -1.0,
	0.5, 0.5, -0.5, 1.0, 1.0, 0.0, 0.0, -1.0,
	-0.5, 0.5, -0.5, 0.0, 1.0, 0.0, 0.0, -1.0,
	// Back face
	-0.5, -0.5, 0.5, 0.0, 0.0, 0.0, 0.0, 1.0,
	0.5, -0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 1.0,
	0.5, 0.5, 0.5, 1.0, 1.0, 0.0, 0.0, 1.0,
	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0,
	// Left face
	-0.5, -0.5, -0.5, 0.0, 0.0, -1.0, 0.0, 0.0,
	-0.5, -0.5, 0.5, 1.0, 0.0, -1.0, 0.0, 0.0,
	-0.5, 0.5, 0.5, 1.0, 1.0, -1.0, 0.0, 0.0,
	-0.5, 0.5, -0.5, 0.0, 1.0, -1.0, 0.0, 0.0,
	// Right face
	0.5, -0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 0.0, 1.0, 0.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 1.0, 1.0, 0.0, 0.0,
	0.5, 0.5, -0.5, 0.0, 1.0, 1.0, 0.0, 0.0,
	// Top face
	-0.5, 0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0,
	0.5, 0.5, -0.5, 1.0, 0.0, 0.0, 1.0, 0.0,
	0.5, 0.5, 0.5, 1.0, 1.0, 0.0, 1.0, 0.0,
	-0.5, 0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0,
	// Bottom face
	-0.5, -0.5, -0.5, 0.0, 0.0, 0.0, -1.0, 0.0,
	0.5, -0.5, -0.5, 1.0, 0.0, 0.0, -1.0, 0.0,
	0.5, -0.5, 0.5, 1.0, 1.0, 0.0, -1.0, 0.0,
	-0.5, -0.5, 0.5, 0.0, 1.0, 0.0, -1.0, 0.0,
}

var defaultIndices = []uint32{
	// Front face
	0, 1, 2, 0, 2, 3,
	// Back face
	4, 5, 6, 4, 6, 7,
	// Left face
	8, 9, 10, 8, 10, 11,
	// Right face
	12, 13, 14, 12, 14, 15,
	// Top face
	16, 17, 18, 16, 18, 19,
	// Bottom face
	20, 21, 22, 20, 22, 23,
}

type NewCubeOptions struct {
	TexturePath string
	Vertices    []float32
	Indicies    []uint32
}

func (es *EntityStore) NewCubeEntity(position mgl32.Vec3, texturePath string) *Entity {
	entity := es.NewEntity()

	transform := components.NewTransformComponent(position)
	es.AddComponent(entity, transform)

	mesh := components.NewMeshComponent(defaultVertices, defaultIndices)
	es.AddComponent(entity, mesh)

	texture, err := components.NewTextureComponent(texturePath)
	if err != nil {
		log.Printf("Error creating texture component: %v", err)
	}
	es.AddComponent(entity, texture)

	buffers := components.NewBufferComponent(defaultVertices, defaultIndices)
	es.AddComponent(entity, buffers)

	return &entity
}
