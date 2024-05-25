package components

import (
	"github.com/go-gl/mathgl/mgl32"
)

type BoxColliderComponent struct {
	Center      mgl32.Vec3
	Size        mgl32.Vec3
	Friction    float32
	Restitution float32
}

func NewBoxColliderComponent(center, size mgl32.Vec3, friction, restitution float32) *BoxColliderComponent {
	return &BoxColliderComponent{
		Center:      center,
		Size:        size,
		Friction:    friction,
		Restitution: restitution,
	}
}
