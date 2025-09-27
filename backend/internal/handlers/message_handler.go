package handlers

import (
	"chat-app/internal/models"
	"chat-app/internal/services"
	"chat-app/pkg/validator"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type MessageHandler struct {
	messageService services.MessageService
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(messageService services.MessageService) *MessageHandler {
	return &MessageHandler{
		messageService: messageService,
	}
}

// CreateMessage handles POST /api/rooms/:id/messages
func (h *MessageHandler) CreateMessage(c echo.Context) error {
	roomID := c.Param("id")
	if roomID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Room ID is required",
			Code:  http.StatusBadRequest,
		})
	}

	var req models.CreateMessageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid request body",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	// Override room ID from URL parameter
	req.RoomID = roomID

	if err := validator.ValidateStruct(&req); err != nil {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
	}

	message, err := h.messageService.CreateMessage(c.Request().Context(), &req)
	if err != nil {
		if err == models.ErrRoomNotFound {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Room not found",
				Code:  http.StatusNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to create message",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	response := models.MessageResponse{
		ID:        message.ID,
		RoomID:    message.RoomID,
		UserID:    message.UserID,
		Username:  message.Username,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}

	return c.JSON(http.StatusCreated, models.SuccessResponse{
		Data:    response,
		Message: "Message created successfully",
	})
}

// GetRoomMessages handles GET /api/rooms/:id/messages
func (h *MessageHandler) GetRoomMessages(c echo.Context) error {
	roomID := c.Param("id")
	if roomID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Room ID is required",
			Code:  http.StatusBadRequest,
		})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}

	messages, pagination, err := h.messageService.GetRoomMessages(c.Request().Context(), roomID, page, limit)
	if err != nil {
		if err == models.ErrRoomNotFound {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Room not found",
				Code:  http.StatusNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get messages",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	// Convert to response format
	responses := make([]models.MessageResponse, len(messages))
	for i, msg := range messages {
		responses[i] = models.MessageResponse{
			ID:        msg.ID,
			RoomID:    msg.RoomID,
			UserID:    msg.UserID,
			Username:  msg.Username,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		}
	}

	pagination.Data = responses
	return c.JSON(http.StatusOK, pagination)
}

// GetRecentMessages handles GET /api/rooms/:id/messages/recent
func (h *MessageHandler) GetRecentMessages(c echo.Context) error {
	roomID := c.Param("id")
	if roomID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Room ID is required",
			Code:  http.StatusBadRequest,
		})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 {
		limit = 20
	}

	messages, err := h.messageService.GetRecentMessages(c.Request().Context(), roomID, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get recent messages",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	// Convert to response format
	responses := make([]models.MessageResponse, len(messages))
	for i, msg := range messages {
		responses[i] = models.MessageResponse{
			ID:        msg.ID,
			RoomID:    msg.RoomID,
			UserID:    msg.UserID,
			Username:  msg.Username,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		}
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Data: responses,
	})
}

// GetMessage handles GET /api/messages/:id
func (h *MessageHandler) GetMessage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Message ID is required",
			Code:  http.StatusBadRequest,
		})
	}

	message, err := h.messageService.GetMessage(c.Request().Context(), id)
	if err != nil {
		if err == models.ErrMessageNotFound {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Message not found",
				Code:  http.StatusNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get message",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	response := models.MessageResponse{
		ID:        message.ID,
		RoomID:    message.RoomID,
		UserID:    message.UserID,
		Username:  message.Username,
		Content:   message.Content,
		CreatedAt: message.CreatedAt,
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Data: response,
	})
}

// DeleteMessage handles DELETE /api/messages/:id
func (h *MessageHandler) DeleteMessage(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Message ID is required",
			Code:  http.StatusBadRequest,
		})
	}

	err := h.messageService.DeleteMessage(c.Request().Context(), id)
	if err != nil {
		if err == models.ErrMessageNotFound {
			return c.JSON(http.StatusNotFound, models.ErrorResponse{
				Error: "Message not found",
				Code:  http.StatusNotFound,
			})
		}
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to delete message",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	return c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Message deleted successfully",
	})
}

// GetUserMessages handles GET /api/users/:userId/messages
func (h *MessageHandler) GetUserMessages(c echo.Context) error {
	userID := c.Param("userId")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "User ID is required",
			Code:  http.StatusBadRequest,
		})
	}

	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 50
	}

	messages, pagination, err := h.messageService.GetUserMessages(c.Request().Context(), userID, page, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "Failed to get user messages",
			Message: err.Error(),
			Code:    http.StatusInternalServerError,
		})
	}

	// Convert to response format
	responses := make([]models.MessageResponse, len(messages))
	for i, msg := range messages {
		responses[i] = models.MessageResponse{
			ID:        msg.ID,
			RoomID:    msg.RoomID,
			UserID:    msg.UserID,
			Username:  msg.Username,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
		}
	}

	pagination.Data = responses
	return c.JSON(http.StatusOK, pagination)
}
