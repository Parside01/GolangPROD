package redisdb

import (
	"context"
	"encoding/json"
	"fmt"
	"solution/models"
	"time"
)

func (r *RedisDB) WriteSession(ctx context.Context, session *models.Session) error {
	name := getSessionNameByGUID(session.SessionGUID)

	if err := r.сlient.Set(ctx, name, session, time.Duration(r.sessionTTL)*time.Minute).Err(); err != nil {
		r.logger.Error("redisdb.WriteSession", "error", err)
		return err
	}
	return nil
}

func (r *RedisDB) UpdateSessionByGUID(ctx context.Context, session *models.Session, SessionGUID string) error {
	err := r.DeleteSessionByGUID(ctx, SessionGUID)
	if err != nil {
		r.logger.Error("redis.UpdateSessionByGUID: failed to delete session", "error", err)
		return err
	}

	err = r.WriteSession(ctx, session)
	if err != nil {
		r.logger.Error("redis.UpdateSessionByGUID: failed to write session", "error", err)
		return err
	}
	return nil
}

func (r *RedisDB) FindSessionByGUID(ctx context.Context, SessionGUID string) (*models.Session, error) {
	name := getSessionNameByGUID(SessionGUID)

	res, err := r.сlient.Get(ctx, name).Result()
	if err != nil {
		r.logger.Error("redisdb.FindSessionByID: failed to find session by id", "error", err)
		return nil, err
	}
	var session *models.Session
	err = json.Unmarshal([]byte(res), &session)
	if err != nil {
		r.logger.Error("redisdb.FindSessionByID: failed to decode session", "error", err)
		return nil, err
	}

	return session, nil
}

func (r *RedisDB) DeleteSessionByGUID(ctx context.Context, SessionGUID string) error {
	name := getSessionNameByGUID(SessionGUID)
	_, err := r.сlient.Del(ctx, name).Result()
	if err != nil {
		r.logger.Error("redisdb.DeleteSession: failed to delete session", "error", err)
		return err
	}
	return nil
}

func getSessionNameByGUID(SessionGUID string) string {
	return fmt.Sprintf("session-%s", SessionGUID)
}
