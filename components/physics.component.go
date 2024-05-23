package components

import "github.com/go-gl/mathgl/mgl32"

type PhysicsComponent struct {
	// Rigid Body
	Velocity     mgl32.Vec3
	Acceleration mgl32.Vec3
	Mass         float32

	// Collision Detection
	BoundingBoxHalfSize mgl32.Vec3
	Grounded            bool

	Static bool
}

// Angular velocity & intertia tensor for rotation?
// Bounding volume
// Collision Mesh
// Restitution
// Friction

func NewPhysicsComponent(velocity mgl32.Vec3, acceleration mgl32.Vec3, mass float32, boundingBoxHalfsize mgl32.Vec3, static bool) *PhysicsComponent {
	return &PhysicsComponent{
		Velocity:            velocity,
		Acceleration:        acceleration,
		Mass:                mass,
		BoundingBoxHalfSize: boundingBoxHalfsize,
		Static:              static,
		Grounded:            false,
	}
}
