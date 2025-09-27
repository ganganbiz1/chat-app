package services

import (
	"chat-app/internal/models"
	"chat-app/internal/repositories/interfaces"
	"context"
)

// RoomService defines the interface for room business logic
type RoomService interface {
	CreateRoom(ctx context.Context, req *models.CreateRoomRequest) (*models.Room, error)
	GetRoom(ctx context.Context, id string) (*models.Room, error)
	GetRooms(ctx context.Context, page, limit int) ([]*models.Room, *models.PaginationResponse, error)
	UpdateRoom(ctx context.Context, id string, req *models.UpdateRoomRequest) (*models.Room, error)
	DeleteRoom(ctx context.Context, id string) error
}

type roomService struct {
	roomRepo  interfaces.RoomRepository
	cacheRepo interfaces.CacheRepository
}

// NewRoomService creates a new room service
func NewRoomService(roomRepo interfaces.RoomRepository, cacheRepo interfaces.CacheRepository) RoomService {
	return &roomService{
		roomRepo:  roomRepo,
		cacheRepo: cacheRepo,
	}
}

// CreateRoom creates a new room
func (s *roomService) CreateRoom(ctx context.Context, req *models.CreateRoomRequest) (*models.Room, error) {
	room := &models.Room{
		Name:        req.Name,
		Description: req.Description,
	}

	err := s.roomRepo.Create(ctx, room)
	if err != nil {
		return nil, err
	}

	return room, nil
}

// GetRoom retrieves a room by ID
func (s *roomService) GetRoom(ctx context.Context, id string) (*models.Room, error) {
	return s.roomRepo.GetByID(ctx, id)
}

// GetRooms retrieves rooms with pagination
func (s *roomService) GetRooms(ctx context.Context, page, limit int) ([]*models.Room, *models.PaginationResponse, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	offset := (page - 1) * limit

	rooms, err := s.roomRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	total, err := s.roomRepo.Count(ctx)
	if err != nil {
		return nil, nil, err
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	pagination := &models.PaginationResponse{
		Data:       rooms,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return rooms, pagination, nil
}

// UpdateRoom updates a room
func (s *roomService) UpdateRoom(ctx context.Context, id string, req *models.UpdateRoomRequest) (*models.Room, error) {
	// Check if room exists
	existingRoom, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != "" {
		existingRoom.Name = req.Name
	}
	if req.Description != "" {
		existingRoom.Description = req.Description
	}

	err = s.roomRepo.Update(ctx, id, existingRoom)
	if err != nil {
		return nil, err
	}

	// Invalidate cache
	s.cacheRepo.InvalidateRoomCache(ctx, id)

	return existingRoom, nil
}

// DeleteRoom deletes a room
func (s *roomService) DeleteRoom(ctx context.Context, id string) error {
	// Check if room exists
	_, err := s.roomRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete room
	err = s.roomRepo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache
	s.cacheRepo.InvalidateRoomCache(ctx, id)

	return nil
}
