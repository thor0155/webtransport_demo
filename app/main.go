package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Item struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"size:255"`
}

type Job struct {
	ID        string
	CreatedAt time.Time
}

const (
	MaxJobs = 5
)

var (
	jobQueue = make(chan *Job, MaxJobs+1) // A buffered channel to hold items
	logger   = zap.NewExample().Named(os.Getenv("POD_NAME"))
)

func main() {
	defer close(jobQueue)
	go func() {
		logger.Info("Worker started")
		for job := range jobQueue {
			logger.Info("Processing job:", zap.String("id", job.ID))
			<-time.After(5 * time.Minute)
			logger.Info("Finished job:", zap.String("id", job.ID))
		}
	}()

	dsn := os.Getenv("MYSQL_DSN") // e.g. "root:password@tcp(mysql:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	logger.Info("mysql connect:", zap.String("dsn", dsn))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect database:", zap.Error(err))
	}

	db.AutoMigrate(&Item{})

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if len(jobQueue) >= MaxJobs {
			w.WriteHeader(http.StatusServiceUnavailable)
			logger.Warn("healthz check failed: too many jobs in queue", zap.Int("queue_length", len(jobQueue)))
			return
		}
		w.Write([]byte("ok"))
	})

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		select {
		case jobQueue <- &Job{ID: uuid.New().String(), CreatedAt: time.Now()}:
			w.WriteHeader(http.StatusOK)
		default:
			w.WriteHeader(http.StatusTooManyRequests)
		}
	})

	http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			var items []Item
			db.Find(&items)
			json.NewEncoder(w).Encode(items)

		case "POST":
			var item Item
			json.NewDecoder(r.Body).Decode(&item)
			db.Create(&item)
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	logger.Info("server running on :8080")
	http.ListenAndServe(":8080", nil)
	logger.Info("server down")
}
