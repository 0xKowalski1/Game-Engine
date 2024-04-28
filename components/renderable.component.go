package components

type RenderableComponent struct {
	MeshComponent      *MeshComponent
	BufferComponent    *BufferComponent
	TransformComponent *TransformComponent
	TextureComponent   *TextureComponent
}

func NewRenderableComponent(meshComponent *MeshComponent, bufferComponent *BufferComponent, transformComponent *TransformComponent, textureComponent *TextureComponent) *RenderableComponent {

	return &RenderableComponent{
		MeshComponent:      meshComponent,
		BufferComponent:    bufferComponent,
		TransformComponent: transformComponent,
		TextureComponent:   textureComponent,
	}
}
