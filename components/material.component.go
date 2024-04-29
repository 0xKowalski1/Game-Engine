package components

import "github.com/go-gl/mathgl/mgl32"

type MaterialComponent struct {
	Ambient   mgl32.Vec3
	Diffuse   mgl32.Vec3
	Specular  mgl32.Vec3
	Shininess float32
}

func NewMaterialComponent(ambient, diffuse, specular mgl32.Vec3, shininess float32) *MaterialComponent {
	return &MaterialComponent{
		Ambient:   ambient,
		Diffuse:   diffuse,
		Specular:  specular,
		Shininess: shininess,
	}
}
