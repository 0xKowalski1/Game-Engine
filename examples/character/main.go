package main

import (
	"0xKowalski/game/components"
	"0xKowalski/game/engine"
	"0xKowalski/game/entities"
	"log"
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Player struct {
	Entity             *entities.Entity
	TransformComponent *components.TransformComponent
	CameraComponent    *components.CameraComponent
	PhysicsComponent   *components.PhysicsComponent
}

func (player *Player) Rotate(yawIncr, pitchIncr float32) {
	cam := player.CameraComponent

	cam.Yaw += yawIncr
	cam.Pitch -= pitchIncr

	// Limit pitch to prevent gimbal lock
	if cam.Pitch > 89.0 {
		cam.Pitch = 89.0
	} else if cam.Pitch < -89.0 {
		cam.Pitch = -89.0
	}

	player.updateCameraVectors()
}

func (player *Player) updateCameraVectors() {
	cam := player.CameraComponent

	front := mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(cam.Yaw))) * math.Cos(float64(mgl32.DegToRad(cam.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(cam.Pitch)))),
		float32(math.Sin(float64(mgl32.DegToRad(cam.Yaw))) * math.Cos(float64(mgl32.DegToRad(cam.Pitch)))),
	}.Normalize()

	cam.Front = front
	cam.Right = cam.Front.Cross(cam.WorldUp).Normalize()
	cam.Up = cam.Right.Cross(cam.Front).Normalize()
}

type Game struct {
	Engine *engine.Engine
	Player *Player
}

func (g *Game) MainLoop() {

}

func main() {

	game := Game{}

	eng, err := engine.InitEngine()
	if err != nil {
		log.Printf("Error starting engine: %v", err)
		panic(err)
	}

	game.Engine = eng

	// Player

	entity := game.Engine.EntityStore.NewEntity()

	transformComponent := components.NewTransformComponent(mgl32.Vec3{0, 3, 0})
	game.Engine.EntityStore.AddComponent(entity, transformComponent)

	cameraComp := components.NewCameraComponent(
		mgl32.Vec3{0, 1, 0}, // WorldUp: The up vector of the world, typically Y-axis is up
		-90.0,               // Yaw: Initial yaw angle, facing forward along the Z-axis
		0.0,                 // Pitch: Initial pitch angle, looking straight at the horizon
		45.0,                // Field of view in degrees
		800.0/600.0,         // Aspect ratio: width divided by height of the viewport
		0.1,                 // Near clipping plane: the closest distance the camera can see
		100.0,               // Far clipping plane: the farthest distance the camera can see
	)
	game.Engine.EntityStore.AddComponent(entity, cameraComp)

	physicsComponent := components.NewPhysicsComponent(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 0, 0}, 1, mgl32.Vec3{1, 1, 1}, false)
	game.Engine.EntityStore.AddComponent(entity, physicsComponent)

	game.Player = &Player{
		Entity:             &entity,
		TransformComponent: transformComponent,
		CameraComponent:    cameraComp,
		PhysicsComponent:   physicsComponent,
	}

	// Floor
	floorEntity := game.Engine.EntityStore.NewPlaneEntity(mgl32.Vec3{0, 0, 0})

	floorPhysicsComponent := components.NewPhysicsComponent(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0.0, 0.0, 0.0}, 1, mgl32.Vec3{25, 1, 25}, true)
	game.Engine.EntityStore.AddComponent(*floorEntity, floorPhysicsComponent)

	// LIGHTING

	// Ambient
	ambientLightEntity := game.Engine.EntityStore.NewEntity()
	ambientLightComponent := components.NewAmbientLightComponent(mgl32.Vec3{1.0, 1.0, 1.0}, 0.1)
	game.Engine.EntityStore.AddComponent(ambientLightEntity, ambientLightComponent)

	// Directional
	directionalLightEntity := game.Engine.EntityStore.NewEntity()
	directionalLightComponent := components.NewDirectionalLightComponent(mgl32.Vec3{-0.2, -1.0, -0.3}, mgl32.Vec3{1.0, 1.0, 1.0}, 1)
	game.Engine.EntityStore.AddComponent(directionalLightEntity, directionalLightComponent)

	// END LIGHTING

	// Inputs
	const (
		CloseApp = iota

		MoveForward
		MoveBackward
		StrafeRight
		StrafeLeft
	)

	currentFrame := glfw.GetTime()
	deltaTime := currentFrame - game.Engine.LastFrame
	game.Engine.LastFrame = currentFrame
	playerSpeed := float32(3 * deltaTime) // Define the player's movement s

	// Key Inputs
	game.Engine.InputManager.RegisterKeyAction(glfw.KeyEscape, CloseApp, func() { game.Engine.Window.GlfwWindow.SetShouldClose(true) })

	game.Engine.InputManager.RegisterKeyAction(glfw.KeyW, MoveForward, func() {
		// Calculate the new forward movement only in the X/Z plane
		newPosition := game.Player.TransformComponent.Position.Add(mgl32.Vec3{game.Player.CameraComponent.Front.X(), 0, game.Player.CameraComponent.Front.Z()}.Mul(playerSpeed))
		game.Player.TransformComponent.Position = newPosition
	})

	game.Engine.InputManager.RegisterKeyAction(glfw.KeyS, MoveBackward, func() {
		// Calculate the new backward movement only in the X/Z plane
		newPosition := game.Player.TransformComponent.Position.Sub(mgl32.Vec3{game.Player.CameraComponent.Front.X(), 0, game.Player.CameraComponent.Front.Z()}.Mul(playerSpeed))
		game.Player.TransformComponent.Position = newPosition
	})

	game.Engine.InputManager.RegisterKeyAction(glfw.KeyD, StrafeRight, func() {
		// Calculate the new right strafe movement only in the X/Z plane
		newPosition := game.Player.TransformComponent.Position.Add(mgl32.Vec3{game.Player.CameraComponent.Right.X(), 0, game.Player.CameraComponent.Right.Z()}.Mul(playerSpeed))
		game.Player.TransformComponent.Position = newPosition
	})

	game.Engine.InputManager.RegisterKeyAction(glfw.KeyA, StrafeLeft, func() {
		// Calculate the new left strafe movement only in the X/Z plane
		newPosition := game.Player.TransformComponent.Position.Sub(mgl32.Vec3{game.Player.CameraComponent.Right.X(), 0, game.Player.CameraComponent.Right.Z()}.Mul(playerSpeed))
		game.Player.TransformComponent.Position = newPosition
	})

	// Mouse Inputs
	game.Engine.InputManager.RegisterMouseMoveHandler(func(xpos, ypos float64) {
		if game.Engine.InputManager.FirstMouse { // Handle the initial jump in mouse position
			game.Engine.InputManager.LastX = xpos
			game.Engine.InputManager.LastY = ypos
			game.Engine.InputManager.FirstMouse = false
		}

		xOffset := float32(xpos - game.Engine.InputManager.LastX)
		yOffset := float32(ypos - game.Engine.InputManager.LastY)

		game.Player.Rotate(xOffset*0.05, yOffset*0.05)
		game.Engine.InputManager.LastX = xpos
		game.Engine.InputManager.LastY = ypos
	})

	game.Engine.InputManager.RegisterMouseScrollHandler(func(xoffset, yoffset float64) {
		fov := cameraComp.FieldOfView
		cameraComp.FieldOfView = fov - float32(yoffset)

		if fov < 1.0 {
			cameraComp.FieldOfView = 1.0
		} else if fov > 45.0 {
			cameraComp.FieldOfView = 45.0
		}
	})

	// Loop
	game.Engine.Run(game.MainLoop)

}
