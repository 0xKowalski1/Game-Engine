package components

type RenderableComponent struct {
	TransformComponent *TransformComponent
	ModelComponent     *ModelComponent
}

func NewRenderableComponent(transformComponent *TransformComponent, modelComponent *ModelComponent) *RenderableComponent {

	return &RenderableComponent{
		TransformComponent: transformComponent,
		ModelComponent:     modelComponent,
	}
}
