package entities

import (
	"0xKowalski/game/components"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

func generateSphere(stacks, slices int, radius float32) ([]components.Vertex, []uint32) {
	const PI = float32(math.Pi)
	vertices := []components.Vertex{}
	indices := []uint32{}

	for i := 0; i <= stacks; i++ {
		V := float32(i) / float32(stacks)
		phi := V * PI

		for j := 0; j <= slices; j++ {
			U := float32(j) / float32(slices)
			theta := U * (PI * 2)

			x := radius * float32(math.Cos(float64(theta))*math.Sin(float64(phi)))
			y := radius * float32(math.Cos(float64(phi)))
			z := radius * float32(math.Sin(float64(theta))*math.Sin(float64(phi)))

			// Convert to mgl32.Vec3 for Position and add a Vertex to the slice
			position := mgl32.Vec3{x, y, z}
			vertices = append(vertices, components.Vertex{
				Position:  position,
				TexCoords: mgl32.Vec2{U, V},
				Normal:    position.Normalize(),
			})
		}
	}

	// Calculate the indices
	for i := 0; i < slices*stacks+slices; i++ {
		indices = append(indices, uint32(i))
		indices = append(indices, uint32(i+slices+1))
		indices = append(indices, uint32(i+slices))

		indices = append(indices, uint32(i+slices+1))
		indices = append(indices, uint32(i))
		indices = append(indices, uint32(i+1))
	}

	return vertices, indices
}

func (es *EntityStore) NewSphereEntity(position mgl32.Vec3, radius float32, segments int, rings int) *Entity {
	entity := es.NewEntity()

	transform := components.NewTransformComponent(position)
	es.AddComponent(entity, transform)

	// Generate vertices and indices for the sphere
	vertices, indices := generateSphere(segments, rings, radius)

	meshComponents := make([]*components.MeshComponent, 1)
	mesh := components.NewMeshComponent(vertices, indices)
	meshComponents[0] = mesh

	bufferComponents := make([]*components.BufferComponent, 1)
	buffer := components.NewBufferComponent(vertices, indices)
	bufferComponents[0] = buffer

	materialComponents := make([]*components.MaterialComponent, 1)
	material := components.NewMaterialComponent(
		"assets/textures/container.png",
		"assets/textures/container_specular.png",
		32.0)

	materialComponents[0] = material

	modelComponent := &components.ModelComponent{
		MeshComponents:     meshComponents,
		MaterialComponents: materialComponents,
		BufferComponents:   bufferComponents,
	}
	es.AddComponent(entity, modelComponent)

	renderable := components.NewRenderableComponent(transform, modelComponent)
	es.AddComponent(entity, renderable)

	return &entity
}
