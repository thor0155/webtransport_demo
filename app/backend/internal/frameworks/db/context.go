package db

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

type Context struct {
	context.Context
	gorm          *gorm.DB
	mongodbClient *mongo.Client
	redisClient   redis.UniversalClient
}

type gormKey struct{}
type mongodbKey struct{}
type redisKey struct{}

func NewContext(parent context.Context) *Context {
	return &Context{
		Context: parent,
	}
}

func (c *Context) WithGorm(db *gorm.DB) *Context {
	c.gorm = db.WithContext(c.Context)
	return c
}

func (c *Context) WithMongodb(mongodb *mongo.Client) *Context {
	c.mongodbClient = mongodb
	return c
}

func (c *Context) WithRedis(redisClient redis.UniversalClient) *Context {
	c.redisClient = redisClient
	return c
}

func (c *Context) GetGorm() *gorm.DB {
	return c.gorm
}

func (c *Context) GetMongodb() *mongo.Client {
	return c.mongodbClient
}

func (c *Context) GetRedis() redis.UniversalClient {
	return c.redisClient
}

func WithGorm(parent context.Context, db *gorm.DB) context.Context {
	return context.WithValue(parent, gormKey{}, db.WithContext(parent))
}

func WithMongodb(parent context.Context, mongodb *mongo.Client) context.Context {
	return context.WithValue(parent, mongodbKey{}, mongodb)
}

func WithRedis(parent context.Context, redisClient redis.UniversalClient) context.Context {
	return context.WithValue(parent, redisKey{}, redisClient)
}

func GetGorm(ctx context.Context) *gorm.DB {
	if dbCtx, ok := ctx.(*Context); ok {
		return dbCtx.gorm
	} else {
		return ctx.Value(gormKey{}).(*gorm.DB)
	}
}

func GetMongodb(ctx context.Context) *mongo.Client {
	if dbCtx, ok := ctx.(*Context); ok {
		return dbCtx.mongodbClient
	} else {
		return ctx.Value(mongodbKey{}).(*mongo.Client)
	}
}

func GetRedis(ctx context.Context) redis.UniversalClient {
	if dbCtx, ok := ctx.(*Context); ok {
		return dbCtx.redisClient
	} else {
		return ctx.Value(redisKey{}).(redis.UniversalClient)
	}
}
