package store

import (
	"errors"
	"testing"

	"goz/internal/domain"
)

func seed() []domain.Task {
	return []domain.Task{
		{ID: "a", Text: "alpha", SortOrder: 1},
		{ID: "b", Text: "beta", SortOrder: 2},
		{ID: "c", Text: "gamma", SortOrder: 3},
	}
}

func ids(tasks []domain.Task) []string {
	out := make([]string, len(tasks))
	for i, t := range tasks {
		out[i] = t.ID
	}
	return out
}

func eq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestList(t *testing.T) {
	m := NewMemory(seed())
	if got := ids(m.List()); !eq(got, []string{"a", "b", "c"}) {
		t.Fatalf("got %v", got)
	}
}

func TestAdd(t *testing.T) {
	m := NewMemory(seed())
	task, err := m.Add("delta")
	if err != nil {
		t.Fatal(err)
	}
	if task.Text != "delta" {
		t.Fatalf("text=%s", task.Text)
	}
	if got := len(m.List()); got != 4 {
		t.Fatalf("len=%d", got)
	}
}

func TestAddEmpty(t *testing.T) {
	m := NewMemory(seed())
	if _, err := m.Add(""); !errors.Is(err, domain.ErrEmptyText) {
		t.Fatalf("got %v", err)
	}
}

func TestToggle(t *testing.T) {
	m := NewMemory(seed())
	if err := m.Toggle("b"); err != nil {
		t.Fatal(err)
	}
	got, _ := m.Get("b")
	if !got.Done {
		t.Fatal("expected done")
	}
	_ = m.Toggle("b")
	got, _ = m.Get("b")
	if got.Done {
		t.Fatal("expected un-done")
	}
}

func TestToggleMissing(t *testing.T) {
	m := NewMemory(seed())
	if err := m.Toggle("zzz"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("got %v", err)
	}
}

func TestDelete(t *testing.T) {
	m := NewMemory(seed())
	if err := m.Delete("b"); err != nil {
		t.Fatal(err)
	}
	if got := ids(m.List()); !eq(got, []string{"a", "c"}) {
		t.Fatalf("got %v", got)
	}
}

func TestMove(t *testing.T) {
	cases := []struct {
		name  string
		id    string
		delta int
		want  []string
	}{
		{"down 1", "a", 1, []string{"b", "a", "c"}},
		{"up 1", "c", -1, []string{"a", "c", "b"}},
		{"clamp neg", "a", -10, []string{"a", "b", "c"}},
		{"clamp pos", "a", 10, []string{"b", "c", "a"}},
		{"zero", "b", 0, []string{"a", "b", "c"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m := NewMemory(seed())
			if err := m.Move(tc.id, tc.delta); err != nil {
				t.Fatal(err)
			}
			if got := ids(m.List()); !eq(got, tc.want) {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestGetMissing(t *testing.T) {
	m := NewMemory(seed())
	if _, ok := m.Get("zzz"); ok {
		t.Fatal("expected not found")
	}
}
