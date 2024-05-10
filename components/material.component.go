package components

type MaterialComponent struct {
	DiffuseMap  string
	SpecularMap string
	Shininess   float32
}

func NewMaterialComponent(diffuseMapPath string, specularMapPath string, shininess float32) *MaterialComponent {

	return &MaterialComponent{
		DiffuseMap:  diffuseMapPath,
		SpecularMap: specularMapPath,
		Shininess:   shininess,
	}
}
