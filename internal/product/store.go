package product

import (
	"fmt"
)

type Store interface {
	List() []Product
	Create(Product) Product
	GetByID(id int) (Product, error)
	Update(id int, p Product) (Product, error)
	Delete(id int) error
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

func (s *MemoryStore) GetByID(id int) (Product, error) {
	for _, p := range s.products {
		if p.ID == id {
			return p, nil
		}
	}
	return Product{}, fmt.Errorf("product %d not found", id)
}

func (s *MemoryStore) Update(id int, updated Product) (Product, error) {
	for i, p := range s.products {
		if p.ID == id {
			updated.ID = id
			s.products[i] = updated
			return updated, nil
		}
	}
	return Product{}, fmt.Errorf("product %d not found", id)
}

func (s *MemoryStore) Delete(id int) error {
	for i, p := range s.products {
		if p.ID == id {
			s.products = append(s.products[:i], s.products[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("product %d not found", id)
}