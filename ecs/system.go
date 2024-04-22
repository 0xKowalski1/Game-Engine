package ecs

// System interface represents a system that operates on entities
type System interface {
	Update(dt float32)
}
