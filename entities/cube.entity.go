package entities

import (
	"0xKowalski/game/components"

	"github.com/go-gl/mathgl/mgl32"
)

func generateCube(size float32) ([]components.Vertex, []uint32) {
	halfSize := size / 2
	vertices := make([]components.Vertex, 0, 24) // 24 vertices for 6 faces with unique normals
	indices := make([]uint32, 0, 36)             // 6 faces * 2 triangles * 3 vertices

	// Define cube offsets and normals for each face
	faces := []struct {
		normal   mgl32.Vec3
		corners  [4]mgl32.Vec3
		texCoord [4]mgl32.Vec2
	}{
		{mgl32.Vec3{0, 0, -1}, [4]mgl32.Vec3{{-halfSize, -halfSize, -halfSize}, {halfSize, -halfSize, -halfSize}, {halfSize, halfSize, -halfSize}, {-halfSize, halfSize, -halfSize}}, [4]mgl32.Vec2{{0, 0}, {1, 0}, {1, 1}, {0, 1}}},
		{mgl32.Vec3{0, 0, 1}, [4]mgl32.Vec3{{halfSize, -halfSize, halfSize}, {-halfSize, -halfSize, halfSize}, {-halfSize, halfSize, halfSize}, {halfSize, halfSize, halfSize}}, [4]mgl32.Vec2{{0, 0}, {1, 0}, {1, 1}, {0, 1}}},
		{mgl32.Vec3{-1, 0, 0}, [4]mgl32.Vec3{{-halfSize, -halfSize, halfSize}, {-halfSize, -halfSize, -halfSize}, {-halfSize, halfSize, -halfSize}, {-halfSize, halfSize, halfSize}}, [4]mgl32.Vec2{{0, 0}, {1, 0}, {1, 1}, {0, 1}}},
		{mgl32.Vec3{1, 0, 0}, [4]mgl32.Vec3{{halfSize, -halfSize, -halfSize}, {halfSize, -halfSize, halfSize}, {halfSize, halfSize, halfSize}, {halfSize, halfSize, -halfSize}}, [4]mgl32.Vec2{{0, 0}, {1, 0}, {1, 1}, {0, 1}}},
		{mgl32.Vec3{0, 1, 0}, [4]mgl32.Vec3{{-halfSize, halfSize, -halfSize}, {halfSize, halfSize, -halfSize}, {halfSize, halfSize, halfSize}, {-halfSize, halfSize, halfSize}}, [4]mgl32.Vec2{{0, 0}, {1, 0}, {1, 1}, {0, 1}}},
		{mgl32.Vec3{0, -1, 0}, [4]mgl32.Vec3{{-halfSize, -halfSize, halfSize}, {halfSize, -halfSize, halfSize}, {halfSize, -halfSize, -halfSize}, {-halfSize, -halfSize, -halfSize}}, [4]mgl32.Vec2{{0, 0}, {1, 0}, {1, 1}, {0, 1}}},
	}

	// Generate vertices and indices for each face
	for _, face := range faces {
		baseIndex := uint32(len(vertices))
		for j, corner := range face.corners {
			vertices = append(vertices, components.Vertex{
				Position:  corner,
				TexCoords: face.texCoord[j],
				Normal:    face.normal,
			})
		}
		indices = append(indices,
			baseIndex, baseIndex+1, baseIndex+2,
			baseIndex, baseIndex+2, baseIndex+3,
		)
	}

	return vertices, indices
}

type CubeOption func(*EntityStore, *Entity)

func (es *EntityStore) NewCubeEntity(position mgl32.Vec3, size float32, opts ...CubeOption) *Entity {
	entity := es.NewEntity()

	vertices, indices := generateCube(size)

	transform := components.NewTransformComponent(position)
	es.AddComponent(entity, transform)

	meshComponents := make([]*components.MeshComponent, 1)
	mesh := components.NewMeshComponent(vertices, indices)
	meshComponents[0] = mesh

	// Similarly, initialize the buffer components slice
	bufferComponents := make([]*components.BufferComponent, 1)
	buffer := components.NewBufferComponent(vertices, indices)
	bufferComponents[0] = buffer

	// Initialize the material components slice
	materialComponents := make([]*components.MaterialComponent, 1)
	material := components.NewMaterialComponent(
		"assets/textures/container.png",
		"assets/textures/container_specular.png",
		32.0)

	materialComponents[0] = material

	modelComponent := &components.ModelComponent{
		MeshComponents:     meshComponents,
		MaterialComponents: materialComponents,
		BufferComponents:   bufferComponents}
	es.AddComponent(entity, modelComponent)

	// Apply any additional options
	for _, opt := range opts {
		opt(es, &entity)
	}

	renderable := components.NewRenderableComponent(transform, modelComponent)
	es.AddComponent(entity, renderable)

	return &entity
}
