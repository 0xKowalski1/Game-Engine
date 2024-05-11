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
			continue // Skip this entity if its physics component is not accessible
		}

		if !physicsComponent.Static {
			physicsComponent.Velocity = physicsComponent.Velocity.Add(ps.Gravity.Mul(dt))

			transformComponent, transformComponentOk := ps.EntityStore.GetComponent(entity, &components.TransformComponent{}).(*components.TransformComponent)
			if !transformComponentOk {
				log.Println("Error converting transform component interface into component")
				continue
			}

			newPosition := transformComponent.Position.Add(physicsComponent.Velocity.Mul(dt))

			hasCollided := false
			for _, otherEntity := range entities {
				if entity.ID == otherEntity.ID {
					continue
				}

				otherTransformComponent, otherTransformComponentOk := ps.EntityStore.GetComponent(otherEntity, &components.TransformComponent{}).(*components.TransformComponent)
				if !otherTransformComponentOk {
					log.Println("Error converting other transform component interface into component")

					continue
				}

				collides, collisionNormal := ps.DetectCollision(transformComponent, otherTransformComponent)
				if collides {
					physicsComponent.Velocity = ps.StopVelocityAlongNormal(physicsComponent.Velocity, collisionNormal)
					hasCollided = true
				}
			}

			if !hasCollided {
				transformComponent.Position = newPosition
			}
		}
	}
}

func (ps *PhysicsSystem) DetectCollision(transformA, transformB *components.TransformComponent) (bool, mgl32.Vec3) {
	// Simple AABB collision detection
	aMin := transformA.Position.Sub(mgl32.Vec3{0.5, 0.5, 0.5})
	aMax := transformA.Position.Add(mgl32.Vec3{0.5, 0.5, 0.5})
	bMin := transformB.Position.Sub(mgl32.Vec3{0.5, 0.5, 0.5})
	bMax := transformB.Position.Add(mgl32.Vec3{0.5, 0.5, 0.5})

	if aMax.X() < bMin.X() || aMin.X() > bMax.X() ||
		aMax.Y() < bMin.Y() || aMin.Y() > bMax.Y() ||
		aMax.Z() < bMin.Z() || aMin.Z() > bMax.Z() {
		return false, mgl32.Vec3{0, 0, 0}
	}

	collisionNormal := transformA.Position.Sub(transformB.Position).Normalize()
	return true, collisionNormal
}

func (ps *PhysicsSystem) StopVelocityAlongNormal(velocity, collisionNormal mgl32.Vec3) mgl32.Vec3 {
	velocityProjection := velocity.Dot(collisionNormal)
	return velocity.Sub(collisionNormal.Mul(velocityProjection))
}
