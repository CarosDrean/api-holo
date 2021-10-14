package model

import "time"

type Service struct {
	ID               string    `json:"id"`
	ProtocolID       string    `json:"protocol_id"`
	PersonID         string    `json:"person_id"`
	StatusID         int       `json:"status_id"`
	ServiceDate      string    `json:"service_date"`
	AptitudeStatusID int       `json:"aptitude_status_id"`
	OrganizationID   string    `json:"organization_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type Services []Service
