package mysql

import (
	"chat-app/internal/models"
	"chat-app/internal/repositories/interfaces"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type roomRepository struct {
	db *gorm.DB
}

// NewRoomRepository creates a new room repository
func NewRoomRepository(db *gorm.DB) interfaces.RoomRepository {
	return &roomRepository{
		db: db,
	}
}

// Create creates a new room
func (r *roomRepository) Create(ctx context.Context, room *models.Room) error {
	if room.ID == "" {
		room.ID = uuid.New().String()
	}

	return r.db.WithContext(ctx).Create(room).Error
}

// GetByID retrieves a room by ID
func (r *roomRepository) GetByID(ctx context.Context, id string) (*models.Room, error) {
	var room models.Room
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&room).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrRoomNotFound
		}
		return nil, err
	}
	return &room, nil
}

// GetAll retrieves all rooms with pagination
func (r *roomRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Room, error) {
	var rooms []*models.Room
	query := r.db.WithContext(ctx).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&rooms).Error
	return rooms, err
}

// Update updates a room
func (r *roomRepository) Update(ctx context.Context, id string, room *models.Room) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Updates(room)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrRoomNotFound
	}
	return nil
}

// Delete deletes a room
func (r *roomRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Room{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrRoomNotFound
	}
	return nil
}

// Count returns the total number of rooms
func (r *roomRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Room{}).Count(&count).Error
	return count, err
}
