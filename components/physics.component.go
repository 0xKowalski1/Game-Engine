package components

import "github.com/go-gl/mathgl/mgl32"

type PhysicsComponent struct {
	Velocity        mgl32.Vec3
	Acceleration    mgl32.Vec3
	AngularVelocity mgl32.Vec3
	InertiaTensor   mgl32.Vec3
	Mass            float32

	Static bool
}

func NewPhysicsComponent(velocity, acceleration, angularVelocity, inertiaTensor mgl32.Vec3, mass float32, static bool) *PhysicsComponent {
	return &PhysicsComponent{
		Velocity:        velocity,
		Acceleration:    acceleration,
		AngularVelocity: angularVelocity,
		InertiaTensor:   inertiaTensor,
		Mass:            mass,
		Static:          static,
	}
}
