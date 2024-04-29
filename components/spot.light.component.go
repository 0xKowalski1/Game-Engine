package components

import "github.com/go-gl/mathgl/mgl32"

type SpotLightComponent struct {
	Position    mgl32.Vec3
	Color       mgl32.Vec3
	Direction   mgl32.Vec3
	CutOff      float32
	OuterCutOff float32
	Intensity   float32
	Constant    float32
	Linear      float32
	Quadratic   float32
}

func NewSpotLightComponent(position, color, direction mgl32.Vec3, cutOff, outerCutOff, intensity, constant, linear, quadratic float32) *SpotLightComponent {
	return &SpotLightComponent{
		Position:    position,
		Color:       color,
		Direction:   direction,
		CutOff:      cutOff,
		OuterCutOff: outerCutOff,
		Intensity:   intensity,
		Constant:    constant,
		Linear:      linear,
		Quadratic:   quadratic,
	}
}
