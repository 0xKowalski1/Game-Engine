package entities

import (
	"reflect"
)

type EntityStore struct {
	entities    []Entity
	activeCount int
	components  map[reflect.Type]map[uint32]Component
}

func NewEntityStore() *EntityStore {
	return &EntityStore{
		components: make(map[reflect.Type]map[uint32]Component),
	}
}

// Entities

type Entity struct {
	ID     uint32
	Active bool
}

func (store *EntityStore) NewEntity() Entity {
	if store.activeCount < len(store.entities) {
		// Reuse a previously freed entity slot.
		entity := &store.entities[store.activeCount]
		entity.Active = true
		entity.ID = uint32(store.activeCount)
		store.activeCount++
		return *entity
	} else {
		// Expand the array with a new entity.
		id := uint32(len(store.entities))
		entity := Entity{ID: id, Active: true}
		store.entities = append(store.entities, entity)
		store.activeCount++
		return entity
	}
}

// Memory currently grows infinetely, swap for a batched approach in future
func (store *EntityStore) FreeEntity(id uint32) {
	if int(id) < len(store.entities) && store.entities[id].Active {
		store.entities[id].Active = false
		store.activeCount--
		// Swap with the last active entity
		if store.activeCount > 0 && int(id) != store.activeCount {
			store.entities[id], store.entities[store.activeCount] = store.entities[store.activeCount], store.entities[id]
		}
	}
}

func (store *EntityStore) ActiveEntities() []Entity {
	return store.entities[:store.activeCount]
}

// Components

type Component interface{}

func (store *EntityStore) AddComponent(entity Entity, component Component) {
	compType := reflect.TypeOf(component)
	if store.components[compType] == nil {
		store.components[compType] = make(map[uint32]Component)
	}
	store.components[compType][entity.ID] = component
}

func (store *EntityStore) GetComponent(entity Entity, componentType Component) Component {
	compType := reflect.TypeOf(componentType)
	if comps, ok := store.components[compType]; ok {
		return comps[entity.ID]
	}
	return nil
}

func (store *EntityStore) RemoveComponent(entity Entity, componentType Component) {
	compType := reflect.TypeOf(componentType)
	if comps, ok := store.components[compType]; ok {
		delete(comps, entity.ID)
	}
}
