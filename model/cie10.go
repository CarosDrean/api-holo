package model

import "time"

type Cie10 struct {
	ID             string    `json:"id"`
	Description    string    `json:"description"`
	DescriptionTwo string    `json:"description_two"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Cie10s []Cie10
