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
	dt          float32
}

func NewPhysicsSystem(entityStore *entities.EntityStore, gravity mgl32.Vec3) *PhysicsSystem {
	return &PhysicsSystem{
		EntityStore: entityStore,
		Gravity:     gravity,
	}
}

func (ps *PhysicsSystem) Update(dt float32) {
	ps.dt = dt

	entities := ps.EntityStore.GetEntitiesWithComponentType(&components.PhysicsComponent{})

	for _, entity := range entities {
		physicsComponent, physicsComponentOk := ps.EntityStore.GetComponent(entity, &components.PhysicsComponent{}).(*components.PhysicsComponent)
		if !physicsComponentOk {
			log.Println("Error converting physics component interface into component")
			continue
		}

		if physicsComponent.Static {
			continue
		}

		transformComponent, transformComponentOk := ps.EntityStore.GetComponent(entity, &components.TransformComponent{}).(*components.TransformComponent)
		if !transformComponentOk {
			log.Println("Error converting transform component interface into component")
			continue
		}

		ps.updateForces(physicsComponent)
		ps.applyForces(physicsComponent, transformComponent)
		ps.handleCollisions(entity)
	}
}

func (ps *PhysicsSystem) updateForces(physicsComponent *components.PhysicsComponent) {
	// Apply gravity
	force := ps.Gravity.Mul(physicsComponent.Mass)
	physicsComponent.Acceleration = force.Mul(1 / physicsComponent.Mass) // a = F/m

}

func (ps *PhysicsSystem) applyForces(physicsComponent *components.PhysicsComponent, transformComponent *components.TransformComponent) {
	//Update velocity based on acceleration
	physicsComponent.Velocity = physicsComponent.Velocity.Add(physicsComponent.Acceleration.Mul(ps.dt))

	// Update position based on velocity
	displacement := physicsComponent.Velocity.Mul(ps.dt)
	transformComponent.Position = transformComponent.Position.Add(displacement)
}

func (ps *PhysicsSystem) handleCollisions(entity entities.Entity) {
	entities := ps.EntityStore.GetEntitiesWithComponentType(&components.PhysicsComponent{})

	for _, entityToCheck := range entities {
		if entityToCheck.ID == entity.ID {
			continue
		}

		broadPhaseCollision := ps.broadPhaseCollisionCheck(entity, entityToCheck)
		log.Println(broadPhaseCollision)
		if !broadPhaseCollision {
			continue
		}

		narrowPhaseCollision := ps.narrowPhaseCollisionCheck(entity, entityToCheck)

		if !narrowPhaseCollision {
			continue
		}

		ps.resolveCollision(entity, entityToCheck)
	}
}

func (ps *PhysicsSystem) broadPhaseCollisionCheck(entity, entityToCheck entities.Entity) bool {
	boxCollider1, ok1 := ps.EntityStore.GetComponent(entity, &components.BoxColliderComponent{}).(*components.BoxColliderComponent)
	if !ok1 {
		log.Println("Error converting box collider component interface into component for the first entity")
		return false
	}
	transform1, ok1 := ps.EntityStore.GetComponent(entity, &components.TransformComponent{}).(*components.TransformComponent)
	if !ok1 {
		log.Println("Error converting transform component interface into component for the first entity")
		return false
	}

	boxCollider2, ok2 := ps.EntityStore.GetComponent(entityToCheck, &components.BoxColliderComponent{}).(*components.BoxColliderComponent)
	if !ok2 {
		log.Println("Error converting box collider component interface into component for the second entity")
		return false
	}
	transform2, ok2 := ps.EntityStore.GetComponent(entityToCheck, &components.TransformComponent{}).(*components.TransformComponent)
	if !ok2 {
		log.Println("Error converting transform component interface into component for the second entity")
		return false
	}

	min1 := transform1.Position.Add(boxCollider1.Center.Sub(boxCollider1.Size.Mul(0.5)))
	max1 := transform1.Position.Add(boxCollider1.Center.Add(boxCollider1.Size.Mul(0.5)))

	min2 := transform2.Position.Add(boxCollider2.Center.Sub(boxCollider2.Size.Mul(0.5)))
	max2 := transform2.Position.Add(boxCollider2.Center.Add(boxCollider2.Size.Mul(0.5)))

	return min1[0] <= max2[0] && max1[0] >= min2[0] &&
		min1[1] <= max2[1] && max1[1] >= min2[1] &&
		min1[2] <= max2[2] && max1[2] >= min2[2]
}

func (ps *PhysicsSystem) narrowPhaseCollisionCheck(entity, entityToCheck entities.Entity) bool {
	return true
}

func (ps *PhysicsSystem) resolveCollision(entity, entityToCheck entities.Entity) {
	physics1, ok1 := ps.EntityStore.GetComponent(entity, &components.PhysicsComponent{}).(*components.PhysicsComponent)
	transform1, okT1 := ps.EntityStore.GetComponent(entity, &components.TransformComponent{}).(*components.TransformComponent)
	boxCollider1, okB1 := ps.EntityStore.GetComponent(entity, &components.BoxColliderComponent{}).(*components.BoxColliderComponent)
	if !ok1 || !okT1 || !okB1 {
		log.Println("Error retrieving components for the first entity")
		return
	}

	physics2, ok2 := ps.EntityStore.GetComponent(entityToCheck, &components.PhysicsComponent{}).(*components.PhysicsComponent)
	transform2, okT2 := ps.EntityStore.GetComponent(entityToCheck, &components.TransformComponent{}).(*components.TransformComponent)
	boxCollider2, okB2 := ps.EntityStore.GetComponent(entityToCheck, &components.BoxColliderComponent{}).(*components.BoxColliderComponent)
	if !ok2 || !okT2 || !okB2 {
		log.Println("Error retrieving components for the second entity")
		return
	}

	center1 := transform1.Position.Add(boxCollider1.Center)
	center2 := transform2.Position.Add(boxCollider2.Center)
	distance := center1.Sub(center2)
	displacement := distance.Normalize().Mul(0.5)

	// Update positions to resolve penetration
	transform1.Position = transform1.Position.Add(displacement)
	if !physics2.Static {
		transform2.Position = transform2.Position.Sub(displacement)
	}
	// Update velocities based on a very basic elastic collision response
	v1 := physics1.Velocity
	v2 := physics2.Velocity
	m1 := physics1.Mass
	m2 := physics2.Mass

	// New velocities based on perfect elastic collision
	newV1 := v1.Sub(v2).Mul(2 * m2 / (m1 + m2)).Add(v2)
	newV2 := v2.Sub(v1).Mul(2 * m1 / (m1 + m2)).Add(v1)

	physics1.Velocity = newV1
	if !physics2.Static {
		physics2.Velocity = newV2
	}
}
