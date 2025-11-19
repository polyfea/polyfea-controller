package repository

// FilterFunc defines a predicate function for filtering repository items.
type FilterFunc[Item interface{ GetName() string }] func(item Item) bool

// Repository is a generic interface for storing and retrieving items.
// Item must implement GetName() string.
type Repository[Item interface{ GetName() string }] interface {
	Store(item Item) error
	Get(item Item) (Item, error)
	List(filter FilterFunc[Item]) ([]Item, error)
	Delete(item Item) error
}

// InMemoryRepository is an in-memory implementation of Repository.
type InMemoryRepository[Item interface{ GetName() string }] struct {
	items map[string]Item
}

// NewInMemoryRepository creates a new in-memory repository.
func NewInMemoryRepository[Item interface{ GetName() string }]() *InMemoryRepository[Item] {
	return &InMemoryRepository[Item]{
		items: make(map[string]Item),
	}
}

// Store adds or updates an item in the repository.
func (r *InMemoryRepository[Item]) Store(item Item) error {
	r.items[item.GetName()] = item
	return nil
}

// Get retrieves an item by its name.
func (r *InMemoryRepository[Item]) Get(item Item) (Item, error) {
	return r.items[item.GetName()], nil
}

// List returns all items matching the filter.
func (r *InMemoryRepository[Item]) List(filter FilterFunc[Item]) ([]Item, error) {
	var result []Item
	for _, item := range r.items {
		if filter(item) {
			result = append(result, item)
		}
	}
	return result, nil
}

// Delete removes an item from the repository.
func (r *InMemoryRepository[Item]) Delete(item Item) error {
	delete(r.items, item.GetName())
	return nil
}
