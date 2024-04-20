package main

import (
	"log"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	shaderProgram uint32
	vao           uint32
)

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(800, 600, "Game", nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}
	window.MakeContextCurrent()

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize opengl:", err)
	}

	// Set initial viewport
	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	// Load Shaders
	vertexShaderSource, err := loadShaderFile("./shaders/vertex.glsl")
	if err != nil {
		log.Fatalf("failed to load vertex shader: %v", err)
	}
	fragmentShaderSource, err := loadShaderFile("./shaders/fragment.glsl")
	if err != nil {
		log.Fatalf("failed to load fragment shader: %v", err)
	}

	// Compile shaders, setup VAO, VBO
	setupOpenGL(vertexShaderSource, fragmentShaderSource)

	for !window.ShouldClose() {
		draw()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func loadShaderFile(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func setupOpenGL(vertexShaderSource, fragmentShaderSource string) {
	vertexShader, fragmentShader := compileShaders(vertexShaderSource, fragmentShaderSource)
	shaderProgram = gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	gl.DeleteShader(vertexShader)   // Delete after linking
	gl.DeleteShader(fragmentShader) // Delete after linking

	prepareTriangle()
}

func prepareTriangle() {
	var vertices = []float32{
		-0.5, -0.5, 0.0, // Vertex 1
		0.5, -0.5, 0.0, // Vertex 2
		0.0, 0.5, 0.0, // Vertex 3
	}

	var vbo uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)
}

func compileShaders(vertexShaderSource, fragmentShaderSource string) (uint32, uint32) {
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	csources, free := gl.Strs(vertexShaderSource + "\x00")
	gl.ShaderSource(vertexShader, 1, csources, nil)
	free()
	gl.CompileShader(vertexShader)

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	csources, free = gl.Strs(fragmentShaderSource + "\x00")
	gl.ShaderSource(fragmentShader, 1, csources, nil)
	free()
	gl.CompileShader(fragmentShader)

	return vertexShader, fragmentShader
}

func draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)

	gl.UseProgram(shaderProgram)
	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}
