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

func main() {
	eng, err := engine.InitEngine()
	if err != nil {
		log.Printf("Error starting engine: %v", err)
		panic(err)
	}

	vertices := []float32{
		// positions          // colors           // texture coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // top left
	}
	indices := []uint32{
		0, 1, 2,
		2, 3, 0,
	}

	entity := eng.EntityStore.NewEntity()

	mesh := components.NewMeshComponent(vertices, indices)
	eng.ComponentStore.AddComponent(entity, mesh)

	texture, err := components.NewTextureComponent("assets/textures/wall.jpg")
	if err != nil {
		log.Printf("Error creating texture component: %v", err)
	}
	eng.ComponentStore.AddComponent(entity, texture)

	transform := components.NewTransformComponent()
	// Apply rotation around the Z-axis
	transform.SetRotation(mgl32.QuatRotate(mgl32.DegToRad(45.0), mgl32.Vec3{0, 0, 1}))
	// Apply scaling
	transform.SetScale(0.5, 0.5, 0.5) // Scale down by 50%
	eng.ComponentStore.AddComponent(entity, transform)

	buffers := components.NewBufferComponent(vertices, indices)
	eng.ComponentStore.AddComponent(entity, buffers)

	// Start the main loop
	eng.Run()
}
