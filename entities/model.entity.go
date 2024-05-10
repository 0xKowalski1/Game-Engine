package entities

import (
	"0xKowalski/game/components"
	"github.com/go-gl/mathgl/mgl32"
)

type ModelOption func(*EntityStore, *Entity)

func (es *EntityStore) NewModelEntity(position mgl32.Vec3, objPath string, mtlPath string, opts ...ModelOption) *Entity {
	entity := es.NewEntity()

	modelComponent := components.NewModelComponent(objPath, mtlPath)
	es.AddComponent(entity, modelComponent)

	transformComponent := components.NewTransformComponent(position)
	es.AddComponent(entity, transformComponent)

	renderComponent := components.NewRenderableComponent(transformComponent, modelComponent)
	es.AddComponent(entity, renderComponent)

	return &entity
}
