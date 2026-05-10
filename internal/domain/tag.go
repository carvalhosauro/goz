package domain

type Tag string

const (
	TagNone     Tag = ""
	TagEng      Tag = "eng"
	TagDoc      Tag = "doc"
	TagWork     Tag = "work"
	TagPersonal Tag = "personal"
	TagPeople   Tag = "people"
	TagLearn    Tag = "learn"
)

var KnownTags = []Tag{TagEng, TagDoc, TagWork, TagPersonal, TagPeople, TagLearn}
