package components

import (
	"0xKowalski/game/entities"
	"reflect"
)

// Component is the base interface for all components, it could be empty as all data handling is type-specific.
type Component interface{}

// ComponentStore holds components mapped by entity ID and component type.
type ComponentStore struct {
	components map[reflect.Type]map[uint32]Component
}

// NewComponentStore creates a new component store.
func NewComponentStore() *ComponentStore {
	return &ComponentStore{
		components: make(map[reflect.Type]map[uint32]Component),
	}
}

// AddComponent adds a component to an entity.
func (store *ComponentStore) AddComponent(entity entities.Entity, component Component) {
	compType := reflect.TypeOf(component)
	if store.components[compType] == nil {
		store.components[compType] = make(map[uint32]Component)
	}
	store.components[compType][entity.ID] = component
}

// GetComponent retrieves a component attached to an entity by ID, returning nil if no component of that type exists.
func (store *ComponentStore) GetComponent(entity entities.Entity, componentType Component) Component {
	compType := reflect.TypeOf(componentType)
	if comps, ok := store.components[compType]; ok {
		return comps[entity.ID]
	}
	return nil
}

// RemoveComponent removes a component from an entity.
func (store *ComponentStore) RemoveComponent(entity entities.Entity, componentType Component) {
	compType := reflect.TypeOf(componentType)
	if comps, ok := store.components[compType]; ok {
		delete(comps, entity.ID)
	}
}
