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
	TextureStore  *TextureStore
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
	rs.TextureStore = NewTextureStore()

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

func (rs *RenderSystem) renderEntity(comp *components.RenderableComponent) {
	modelMatrix := comp.TransformComponent.GetModelMatrix()
	rs.SetShaderUniformMat4("model", modelMatrix)

	if comp.TransformComponent == nil || comp.ModelComponent == nil {
		log.Println("Mesh, buffer, transform or material component is nil, cannot render entity")
		return
	}

	for i, meshComponent := range comp.ModelComponent.MeshComponents {
		materialComponent := comp.ModelComponent.MaterialComponents[i]
		bufferComponent := comp.ModelComponent.BufferComponents[i]

		diffuseMap, err := rs.TextureStore.GetTexture(materialComponent.DiffuseMap)
		if err != nil {
			log.Printf("Error getting diffusemap texture: %v", err)
		}
		specularMap, err := rs.TextureStore.GetTexture(materialComponent.SpecularMap)
		if err != nil {
			log.Printf("Error getting specularmap texture: %v", err)
		}

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, diffuseMap)
		rs.SetShaderUniformInt("material.diffuseMap", 0)

		gl.ActiveTexture(gl.TEXTURE1)
		gl.BindTexture(gl.TEXTURE_2D, specularMap)
		rs.SetShaderUniformInt("material.specularMap", 1)

		rs.SetShaderUniformFloat("material.shininess", materialComponent.Shininess)

		gl.BindVertexArray(bufferComponent.VAO)
		gl.DrawElements(gl.TRIANGLES, int32(len(meshComponent.Indices)), gl.UNSIGNED_INT, gl.Ptr(nil))
		gl.BindVertexArray(0)
	}

}

func (rs *RenderSystem) Update() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear the color and depth buffers
	gl.ClearColor(0.0, 0.0, 0.1, 0.0)                   // Set background color to black

	rs.ShaderProgram.Use()

	cameraEntity := rs.EntityStore.GetEntityWithComponentType(&components.CameraComponent{})

	// get camera component
	cameraComponentInterface := rs.EntityStore.GetComponent(cameraEntity, &components.CameraComponent{})
	cameraComponent, cameraOk := cameraComponentInterface.(*components.CameraComponent)

	transformComponentInterface := rs.EntityStore.GetComponent(cameraEntity, &components.TransformComponent{})
	transformComponent, transformComponentOk := transformComponentInterface.(*components.TransformComponent)

	if !cameraOk || !transformComponentOk {
		log.Println(cameraEntity, cameraComponentInterface, transformComponentInterface)
		log.Fatalf("Failed to get camera component")
	}

	rs.SetShaderUniformVec3("viewPos", transformComponent.Position)

	viewMatrix := cameraComponent.GetViewMatrix(transformComponent.Position)
	rs.SetShaderUniformMat4("view", viewMatrix)

	projectionMatrix := cameraComponent.GetProjectionMatrix()
	rs.SetShaderUniformMat4("projection", projectionMatrix)

	// Get renderable components and render them
	renderableComponents := rs.EntityStore.GetAllComponents(&components.RenderableComponent{})
	for _, renderableComponent := range renderableComponents {
		comp, ok := renderableComponent.(*components.RenderableComponent)

		if ok {
			rs.renderEntity(comp)
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
			rs.SetShaderUniformVec3("ambientLight.color", ambientLightComponent.Color)
			rs.SetShaderUniformFloat("ambientLight.intensity", ambientLightComponent.Intensity)
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
			rs.SetShaderUniformVec3("directionalLight.direction", directionalLightComponent.Direction)
			rs.SetShaderUniformVec3("directionalLight.color", directionalLightComponent.Color)
			rs.SetShaderUniformFloat("directionalLight.intensity", directionalLightComponent.Intensity)
		}
	}

	// Point Lights (idk max yet)
	pointLightComponents := rs.EntityStore.GetAllComponents(&components.PointLightComponent{})
	rs.SetShaderUniformInt("pointLightsCount", int32(len(pointLightComponents)))
	if len(pointLightComponents) > 0 {
		for index, pointLightComponentInterface := range pointLightComponents {
			pointLightComponent, pointLightOk := pointLightComponentInterface.(*components.PointLightComponent)

			if pointLightOk {
				rs.SetShaderUniformVec3(fmt.Sprintf("pointLights[%d].position", index), pointLightComponent.Position)
				rs.SetShaderUniformVec3(fmt.Sprintf("pointLights[%d].color", index), pointLightComponent.Color)
				rs.SetShaderUniformFloat(fmt.Sprintf("pointLights[%d].intensity", index), pointLightComponent.Intensity)
				rs.SetShaderUniformFloat(fmt.Sprintf("pointLights[%d].constant", index), pointLightComponent.Constant)
				rs.SetShaderUniformFloat(fmt.Sprintf("pointLights[%d].linear", index), pointLightComponent.Linear)
				rs.SetShaderUniformFloat(fmt.Sprintf("pointLights[%d].quadratic", index), pointLightComponent.Quadratic)
			}
		}
	}

	// Spot Lights (idk max yet)
	spotLightComponents := rs.EntityStore.GetAllComponents(&components.SpotLightComponent{})
	rs.SetShaderUniformInt("spotLightsCount", int32(len(spotLightComponents)))
	if len(spotLightComponents) > 0 {
		for index, spotLightComponentInterface := range spotLightComponents {
			spotLightComponent, spotLightOk := spotLightComponentInterface.(*components.SpotLightComponent)

			if spotLightOk {
				rs.SetShaderUniformVec3(fmt.Sprintf("spotLights[%d].position", index), spotLightComponent.Position)
				rs.SetShaderUniformVec3(fmt.Sprintf("spotLights[%d].color", index), spotLightComponent.Color)
				rs.SetShaderUniformVec3(fmt.Sprintf("spotLights[%d].direction", index), spotLightComponent.Direction)
				rs.SetShaderUniformFloat(fmt.Sprintf("spotLights[%d].cutOff", index), spotLightComponent.CutOff)
				rs.SetShaderUniformFloat(fmt.Sprintf("spotLights[%d].outerCutOff", index), spotLightComponent.OuterCutOff)

				rs.SetShaderUniformFloat(fmt.Sprintf("spotLights[%d].intensity", index), spotLightComponent.Intensity)
				rs.SetShaderUniformFloat(fmt.Sprintf("spotLights[%d].constant", index), spotLightComponent.Constant)
				rs.SetShaderUniformFloat(fmt.Sprintf("spotLights[%d].linear", index), spotLightComponent.Linear)
				rs.SetShaderUniformFloat(fmt.Sprintf("spotLights[%d].quadratic", index), spotLightComponent.Quadratic)

			}
		}
	}
}
