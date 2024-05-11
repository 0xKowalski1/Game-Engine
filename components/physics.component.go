package components

import "github.com/go-gl/mathgl/mgl32"

type PhysicsComponent struct {
	Velocity mgl32.Vec3
	Mass     float32
	Static   bool
}

func NewPhysicsComponent(velocity mgl32.Vec3, mass float32, static bool) *PhysicsComponent {
	return &PhysicsComponent{
		Velocity: velocity,
		Mass:     mass,
		Static:   static,
	}
}
