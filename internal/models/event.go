package models

import (
	"errors"
	"time"
)

type EventStatus string

const (
	StatusDraft     EventStatus = "draft"
	StatusPublished EventStatus = "published"
	StatusOngoing   EventStatus = "ongoing"
	StatusCompleted EventStatus = "completed"
	StatusCancelled EventStatus = "cancelled"
	StatusPostponed EventStatus = "postponed"
)

func (s EventStatus) IsValid() bool {
	switch s {
	case StatusDraft, StatusPublished, StatusOngoing, StatusCompleted, StatusCancelled, StatusPostponed:
		return true
	default:
		return false
	}
}

type EventCreateRequest struct {
	Title        string      `json:"title" binding:"required,min=3,max=255"`
	About        string      `json:"description"`
	StartDate    time.Time   `json:"start_date" binding:"required"`
	Location     string      `json:"location"`
	Status       EventStatus `json:"status" binding:"required"`
	MaxAttendees int         `json:"max_attendees" binding:"gt=5,lte=100"`
}

func (r *EventCreateRequest) Validate() error {
	if !r.Status.IsValid() {
		return errors.New("invalid status")
	}
	return nil
}

type EventResponse struct {
	Id           int       `json:"id"`
	Title        string    `json:"title"`
	About        string    `json:"description"`
	StartDate    time.Time `json:"start_date"`
	Location     string    `json:"location"`
	Status       string    `json:"status"`
	MaxAttendees int       `json:"max_attendees"`
	Creator      string    `json:"creator"`
}

type EventUpdateRequest struct {
	Id           *int         `json:"id" binding:"required"`
	Title        *string      `json:"title"`
	About        *string      `json:"description"`
	StartDate    *time.Time   `json:"start_date"`
	Location     *string      `json:"location"`
	Status       *EventStatus `json:"status"`
	MaxAttendees *int         `json:"max_attendees"`
}
