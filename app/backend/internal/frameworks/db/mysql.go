package db

import (
	"api/internal/frameworks/config"
	"api/internal/frameworks/obj"
	"database/sql"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// default
const (
	defaultMysqlUserName     string        = "root"
	defaultMysqlPassword     string        = ""
	defaultMysqlHost         string        = "localhost"
	defaultMysqlPort         int           = 3306
	defaultMysqlDatabase     string        = "test"
	defaultMysqlMaxLifetime  time.Duration = 0
	defaultMysqlMaxIdleTime  time.Duration = 0
	defaultMysqlMaxOpenConns int           = 0
	defaultMysqlMaxIdleConns int           = 0
	defaultMysqlCharset      string        = "utf8mb4"
	defaultMysqlTimeout      string        = "5s"
	defaultMysqlReadTimeout  string        = "10s"
	defaultMysqlWriteTimeout string        = "10s"
)

const DSN_FORMAT = "%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s&timeout=%s&readTimeout=%s&writeTimeout=%s"

type MysqlDB interface {
	obj.Activity
	obj.Component
	obj.Phaseable
	Client() *gorm.DB
	Session() *gorm.DB
}

type mysqlDB struct {
	logger *zap.Logger
	client *gorm.DB
	conn   *sql.DB
	cfg    *MySQLConfig
	name   string
}

func (s *mysqlDB) Active() bool {
	return s.cfg.Active
}

func (m *mysqlDB) Phase() int {
	return 1
}

func (s *mysqlDB) Init() (err error) {
	var username string = defaultMysqlUserName
	var password string = defaultMysqlPassword
	var host string = defaultMysqlHost
	var port int = defaultMysqlPort
	var database string = defaultMysqlDatabase
	var charset string = defaultMysqlCharset
	var maxLifeTime = defaultMysqlMaxLifetime
	var maxIdleTime = defaultMysqlMaxIdleTime
	var maxIdleConns int = defaultMysqlMaxIdleConns
	var maxOpenConns int = defaultMysqlMaxOpenConns
	var parseTime bool
	var loc string
	var timeout string = defaultMysqlTimeout
	var readTimeout = defaultMysqlReadTimeout
	var writeTimeout = defaultMysqlWriteTimeout

	if s.cfg != nil {
		username = s.cfg.User
		password = s.cfg.Password
		host = s.cfg.Host
		port = s.cfg.Port
		database = s.cfg.DBName
		charset = s.cfg.Charset
		maxIdleConns = s.cfg.MaxIdleConns
		maxOpenConns = s.cfg.MaxOpenConns
		parseTime = s.cfg.ParseTime
		loc = s.cfg.Loc
		timeout = s.cfg.Timeout
		readTimeout = s.cfg.ReadTimeout
		writeTimeout = s.cfg.WriteTimeout
		if s.cfg.ConnMaxLifetime != "" {
			if maxLifeTime, err = time.ParseDuration(s.cfg.ConnMaxLifetime); err != nil {
				err = errors.WithStack(err)
				return
			}
		}
		if s.cfg.ConnMaxIdleTime != "" {
			if maxIdleTime, err = time.ParseDuration(s.cfg.ConnMaxIdleTime); err != nil {
				err = errors.WithStack(err)
				return
			}
		}
	}

	dsn := fmt.Sprintf(DSN_FORMAT,
		username, password, host, port, database, charset, parseTime, loc, timeout, readTimeout, writeTimeout)
	s.logger.Sugar().Debugf("dsn: %s", dsn)

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// SkipDefaultTransaction: true,
		PrepareStmt: true,
	})
	if err != nil {
		return errors.WithStack(err)
	}
	if config.IsDebugging() {
		db = db.Debug()
	}
	raw, err := db.DB()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	raw.SetConnMaxLifetime(maxLifeTime)
	s.logger.Sugar().Debugf("maxLifeTime: %v", maxLifeTime)
	raw.SetConnMaxIdleTime(maxIdleTime)
	s.logger.Sugar().Debugf("maxIdleTime: %v", maxIdleTime)
	raw.SetMaxIdleConns(maxIdleConns)
	s.logger.Sugar().Debugf("maxIdleConns: %d", maxIdleConns)
	if maxOpenConns > 0 {
		raw.SetMaxOpenConns(maxOpenConns)
		s.logger.Sugar().Debugf("maxOpenConns: %d", maxOpenConns)
	}
	s.client = db
	return
}

// create a new session mode
func (s *mysqlDB) Session() *gorm.DB {
	return s.client.Session(&gorm.Session{NewDB: true})
}

func (s *mysqlDB) Client() *gorm.DB {
	return s.client
}

func (s *mysqlDB) Name() string {
	return s.name
}

// new mysql db
func NewMysqlDB(cfg *MySQLConfig, logger *zap.Logger) MysqlDB {
	name := "mysql"
	svc := &mysqlDB{
		cfg:    cfg,
		name:   name,
		logger: logger.Named(name),
	}
	return svc
}

func ProvideMysqlGorm(mysql MysqlDB) *gorm.DB {
	return mysql.Client()
}
