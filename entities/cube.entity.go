package entities

import (
	"0xKowalski/game/components"
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

var defaultVertices = []float32{
	// Front face
	// 3*Position     2*Tex     3*Normals
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

type CubeOption func(*EntityStore, *Entity)

func WithTexture(texturePath string) CubeOption {
	return func(es *EntityStore, e *Entity) {
		texture, err := components.NewTextureComponent(texturePath)
		if err != nil {
			log.Printf("Error creating texture component: %v", err)
		} else {
			es.AddComponent(*e, texture)
		}
	}
}

func (es *EntityStore) NewCubeEntity(position mgl32.Vec3, opts ...CubeOption) *Entity {
	entity := es.NewEntity()

	transform := components.NewTransformComponent(position)
	es.AddComponent(entity, transform)

	// Might want to allow vertices/indicies in params in future
	mesh := components.NewMeshComponent(defaultVertices, defaultIndices)
	es.AddComponent(entity, mesh)

	buffer := components.NewBufferComponent(defaultVertices, defaultIndices)
	es.AddComponent(entity, buffer)

	// Apply any additional options
	for _, opt := range opts {
		opt(es, &entity)
	}

	texture, _ := es.GetComponent(entity, &components.TextureComponent{}).(*components.TextureComponent) // Can be nil
	renderable := components.NewRenderableComponent(mesh, buffer, transform, texture)
	es.AddComponent(entity, renderable)

	return &entity
}
