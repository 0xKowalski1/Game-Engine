package components

type RenderableComponent struct {
	MeshComponent      *MeshComponent
	BufferComponent    *BufferComponent
	TransformComponent *TransformComponent
	TextureComponent   *TextureComponent
	MaterialComponent  *MaterialComponent
}

func NewRenderableComponent(meshComponent *MeshComponent, bufferComponent *BufferComponent, transformComponent *TransformComponent, textureComponent *TextureComponent, materialComponent *MaterialComponent) *RenderableComponent {

	return &RenderableComponent{
		MeshComponent:      meshComponent,
		BufferComponent:    bufferComponent,
		TransformComponent: transformComponent,
		TextureComponent:   textureComponent,
		MaterialComponent:  materialComponent,
	}
}
