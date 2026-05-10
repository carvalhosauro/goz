package tui

import (
	"time"

	"goz/internal/domain"
)

type tasksLoadedMsg []domain.Task

type taskAddedMsg struct{ Task domain.Task }

type taskMutatedMsg struct{}

type errMsg struct{ Err error }

type blinkMsg time.Time

type clockMsg time.Time
