package components

import "github.com/go-gl/mathgl/mgl32"

type PointLightComponent struct {
	Position  mgl32.Vec3
	Color     mgl32.Vec3
	Intensity float32
	Constant  float32
	Linear    float32
	Quadratic float32
}

func NewPointLightComponent(position, color mgl32.Vec3, intensity, constant, linear, quadratic float32) *PointLightComponent {
	return &PointLightComponent{
		Position:  position,
		Color:     color,
		Intensity: intensity,
		Constant:  constant,
		Linear:    linear,
		Quadratic: quadratic,
	}
}
