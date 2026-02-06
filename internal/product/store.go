package product

import (
	"context"
	"fmt"
)

type Store interface {
	List(ctx context.Context) ([]Product, error)
	Create(ctx context.Context, p Product) (Product, error)
	GetByID(ctx context.Context, id int) (Product, error)
	Update(ctx context.Context, id int, p Product) (Product, error)
	Delete(ctx context.Context, id int) error
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

func (s *MemoryStore) List(ctx context.Context) ([]Product, error) {
	select {
    case <-ctx.Done():
        return nil, ctx.Err()
    default:
	}

	return s.products, nil
}

func (s *MemoryStore) Create(ctx context.Context, p Product) (Product, error) {
	select {
    case <-ctx.Done():
        return Product{}, ctx.Err()
    default:
    }

	p.ID = s.nextID
	s.nextID++
	s.products = append(s.products, p)
	return p, nil
}

func (s *MemoryStore) GetByID(ctx context.Context, id int) (Product, error) {
	select {
    case <-ctx.Done():
        return Product{}, ctx.Err()
    default:
    }

	for _, p := range s.products {
		if p.ID == id {
			return p, nil
		}
	}
	return Product{}, fmt.Errorf("product %d not found", id)
}

func (s *MemoryStore) Update(ctx context.Context, id int, updated Product) (Product, error) {
	select {
    case <-ctx.Done():
        return Product{}, ctx.Err()
    default:
    }

	for i, p := range s.products {
		if p.ID == id {
			updated.ID = id
			s.products[i] = updated
			return updated, nil
		}
	}
	return Product{}, fmt.Errorf("product %d not found", id)
}

func (s *MemoryStore) Delete(ctx context.Context, id int) error {
	select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

	for i, p := range s.products {
		if p.ID == id {
			s.products = append(s.products[:i], s.products[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("product %d not found", id)
}