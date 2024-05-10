package entities

import (
	"0xKowalski/game/components"

	"github.com/go-gl/mathgl/mgl32"
)

var defaultPlaneVertices = []components.Vertex{
	// First triangle
	{Position: mgl32.Vec3{-0.5, 0.0, -0.5}, TexCoords: mgl32.Vec2{0.0, 0.0}, Normal: mgl32.Vec3{0.0, 1.0, 0.0}},
	{Position: mgl32.Vec3{0.5, 0.0, -0.5}, TexCoords: mgl32.Vec2{1.0, 0.0}, Normal: mgl32.Vec3{0.0, 1.0, 0.0}},
	{Position: mgl32.Vec3{0.5, 0.0, 0.5}, TexCoords: mgl32.Vec2{1.0, 1.0}, Normal: mgl32.Vec3{0.0, 1.0, 0.0}},

	// Second triangle
	{Position: mgl32.Vec3{-0.5, 0.0, -0.5}, TexCoords: mgl32.Vec2{0.0, 0.0}, Normal: mgl32.Vec3{0.0, 1.0, 0.0}},
	{Position: mgl32.Vec3{0.5, 0.0, 0.5}, TexCoords: mgl32.Vec2{1.0, 1.0}, Normal: mgl32.Vec3{0.0, 1.0, 0.0}},
	{Position: mgl32.Vec3{-0.5, 0.0, 0.5}, TexCoords: mgl32.Vec2{0.0, 1.0}, Normal: mgl32.Vec3{0.0, 1.0, 0.0}},
}

var defaultPlaneIndices = []uint32{
	0, 1, 2, 3, 4, 5,
}

func (es *EntityStore) NewPlaneEntity(position mgl32.Vec3) *Entity {
	entity := es.NewEntity()

	transform := components.NewTransformComponent(position)
	es.AddComponent(entity, transform)
	transform.SetScale(50, 50, 50)

	// Mesh component for the plane
	meshComponents := make([]*components.MeshComponent, 1)
	mesh := components.NewMeshComponent(defaultPlaneVertices, defaultPlaneIndices)
	meshComponents[0] = mesh

	// Buffer component for the plane
	bufferComponents := make([]*components.BufferComponent, 1)
	buffer := components.NewBufferComponent(defaultPlaneVertices, defaultPlaneIndices)
	bufferComponents[0] = buffer

	// Material component for the plane
	materialComponents := make([]*components.MaterialComponent, 1)
	material := components.NewMaterialComponent(
		"assets/textures/container.png",
		"assets/textures/container_specular.png",
		32.0)
	materialComponents[0] = material

	// Model component that aggregates mesh, material, and buffer components
	modelComponent := &components.ModelComponent{
		MeshComponents:     meshComponents,
		MaterialComponents: materialComponents,
		BufferComponents:   bufferComponents,
	}
	es.AddComponent(entity, modelComponent)

	// Renderable component to integrate with the rendering system
	renderable := components.NewRenderableComponent(transform, modelComponent)
	es.AddComponent(entity, renderable)

	return &entity
}
