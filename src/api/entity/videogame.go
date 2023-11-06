// Package entities holds all the entities that are shared across all subdomains
package entities

import (
	"github.com/google/uuid"
)

// VideoGame represents a video game in the system
type VideoGame struct {
	ID       uuid.UUID
	Name     string
	Platform string
}
