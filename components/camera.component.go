package components

import (
	"github.com/go-gl/mathgl/mgl32"
)

type CameraComponent struct {
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

func NewCameraComponent(worldUp mgl32.Vec3, yaw, pitch, fov, aspect, nearClip, farClip float32) *CameraComponent {
	cam := &CameraComponent{
		WorldUp:     worldUp,
		Yaw:         yaw,
		Pitch:       pitch,
		FieldOfView: fov,
		AspectRatio: aspect,
		NearClip:    nearClip,
		FarClip:     farClip,
	}
	return cam
}

func (cam *CameraComponent) GetViewMatrix(camPosition mgl32.Vec3) mgl32.Mat4 {
	target := camPosition.Add(cam.Front)
	return mgl32.LookAtV(camPosition, target, cam.Up)
}

func (cam *CameraComponent) GetProjectionMatrix() mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(cam.FieldOfView), cam.AspectRatio, cam.NearClip, cam.FarClip)
}
