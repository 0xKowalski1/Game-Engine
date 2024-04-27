package components

import "github.com/go-gl/mathgl/mgl32"

type DirectionalLightComponent struct {
	Direction mgl32.Vec3
	Color     mgl32.Vec3
	Intensity float32
}

func NewDirectionalLightComponent(direction mgl32.Vec3, color mgl32.Vec3, intensity float32) *DirectionalLightComponent {
	return &DirectionalLightComponent{
		Direction: direction,
		Color:     color,
		Intensity: intensity,
	}
}
