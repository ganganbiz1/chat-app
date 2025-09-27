package interfaces

import (
	"chat-app/internal/models"
	"context"
)

// RoomRepository defines the interface for room data operations
type RoomRepository interface {
	Create(ctx context.Context, room *models.Room) error
	GetByID(ctx context.Context, id string) (*models.Room, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Room, error)
	Update(ctx context.Context, id string, room *models.Room) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
}
