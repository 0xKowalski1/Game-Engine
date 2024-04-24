package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type InputManager struct {
	GlfwWindow       *glfw.Window
	keyMap           map[glfw.Key]int
	actionState      map[int]bool
	actionHandlers   map[int]func()
	mouseMoveHandler func(xpos, ypos float64)
	LastX, LastY     float64
	FirstMouse       bool
}

func NewInputManager(glfwWindow *glfw.Window) *InputManager {
	glfwWindow.SetInputMode(glfw.CursorMode, glfw.CursorDisabled) // Disable cursor and capture inputs
	return &InputManager{
		GlfwWindow:     glfwWindow,
		keyMap:         make(map[glfw.Key]int),
		actionState:    make(map[int]bool),
		actionHandlers: make(map[int]func()),
		FirstMouse:     true,
	}
}

func (im *InputManager) RegisterKeyAction(key glfw.Key, action int, handler func()) {
	im.keyMap[key] = action
	im.actionHandlers[action] = handler
}

func (im *InputManager) RegisterMouseMoveHandler(handler func(xpos, ypos float64)) {
	im.mouseMoveHandler = handler
	im.GlfwWindow.SetCursorPosCallback(im.onMouseMove)
}

func (im *InputManager) Update() {
	im.GlfwWindow.SetKeyCallback(im.onKey)
	for action, active := range im.actionState {
		if active && im.actionHandlers[action] != nil {
			im.actionHandlers[action]()
		}
	}
}

func (im *InputManager) onKey(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if act, ok := im.keyMap[key]; ok {
		if action == glfw.Press {
			im.actionState[act] = true
		} else if action == glfw.Release {
			im.actionState[act] = false
		}

	}
}

func (im *InputManager) onMouseMove(w *glfw.Window, xpos, ypos float64) {
	if im.FirstMouse {
		im.LastX = xpos
		im.LastY = ypos
		im.FirstMouse = false
	}

	if im.mouseMoveHandler != nil {
		im.mouseMoveHandler(xpos, ypos)
	}

	im.LastX = xpos
	im.LastY = ypos
}
