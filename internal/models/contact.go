package models

import (
	"github.com/docker/distribution/uuid"
	"time"
)

type ContactRequest struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	Email       string
	Description string
	Status      string
	CreatedAt   time.Time
}
