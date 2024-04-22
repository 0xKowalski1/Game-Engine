package window

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
)

type WindowConfig struct {
	Title  string
	Width  int
	Height int
}

type Window struct {
	GlfwWindow   *glfw.Window
	WindowConfig WindowConfig
}

func InitWindow(windowConfig WindowConfig) (*Window, error) {
	if err := glfw.Init(); err != nil {
		log.Fatalf("Failed to initialize GLFW: %v", err)
		return nil, err
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	// Create a GLFW window.
	glfwWindow, err := glfw.CreateWindow(windowConfig.Width, windowConfig.Height, windowConfig.Title, nil, nil)
	if err != nil {
		glfw.Terminate()
		log.Fatalf("Failed to create GLFW window: %v", err)
		return nil, err
	}

	glfwWindow.MakeContextCurrent()

	window := &Window{
		GlfwWindow:   glfwWindow,
		WindowConfig: windowConfig,
	}

	return window, nil
}

func (w *Window) Terminate() {
	glfw.Terminate()
}

func (w *Window) GetWidthAndHeight() (int32, int32) {
	width, height := w.GlfwWindow.GetFramebufferSize()

	return int32(width), int32(height)
}
