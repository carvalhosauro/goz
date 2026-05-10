package domain

type Priority string

const (
	PNone Priority = ""
	P1    Priority = "P1"
	P2    Priority = "P2"
	P3    Priority = "P3"
)

func (p Priority) Valid() bool {
	switch p {
	case PNone, P1, P2, P3:
		return true
	}
	return false
}

func ParsePriority(s string) (Priority, bool) {
	p := Priority(s)
	return p, p.Valid()
}
