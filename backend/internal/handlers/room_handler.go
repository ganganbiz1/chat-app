package handlers

import (
	"chat-app/internal/models"
	"chat-app/internal/services"
	"chat-app/pkg/validator"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type RoomHandler struct {
	roomService services.RoomService
}

// NewRoomHandler creates a new room handler
func NewRoomHandler(roomService services.RoomService) *RoomHandler {
	return &RoomHandler{
		roomService: roomService,
	}
}

// CreateRoom handles POST /api/rooms
func (h *RoomHandler) CreateRoom(c echo.Context) error {
	var req models.CreateRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	room, err := h.roomService.CreateRoom(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create room",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusCreated, models.SuccessResponse{
		Data:    room,
		Message: "Room created successfully",
	})
}

// GetRoom handles GET /api/rooms/:id
func (h *RoomHandler) GetRoom(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Room ID is required",
			Code:  http.StatusBadRequest,
		})
	}

	room, err := h.roomService.GetRoom(c.Request().Context(), id)
	if err != nil {
		if err == models.ErrRoomNotFound {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Room not found",
				Code:  http.StatusNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get room",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Data: room,
	})
}

// GetRooms handles GET /api/rooms
func (h *RoomHandler) GetRooms(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	rooms, pagination, err := h.roomService.GetRooms(c.Request().Context(), page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get rooms",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	pagination.Data = rooms
	return c.JSON(http.StatusOK, pagination)
}

// UpdateRoom handles PUT /api/rooms/:id
func (h *RoomHandler) UpdateRoom(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Room ID is required",
			Code:  http.StatusBadRequest,
		})
	}

	var req models.UpdateRoomRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	if err := validator.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	room, err := h.roomService.UpdateRoom(c.Request().Context(), id, &req)
	if err != nil {
		if err == models.ErrRoomNotFound {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Room not found",
				Code:  http.StatusNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to update room",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Data:    room,
		Message: "Room updated successfully",
	})
}

// DeleteRoom handles DELETE /api/rooms/:id
func (h *RoomHandler) DeleteRoom(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Room ID is required",
			Code:  http.StatusBadRequest,
		})
	}

	err := h.roomService.DeleteRoom(c.Request().Context(), id)
	if err != nil {
		if err == models.ErrRoomNotFound {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Room not found",
				Code:  http.StatusNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete room",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Room deleted successfully",
	})
}
