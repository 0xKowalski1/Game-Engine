package input

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type InputManager struct {
	GlfwWindow     *glfw.Window
	keyMap         map[glfw.Key]int
	actionState    map[int]bool
	actionHandlers map[int]func()
}

func NewInputManager(glfwWindow *glfw.Window) *InputManager {
	return &InputManager{
		GlfwWindow:     glfwWindow,
		keyMap:         make(map[glfw.Key]int),
		actionState:    make(map[int]bool),
		actionHandlers: make(map[int]func()),
	}
}

func (im *InputManager) RegisterKeyAction(key glfw.Key, action int, handler func()) {
	im.keyMap[key] = action
	im.actionHandlers[action] = handler
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
