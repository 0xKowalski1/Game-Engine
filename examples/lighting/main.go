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
		mgl32.Vec3{0, 0, 10}, // Position: Initial position of the camera in the world
		mgl32.Vec3{0, 1, 0},  // WorldUp: The up vector of the world, typically Y-axis is up
		-90.0,                // Yaw: Initial yaw angle, facing forward along the Z-axis
		0.0,                  // Pitch: Initial pitch angle, looking straight at the horizon
		45.0,                 // Field of view in degrees
		800.0/600.0,          // Aspect ratio: width divided by height of the viewport
		0.1,                  // Near clipping plane: the closest distance the camera can see
		100.0,                // Far clipping plane: the farthest distance the camera can see
	)

	game.Engine.EntityStore.AddComponent(cameraEntity, cameraComp)
	game.Camera = cameraEntity

	game.Engine.EntityStore.NewCubeEntity(mgl32.Vec3{1.0, 1.0, 1.0})

	// LIGHTING

	// END LIGHTING

	// Inputs
	const (
		CloseApp = iota
	)
	game.Engine.InputManager.RegisterKeyAction(glfw.KeyEscape, CloseApp, func() { game.Engine.Window.GlfwWindow.SetShouldClose(true) })

	// Loop
	game.Engine.Run(game.MainLoop)
}
