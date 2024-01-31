package repository

// PolyfeaRepositoryFilterFunc is a function that takes an item of type ItemType and returns a boolean.
type PolyfeaRepositoryFilterFunc[ItemType interface{}] func(mf ItemType) bool

// PolyfeaRepository is a generic interface for storing and retrieving items of type ItemType.
// The type ItemType must implement the GetName() method.
type PolyfeaRepository[ItemType interface{ GetName() string }] interface {
	StoreItem(item ItemType) error

	GetItem(item ItemType) (ItemType, error)

	GetItems(filter PolyfeaRepositoryFilterFunc[ItemType]) ([]ItemType, error)

	DeleteItem(item ItemType) error
}

type InMemoryPolyfeaRepository[ItemType interface{ GetName() string }] struct {
	items map[string]ItemType
}

func NewInMemoryPolyfeaRepository[ItemType interface{ GetName() string }]() *InMemoryPolyfeaRepository[ItemType] {
	return &InMemoryPolyfeaRepository[ItemType]{
		items: make(map[string]ItemType),
	}
}

func (r *InMemoryPolyfeaRepository[ItemType]) StoreItem(item ItemType) error {
	r.items[item.GetName()] = item
	return nil
}

func (r *InMemoryPolyfeaRepository[ItemType]) GetItem(item ItemType) (ItemType, error) {
	return r.items[item.GetName()], nil
}

func (r *InMemoryPolyfeaRepository[ItemType]) GetItems(filter PolyfeaRepositoryFilterFunc[ItemType]) ([]ItemType, error) {
	var result []ItemType
	for _, item := range r.items {
		if filter(item) {
			result = append(result, item)
		}
	}
	return result, nil
}

func (r *InMemoryPolyfeaRepository[ItemType]) DeleteItem(item ItemType) error {
	delete(r.items, item.GetName())
	return nil
}
