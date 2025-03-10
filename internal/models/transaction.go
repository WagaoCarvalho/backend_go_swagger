package models

import "time"

type Transaction struct {
	UID       uint64    `json:"id"`
	Origin    Wallet    `json:"origin"`
	Target    Wallet    `json:"target"`
	Cash      float32   `json:"cash"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
