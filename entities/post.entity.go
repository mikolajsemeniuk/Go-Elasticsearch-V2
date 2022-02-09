package entities

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Id      uuid.UUID  `json:"id"`
	Title   string     `json:"title"`
	Done    bool       `json:"done"`
	Created time.Time  `json:"created"`
	Updated *time.Time `json:"updated"`
}
