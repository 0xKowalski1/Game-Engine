package components

type MeshComponent struct {
	Vertices []float32
	Indices  []uint32
}

func NewMeshComponent(vertices []float32, indices []uint32) *MeshComponent {
	return &MeshComponent{
		Vertices: vertices,
		Indices:  indices,
	}
}
