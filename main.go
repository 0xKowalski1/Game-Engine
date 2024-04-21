package main

import (
	"fmt"
	"image"
	drawImage "image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	shaderProgram uint32 // Global variable to store the shader program ID
	vao           uint32 // Global variable to store the Vertex Array Object ID
	texture1      uint32
	texture2      uint32
)

var vertices = []float32{
	// positions          // colors           // texture coords
	0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
	0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom right
	-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom left
	-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // top left
}

var indices = []uint32{
	0, 1, 3,
	1, 2, 3,
}

var texCoords = []float32{
	0.0, 0.0, // lower-left corner
	1.0, 0.0, // lower-right corner
	0.5, 1.0, // top-center corner
}

func main() {
	// Initialize GLFW
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// Set the required options for OpenGL version and compatibility
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create Window
	window, err := glfw.CreateWindow(800, 600, "Game", nil, nil)
	if err != nil {
		log.Fatalln("failed to create window:", err)
		glfw.Terminate()
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

	// Texture 1
	gl.GenTextures(1, &texture1)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture1)

	img1, err := loadImage("wall.jpg")
	if err != nil {
		log.Fatalf("Failed to load texture: %v", err)
	}

	rgba := imageToRGBA(img1)
	width, height := rgba.Rect.Size().X, rgba.Rect.Size().Y
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	// Texture 2
	gl.GenTextures(1, &texture2)
	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, texture2)

	img2, err := loadImage("smiley.png")
	if err != nil {
		log.Fatalf("Failed to load texture: %v", err)
	}

	rgbaFlipped := imageToRGBA(img2)
	rgba = flipImageVertically(rgbaFlipped)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR_MIPMAP_LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.UseProgram(shaderProgram)
	textureUniform1 := gl.GetUniformLocation(shaderProgram, gl.Str("texture1\x00"))
	gl.Uniform1i(textureUniform1, 0)

	textureUniform2 := gl.GetUniformLocation(shaderProgram, gl.Str("texture2\x00"))
	gl.Uniform1i(textureUniform2, 1)

	prepareTriangle() // Setup the vertex data for a triangle
}

func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Determine the file format based on the file extension
	if strings.HasSuffix(strings.ToLower(filename), ".jpg") || strings.HasSuffix(strings.ToLower(filename), ".jpeg") {
		img, err := jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
		return img, nil
	} else if strings.HasSuffix(strings.ToLower(filename), ".png") {
		img, err := png.Decode(file)
		if err != nil {
			return nil, err
		}
		return img, nil
	}

	return nil, fmt.Errorf("unsupported file format for %v", filename)
}

func imageToRGBA(img image.Image) *image.RGBA {
	rgba := image.NewRGBA(img.Bounds())
	drawImage.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, drawImage.Src)
	return rgba
}

func flipImageVertically(img *image.RGBA) *image.RGBA {
	src := img.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, src.Dx(), src.Dy()))
	for y := 0; y < src.Dy(); y++ {
		for x := 0; x < src.Dx(); x++ {
			dst.Set(x, src.Dy()-y-1, img.At(x, y))
		}
	}
	return dst
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

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.Ptr(nil)) // Describe the vertex data layout
	gl.EnableVertexAttribArray(0)                                   // Enable the vertex attribute array

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.Ptr(uintptr(3*4)))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.Ptr(uintptr(6*4)))
	gl.EnableVertexAttribArray(2)

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
	gl.UseProgram(shaderProgram)                        // Use the shader program

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture1)

	gl.ActiveTexture(gl.TEXTURE1)
	gl.BindTexture(gl.TEXTURE_2D, texture2)

	gl.BindVertexArray(vao)                                                  // Bind the VAO
	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, nil) // Draw the triangle
}
