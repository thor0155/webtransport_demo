package db

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
)

type Context struct {
	context.Context
	gorm          *gorm.DB
	mongodbClient *mongo.Client
}

type gormKey struct{}
type mongodbKey struct{}

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

func (c *Context) GetGorm() *gorm.DB {
	return c.gorm
}

func (c *Context) GetMongodb() *mongo.Client {
	return c.mongodbClient
}

func WithGorm(parent context.Context, db *gorm.DB) context.Context {
	return context.WithValue(parent, gormKey{}, db.WithContext(parent))
}

func WithMongodb(parent context.Context, mongodb *mongo.Client) context.Context {
	return context.WithValue(parent, mongodbKey{}, mongodb)
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
