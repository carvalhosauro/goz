package domain

// Seed returns the fixture task list used in Phase 0 (in-memory only).
// Mirrors the design mock data, translated to English.
func Seed() []Task {
	return []Task{
		{ID: "t1", Text: "Review Marcos's PR #482 — auth flow", Priority: P1, Estimate: "20m", Due: "today", Tag: TagEng, SortOrder: 1},
		{ID: "t2", Text: "Refactor billing module to Stripe Checkout", Tag: TagEng, SortOrder: 2},
		{ID: "t3", Text: "Write ADR on storage v3", Priority: P2, Estimate: "1h", Due: "fri", Tag: TagDoc, SortOrder: 3},
		{ID: "t4", Text: "Call accountant about taxes", Priority: P1, Estimate: "15m", Due: "yesterday", Tag: TagPersonal, SortOrder: 4, Overdue: true},
		{ID: "t5", Text: "Buy Júlia's birthday gift (saturday)", Priority: P2, Estimate: "30m", Due: "thu", Tag: TagPersonal, SortOrder: 5},
		{ID: "t6", Text: "Investor deck for thursday meeting", Priority: P1, Estimate: "2h", Due: "wed", Tag: TagWork, SortOrder: 6},
		{ID: "t7", Text: "Onboard João (new eng)", Priority: P2, Estimate: "45m", Due: "today", Tag: TagPeople, SortOrder: 7},
		{ID: "t8", Text: "Schedule dentist", Priority: P3, Estimate: "5m", Tag: TagPersonal, SortOrder: 8},
		{ID: "t9", Text: "Reply HR email about benefits", Priority: P3, Estimate: "10m", Tag: TagWork, SortOrder: 9, Done: true},
		{ID: "t10", Text: "Study for AWS SA certification", Tag: TagLearn, SortOrder: 10},
	}
}
