package components

import (
	"github.com/go-gl/mathgl/mgl32"
)

type TransformComponent struct {
	Position mgl32.Vec3
	Rotation mgl32.Quat
	Scale    mgl32.Vec3
}

func NewTransformComponent(position mgl32.Vec3) *TransformComponent {
	return &TransformComponent{
		Position: position,
		Rotation: mgl32.QuatIdent(),
		Scale:    mgl32.Vec3{1, 1, 1},
	}
}

func (t *TransformComponent) SetPosition(x, y, z float32) {
	t.Position = mgl32.Vec3{x, y, z}
}

func (t *TransformComponent) SetRotation(quat mgl32.Quat) {
	t.Rotation = quat
}

func (t *TransformComponent) SetScale(x, y, z float32) {
	t.Scale = mgl32.Vec3{x, y, z}
}

func (t *TransformComponent) GetModelMatrix() mgl32.Mat4 {
	translateMat := mgl32.Translate3D(t.Position.X(), t.Position.Y(), t.Position.Z())
	scaleMat := mgl32.Scale3D(t.Scale.X(), t.Scale.Y(), t.Scale.Z())
	rotateMat := t.Rotation.Mat4()

	return translateMat.Mul4(rotateMat).Mul4(scaleMat)
}
