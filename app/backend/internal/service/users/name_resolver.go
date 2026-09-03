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
	// Invalidate Call this function when a user changes their name, clear the old cache, and ensure the next query retrieves the new name from the database.
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

	r.logger.Debug("resolve names parameters", zap.Strings("userIds", userIds))

	result := make(map[string]string, len(userIds))
	keys := make([]string, len(userIds))
	for i, id := range userIds {
		keys[i] = r.fetchKey(id)
	}

	// Get data in batches using [MGET], instead of getting data one transaction at a time.
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

	r.logger.Debug("resolve names missing", zap.Strings("userIds", missing))

	fetched, err := r.userDao.GetNameMapByIds(gorm, missing)
	if err != nil {
		return result, errors.WithStack(err)
	}

	pipe := redisClient.Pipeline()
	for id, name := range fetched {
		result[id] = name
		pipe.Set(ctx, r.fetchKey(id), name, cachetime.WithJitter(nameCacheBaseTTL, nameCacheJitter))
	}
	if _, err := pipe.Exec(ctx); err != nil {
		// A failure to write to the cache shouldn't cause the entire query to fail;
		// the retrieved name should still be returned, it just won't be found in the next cache cache.
		r.logger.Error("write name cache failed", zap.Error(err))
	}
	return result, nil
}

func (r *redisNameResolver) Invalidate(ctx context.Context, userId string) error {
	return db.GetRedis(ctx).Del(ctx, r.fetchKey(userId)).Err()
}

func (*redisNameResolver) fetchKey(id string) string {
	return nameCacheKeyPrefix + id
}
