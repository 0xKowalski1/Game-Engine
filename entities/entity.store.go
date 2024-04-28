package entities

type Entity struct {
	ID     uint32
	Active bool
}

type EntityStore struct {
	entities    []Entity
	activeCount int
}

func NewEntityStore() *EntityStore {
	return &EntityStore{}
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
