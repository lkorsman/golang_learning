package auth

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	Create(ctx context.Context, email, password string) (User, error) 
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id int) (User, error)
}

type MemoryUserStore struct {
	users map[int]User
	emails map[string]int
	nextID int
	mu sync.RWMutex
}

func NewMemoryUserStore() *MemoryUserStore {
	return &MemoryUserStore{
		users: make(map[int]User),
		emails: make(map[string]int),
		nextID: 1,
	}
}

func (s *MemoryUserStore) Create(ctx context.Context, email, password string) (User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.emails[email]; exists {
		return User{}, fmt.Errorf("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err 
	}

	user := User{
		ID: s.nextID,
		Email: email,
		Password: string(hashedPassword),
	}

	s.users[user.ID] = user
	s.emails[email] = user.ID
	s.nextID++

	return user, nil
}

func (s *MemoryUserStore) GetByEmail(ctx context.Context, email string) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	userID, exists := s.emails[email]
	if !exists {
		return User{}, fmt.Errorf("user not found")
	}

	return s.users[userID], nil
}

func (s *MemoryUserStore) GetByID(ctx context.Context, id int) (User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.users[id]
	if !exists {
		return User{}, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *MemoryUserStore) ValidatePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}