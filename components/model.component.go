package components

type ModelComponent struct {
	MeshComponents     []*MeshComponent
	MaterialComponents []*MaterialComponent
	BufferComponents   []*BufferComponent
}

func NewModelComponent(meshComponents []*MeshComponent, materialComponents []*MaterialComponent, bufferComponents []*BufferComponent) *ModelComponent {
	return &ModelComponent{
		MeshComponents:     meshComponents,
		MaterialComponents: materialComponents,
		BufferComponents:   bufferComponents,
	}
}
