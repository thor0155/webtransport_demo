package users

import (
	"api/internal/dao/userd"
	"api/internal/frameworks/db"
	"api/internal/frameworks/utils/cachetime"
	"api/internal/utils/types"
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
)

const (
	nameCacheKeyPrefix = "user:name:"
	nameCacheBaseTTL   = 6 * time.Hour
	nameCacheJitter    = 2 * time.Hour
)

type NameResolver interface {
	ResolveNames(ctx context.Context, userIds ...string) (map[types.UserId]types.UserName, error)
	Invalidate(ctx context.Context, userId string) error
}

type redisNameResolver struct {
	logger  *zap.Logger
	userDao userd.UserDao
}

func NewRedisNameResolver(logger *zap.Logger, userDao userd.UserDao) NameResolver {
	return &redisNameResolver{
		logger:  logger.Named("name-resolver"),
		userDao: userDao}
}

func (r *redisNameResolver) ResolveNames(ctx context.Context, userIds ...string) (map[types.UserId]types.UserName, error) {

	gorm := db.GetGorm(ctx)
	redisClient := db.GetRedis(ctx)

	result := make(map[string]string, len(userIds))
	keys := make([]string, len(userIds))
	for i, id := range userIds {
		keys[i] = r.fetchKey(id)
	}

	// 一次 MGET 批次拿,而不是逐筆 GET
	values, err := redisClient.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var missing []string
	for i, v := range values {
		if v == nil {
			missing = append(missing, userIds[i])
			continue
		}
		name, ok := v.(string)
		if !ok {
			missing = append(missing, userIds[i])
			continue
		}
		result[userIds[i]] = name
	}

	if len(missing) == 0 {
		return result, nil
	}

	fetched, err := r.userDao.GetNameMapByIds(gorm, missing)
	if err != nil {
		return result, errors.WithStack(err)
	}

	// 批次寫回 Redis,用 pipeline 減少往返次數
	pipe := redisClient.Pipeline()
	for id, name := range fetched {
		result[id] = name
		pipe.Set(ctx, r.fetchKey(id), name, cachetime.WithJitter(nameCacheBaseTTL, nameCacheJitter))
	}
	if _, err := pipe.Exec(ctx); err != nil {
		// 寫快取失敗不該讓整個查詢失敗,查到的名稱還是要回傳,只是下次快取沒命中
		// 記 log 即可
		r.logger.Error("write name cache failed", zap.Error(err))
	}
	return result, nil
}

// Invalidate 在使用者改名時呼叫,清掉舊快取,讓下次查詢重新從 DB 取得新名稱
func (r *redisNameResolver) Invalidate(ctx context.Context, userId string) error {
	return db.GetRedis(ctx).Del(ctx, r.fetchKey(userId)).Err()
}

func (*redisNameResolver) fetchKey(id string) string {
	return nameCacheKeyPrefix + id
}
