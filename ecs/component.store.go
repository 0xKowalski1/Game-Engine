package ecs

import (
	"reflect"
)

// Component is the base interface for all components, it could be empty as all data handling is type-specific.
type Component interface{}

// ComponentStore holds components mapped by entity and component type.
type ComponentStore struct {
	components map[reflect.Type]map[Entity]Component
}

// NewComponentStore creates a new component store.
func NewComponentStore() *ComponentStore {
	return &ComponentStore{
		components: make(map[reflect.Type]map[Entity]Component),
	}
}

// AddComponent adds a component to an entity.
func (store *ComponentStore) AddComponent(entity Entity, component Component) {
	compType := reflect.TypeOf(component)
	if store.components[compType] == nil {
		store.components[compType] = make(map[Entity]Component)
	}
	store.components[compType][entity] = component
}

// GetComponent retrieves a component attached to an entity, returning nil if no component of that type exists.
func (store *ComponentStore) GetComponent(entity Entity, componentType Component) Component {
	compType := reflect.TypeOf(componentType)
	if comps, ok := store.components[compType]; ok {
		return comps[entity]
	}
	return nil
}

// RemoveComponent removes a component from an entity.
func (store *ComponentStore) RemoveComponent(entity Entity, componentType Component) {
	compType := reflect.TypeOf(componentType)
	if comps, ok := store.components[compType]; ok {
		delete(comps, entity)
	}
}

// HasComponent checks if an entity has a component of a specific type.
func (store *ComponentStore) HasComponent(entity Entity, componentType Component) bool {
	compType := reflect.TypeOf(componentType)
	if comps, ok := store.components[compType]; ok {
		_, exists := comps[entity]
		return exists
	}
	return false
}

// GetAllComponentsOfType returns a map of all components of a specific type.
func (store *ComponentStore) GetAllComponentsOfType(componentType Component) map[Entity]Component {
	compType := reflect.TypeOf(componentType)
	if comps, ok := store.components[compType]; ok {
		return comps
	}
	return nil
}
