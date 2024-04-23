package main

import (
	"0xKowalski/game/components"
	"0xKowalski/game/engine"
	"log"
	"runtime"

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

type Game struct {
	Engine    *engine.Engine
	TestCubes []TestCube
}

func (g *Game) MainLoop() {
	for _, testCube := range g.TestCubes {
		g.RotateTestCube(testCube.ID)
	}
}

func (g *Game) NewTestCube(position mgl32.Vec3) {
	var vertices = []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}

	var indices = []uint32{}

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

	// Rotation amount in radians for each frame
	rotationAmount := mgl32.DegToRad(1.0) // Adjust this value as needed

	// Update the rotation of the cube
	// This will rotate the cube around the x-axis; for a different axis, change the vector
	currentRotation := transformComponent.Rotation
	deltaRotation := mgl32.QuatRotate(rotationAmount, mgl32.Vec3{1, 1, 0})
	newRotation := currentRotation.Mul(deltaRotation)

	// Set the new rotation back to the transform component
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

	cameraEntity := eng.EntityStore.NewEntity()

	cameraComp := components.NewCameraComponent(
		mgl32.Vec3{0, 0, 5}, // Position
		mgl32.Vec3{0, 0, 0}, // Target
		mgl32.Vec3{0, 1, 0}, // Up vector
		45.0,                // Field of view in degrees
		800.0/600.0,         // Aspect ratio, should get this from elsewhere
		0.1,                 // Near clipping plane
		100.0,               // Far clipping plane
	)
	eng.ComponentStore.AddComponent(cameraEntity, cameraComp)

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

	eng.Run(game.MainLoop)
}
