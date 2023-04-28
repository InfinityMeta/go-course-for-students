package ads

import "time"

type Ad struct {
	ID           int64
	Title        string `validate:"max:100"`
	Text         string `validate:"max:500"`
	AuthorID     int64
	Published    bool
	CreationDate time.Time
	UpdateDate   time.Time
}
