package systems

import (
	"0xKowalski/game/components"
	"0xKowalski/game/entities"
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

type PhysicsSystem struct {
	EntityStore *entities.EntityStore
	Gravity     mgl32.Vec3
}

func NewPhysicsSystem(entityStore *entities.EntityStore, gravity mgl32.Vec3) *PhysicsSystem {
	return &PhysicsSystem{
		EntityStore: entityStore,
		Gravity:     gravity,
	}
}

func (ps *PhysicsSystem) Update(dt float32) {
	entities := ps.EntityStore.GetEntitiesWithComponentType(&components.PhysicsComponent{})

	for _, entity := range entities {
		physicsComponent, physicsComponentOk := ps.EntityStore.GetComponent(entity, &components.PhysicsComponent{}).(*components.PhysicsComponent)
		if !physicsComponentOk {
			log.Println("Error converting physics component interface into component")
		}

		if !physicsComponent.Static {
			// Apply Gravity
			physicsComponent.Velocity = physicsComponent.Velocity.Add(ps.Gravity.Mul(dt))
			//Update Position
			// Get transform component
			transformComponent, transformComponentOk := ps.EntityStore.GetComponent(entity, &components.TransformComponent{}).(*components.TransformComponent)
			if !transformComponentOk {
				log.Println("Error converting transform component interface into component")
			}
			transformComponent.Position = transformComponent.Position.Add(physicsComponent.Velocity.Mul(dt))
		}
	}
}
