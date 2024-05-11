package entities

import (
	"0xKowalski/game/components"
	"math"

	"github.com/go-gl/mathgl/mgl32"
)

type Freecam struct {
	Entity             *Entity
	TransformComponent *components.TransformComponent
	CameraComponent    *components.CameraComponent
}

func (es *EntityStore) NewFreecamEntity(position mgl32.Vec3) *Freecam {
	entity := es.NewEntity()

	transformComponent := components.NewTransformComponent(position)
	es.AddComponent(entity, transformComponent)

	cameraComponent := components.NewCameraComponent(
		mgl32.Vec3{0, 1, 0}, // WorldUp: The up vector of the world, typically Y-axis is up
		-90.0,               // Yaw: Initial yaw angle, facing forward along the Z-axis
		0.0,                 // Pitch: Initial pitch angle, looking straight at the horizon
		45.0,                // Field of view in degrees
		800.0/600.0,         // Aspect ratio: width divided by height of the viewport
		0.1,                 // Near clipping plane: the closest distance the camera can see
		100.0,               // Far clipping plane: the farthest distance the camera can see
	)
	es.AddComponent(entity, cameraComponent)

	freecam := &Freecam{
		Entity:             &entity,
		TransformComponent: transformComponent,
		CameraComponent:    cameraComponent,
	}

	freecam.updateCameraVectors()

	return freecam
}

func (camEntity *Freecam) Move(direction mgl32.Vec3, amount float32) {
	cam := camEntity.TransformComponent
	cam.Position = cam.Position.Add(direction.Mul(amount))
}

func (camEntity *Freecam) Rotate(yawIncr, pitchIncr float32) {
	cam := camEntity.CameraComponent

	cam.Yaw += yawIncr
	cam.Pitch -= pitchIncr

	// Limit pitch to prevent gimbal lock
	if cam.Pitch > 89.0 {
		cam.Pitch = 89.0
	} else if cam.Pitch < -89.0 {
		cam.Pitch = -89.0
	}

	camEntity.updateCameraVectors()
}

func (camEntity *Freecam) updateCameraVectors() {
	cam := camEntity.CameraComponent

	front := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(cam.Yaw))) * math.Cos(float64(mgl32.DegToRad(cam.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(cam.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(cam.Yaw))) * math.Cos(float64(mgl32.DegToRad(cam.Pitch)))),
	}.Normalize()

	cam.Front = front
	cam.Right = cam.Front.Cross(cam.WorldUp).Normalize()
	cam.Up = cam.Right.Cross(cam.Front).Normalize()
}
