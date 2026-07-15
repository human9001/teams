package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type TaskListCache struct {
	rdb *redis.Client
	ttl time.Duration
}

func NewTaskListCache(rdb *redis.Client, ttl time.Duration) *TaskListCache {
	return &TaskListCache{rdb: rdb, ttl: ttl}
}

func (c *TaskListCache) versionKey(teamID int64) string {
	return fmt.Sprintf("tasks:version:team:%d", teamID)
}

func (c *TaskListCache) ensureVersion(ctx context.Context, teamID int64) (int64, error) {
	key := c.versionKey(teamID)
	v, err := c.rdb.Get(ctx, key).Int64()
	if errors.Is(err, redis.Nil) {
		if err := c.rdb.Set(ctx, key, 1, 0).Err(); err != nil {
			return 0, err
		}
		return 1, nil
	}
	return v, err
}

func (c *TaskListCache) listKey(teamID, version int64, status, assigneeID string, page, limit int64) string {
	return fmt.Sprintf(
		"tasks:list:v%d:team:%d:status:%s:assignee:%s:page:%d:limit:%d",
		version, teamID, status, assigneeID, page, limit,
	)
}

func (c *TaskListCache) Get(ctx context.Context, teamID int64, status, assigneeID string, page, limit int64, dst any) (bool, error) {
	version, err := c.ensureVersion(ctx, teamID)
	if err != nil {
		return false, err
	}

	key := c.listKey(teamID, version, status, assigneeID, page, limit)
	b, err := c.rdb.Get(ctx, key).Bytes()
	if errors.Is(err, redis.Nil) {
		slog.Info("task from cache", "info", "empty")
		return false, nil
	}
	if err != nil {
		slog.Error("task from cache", "error", err)
		return false, err
	}

	return true, json.Unmarshal(b, dst)
}

func (c *TaskListCache) Set(ctx context.Context, teamID int64, status, assigneeID string, page, limit int64, value any) error {
	version, err := c.ensureVersion(ctx, teamID)
	if err != nil {
		return err
	}

	key := c.listKey(teamID, version, status, assigneeID, page, limit)
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.rdb.Set(ctx, key, b, c.ttl).Err()
}

func (c *TaskListCache) BumpVersion(ctx context.Context, teamID int64) (int64, error) {
	return c.rdb.Incr(ctx, c.versionKey(teamID)).Result()
}
