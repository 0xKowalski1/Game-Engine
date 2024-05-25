package main

import (
	"0xKowalski/game/components"
	"0xKowalski/game/engine"
	"log"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Game struct {
	Engine *engine.Engine
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
	freeCam := game.Engine.EntityStore.NewFreecamEntity(mgl32.Vec3{0, 3, 10})

	// Ground plane
	plane := game.Engine.EntityStore.NewPlaneEntity(mgl32.Vec3{0.0, 0.0, 0.0})
	planePhysicsComp := components.NewPhysicsComponent(mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 0.0, 0.0}, 1, true)
	game.Engine.EntityStore.AddComponent(*plane, planePhysicsComp)
	planeColliderComp := components.NewBoxColliderComponent(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 0, 1}, 0.5, 0.3)
	game.Engine.EntityStore.AddComponent(*plane, planeColliderComp)

	// Physics Cube
	cube := game.Engine.EntityStore.NewCubeEntity(mgl32.Vec3{0.0, 5.0, 0.0}, 1)
	physicsComp := components.NewPhysicsComponent(mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 0.0, 0.0}, mgl32.Vec3{0.0, 0.0, 0.0}, 1, false)
	game.Engine.EntityStore.AddComponent(*cube, physicsComp)
	boxColliderComp := components.NewBoxColliderComponent(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{1, 1, 1}, 0.5, 0.3)
	game.Engine.EntityStore.AddComponent(*cube, boxColliderComp)

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
		freeCam.Move(freeCam.CameraComponent.Front, cameraSpeed)
	})
	game.Engine.InputManager.RegisterKeyAction(glfw.KeyS, MoveBackward, func() {
		freeCam.Move(freeCam.CameraComponent.Front.Mul(-1), cameraSpeed)
	})
	game.Engine.InputManager.RegisterKeyAction(glfw.KeyD, StrafeRight, func() {
		freeCam.Move(freeCam.CameraComponent.Right, cameraSpeed)

	})
	game.Engine.InputManager.RegisterKeyAction(glfw.KeyA, StrafeLeft, func() {
		freeCam.Move(freeCam.CameraComponent.Right.Mul(-1), cameraSpeed)

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

		freeCam.Rotate(xOffset*0.05, yOffset*0.05)
		game.Engine.InputManager.LastX = xpos
		game.Engine.InputManager.LastY = ypos
	})

	game.Engine.InputManager.RegisterMouseScrollHandler(func(xoffset, yoffset float64) {
		fov := freeCam.CameraComponent.FieldOfView
		freeCam.CameraComponent.FieldOfView = fov - float32(yoffset)

		if fov < 1.0 {
			freeCam.CameraComponent.FieldOfView = 1.0
		} else if fov > 45.0 {
			freeCam.CameraComponent.FieldOfView = 45.0
		}
	})

	// Loop
	game.Engine.Run(game.MainLoop)
}
