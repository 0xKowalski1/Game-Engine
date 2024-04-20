package main

import (
	"log"
	"os"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	shaderProgram uint32 // Global variable to store the shader program ID
	vao           uint32 // Global variable to store the Vertex Array Object ID
)

var vertices = []float32{
	0.5, -0.5, 0.0, 1.0, 0.0, 0.0,
	-0.5, -0.5, 0.0, 0.0, 1.0, 0.0,
	0.0, 0.5, 0.0, 0.0, 0.0, 1.0,
}

var indices = []uint32{
	0, 1, 2, // First triangle
}

func main() {
	// Initialize GLFW
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// Set the required options for OpenGL version and compatibility
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create Window
	window, err := glfw.CreateWindow(800, 600, "Game", nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
	}
	window.MakeContextCurrent() // Make the OpenGL context current on the created window

	// Initialize OpenGL
	if err := gl.Init(); err != nil {
		log.Fatalln("failed to initialize opengl:", err)
	}

	// Set initial viewport to match the window size
	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))
	window.SetFramebufferSizeCallback(func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height)) // Adjust viewport on window resize
	})

	// Load Shaders from files
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

	// Main rendering loop
	for !window.ShouldClose() {
		// Process input

		// Render
		draw() // Call the draw function to render graphics

		window.SwapBuffers() // Swap the front and back buffers
		glfw.PollEvents()    // Poll for and process events like keyboard and mouse input
	}

	glfw.Terminate()
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

	// Delete shaders after linking as they are no longer needed

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	prepareTriangle() // Setup the vertex data for a triangle
}

func prepareTriangle() {
	// Generate a Vertex Array Object
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	// Generate a Vertex Buffer Object
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW) // Upload vertex data to the buffer

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 6*4, gl.Ptr(nil)) // Describe the vertex data layout
	gl.EnableVertexAttribArray(0)                                   // Enable the vertex attribute array

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 6*4, gl.Ptr(uintptr(3*4)))
	gl.EnableVertexAttribArray(1)

	// Generate a Element Buffer Object
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

}

func compileShaders(vertexShaderSource, fragmentShaderSource string) (uint32, uint32) {
	// Create & compile a vertex shader
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	csources, free := gl.Strs(vertexShaderSource + "\x00")
	gl.ShaderSource(vertexShader, 1, csources, nil)
	free()
	gl.CompileShader(vertexShader)

	// Create & compile a fragment shader
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	csources, free = gl.Strs(fragmentShaderSource + "\x00")
	gl.ShaderSource(fragmentShader, 1, csources, nil)
	free()
	gl.CompileShader(fragmentShader)

	return vertexShader, fragmentShader
}

func draw() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // Clear the color and depth buffers
	gl.ClearColor(0.0, 0.0, 0.4, 0.0)                   // Set the clear color to a dark blue

	gl.UseProgram(shaderProgram) // Use the shader program

	gl.BindVertexArray(vao)                                                  // Bind the VAO
	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil) // Draw the triangle
}
