package main

import (
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

	// Start the main loop
	eng.Run()
}
