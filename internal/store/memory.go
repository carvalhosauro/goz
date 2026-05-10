package store

import (
	"crypto/rand"
	"encoding/hex"
	"sort"
	"sync"

	"goz/internal/domain"
)

// Memory is an in-memory Store backed by a map + ordered slice for manual sort.
type Memory struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
	order []string
}

func NewMemory(seed []domain.Task) *Memory {
	m := &Memory{tasks: make(map[string]*domain.Task, len(seed))}
	for i := range seed {
		t := seed[i]
		m.tasks[t.ID] = &t
		m.order = append(m.order, t.ID)
	}
	m.normalize()
	return m
}

func (m *Memory) normalize() {
	sort.SliceStable(m.order, func(i, j int) bool {
		return m.tasks[m.order[i]].SortOrder < m.tasks[m.order[j]].SortOrder
	})
	for i, id := range m.order {
		m.tasks[id].SortOrder = i
	}
}

func (m *Memory) List() []domain.Task {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]domain.Task, 0, len(m.order))
	for _, id := range m.order {
		out = append(out, *m.tasks[id])
	}
	return out
}

func (m *Memory) Get(id string) (domain.Task, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.tasks[id]
	if !ok {
		return domain.Task{}, false
	}
	return *t, true
}

func (m *Memory) Add(text string) (domain.Task, error) {
	if text == "" {
		return domain.Task{}, domain.ErrEmptyText
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	id := newID()
	t := domain.Task{ID: id, Text: text, SortOrder: len(m.order)}
	m.tasks[id] = &t
	m.order = append(m.order, id)
	return t, nil
}

func (m *Memory) Toggle(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	t, ok := m.tasks[id]
	if !ok {
		return ErrNotFound
	}
	t.Done = !t.Done
	return nil
}

func (m *Memory) Delete(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(m.tasks, id)
	for i, x := range m.order {
		if x == id {
			m.order = append(m.order[:i], m.order[i+1:]...)
			break
		}
	}
	m.normalize()
	return nil
}

func (m *Memory) Move(id string, delta int) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.tasks[id]; !ok {
		return ErrNotFound
	}
	idx := -1
	for i, x := range m.order {
		if x == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return ErrNotFound
	}
	next := idx + delta
	if next < 0 {
		next = 0
	}
	if next > len(m.order)-1 {
		next = len(m.order) - 1
	}
	if next == idx {
		return nil
	}

	without := make([]string, 0, len(m.order)-1)
	without = append(without, m.order[:idx]...)
	without = append(without, m.order[idx+1:]...)

	rebuilt := make([]string, 0, len(m.order))
	rebuilt = append(rebuilt, without[:next]...)
	rebuilt = append(rebuilt, id)
	rebuilt = append(rebuilt, without[next:]...)
	m.order = rebuilt

	for i, oid := range m.order {
		m.tasks[oid].SortOrder = i
	}
	return nil
}

func newID() string {
	var b [4]byte
	_, _ = rand.Read(b[:])
	return "n" + hex.EncodeToString(b[:])
}
