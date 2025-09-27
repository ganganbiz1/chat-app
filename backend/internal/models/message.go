package models

import (
	"time"

	"gorm.io/gorm"
)

// Message represents a chat message
type Message struct {
	ID        string         `json:"id" gorm:"type:char(36);primary_key;default:(UUID())"`
	RoomID    string         `json:"room_id" gorm:"type:char(36);not null;index:idx_room_created,priority:1"`
	UserID    string         `json:"user_id" gorm:"type:varchar(255);not null;index"`
	Username  string         `json:"username" gorm:"type:varchar(255);not null"`
	Content   string         `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime;index:idx_room_created,priority:2"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Note: Room relationship removed to avoid circular imports
	// Use joins in queries when room data is needed
}

// TableName specifies the table name for Message model
func (Message) TableName() string {
	return "messages"
}

// CreateMessageRequest represents the request body for creating a message
type CreateMessageRequest struct {
	RoomID   string `json:"room_id" validate:"required,uuid"`
	UserID   string `json:"user_id" validate:"required,min=1,max=255"`
	Username string `json:"username" validate:"required,min=1,max=255"`
	Content  string `json:"content" validate:"required,min=1,max=5000"`
}

// MessageResponse represents the response for message operations
type MessageResponse struct {
	ID        string    `json:"id"`
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
