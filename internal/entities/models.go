package entities

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID        int        `json:"id" db:"id"`
	Name      string     `json:"service_name" db:"name"`
	Price     int        `json:"price" db:"price"`
	UserID    uuid.UUID  `json:"user_id" db:"user_id"`
	StartDate time.Time  `json:"start_date" db:"start_date"`
	EndDate   *time.Time `json:"end_date" db:"end_date"`
}

type FrontendSubscription struct {
	ID        int    `json:"id" db:"id"`
	Name      string `json:"service_name" db:"name"`
	Price     int    `json:"price" db:"price"`
	UserID    string `json:"user_id" db:"user_id"`
	StartDate string `json:"start_date" db:"start_date"`
	EndDate   string `json:"end_date" db:"end_date"`
}
