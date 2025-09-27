package services

import (
	"chat-app/internal/models"
	"chat-app/internal/repositories/interfaces"
	"context"
)

// MessageService defines the interface for message business logic
type MessageService interface {
	CreateMessage(ctx context.Context, req *models.CreateMessageRequest) (*models.Message, error)
	GetMessage(ctx context.Context, id string) (*models.Message, error)
	GetRoomMessages(ctx context.Context, roomID string, page, limit int) ([]*models.Message, *models.PaginationResponse, error)
	GetUserMessages(ctx context.Context, userID string, page, limit int) ([]*models.Message, *models.PaginationResponse, error)
	DeleteMessage(ctx context.Context, id string) error
	GetRecentMessages(ctx context.Context, roomID string, limit int) ([]*models.Message, error)
}

type messageService struct {
	messageRepo interfaces.MessageRepository
	roomRepo    interfaces.RoomRepository
	cacheRepo   interfaces.CacheRepository
}

// NewMessageService creates a new message service
func NewMessageService(
	messageRepo interfaces.MessageRepository,
	roomRepo interfaces.RoomRepository,
	cacheRepo interfaces.CacheRepository,
) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		roomRepo:    roomRepo,
		cacheRepo:   cacheRepo,
	}
}

// CreateMessage creates a new message
func (s *messageService) CreateMessage(ctx context.Context, req *models.CreateMessageRequest) (*models.Message, error) {
	// Verify room exists
	_, err := s.roomRepo.GetByID(ctx, req.RoomID)
	if err != nil {
		return nil, err
	}

	message := &models.Message{
		RoomID:   req.RoomID,
		UserID:   req.UserID,
		Username: req.Username,
		Content:  req.Content,
	}

	err = s.messageRepo.Create(ctx, message)
	if err != nil {
		return nil, err
	}

	// Update cache with recent messages
	go func() {
		// Get recent messages from DB
		recentMessages, err := s.messageRepo.GetByRoomID(context.Background(), req.RoomID, 50, 0)
		if err == nil {
			s.cacheRepo.SetRecentMessages(context.Background(), req.RoomID, recentMessages)
		}
	}()

	return message, nil
}

// GetMessage retrieves a message by ID
func (s *messageService) GetMessage(ctx context.Context, id string) (*models.Message, error) {
	return s.messageRepo.GetByID(ctx, id)
}

// GetRoomMessages retrieves messages for a room with pagination
func (s *messageService) GetRoomMessages(ctx context.Context, roomID string, page, limit int) ([]*models.Message, *models.PaginationResponse, error) {
	// Verify room exists
	_, err := s.roomRepo.GetByID(ctx, roomID)
	if err != nil {
		return nil, nil, err
	}

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}

	// Try to get from cache first for recent messages
	if page == 1 {
		cachedMessages, err := s.cacheRepo.GetRecentMessages(ctx, roomID, limit)
		if err == nil && len(cachedMessages) > 0 {
			total, _ := s.messageRepo.CountByRoomID(ctx, roomID)
			totalPages := int(total) / limit
			if int(total)%limit != 0 {
				totalPages++
			}

			pagination := &models.PaginationResponse{
				Data:       cachedMessages,
				Page:       page,
				Limit:      limit,
				Total:      total,
				TotalPages: totalPages,
			}

			return cachedMessages, pagination, nil
		}
	}

	offset := (page - 1) * limit

	messages, err := s.messageRepo.GetByRoomID(ctx, roomID, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	total, err := s.messageRepo.CountByRoomID(ctx, roomID)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	pagination := &models.PaginationResponse{
		Data:       messages,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	// Cache recent messages if this is the first page
	if page == 1 {
		go s.cacheRepo.SetRecentMessages(context.Background(), roomID, messages)
	}

	return messages, pagination, nil
}

// GetUserMessages retrieves messages for a user with pagination
func (s *messageService) GetUserMessages(ctx context.Context, userID string, page, limit int) ([]*models.Message, *models.PaginationResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}

	offset := (page - 1) * limit

	messages, err := s.messageRepo.GetByUserID(ctx, userID, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	// For simplicity, we're not implementing total count for user messages
	// In a real application, you might want to add this functionality
	pagination := &models.PaginationResponse{
		Data:       messages,
		Page:       page,
		Limit:      limit,
		Total:      int64(len(messages)),
		TotalPages: 1,
	}

	return messages, pagination, nil
}

// DeleteMessage deletes a message
func (s *messageService) DeleteMessage(ctx context.Context, id string) error {
	// Get message to find room ID for cache invalidation
	message, err := s.messageRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.messageRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache for the room
	go s.cacheRepo.InvalidateRoomCache(context.Background(), message.RoomID)

	return nil
}

// GetRecentMessages retrieves recent messages from cache
func (s *messageService) GetRecentMessages(ctx context.Context, roomID string, limit int) ([]*models.Message, error) {
	if limit <= 0 {
		limit = 20
	}

	// Try cache first
	messages, err := s.cacheRepo.GetRecentMessages(ctx, roomID, limit)
	if err == nil && len(messages) > 0 {
		return messages, nil
	}

	// Fallback to database
	messages, err = s.messageRepo.GetByRoomID(ctx, roomID, limit, 0)
	if err != nil {
		return nil, err
	}

	// Update cache
	go s.cacheRepo.SetRecentMessages(context.Background(), roomID, messages)

	return messages, nil
}
