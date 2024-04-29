package components

type RenderableComponent struct {
	MeshComponent      *MeshComponent
	BufferComponent    *BufferComponent
	TransformComponent *TransformComponent
	MaterialComponent  *MaterialComponent
}

func NewRenderableComponent(meshComponent *MeshComponent, bufferComponent *BufferComponent, transformComponent *TransformComponent, materialComponent *MaterialComponent) *RenderableComponent {

	return &RenderableComponent{
		MeshComponent:      meshComponent,
		BufferComponent:    bufferComponent,
		TransformComponent: transformComponent,
		MaterialComponent:  materialComponent,
	}
}
