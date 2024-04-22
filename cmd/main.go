package main

import (
	"0xKowalski/game/engine"
	"log"
	"runtime"

	"github.com/go-gl/glfw/v3.3/glfw"
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

	// Start the main loop
	for !eng.Window.GlfwWindow.ShouldClose() {
		// Game rendering or updates can be handled here

		// Swap buffers and poll IO events (keys pressed/released, mouse moved, etc.)
		eng.Window.GlfwWindow.SwapBuffers()
		glfw.PollEvents()
	}

	// Cleanup resources
	eng.Cleanup()
}
