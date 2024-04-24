package components

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraComponent struct {
	Position    mgl32.Vec3
	Orientation mgl32.Quat
	FieldOfView float32
	AspectRatio float32
	NearClip    float32
	FarClip     float32
}

func NewCameraComponent(position mgl32.Vec3, orientation mgl32.Quat, fov, aspect, nearClip, farClip float32) *CameraComponent {
	if orientation.ApproxEqual(mgl32.Quat{}) {
		orientation = mgl32.QuatIdent()
	}
	return &CameraComponent{
		Position:    position,
		Orientation: orientation,
		FieldOfView: fov,
		AspectRatio: aspect,
		NearClip:    nearClip,
		FarClip:     farClip,
	}
}

func (cam *CameraComponent) Front() mgl32.Vec3 {
	return cam.Orientation.Rotate(mgl32.Vec3{0, 0, -1})
}

func (cam *CameraComponent) Right() mgl32.Vec3 {
	up := cam.Orientation.Rotate(mgl32.Vec3{0, 1, 0})
	return cam.Front().Cross(up).Normalize()
}

func (cam *CameraComponent) Move(direction mgl32.Vec3, amount float32) {
	moveVector := direction.Mul(amount)
	cam.Position = cam.Position.Add(moveVector)
}

func (cam *CameraComponent) Rotate(pitch, yaw float32) {
	pitchRad := mgl32.DegToRad(pitch)
	yawRad := mgl32.DegToRad(yaw)

	pitchQuat := angleAxis(-pitchRad, cam.Right())
	yawQuat := angleAxis(yawRad, mgl32.Vec3{0, 1, 0})

	cam.Orientation = yawQuat.Mul(cam.Orientation).Mul(pitchQuat)
}

func angleAxis(angle float32, axis mgl32.Vec3) mgl32.Quat {
	normalizedAxis := axis.Normalize()
	sinHalfAngle := float32(math.Sin(float64(angle / 2.0)))
	cosHalfAngle := float32(math.Cos(float64(angle / 2.0)))

	x := normalizedAxis.X() * sinHalfAngle
	y := normalizedAxis.Y() * sinHalfAngle
	z := normalizedAxis.Z() * sinHalfAngle
	w := cosHalfAngle

	return mgl32.Quat{W: w, V: mgl32.Vec3{x, y, z}}
}

func (cam *CameraComponent) GetViewMatrix() mgl32.Mat4 {
	front := cam.Orientation.Rotate(mgl32.Vec3{0, 0, -1})
	up := cam.Orientation.Rotate(mgl32.Vec3{0, 1, 0})
	return mgl32.LookAtV(cam.Position, cam.Position.Add(front), up)
}

func (cam *CameraComponent) GetProjectionMatrix() mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(cam.FieldOfView), cam.AspectRatio, cam.NearClip, cam.FarClip)
}
