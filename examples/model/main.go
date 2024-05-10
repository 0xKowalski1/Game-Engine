package main

import (
	"0xKowalski/game/components"
	"0xKowalski/game/engine"
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Game struct {
	Engine     *engine.Engine
	CameraComp *components.CameraComponent
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

	// Camera
	cameraEntity := game.Engine.EntityStore.NewEntity()

	cameraComp := components.NewCameraComponent(
		mgl32.Vec3{0, 0, 5}, // Position: Initial position of the camera in the world
		mgl32.Vec3{0, 1, 0}, // WorldUp: The up vector of the world, typically Y-axis is up
		-90.0,               // Yaw: Initial yaw angle, facing forward along the Z-axis
		0.0,                 // Pitch: Initial pitch angle, looking straight at the horizon
		45.0,                // Field of view in degrees
		800.0/600.0,         // Aspect ratio: width divided by height of the viewport
		0.1,                 // Near clipping plane: the closest distance the camera can see
		100.0,               // Far clipping plane: the farthest distance the camera can see
	)
	game.Engine.EntityStore.AddComponent(cameraEntity, cameraComp)
	game.CameraComp = cameraComp

	// Model
	_ = game.Engine.EntityStore.NewModelEntity(mgl32.Vec3{0.0, 0.0, 0.0}, "assets/models/backpack/backpack.obj", "assets/models/backpack/backpack.mtl")

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

	// Key Inputs
	game.Engine.InputManager.RegisterKeyAction(glfw.KeyEscape, CloseApp, func() { game.Engine.Window.GlfwWindow.SetShouldClose(true) })

	currentFrame := glfw.GetTime()
	deltaTime := currentFrame - game.Engine.LastFrame
	cameraSpeed := float32(1 * deltaTime)

	game.Engine.InputManager.RegisterKeyAction(glfw.KeyW, MoveForward, func() {
		cameraComp.Move(cameraComp.Front, cameraSpeed)
	})
	game.Engine.InputManager.RegisterKeyAction(glfw.KeyS, MoveBackward, func() {
		cameraComp.Move(cameraComp.Front.Mul(-1), cameraSpeed)
	})
	game.Engine.InputManager.RegisterKeyAction(glfw.KeyD, StrafeRight, func() {
		cameraComp.Move(cameraComp.Right, cameraSpeed)

	})
	game.Engine.InputManager.RegisterKeyAction(glfw.KeyA, StrafeLeft, func() {
		cameraComp.Move(cameraComp.Right.Mul(-1), cameraSpeed)

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

		cameraComp.Rotate(xOffset*0.05, yOffset*0.05)
		game.Engine.InputManager.LastX = xpos
		game.Engine.InputManager.LastY = ypos
	})

	// Loop
	game.Engine.Run(game.MainLoop)
}
