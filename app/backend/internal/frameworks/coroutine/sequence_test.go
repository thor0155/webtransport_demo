package coroutine

import (
	"context"
	"testing"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mysequenceTask struct {
	logger *zap.Logger
	name   string
	sec    int
}

func TestSequence(t *testing.T) {

	logger := zap.NewExample()

	corotine, err := NewSequenceCoroutine(func(m *mysequenceTask) error {
		return m.Run()
	})

	if err != nil {
		panic(err)
	}

	corotine.Add(0, &mysequenceTask{
		name:   "c",
		logger: logger.Named("c"),
		sec:    3,
	})
	corotine.Add(1, &mysequenceTask{
		name:   "b",
		logger: logger.Named("b"),
		sec:    3,
	})
	corotine.Add(2, &mysequenceTask{
		name:   "a",
		logger: logger.Named("a"),
		sec:    3,
	})
	corotine.Add(2, &mysequenceTask{
		name:   "a-a",
		logger: logger.Named("a-a"),
		sec:    3,
	})

	t.Run("normal", func(t *testing.T) {

		timeoutCtx, timeoutCancel := context.WithTimeout(t.Context(), time.Second*20)
		defer timeoutCancel()

		logger.Info("test start")
		assert.NoError(t, corotine.Invoke(timeoutCtx))
		logger.Info("test done")
	})

	t.Run("timeout", func(t *testing.T) {
		timeoutCtx, timeoutCancel := context.WithTimeout(t.Context(), time.Second*5)
		defer timeoutCancel()
		logger.Info("test start")
		err := corotine.Invoke(timeoutCtx)
		assert.True(t, errors.Is(err, context.DeadlineExceeded))
		logger.Error(err.Error())
		logger.Info("test done")
	})

	t.Run("some error happened", func(t *testing.T) {
		corotine.Add(1, &mysequenceTask{
			name:   "b-a",
			logger: logger.Named("b-a"),
			sec:    20,
		})

		timeoutCtx, timeoutCancel := context.WithTimeout(t.Context(), time.Second*20)
		defer timeoutCancel()
		logger.Info("test start")
		err := corotine.Invoke(timeoutCtx)
		assert.Error(t, err)
		logger.Error(err.Error())
		logger.Info("test done")
	})

}

func (m *mysequenceTask) Run() error {
	if m.sec > 10 {
		return errors.New("crash")
	}
	for i := 0; i < m.sec; i++ {
		<-time.After(time.Second)
		m.logger.Info("after", zap.Int("sec", i+1))
	}
	m.logger.Info("done")
	return nil
}
