package db

import (
	"api/internal/frameworks/obj"
	"api/internal/frameworks/utils"
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
	"go.uber.org/zap"
)

type MongoDB interface {
	obj.Component
	obj.Closeable
	obj.Phaseable
	Client() *mongo.Client
	Database() *mongo.Database
}

type mongoDB struct {
	logger          *zap.Logger
	client          *mongo.Client
	cfg             *MongoConfig
	name            string
	defaultDatabase *mongo.Database
}

// Database implements [MongoDB].
func (m *mongoDB) Database() *mongo.Database {
	return m.defaultDatabase
}

// Close implements [MongoDB].
func (m *mongoDB) Close(ctx context.Context) error {
	if m.client == nil {
		return nil
	}
	return errors.WithStack(m.client.Disconnect(ctx))
}

// Active implements [MongoDB].
func (m *mongoDB) Active() bool {
	return m.cfg.Active
}

func (m *mongoDB) Client() *mongo.Client {
	return m.client
}

func (m *mongoDB) Phase() int {
	return DatabasePhase
}

// Init implements [MongoDB].
func (m *mongoDB) Init() error {
	opts := options.Client().
		SetHosts(m.cfg.Hosts).
		SetMaxPoolSize(m.cfg.MaxPoolSize).
		SetMinPoolSize(m.cfg.MinPoolSize).
		SetRetryWrites(m.cfg.RetryWrites).
		SetRetryReads(m.cfg.RetryReads)

	if m.cfg.Log.Enabled {
		level := options.LogLevelInfo
		if m.cfg.Log.Debug {
			level = options.LogLevelDebug
		}
		opts.SetLoggerOptions(options.Logger().
			SetSink(NewZapLogSink(m.logger)).
			SetComponentLevel(options.LogComponentCommand, level))
	}

	var err error
	if m.cfg.Username != "" {
		opts.SetAuth(options.Credential{
			Username: m.cfg.Username,
			Password: m.cfg.Password,
		})
	}
	if m.cfg.ReplicaSet != "" {
		opts.SetReplicaSet(m.cfg.ReplicaSet)
	}
	if m.cfg.ConnectTimeoutMS != 0 {
		opts.SetConnectTimeout(time.Duration(m.cfg.ConnectTimeoutMS) * time.Millisecond)
	}
	if m.cfg.ServerSelectionTimeout != 0 {
		opts.SetServerSelectionTimeout(time.Duration(m.cfg.ServerSelectionTimeout) * time.Millisecond)
	}
	if m.cfg.HeartbeatIntervalMS != 0 {
		opts.SetHeartbeatInterval(time.Duration(m.cfg.HeartbeatIntervalMS) * time.Millisecond)
	}
	if m.cfg.WriteConcern != "" {
		opts.SetWriteConcern(writeconcern.Custom(m.cfg.WriteConcern))
	}
	if m.cfg.ReadPreference != "" {
		var mode readpref.Mode
		var rp *readpref.ReadPref
		mode, err = readpref.ModeFromString(m.cfg.ReadPreference)
		if err != nil {
			return errors.WithStack(err)
		}
		rp, err = readpref.New(mode)
		if err != nil {
			return errors.WithStack(err)
		}
		opts.SetReadPreference(rp)
	}
	if m.client, err = mongo.Connect(opts); err != nil {
		return errors.WithStack(err)
	}

	m.defaultDatabase = m.client.Database(m.cfg.DbName)

	m.logger.With(zap.Namespace("connect")).Debug("init", utils.ObjectToZapFields(m.cfg)...)
	return nil
}

// Name implements [MongoDB].
func (m *mongoDB) Name() string {
	return m.name
}

func NewMongoDB(cfg *MongoConfig, logger *zap.Logger) MongoDB {
	name := "mongodb"
	return &mongoDB{
		cfg:    cfg,
		name:   name,
		logger: logger.Named(name),
	}
}
