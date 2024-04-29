package main

import (
	"0xKowalski/game/components"
	"0xKowalski/game/engine"
	"0xKowalski/game/entities"
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Game struct {
	Engine *engine.Engine
	Camera entities.Entity
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
	game.Camera = cameraEntity

	// Cubes
	game.Engine.EntityStore.NewCubeEntity(mgl32.Vec3{-1.0, -1.0, -2.0})

	// LIGHTING
	// Ambient
	ambientLightEntity := game.Engine.EntityStore.NewEntity()
	ambientLightComponent := components.NewAmbientLightComponent(mgl32.Vec3{1.0, 1.0, 1.0}, 0.1)
	game.Engine.EntityStore.AddComponent(ambientLightEntity, ambientLightComponent)

	pointLightEntity := game.Engine.EntityStore.NewEntity()
	pointLightComponent := components.NewPointLightComponent(mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{1.0, 0.8, 0.7}, 1.0, 1.0, 0.09, 0.032)
	game.Engine.EntityStore.AddComponent(pointLightEntity, pointLightComponent)

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
