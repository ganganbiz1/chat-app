package mysql

import (
	"chat-app/internal/models"
	"chat-app/internal/repositories/interfaces"
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type messageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) interfaces.MessageRepository {
	return &messageRepository{
		db: db,
	}
}

// Create creates a new message
func (r *messageRepository) Create(ctx context.Context, message *models.Message) error {
	if message.ID == "" {
		message.ID = uuid.New().String()
	}

	return r.db.WithContext(ctx).Create(message).Error
}

// GetByID retrieves a message by ID
func (r *messageRepository) GetByID(ctx context.Context, id string) (*models.Message, error) {
	var message models.Message
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&message).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, models.ErrMessageNotFound
		}
		return nil, err
	}
	return &message, nil
}

// GetByRoomID retrieves messages by room ID with pagination
func (r *messageRepository) GetByRoomID(ctx context.Context, roomID string, limit, offset int) ([]*models.Message, error) {
	var messages []*models.Message
	query := r.db.WithContext(ctx).Where("room_id = ?", roomID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&messages).Error
	return messages, err
}

// GetByUserID retrieves messages by user ID with pagination
func (r *messageRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*models.Message, error) {
	var messages []*models.Message
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&messages).Error
	return messages, err
}

// Delete deletes a message
func (r *messageRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Message{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return models.ErrMessageNotFound
	}
	return nil
}

// CountByRoomID returns the total number of messages in a room
func (r *messageRepository) CountByRoomID(ctx context.Context, roomID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Message{}).Where("room_id = ?", roomID).Count(&count).Error
	return count, err
}
