package models

import "time"

type Account struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	PIN       string    `json:"-"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}
