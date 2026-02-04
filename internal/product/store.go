package product

type Store interface {
	List() []Product
	Create(Product) Product
}

type MemoryStore struct {
	products []Product
	nextID   int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		nextID: 1,
	}
}

func (s *MemoryStore) List() []Product {
	return s.products
}

func (s *MemoryStore) Create(p Product) Product {
	p.ID = s.nextID
	s.nextID++
	s.products = append(s.products, p)
	return p
}
