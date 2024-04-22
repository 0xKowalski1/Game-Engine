package main

import (
	"0xKowalski/game/components"
	"0xKowalski/game/engine"
	"log"
	"runtime"
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

	buffers := components.NewBufferComponent(vertices, indices)
	eng.ComponentStore.AddComponent(entity, buffers)

	// Start the main loop
	eng.Run()
}
