package chatd

import (
	"api/internal/frameworks/db"
	"context"
	"log"
	"os"
	"testing"
)

type tester struct {
	mongodb        db.MongoDB
	chatMessageDao ChatMessageDao
}

var test *tester

func TestMain(m *testing.M) {

	mongodb := db.NewMockMongoDB()
	test = &tester{
		mongodb:        mongodb,
		chatMessageDao: NewChatMessageDao(mongodb),
	}

	if err := mongodb.Init(); err != nil {
		panic(err)
	}
	if err := test.chatMessageDao.Init(); err != nil {
		panic(err)
	}

	log.Println("test start")
	code := m.Run()

	test.mongodb.Close(context.Background())

	os.Exit(code)
}
