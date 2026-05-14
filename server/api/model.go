package api

import "time"

type Agent struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	LastSeen time.Time `json:"last_seen"`
	Online   bool      `json:"online"`
}
