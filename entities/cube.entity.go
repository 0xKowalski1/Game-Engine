package entities

import (
	"0xKowalski/game/components"
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

var defaultVertices = []components.Vertex{
	// Front face
	{Position: mgl32.Vec3{-0.5, -0.5, -0.5}, TexCoords: mgl32.Vec2{0.0, 0.0}, Normal: mgl32.Vec3{0.0, 0.0, -1.0}},
	{Position: mgl32.Vec3{0.5, -0.5, -0.5}, TexCoords: mgl32.Vec2{1.0, 0.0}, Normal: mgl32.Vec3{0.0, 0.0, -1.0}},
	{Position: mgl32.Vec3{0.5, 0.5, -0.5}, TexCoords: mgl32.Vec2{1.0, 1.0}, Normal: mgl32.Vec3{0.0, 0.0, -1.0}},
	{Position: mgl32.Vec3{-0.5, 0.5, -0.5}, TexCoords: mgl32.Vec2{0.0, 1.0}, Normal: mgl32.Vec3{0.0, 0.0, -1.0}},

	// Back face
	{Position: mgl32.Vec3{-0.5, -0.5, 0.5}, TexCoords: mgl32.Vec2{0.0, 0.0}, Normal: mgl32.Vec3{0.0, 0.0, 1.0}},
	{Position: mgl32.Vec3{0.5, -0.5, 0.5}, TexCoords: mgl32.Vec2{1.0, 0.0}, Normal: mgl32.Vec3{0.0, 0.0, 1.0}},
	{Position: mgl32.Vec3{0.5, 0.5, 0.5}, TexCoords: mgl32.Vec2{1.0, 1.0}, Normal: mgl32.Vec3{0.0, 0.0, 1.0}},
	{Position: mgl32.Vec3{-0.5, 0.5, 0.5}, TexCoords: mgl32.Vec2{0.0, 1.0}, Normal: mgl32.Vec3{0.0, 0.0, 1.0}},

	// Left face
	{Position: mgl32.Vec3{-0.5, -0.5, -0.5}, TexCoords: mgl32.Vec2{0.0, 0.0}, Normal: mgl32.Vec3{-1.0, 0.0, 0.0}},
	{Position: mgl32.Vec3{-0.5, -0.5, 0.5}, TexCoords: mgl32.Vec2{1.0, 0.0}, Normal: mgl32.Vec3{-1.0, 0.0, 0.0}},
	{Position: mgl32.Vec3{-0.5, 0.5, 0.5}, TexCoords: mgl32.Vec2{1.0, 1.0}, Normal: mgl32.Vec3{-1.0, 0.0, 0.0}},
	{Position: mgl32.Vec3{-0.5, 0.5, -0.5}, TexCoords: mgl32.Vec2{0.0, 1.0}, Normal: mgl32.Vec3{-1.0, 0.0, 0.0}},

	// Right face
	{Position: mgl32.Vec3{0.5, -0.5, -0.5}, TexCoords: mgl32.Vec2{0.0, 0.0}, Normal: mgl32.Vec3{1.0, 0.0, 0.0}},
	{Position: mgl32.Vec3{0.5, -0.5, 0.5}, TexCoords: mgl32.Vec2{1.0, 0.0}, Normal: mgl32.Vec3{1.0, 0.0, 0.0}},
	{Position: mgl32.Vec3{0.5, 0.5, 0.5}, TexCoords: mgl32.Vec2{1.0, 1.0}, Normal: mgl32.Vec3{1.0, 0.0, 0.0}},
	{Position: mgl32.Vec3{0.5, 0.5, -0.5}, TexCoords: mgl32.Vec2{0.0, 1.0}, Normal: mgl32.Vec3{1.0, 0.0, 0.0}},

	// Top face
	{Position: mgl32.Vec3{-0.5, 0.5, -0.5}, TexCoords: mgl32.Vec2{0.0, 0.0}, Normal: mgl32.Vec3{0.0, 1.0, 0.0}},
	{Position: mgl32.Vec3{0.5, 0.5, -0.5}, TexCoords: mgl32.Vec2{1.0, 0.0}, Normal: mgl32.Vec3{0.0, 1.0, 0.0}},
	{Position: mgl32.Vec3{0.5, 0.5, 0.5}, TexCoords: mgl32.Vec2{1.0, 1.0}, Normal: mgl32.Vec3{0.0, 1.0, 0.0}},
	{Position: mgl32.Vec3{-0.5, 0.5, 0.5}, TexCoords: mgl32.Vec2{0.0, 1.0}, Normal: mgl32.Vec3{0.0, 1.0, 0.0}},

	// Bottom face
	{Position: mgl32.Vec3{-0.5, -0.5, -0.5}, TexCoords: mgl32.Vec2{0.0, 0.0}, Normal: mgl32.Vec3{0.0, -1.0, 0.0}},
	{Position: mgl32.Vec3{0.5, -0.5, -0.5}, TexCoords: mgl32.Vec2{1.0, 0.0}, Normal: mgl32.Vec3{0.0, -1.0, 0.0}},
	{Position: mgl32.Vec3{0.5, -0.5, 0.5}, TexCoords: mgl32.Vec2{1.0, 1.0}, Normal: mgl32.Vec3{0.0, -1.0, 0.0}},
	{Position: mgl32.Vec3{-0.5, -0.5, 0.5}, TexCoords: mgl32.Vec2{0.0, 1.0}, Normal: mgl32.Vec3{0.0, -1.0, 0.0}},
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

func (es *EntityStore) NewCubeEntity(position mgl32.Vec3, opts ...CubeOption) *Entity {
	entity := es.NewEntity()

	transform := components.NewTransformComponent(position)
	es.AddComponent(entity, transform)

	meshComponents := make([]*components.MeshComponent, 1)
	mesh := components.NewMeshComponent(defaultVertices, defaultIndices)
	meshComponents[0] = mesh

	// Similarly, initialize the buffer components slice
	bufferComponents := make([]*components.BufferComponent, 1)
	buffer := components.NewBufferComponent(defaultVertices, defaultIndices)
	bufferComponents[0] = buffer

	// Initialize the material components slice
	materialComponents := make([]*components.MaterialComponent, 1)
	material, err := components.NewMaterialComponent(
		"assets/textures/container.png",
		"assets/textures/container_specular.png",
		32.0)
	if err != nil {
		log.Printf("Error creating material component for cube: %v", err)
	}
	materialComponents[0] = material

	modelComponent := components.NewModelComponent(meshComponents, materialComponents, bufferComponents)
	es.AddComponent(entity, modelComponent)

	// Apply any additional options
	for _, opt := range opts {
		opt(es, &entity)
	}

	renderable := components.NewRenderableComponent(transform, modelComponent)
	es.AddComponent(entity, renderable)

	return &entity
}
