package components

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type CameraComponent struct {
	Position    mgl32.Vec3
	Front       mgl32.Vec3
	Up          mgl32.Vec3
	Right       mgl32.Vec3
	WorldUp     mgl32.Vec3
	Yaw         float32
	Pitch       float32
	FieldOfView float32
	AspectRatio float32
	NearClip    float32
	FarClip     float32
}

func NewCameraComponent(position, worldUp mgl32.Vec3, yaw, pitch, fov, aspect, nearClip, farClip float32) *CameraComponent {
	cam := &CameraComponent{
		Position:    position,
		WorldUp:     worldUp,
		Yaw:         yaw,
		Pitch:       pitch,
		FieldOfView: fov,
		AspectRatio: aspect,
		NearClip:    nearClip,
		FarClip:     farClip,
	}
	cam.updateCameraVectors()
	return cam
}

func (cam *CameraComponent) Move(direction mgl32.Vec3, amount float32) {
	cam.Position = cam.Position.Add(direction.Mul(amount))
}

func (cam *CameraComponent) Rotate(yawIncr, pitchIncr float32) {
	cam.Yaw += yawIncr
	cam.Pitch -= pitchIncr

	// Limit pitch to prevent gimbal lock
	if cam.Pitch > 89.0 {
		cam.Pitch = 89.0
	} else if cam.Pitch < -89.0 {
		cam.Pitch = -89.0
	}

	cam.updateCameraVectors()
}

func (cam *CameraComponent) updateCameraVectors() {
	front := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(cam.Yaw))) * math.Cos(float64(mgl32.DegToRad(cam.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(cam.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(cam.Yaw))) * math.Cos(float64(mgl32.DegToRad(cam.Pitch)))),
	}.Normalize()

	cam.Front = front
	cam.Right = cam.Front.Cross(cam.WorldUp).Normalize()
	cam.Up = cam.Right.Cross(cam.Front).Normalize()
}

func (cam *CameraComponent) GetViewMatrix() mgl32.Mat4 {
	target := cam.Position.Add(cam.Front)
	return mgl32.LookAtV(cam.Position, target, cam.Up)
}

func (cam *CameraComponent) GetProjectionMatrix() mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(cam.FieldOfView), cam.AspectRatio, cam.NearClip, cam.FarClip)
}
