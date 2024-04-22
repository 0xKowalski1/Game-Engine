package ecs

type Entity uint32

var nextID Entity = 1
var freeList []Entity

// NewEntity creates a new unique entity
func NewEntity() Entity {
	if len(freeList) > 0 {
		id := freeList[len(freeList)-1]
		freeList = freeList[:len(freeList)-1]
		return id
	}
	id := nextID
	nextID++
	return id
}

// FreeEntity marks an entity as free to be reused
func FreeEntity(id Entity) {
	freeList = append(freeList, id)
}
