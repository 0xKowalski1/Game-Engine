package systems

import (
	"0xKowalski/game/components"
	"0xKowalski/game/entities"
	"log"
	"math"

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
			continue
		}

		transformComponent, transformComponentOk := ps.EntityStore.GetComponent(entity, &components.TransformComponent{}).(*components.TransformComponent)
		if !transformComponentOk {
			log.Println("Error converting transform component interface into component")
			continue
		}

		if !physicsComponent.Static {
			// Handle Gravity
			if !physicsComponent.Grounded {
				force := ps.Gravity.Mul(physicsComponent.Mass)
				physicsComponent.Acceleration = force.Mul(1 / physicsComponent.Mass) // a = F/m

				//Update velocity based on acceleration
				physicsComponent.Velocity = physicsComponent.Velocity.Add(physicsComponent.Acceleration.Mul(dt))

				// Update position based on velocity
				displacement := physicsComponent.Velocity.Mul(dt)
				transformComponent.Position = transformComponent.Position.Add(displacement)
			}

			// Handle Collisions
			for _, entityToCheck := range entities {
				if entityToCheck.ID == entity.ID {
					continue
				}

				physicsComponentToCheck, physicsComponentToCheckOk := ps.EntityStore.GetComponent(entityToCheck, &components.PhysicsComponent{}).(*components.PhysicsComponent)
				if !physicsComponentToCheckOk {
					log.Println("Error converting physics component interface into component")
					continue
				}

				transformComponentToCheck, transformComponentToCheckOk := ps.EntityStore.GetComponent(entityToCheck, &components.TransformComponent{}).(*components.TransformComponent)
				if !transformComponentToCheckOk {
					log.Println("Error converting transform component interface into component")
					continue
				}

				// Check for early stage collision through bounding volumes
				if checkEarlyStageCollision(transformComponent, transformComponentToCheck, physicsComponent, physicsComponentToCheck) {
					resolveCollision(transformComponent, transformComponentToCheck, physicsComponent, physicsComponentToCheck)
					physicsComponent.Grounded = true
				} else {
					physicsComponent.Grounded = false
				}
			}
		}
	}
}

func checkEarlyStageCollision(transformComponent, transformComponentToCheck *components.TransformComponent, physicsComponent, physicsComponentToCheck *components.PhysicsComponent) bool {
	minA := transformComponent.Position.Sub(physicsComponent.BoundingBoxHalfSize)
	maxA := transformComponent.Position.Add(physicsComponent.BoundingBoxHalfSize)

	minB := transformComponentToCheck.Position.Sub(physicsComponentToCheck.BoundingBoxHalfSize)
	maxB := transformComponentToCheck.Position.Add(physicsComponentToCheck.BoundingBoxHalfSize)

	return (minA.X() <= maxB.X() && maxA.X() >= minB.X()) &&
		(minA.Y() <= maxB.Y() && maxA.Y() >= minB.Y()) &&
		(minA.Z() <= maxB.Z() && maxA.Z() >= minB.Z())
}

func resolveCollision(transformA, transformB *components.TransformComponent, physicsA, physicsB *components.PhysicsComponent) {
	// Example for perfectly elastic collisions and equal mass, without rotation
	normal := transformB.Position.Sub(transformA.Position).Normalize()
	relativeVelocity := physicsB.Velocity.Sub(physicsA.Velocity)

	// Calculate the velocity along the normal due to collision
	velocityAlongNormal := relativeVelocity.Dot(normal)
	log.Println(velocityAlongNormal)
	if velocityAlongNormal > 0 {
		return // They are moving away from each other
	}

	// Calculate restitution (elasticity) coefficient, simple example with 1 for perfectly elastic
	restitution := float32(1.0)
	impulseScalar := -(1 + restitution) * velocityAlongNormal
	impulseScalar /= (1 / physicsA.Mass) + (1 / physicsB.Mass)

	// Apply impulse to both entities
	impulse := normal.Mul(impulseScalar)
	physicsA.Velocity = physicsA.Velocity.Sub(impulse.Mul(1 / physicsA.Mass))
	physicsB.Velocity = physicsB.Velocity.Add(impulse.Mul(1 / physicsB.Mass))

	// Positional correction to avoid sinking due to floating point precision errors
	percent := float32(0.2) // usually 20% to 80%
	slop := float32(0.01)   // usually 0.01 to 0.1
	correctionMagnitude := float32(math.Max(float64(transformA.Position.Sub(transformB.Position).Len()-slop)/float64((1/physicsA.Mass)+(1/physicsB.Mass)), 0.0))
	correction := normal.Mul(correctionMagnitude * percent)
	transformA.Position = transformA.Position.Sub(correction.Mul(1 / physicsA.Mass))
	transformB.Position = transformB.Position.Add(correction.Mul(1 / physicsB.Mass))
}
