package systems

import (
	"0xKowalski/game/components"
	"0xKowalski/game/ecs"
	"0xKowalski/game/graphics"
	"0xKowalski/game/window"
	"log"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
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

func (rs *RenderSystem) bindTexture(textureComponent *components.TextureComponent) {
	if textureComponent == nil {
		log.Println("No texture component provided")
		return
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, textureComponent.TextureID)
	texUniform := gl.GetUniformLocation(rs.ShaderProgram.ID, gl.Str("texture1\x00"))
	if texUniform == int32(-1) {
		log.Println("Error getting uniform location for texture1")
		return
	}
	gl.Uniform1i(texUniform, 0)
}

func (rs *RenderSystem) SetShaderUniformVec3(name string, value mgl32.Vec3) {
	loc := gl.GetUniformLocation(rs.ShaderProgram.ID, gl.Str(name+"\x00"))
	if loc == -1 {
		log.Printf("Could not find the '%s' uniform location", name)
		return
	}

	gl.Uniform3f(loc, value.X(), value.Y(), value.Z())
}

func (rs *RenderSystem) SetShaderUniformFloat(name string, value float32) {
	loc := gl.GetUniformLocation(rs.ShaderProgram.ID, gl.Str(name+"\x00"))
	if loc == -1 {
		log.Printf("Could not find the '%s' uniform location", name)
		return
	}

	gl.Uniform1f(loc, value)
}

func (rs *RenderSystem) renderEntity(meshComponent *components.MeshComponent, bufferComponent *components.BufferComponent, transformComponent *components.TransformComponent, cameraComponent *components.CameraComponent) {
	if meshComponent == nil || bufferComponent == nil || transformComponent == nil {
		log.Println("Mesh, buffer or transform component is nil, cannot render entity")
		return
	}

	// Compute the model matrix based on the transform component
	modelMatrix := transformComponent.GetModelMatrix()
	// Get the uniform location for the model matrix in the shader
	modelLoc := gl.GetUniformLocation(rs.ShaderProgram.ID, gl.Str("model\x00"))
	if modelLoc == -1 {
		log.Println("Could not find the 'model' uniform location")
		return
	}
	// Pass the model matrix to the shader
	gl.UniformMatrix4fv(modelLoc, 1, false, &modelMatrix[0])

	viewMatrix := cameraComponent.GetViewMatrix()
	viewLoc := gl.GetUniformLocation(rs.ShaderProgram.ID, gl.Str("view\x00"))
	if viewLoc == -1 {
		log.Println("Could not find the 'view' uniform location")
		return
	}
	gl.UniformMatrix4fv(viewLoc, 1, false, &viewMatrix[0])

	projectionMatrix := cameraComponent.GetProjectionMatrix()
	projectionLoc := gl.GetUniformLocation(rs.ShaderProgram.ID, gl.Str("projection\x00"))
	if projectionLoc == -1 {
		log.Println("Could not find the 'projection' uniform location")
		return
	}
	gl.UniformMatrix4fv(projectionLoc, 1, false, &projectionMatrix[0])

	gl.BindVertexArray(bufferComponent.VAO)
	gl.DrawElements(gl.TRIANGLES, int32(len(meshComponent.Indices)), gl.UNSIGNED_INT, gl.Ptr(nil))
	gl.BindVertexArray(0)
}

func (rs *RenderSystem) Update() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear the color and depth buffers
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)                   // Set the clear color to a dark blue

	rs.ShaderProgram.Use()

	cameraEntity := rs.EntityStore.ActiveEntities()[0] // Camera is always first entity
	cameraComponent, cameraOk := rs.ComponentStore.GetComponent(cameraEntity, &components.CameraComponent{}).(*components.CameraComponent)

	if !cameraOk {
		log.Fatalf("Failed to get camera component")
	}

	for i, entity := range rs.EntityStore.ActiveEntities() {
		// Can skip first entity here as its the camera
		if i == 0 {
			continue
		}

		// Check if ambient light
		ambientLightComponent, ambientLightOk := rs.ComponentStore.GetComponent(entity, &components.AmbientLightComponent{}).(*components.AmbientLightComponent)
		if ambientLightComponent != nil && ambientLightOk {
			rs.SetShaderUniformVec3("ambientLightColor", ambientLightComponent.Color)
			rs.SetShaderUniformFloat("ambientLightIntensity", ambientLightComponent.Intensity)
			continue
		}

		// Check if directional light
		directionalLightComponent, directionalLightOk := rs.ComponentStore.GetComponent(entity, &components.DirectionalLightComponent{}).(*components.DirectionalLightComponent)

		if directionalLightComponent != nil && directionalLightOk {
			rs.SetShaderUniformVec3("directionalLightDirection", directionalLightComponent.Direction)

			rs.SetShaderUniformVec3("directionalLightColor", directionalLightComponent.Color)
			rs.SetShaderUniformFloat("directionalLightIntensity", directionalLightComponent.Intensity)
			continue
		}

		meshComponent, meshOk := rs.ComponentStore.GetComponent(entity, &components.MeshComponent{}).(*components.MeshComponent)
		bufferComponent, bufferOk := rs.ComponentStore.GetComponent(entity, &components.BufferComponent{}).(*components.BufferComponent)
		transformComponent, transformOk := rs.ComponentStore.GetComponent(entity, &components.TransformComponent{}).(*components.TransformComponent)
		textureComponent, _ := rs.ComponentStore.GetComponent(entity, &components.TextureComponent{}).(*components.TextureComponent)

		if !meshOk || !bufferOk || !transformOk {
			log.Println("Failed to get necessary rendering components for entity")
			continue
		}

		rs.bindTexture(textureComponent)
		rs.renderEntity(meshComponent, bufferComponent, transformComponent, cameraComponent)
	}
}
