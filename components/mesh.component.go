package components

import "github.com/go-gl/mathgl/mgl32"

type Vertex struct {
	Position  mgl32.Vec3
	TexCoords mgl32.Vec2
	Normal    mgl32.Vec3
}

type MeshComponent struct {
	Vertices []Vertex
	Indices  []uint32
}

func NewMeshComponent(vertices []Vertex, indices []uint32) *MeshComponent {
	return &MeshComponent{
		Vertices: vertices,
		Indices:  indices,
	}
}
