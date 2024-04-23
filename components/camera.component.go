package components

import "github.com/go-gl/mathgl/mgl32"

type CameraComponent struct {
	Position    mgl32.Vec3
	Target      mgl32.Vec3
	Up          mgl32.Vec3
	FieldOfView float32
	AspectRatio float32
	NearClip    float32
	FarClip     float32
}

func NewCameraComponent(position, target, up mgl32.Vec3, fov, aspect, nearClip, farClip float32) *CameraComponent {
	return &CameraComponent{
		Position:    position,
		Target:      target,
		Up:          up,
		FieldOfView: fov,
		AspectRatio: aspect,
		NearClip:    nearClip,
		FarClip:     farClip,
	}
}

func (cam *CameraComponent) GetViewMatrix() mgl32.Mat4 {
	return mgl32.LookAtV(cam.Position, cam.Target, cam.Up)
}

func (cam *CameraComponent) GetProjectionMatrix() mgl32.Mat4 {
	return mgl32.Perspective(mgl32.DegToRad(cam.FieldOfView), cam.AspectRatio, cam.NearClip, cam.FarClip)
}
