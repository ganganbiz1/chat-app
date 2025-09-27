package models

import (
	"time"

	"gorm.io/gorm"
)

// Room represents a chat room
type Room struct {
	ID          string         `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	Name        string         `json:"name" gorm:"type:varchar(255);not null" validate:"required,min=1,max=255"`
	Description string         `json:"description" gorm:"type:text"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Note: Messages relationship is defined in Message model to avoid circular imports
}

// TableName specifies the table name for Room model
func (Room) TableName() string {
	return "rooms"
}

// CreateRoomRequest represents the request body for creating a room
type CreateRoomRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}

// UpdateRoomRequest represents the request body for updating a room
type UpdateRoomRequest struct {
	Name        string `json:"name" validate:"omitempty,min=1,max=255"`
	Description string `json:"description" validate:"max=1000"`
}
