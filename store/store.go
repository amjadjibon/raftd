package store

import (
	"errors"
	"sync"
)

type Store struct {
	mu sync.Mutex
	kv map[string][]byte
}

func New() *Store {
	return &Store{
		kv: make(map[string][]byte),
	}
}

func (s *Store) Set(key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.kv[key] = value
	return nil
}

func (s *Store) Get(key string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	value, ok := s.kv[key]
	if !ok {
		return nil, errors.New("key not found")
	}
	return value, nil
}

func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.kv, key)
	return nil
}

func (s *Store) Keys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	keys := make([]string, 0, len(s.kv))
	for key := range s.kv {
		keys = append(keys, key)
	}
	return keys
}
