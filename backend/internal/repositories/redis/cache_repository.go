package redis

import (
	"chat-app/internal/models"
	"chat-app/internal/repositories/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type cacheRepository struct {
	client *redis.Client
}

// NewCacheRepository creates a new cache repository
func NewCacheRepository(client *redis.Client) interfaces.CacheRepository {
	return &cacheRepository{
		client: client,
	}
}

// SetRecentMessages caches recent messages for a room
func (r *cacheRepository) SetRecentMessages(ctx context.Context, roomID string, messages []*models.Message) error {
	key := fmt.Sprintf("room:%s:recent_messages", roomID)

	// Clear existing cached messages
	r.client.Del(ctx, key)

	// Add messages to sorted set with timestamp as score
	for _, msg := range messages {
		msgJSON, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		score := float64(msg.CreatedAt.Unix())
		err = r.client.ZAdd(ctx, key, redis.Z{
			Score:  score,
			Member: string(msgJSON),
		}).Err()
		if err != nil {
			return err
		}
	}

	// Set expiration to 1 hour
	r.client.Expire(ctx, key, time.Hour)

	return nil
}

// GetRecentMessages retrieves recent messages from cache
func (r *cacheRepository) GetRecentMessages(ctx context.Context, roomID string, limit int) ([]*models.Message, error) {
	key := fmt.Sprintf("room:%s:recent_messages", roomID)

	// Get messages in descending order (most recent first)
	result, err := r.client.ZRevRange(ctx, key, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]*models.Message, 0, len(result))
	for _, msgStr := range result {
		var msg models.Message
		err := json.Unmarshal([]byte(msgStr), &msg)
		if err != nil {
			continue // Skip malformed messages
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

// InvalidateRoomCache removes cached data for a room
func (r *cacheRepository) InvalidateRoomCache(ctx context.Context, roomID string) error {
	pattern := fmt.Sprintf("room:%s:*", roomID)
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Del(ctx, keys...).Err()
	}

	return nil
}

// SetRoomOnlineUsers sets the list of online users for a room
func (r *cacheRepository) SetRoomOnlineUsers(ctx context.Context, roomID string, userIDs []string) error {
	key := fmt.Sprintf("room:%s:online_users", roomID)

	// Clear existing set
	r.client.Del(ctx, key)

	if len(userIDs) > 0 {
		// Add users to set
		members := make([]interface{}, len(userIDs))
		for i, userID := range userIDs {
			members[i] = userID
		}

		err := r.client.SAdd(ctx, key, members...).Err()
		if err != nil {
			return err
		}

		// Set expiration to 30 minutes
		r.client.Expire(ctx, key, 30*time.Minute)
	}

	return nil
}

// GetRoomOnlineUsers retrieves the list of online users for a room
func (r *cacheRepository) GetRoomOnlineUsers(ctx context.Context, roomID string) ([]string, error) {
	key := fmt.Sprintf("room:%s:online_users", roomID)
	return r.client.SMembers(ctx, key).Result()
}

// AddUserToRoom adds a user to the room's online users set
func (r *cacheRepository) AddUserToRoom(ctx context.Context, roomID, userID string) error {
	key := fmt.Sprintf("room:%s:online_users", roomID)
	err := r.client.SAdd(ctx, key, userID).Err()
	if err != nil {
		return err
	}

	// Extend expiration
	r.client.Expire(ctx, key, 30*time.Minute)
	return nil
}

// RemoveUserFromRoom removes a user from the room's online users set
func (r *cacheRepository) RemoveUserFromRoom(ctx context.Context, roomID, userID string) error {
	key := fmt.Sprintf("room:%s:online_users", roomID)
	return r.client.SRem(ctx, key, userID).Err()
}
