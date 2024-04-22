package systems

import (
	"0xKowalski/game/components"
	"0xKowalski/game/ecs"
	"0xKowalski/game/graphics"
	"0xKowalski/game/window"
	"log"

	"github.com/go-gl/gl/v4.3-core/gl"
)

type RenderSystem struct {
	ShaderProgram  *graphics.ShaderProgram
	EntityStore    *ecs.EntityStore
	ComponentStore *ecs.ComponentStore
}

func NewRenderSystem(win *window.Window, entityStore *ecs.EntityStore, componentStore *ecs.ComponentStore) (*RenderSystem, error) {
	err := graphics.InitOpenGL(win)
	if err != nil {
		log.Printf("Error initializing renderer: %v", err)
		return nil, err
	}

	rs := new(RenderSystem)

	shaderProgram, err := graphics.InitShaderProgram("assets/shaders/vertex.glsl", "assets/shaders/fragment.glsl")
	if err != nil {
		return nil, err
	}

	rs.ShaderProgram = shaderProgram
	rs.EntityStore = entityStore
	rs.ComponentStore = componentStore

	return rs, nil
}

func (rs *RenderSystem) Update() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear the color and depth buffers
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)                   // Set the clear color to a dark blue

	rs.ShaderProgram.Use()

	for _, entity := range rs.EntityStore.ActiveEntities() {
		meshComponent, meshOk := rs.ComponentStore.GetComponent(entity, &components.MeshComponent{}).(*components.MeshComponent)
		bufferComponent, bufferOk := rs.ComponentStore.GetComponent(entity, &components.BufferComponent{}).(*components.BufferComponent)

		if !meshOk || !bufferOk {
			log.Println("Failed to get necessary rendering components for entity")
			continue // Skip this entity if the required components aren't available
		}

		// Proceed with rendering only if both components are present
		if meshComponent != nil && bufferComponent != nil {
			gl.BindVertexArray(bufferComponent.VAO)
			gl.DrawElements(gl.TRIANGLES, int32(len(meshComponent.Indices)), gl.UNSIGNED_INT, nil)
			gl.BindVertexArray(0)
		}
	}
}
