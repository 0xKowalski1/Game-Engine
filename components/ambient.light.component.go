package components

import "github.com/go-gl/mathgl/mgl32"

type AmbientLightComponent struct {
	Color     mgl32.Vec3
	Intensity float32
}

func NewAmbientLightComponent(color mgl32.Vec3, intensity float32) *AmbientLightComponent {
	return &AmbientLightComponent{
		Color:     color,
		Intensity: intensity,
	}
}
