package consumer

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/azagrivin/testProjectCourses/internal/logger"
)

type Consumer struct {
	tick time.Duration
	log  logger.Logger

	done      chan struct{}
	shutdown  int32
	terminate func()
}

func NewConsumer(tick time.Duration, log logger.Logger) *Consumer {
	return &Consumer{tick: tick, log: log}
}

func (c *Consumer) CatchAndServe(f func(ctx context.Context) error) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	c.terminate = cancel

	c.done = make(chan struct{})
	defer close(c.done)

	err := f(ctx)
	if err != nil {
		return err
	}

	timer := time.Tick(c.tick)
	for {
		select {
		case <-timer:
			err = f(ctx)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		default:
		}
	}
}

func (c *Consumer) Shutdown(ctx context.Context) error {
	if c.done == nil {
		return fmt.Errorf("consumer wasn't ran")
	}

	atomic.StoreInt32(&c.shutdown, 1)

	select {
	case <-c.done:
	case <-ctx.Done():
		c.terminate()
	}

	return nil
}

func (c *Consumer) isShutdowning() bool {
	return atomic.LoadInt32(&c.shutdown) == 1
}
