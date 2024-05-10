package main

import (
	"0xKowalski/game/components"
	"0xKowalski/game/engine"
	"log"
	"math"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Game struct {
	Engine        *engine.Engine
	CameraComp    *components.CameraComponent
	SpotLightComp *components.SpotLightComponent
}

func (g *Game) MainLoop() {
	/*
		lightComp := g.Engine.EntityStore.GetAllComponents(&components.PointLightComponent{})[0].(*components.PointLightComponent)
		lightComp.Color = mgl32.Vec3{
			float32(math.Sin(glfw.GetTime() * 2.0)),
			float32(math.Sin(glfw.GetTime() * 0.7)),
			float32(math.Sin(glfw.GetTime() * 1.3)),
		}
	*/
	g.SpotLightComp.Position = g.CameraComp.Position
	g.SpotLightComp.Direction = g.CameraComp.Front
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

	// Cubes
	game.Engine.EntityStore.NewCubeEntity(mgl32.Vec3{-1.0, -1.0, -2.0}, 1)

	// LIGHTING

	// Ambient
	ambientLightEntity := game.Engine.EntityStore.NewEntity()
	ambientLightComponent := components.NewAmbientLightComponent(mgl32.Vec3{1.0, 1.0, 1.0}, 0.1)
	game.Engine.EntityStore.AddComponent(ambientLightEntity, ambientLightComponent)

	// Directional
	directionalLightEntity := game.Engine.EntityStore.NewEntity()
	directionalLightComponent := components.NewDirectionalLightComponent(mgl32.Vec3{-0.2, -1.0, -0.3}, mgl32.Vec3{1.0, 1.0, 1.0}, 1)
	game.Engine.EntityStore.AddComponent(directionalLightEntity, directionalLightComponent)

	// Point
	pointLightEntity := game.Engine.EntityStore.NewEntity()
	pointLightComponent := components.NewPointLightComponent(mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{1.0, 0.8, 0.7}, 1.0, 1.0, 0.09, 0.032)
	game.Engine.EntityStore.AddComponent(pointLightEntity, pointLightComponent)

	// Spot
	spotLightEntity := game.Engine.EntityStore.NewEntity()
	spotLightComponent := components.NewSpotLightComponent(cameraComp.Position, mgl32.Vec3{1.0, 1.0, 1.0}, cameraComp.Front, float32(math.Cos(float64(mgl32.DegToRad(12.5)))), float32(math.Cos(float64(mgl32.DegToRad(17.5)))), 1.0, 1.0, 0.09, 0.032)
	game.Engine.EntityStore.AddComponent(spotLightEntity, spotLightComponent)
	game.SpotLightComp = spotLightComponent

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
