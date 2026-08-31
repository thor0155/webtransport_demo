package db

import (
	"api/internal/frameworks/obj"
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisDB interface {
	obj.Component
	obj.Phaseable
	obj.Closeable
	Client() redis.UniversalClient
}

type redisDB struct {
	logger *zap.Logger
	cfg    *RedisConfig
	client redis.UniversalClient
	name   string
}

func NewRedisDB(cfg *RedisConfig, logger *zap.Logger) RedisDB {
	name := "redis"
	return &redisDB{
		name:   name,
		logger: logger.Named(name),
		cfg:    cfg,
	}
}

// Phase implements [RedisDB].
func (r *redisDB) Phase() int {
	return DatabasePhase
}

func (r *redisDB) Name() string {
	return r.name
}

func (r *redisDB) Active() bool {
	return r.cfg.Active
}

func (r *redisDB) Client() redis.UniversalClient {
	return r.client
}

func (r *redisDB) Init() error {
	var err error
	switch r.cfg.Mode {
	case RedisModeSingle:
		if r.client, err = r.initSingle(); err != nil {
			return err
		}
	case RedisModeCluster:
		if r.client, err = r.initCluster(); err != nil {
			return err
		}
	case RedisModeSentinel:
		if r.client, err = r.initSentinel(); err != nil {
			return err
		}
	default:
		return errors.Errorf("unknown redis mode '%v'", r.cfg.Mode)
	}

	var pong string
	if pong, err = r.client.Ping(context.Background()).Result(); err != nil {
		return errors.WithStack(err)
	} else if pong != "PONG" {
		return errors.Errorf("connecting refuse, conn state is '%v'", pong)
	}
	return nil
}

func (r *redisDB) Close(ctx context.Context) error {
	return r.client.Shutdown(ctx).Err()
}

func (r *redisDB) initSentinel() (*redis.Client, error) {
	var err error
	options := &redis.FailoverOptions{
		MasterName:    r.cfg.MasterName,
		SentinelAddrs: r.cfg.Endpoints,
		Username:      r.cfg.Username,
		Password:      r.cfg.Password,
		DB:            r.cfg.DB,
	}
	if r.cfg.PoolSize > 0 {
		options.PoolSize = r.cfg.PoolSize
	}
	if r.cfg.MaxRetries > 0 {
		options.MaxRetries = r.cfg.MaxRetries
	}
	if r.cfg.MinIdleConns > 0 {
		options.MinIdleConns = r.cfg.MinIdleConns
	}
	if r.cfg.ConnMaxLifetime != "" {
		if options.ConnMaxLifetime, err = time.ParseDuration(r.cfg.ConnMaxLifetime); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if r.cfg.PoolTimeout != "" {
		if options.PoolTimeout, err = time.ParseDuration(r.cfg.PoolTimeout); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if r.cfg.ConnMaxIdleTime != "" {
		if options.ConnMaxIdleTime, err = time.ParseDuration(r.cfg.ConnMaxIdleTime); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	return redis.NewFailoverClient(options), nil

}

func (r *redisDB) initSingle() (*redis.Client, error) {
	var err error
	options := &redis.Options{
		Addr:     r.cfg.Endpoints[0],
		Username: r.cfg.Username,
		Password: r.cfg.Password,
		DB:       r.cfg.DB,
	}
	if r.cfg.PoolSize > 0 {
		options.PoolSize = r.cfg.PoolSize
	}
	if r.cfg.MaxRetries > 0 {
		options.MaxRetries = r.cfg.MaxRetries
	}
	if r.cfg.MinIdleConns > 0 {
		options.MinIdleConns = r.cfg.MinIdleConns
	}
	if r.cfg.ConnMaxLifetime != "" {
		if options.ConnMaxLifetime, err = time.ParseDuration(r.cfg.ConnMaxLifetime); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if r.cfg.PoolTimeout != "" {
		if options.PoolTimeout, err = time.ParseDuration(r.cfg.PoolTimeout); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if r.cfg.ConnMaxIdleTime != "" {
		if options.ConnMaxIdleTime, err = time.ParseDuration(r.cfg.ConnMaxIdleTime); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return redis.NewClient(options), nil

}

func (r *redisDB) initCluster() (*redis.ClusterClient, error) {
	var err error
	options := &redis.ClusterOptions{
		Addrs:    r.cfg.Endpoints,
		Username: r.cfg.Username,
		Password: r.cfg.Password,
	}
	if r.cfg.PoolSize > 0 {
		options.PoolSize = r.cfg.PoolSize
	}
	if r.cfg.MaxRetries > 0 {
		options.MaxRetries = r.cfg.MaxRetries
	}
	if r.cfg.MinIdleConns > 0 {
		options.MinIdleConns = r.cfg.MinIdleConns
	}
	if r.cfg.ConnMaxLifetime != "" {
		if options.ConnMaxLifetime, err = time.ParseDuration(r.cfg.ConnMaxLifetime); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if r.cfg.PoolTimeout != "" {
		if options.PoolTimeout, err = time.ParseDuration(r.cfg.PoolTimeout); err != nil {
			return nil, errors.WithStack(err)
		}
	}
	if r.cfg.ConnMaxIdleTime != "" {
		if options.ConnMaxIdleTime, err = time.ParseDuration(r.cfg.ConnMaxIdleTime); err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return redis.NewClusterClient(options), nil
}
