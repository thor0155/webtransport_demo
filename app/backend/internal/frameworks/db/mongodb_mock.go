package db

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type mockMongoDB struct {
	client          *mongo.Client
	mongoC          *mongodb.MongoDBContainer
	name            string
	defaultDatabase *mongo.Database
}

// Database implements [MongoDB].
func (m *mockMongoDB) Database() *mongo.Database {
	return m.defaultDatabase
}

// Close implements [MongoDB].
func (m *mockMongoDB) Close(ctx context.Context) error {
	if m.client == nil {
		return nil
	}
	defer m.mongoC.Terminate(ctx)
	if err := errors.WithStack(m.client.Disconnect(ctx)); err != nil {
		return err
	}
	return nil
}

// Active implements [MongoDB].
func (m *mockMongoDB) Active() bool {
	return true
}

func (m *mockMongoDB) Client() *mongo.Client {
	return m.client
}

func (m *mockMongoDB) Phase() int {
	return 1
}

// Init implements [MongoDB].
func (m *mockMongoDB) Init() error {
	ctx := context.Background()
	var (
		uri string
		err error
	)
	if m.mongoC, err = mongodb.Run(ctx, "mongo:8"); err != nil {
		return err
	}
	if uri, err = m.mongoC.ConnectionString(ctx); err != nil {
		return err
	}

	m.client, err = mongo.Connect(options.Client().ApplyURI(uri))
	m.defaultDatabase = m.client.Database("test")
	return nil
}

// Name implements [MongoDB].
func (m *mockMongoDB) Name() string {
	return m.name
}

func NewMockMongoDB() MongoDB {
	name := "mock_mongodb"
	return &mockMongoDB{
		name: name,
	}
}
