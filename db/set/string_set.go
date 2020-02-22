package set

type StringSet interface {
	Put(value string) (added bool)
	Contains(value string) bool
	Remove(value string) (removed bool)
	Len() int
}

func NewStringSet() StringSet {
	return &mapStringSet{
		storage: make(map[string]uint8),
	}
}

type mapStringSet struct {
	storage map[string]uint8
}

func (s *mapStringSet) Put(value string) (added bool) {
	contains := s.Contains(value)
	s.storage[value] = 1
	return !contains
}

func (s *mapStringSet) Contains(value string) bool {
	_, contains := s.storage[value]
	return contains
}

func (s *mapStringSet) Remove(value string) (removed bool) {
	if s.Contains(value) {
		delete(s.storage, value)
		return true
	} else {
		return false
	}
}

func (s *mapStringSet) Len() int {
	return len(s.storage)
}
