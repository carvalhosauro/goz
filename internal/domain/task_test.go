package domain

import (
	"errors"
	"testing"
)

func TestTaskValidate(t *testing.T) {
	cases := []struct {
		name    string
		task    Task
		wantErr error
	}{
		{"empty text", Task{}, ErrEmptyText},
		{"valid p1", Task{Text: "x", Priority: P1}, nil},
		{"valid no prio", Task{Text: "x", Priority: PNone}, nil},
		{"invalid prio", Task{Text: "x", Priority: "P9"}, ErrInvalidPriority},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.task.Validate()
			if !errors.Is(err, tc.wantErr) {
				t.Fatalf("got %v want %v", err, tc.wantErr)
			}
		})
	}
}

func TestPriorityValid(t *testing.T) {
	for _, p := range []Priority{PNone, P1, P2, P3} {
		if !p.Valid() {
			t.Fatalf("expected %q valid", p)
		}
	}
	if Priority("X").Valid() {
		t.Fatalf("X should be invalid")
	}
}

func TestSeedNotEmpty(t *testing.T) {
	s := Seed()
	if len(s) == 0 {
		t.Fatal("seed empty")
	}
	for _, x := range s {
		if err := x.Validate(); err != nil {
			t.Fatalf("seed task %s invalid: %v", x.ID, err)
		}
	}
}
