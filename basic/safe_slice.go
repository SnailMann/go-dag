package basic

import "sync"

type SafeSlice[T any] struct {
	data []T
	mu   *sync.Mutex
}

func NewSafeSlice[T any](cap int) *SafeSlice[T] {
	return &SafeSlice[T]{
		data: make([]T, 0, cap),
		mu:   &sync.Mutex{},
	}
}

func NewSafeSliceWithLen[T any](len, cap int) *SafeSlice[T] {
	return &SafeSlice[T]{
		data: make([]T, len, cap),
		mu:   &sync.Mutex{},
	}
}

func (s *SafeSlice[T]) Append(value T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = append(s.data, value)
}

func (s *SafeSlice[T]) AppendAll(value ...T) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data = append(s.data, value...)
}

func (s *SafeSlice[T]) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.data)
}

func (s *SafeSlice[T]) ToSlice() []T {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.data
}
