package systems

import (
	"0xKowalski/game/components"
	"0xKowalski/game/entities"
	"0xKowalski/game/graphics"
	"0xKowalski/game/window"
	"fmt"
	"log"

	"github.com/go-gl/gl/v4.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type RenderSystem struct {
	ShaderProgram *graphics.ShaderProgram
	EntityStore   *entities.EntityStore
}

func NewRenderSystem(win *window.Window, entityStore *entities.EntityStore) (*RenderSystem, error) {
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

	return rs, nil
}

func (rs *RenderSystem) getShaderLoc(name string) (int32, error) {
	loc := gl.GetUniformLocation(rs.ShaderProgram.ID, gl.Str(name+"\x00"))
	if loc == -1 {
		return -1, fmt.Errorf("Could not find the '%s' uniform location", name)
	}
	return loc, nil
}

func (rs *RenderSystem) SetShaderUniformMat4(name string, value mgl32.Mat4) {
	loc, err := rs.getShaderLoc(name)
	if err != nil {
		log.Println(err)
		return
	}

	gl.UniformMatrix4fv(loc, 1, false, &value[0])
}

func (rs *RenderSystem) SetShaderUniformVec3(name string, value mgl32.Vec3) {
	loc, err := rs.getShaderLoc(name)
	if err != nil {
		log.Println(err)
		return
	}

	gl.Uniform3f(loc, value.X(), value.Y(), value.Z())
}

func (rs *RenderSystem) SetShaderUniformFloat(name string, value float32) {
	loc, err := rs.getShaderLoc(name)
	if err != nil {
		log.Println(err)
		return
	}

	gl.Uniform1f(loc, value)
}

func (rs *RenderSystem) SetShaderUniformInt(name string, value int32) {
	loc, err := rs.getShaderLoc(name)
	if err != nil {
		log.Println(err)
		return
	}

	gl.Uniform1i(loc, value)
}

func (rs *RenderSystem) bindTexture(textureComponent *components.TextureComponent) {
	if textureComponent == nil {
		log.Println("No texture component provided")
		return
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, textureComponent.TextureID)

	rs.SetShaderUniformInt("texture1", 0)
}

func (rs *RenderSystem) renderEntity(meshComponent *components.MeshComponent, bufferComponent *components.BufferComponent, transformComponent *components.TransformComponent, cameraComponent *components.CameraComponent) {
	if meshComponent == nil || bufferComponent == nil || transformComponent == nil {
		log.Println("Mesh, buffer or transform component is nil, cannot render entity")
		return
	}

	modelMatrix := transformComponent.GetModelMatrix()
	rs.SetShaderUniformMat4("model", modelMatrix)

	viewMatrix := cameraComponent.GetViewMatrix()
	rs.SetShaderUniformMat4("view", viewMatrix)

	projectionMatrix := cameraComponent.GetProjectionMatrix()
	rs.SetShaderUniformMat4("projection", projectionMatrix)

	gl.BindVertexArray(bufferComponent.VAO)
	gl.DrawElements(gl.TRIANGLES, int32(len(meshComponent.Indices)), gl.UNSIGNED_INT, gl.Ptr(nil))
	gl.BindVertexArray(0)
}

func (rs *RenderSystem) Update() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear the color and depth buffers
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)                   // Set background color to black

	rs.ShaderProgram.Use()

	// get camera component
	cameraComponentInterface := rs.EntityStore.GetAllComponents(&components.CameraComponent{})[0]
	cameraComponent, cameraOk := cameraComponentInterface.(*components.CameraComponent)

	if !cameraOk {
		log.Fatalf("Failed to get camera component")
	}

	// Get renderable components and render them
	renderableComponents := rs.EntityStore.GetAllComponents(&components.RenderableComponent{})
	for _, renderableComponent := range renderableComponents {
		comp, ok := renderableComponent.(*components.RenderableComponent)

		if ok {
			rs.renderEntity(comp.MeshComponent, comp.BufferComponent, comp.TransformComponent, cameraComponent)
			if comp.TextureComponent != nil {
				rs.bindTexture(comp.TextureComponent)
			}
		} else {
			log.Println("Failed to parse render component")
		}

	}

	// Lighting
	// Ambient light (1 max)
	ambientLightComponents := rs.EntityStore.GetAllComponents(&components.AmbientLightComponent{})
	if len(ambientLightComponents) > 1 {
		log.Fatalf("Exceeded max amount of ambient lights - Max: 1, Used: %d", len(ambientLightComponents))

	}
	if len(ambientLightComponents) == 1 {
		ambientLightComponentInterface := ambientLightComponents[0]
		ambientLightComponent, ambientOk := ambientLightComponentInterface.(*components.AmbientLightComponent)
		if ambientOk {
			rs.SetShaderUniformVec3("ambientLightColor", ambientLightComponent.Color)
			rs.SetShaderUniformFloat("ambientLightIntensity", ambientLightComponent.Intensity)
		}
	}

	// Directional Light (1 max) - Might want to increase this, e.g, sun & moon? multiple suns?
	directionalLightComponents := rs.EntityStore.GetAllComponents(&components.DirectionalLightComponent{})
	if len(directionalLightComponents) > 1 {
		log.Fatalf("Exceeded max amount of Directional lights - Max: 1, Used: %d", len(directionalLightComponents))
	}
	if len(directionalLightComponents) == 1 {
		directionalLightComponentInterface := directionalLightComponents[0]
		directionalLightComponent, directionalOk := directionalLightComponentInterface.(*components.DirectionalLightComponent)
		if directionalOk {
			rs.SetShaderUniformVec3("directionalLightDirection", directionalLightComponent.Direction)
			rs.SetShaderUniformVec3("directionalLightColor", directionalLightComponent.Color)
			rs.SetShaderUniformFloat("directionalLightIntensity", directionalLightComponent.Intensity)
		}
	}

	// Point Lights (idk max yet)
	pointLightComponents := rs.EntityStore.GetAllComponents(&components.PointLightComponent{})
	if len(pointLightComponents) > 0 {
		for _, pointLightComponentInterface := range pointLightComponents {
			pointLightComponent, pointLightOk := pointLightComponentInterface.(*components.PointLightComponent)

			if pointLightOk {
				rs.SetShaderUniformVec3("pointLight.position", pointLightComponent.Position)
				rs.SetShaderUniformVec3("pointLight.color", pointLightComponent.Color)
				rs.SetShaderUniformFloat("pointLight.intensity", pointLightComponent.Intensity)
				rs.SetShaderUniformFloat("pointLight.constant", pointLightComponent.Constant)
				rs.SetShaderUniformFloat("pointLight.linear", pointLightComponent.Linear)
				rs.SetShaderUniformFloat("pointLight.quadratic", pointLightComponent.Quadratic)
			}
		}
	}
}
