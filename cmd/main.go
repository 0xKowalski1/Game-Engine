package main

import (
	"0xKowalski/game/components"
	"0xKowalski/game/engine"
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// GLFW event handling must run on the main OS thread.
	runtime.LockOSThread()
}

type TestCube struct {
	ID uint32
}

type Camera struct {
	ID   uint32
	Comp *components.CameraComponent
}

func (g *Game) NewCamera() {
	cameraEntity := g.Engine.EntityStore.NewEntity()

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

	g.Engine.ComponentStore.AddComponent(cameraEntity, cameraComp)

	g.Camera = Camera{ID: cameraEntity.ID, Comp: cameraComp}
}

type Game struct {
	Engine    *engine.Engine
	Camera    Camera
	TestCubes []TestCube
}

func (g *Game) MainLoop() {
	for _, testCube := range g.TestCubes {
		g.RotateTestCube(testCube.ID)
	}
}

func (g *Game) NewTestCube(position mgl32.Vec3) {
	var vertices = []float32{
		// Front face
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		// Back face
		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		// Left face
		-0.5, -0.5, -0.5, 0.0, 0.0,
		-0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		// Right face
		0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 0.0, 1.0,
		// Top face
		-0.5, 0.5, -0.5, 0.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		// Bottom face
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 1.0,
	}

	var indices = []uint32{
		// Front face
		0, 1, 2, 0, 2, 3,
		// Back face
		4, 5, 6, 4, 6, 7,
		// Left face
		8, 9, 10, 8, 10, 11,
		// Right face
		12, 13, 14, 12, 14, 15,
		// Top face
		16, 17, 18, 16, 18, 19,
		// Bottom face
		20, 21, 22, 20, 22, 23,
	}

	entity := g.Engine.EntityStore.NewEntity()

	mesh := components.NewMeshComponent(vertices, indices)
	g.Engine.ComponentStore.AddComponent(entity, mesh)

	texture, err := components.NewTextureComponent("assets/textures/wall.jpg")
	if err != nil {
		log.Printf("Error creating texture component: %v", err)
	}
	g.Engine.ComponentStore.AddComponent(entity, texture)

	transform := components.NewTransformComponent(position)
	g.Engine.ComponentStore.AddComponent(entity, transform)

	buffers := components.NewBufferComponent(vertices, indices)
	g.Engine.ComponentStore.AddComponent(entity, buffers)

	g.TestCubes = append(g.TestCubes, TestCube{ID: entity.ID})
}

func (g *Game) RotateTestCube(testCubeID uint32) {
	cubeEntity := g.Engine.EntityStore.ActiveEntities()[testCubeID]
	transformComponent, _ := g.Engine.ComponentStore.GetComponent(cubeEntity, &components.TransformComponent{}).(*components.TransformComponent)

	rotationAmount := mgl32.DegToRad(1.0)

	currentRotation := transformComponent.Rotation
	deltaRotation := mgl32.QuatRotate(rotationAmount, mgl32.Vec3{1, 1, 0})
	newRotation := currentRotation.Mul(deltaRotation)

	transformComponent.SetRotation(newRotation)
}

func main() {
	game := Game{}

	eng, err := engine.InitEngine()
	if err != nil {
		log.Printf("Error starting engine: %v", err)
		panic(err)
	}

	game.Engine = eng

	game.NewCamera()

	var testCubePositions = []mgl32.Vec3{
		{0.0, 0.0, 0.0},
		{2.0, 5.0, -15.0},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{2.4, -0.4, -3.5},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
	}

	for _, testCubePosition := range testCubePositions {
		game.NewTestCube(testCubePosition)
	}

	// Lighting
	// Ambient
	ambientLightEntity := game.Engine.EntityStore.NewEntity()

	ambientLightComponent := components.NewAmbientLightComponent(mgl32.Vec3{1.0, 1.0, 1.0}, 0.2)
	game.Engine.ComponentStore.AddComponent(ambientLightEntity, ambientLightComponent)

	// Register Inputs
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
		game.Camera.Comp.Move(game.Camera.Comp.Front, cameraSpeed)
	})
	game.Engine.InputManager.RegisterKeyAction(glfw.KeyS, MoveBackward, func() {
		game.Camera.Comp.Move(game.Camera.Comp.Front.Mul(-1), cameraSpeed)
	})
	game.Engine.InputManager.RegisterKeyAction(glfw.KeyD, StrafeRight, func() {
		game.Camera.Comp.Move(game.Camera.Comp.Right, cameraSpeed)

	})
	game.Engine.InputManager.RegisterKeyAction(glfw.KeyA, StrafeLeft, func() {
		game.Camera.Comp.Move(game.Camera.Comp.Right.Mul(-1), cameraSpeed)

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

		game.Camera.Comp.Rotate(xOffset*0.05, yOffset*0.05)
		game.Engine.InputManager.LastX = xpos
		game.Engine.InputManager.LastY = ypos
	})

	game.Engine.InputManager.RegisterMouseScrollHandler(func(xoffset, yoffset float64) {
		fov := game.Camera.Comp.FieldOfView
		game.Camera.Comp.FieldOfView = fov - float32(yoffset)

		if fov < 1.0 {
			game.Camera.Comp.FieldOfView = 1.0
		} else if fov > 45.0 {
			game.Camera.Comp.FieldOfView = 45.0
		}
	})

	eng.Run(game.MainLoop)
}
