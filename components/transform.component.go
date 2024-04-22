package components

import (
	"github.com/go-gl/mathgl/mgl32"
)

type TransformComponent struct {
	Position mgl32.Vec3
	Rotation mgl32.Quat
	Scale    mgl32.Vec3
}

// NewTransformComponent creates a new Transform component with default values.
// Position at origin, no rotation, and uniform scale of 1.
func NewTransformComponent() *TransformComponent {
	return &TransformComponent{
		Position: mgl32.Vec3{0, 0, 0},
		Rotation: mgl32.QuatIdent(),
		Scale:    mgl32.Vec3{1, 1, 1},
	}
}

// SetPosition updates the position vector of the transform.
func (t *TransformComponent) SetPosition(x, y, z float32) {
	t.Position = mgl32.Vec3{x, y, z}
}

// SetRotation updates the rotation quaternion of the transform.
func (t *TransformComponent) SetRotation(quat mgl32.Quat) {
	t.Rotation = quat
}

// SetScale updates the scale vector of the transform.
func (t *TransformComponent) SetScale(x, y, z float32) {
	t.Scale = mgl32.Vec3{x, y, z}
}

// GetModelMatrix calculates and returns the model matrix for this transform.
func (t *TransformComponent) GetModelMatrix() mgl32.Mat4 {
	translateMat := mgl32.Translate3D(t.Position.X(), t.Position.Y(), t.Position.Z())
	scaleMat := mgl32.Scale3D(t.Scale.X(), t.Scale.Y(), t.Scale.Z())
	rotateMat := t.Rotation.Mat4()

	// Model matrix is T * R * S
	return translateMat.Mul4(rotateMat).Mul4(scaleMat)
}
