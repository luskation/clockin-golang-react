package domain

import "time"

type Company struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CNPJ      string    `json:"cnpj"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
