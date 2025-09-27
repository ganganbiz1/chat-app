package interfaces

import (
	"chat-app/internal/models"
	"context"
)

// MessageRepository defines the interface for message data operations
type MessageRepository interface {
	Create(ctx context.Context, message *models.Message) error
	GetByID(ctx context.Context, id string) (*models.Message, error)
	GetByRoomID(ctx context.Context, roomID string, limit, offset int) ([]*models.Message, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.Message, error)
	Delete(ctx context.Context, id string) error
	CountByRoomID(ctx context.Context, roomID string) (int64, error)
}

// CacheRepository defines the interface for cache operations
type CacheRepository interface {
	SetRecentMessages(ctx context.Context, roomID string, messages []*models.Message) error
	GetRecentMessages(ctx context.Context, roomID string, limit int) ([]*models.Message, error)
	InvalidateRoomCache(ctx context.Context, roomID string) error
	SetRoomOnlineUsers(ctx context.Context, roomID string, userIDs []string) error
	GetRoomOnlineUsers(ctx context.Context, roomID string) ([]string, error)
	AddUserToRoom(ctx context.Context, roomID, userID string) error
	RemoveUserFromRoom(ctx context.Context, roomID, userID string) error
}
