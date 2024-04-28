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

type Game struct {
	Engine *engine.Engine
}

func (g *Game) MainLoop() {

}

func main() {
	// Init game and engine
	game := Game{}

	eng, err := engine.InitEngine()
	if err != nil {
		log.Printf("Error starting engine: %v", err)
		panic(err)
	}

	game.Engine = eng
	// End Init

	// Setup Game

	// End setup game

	// Run engine with mainloop
	game.Engine.Run(game.MainLoop)
}
