package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/xamust/go-utils/logger"

	"github.com/segmentio/kafka-go"
)

type Subscriber interface {
	Subscribe(context.Context) error
	Unsubscribe() error
}

type Handler interface {
	Handler(context.Context, *Frame) error
}

type HandlerErr interface {
	HandlerError(ctx context.Context, frame *Frame)
}

type subscriber struct {
	sync.Mutex
	sync.WaitGroup

	reader     *kafka.Reader
	handler    Handler
	HandlerErr HandlerErr
	topic      string
	workers    int
}

// NewSubscriber need testing
func NewSubscriber(cfg Config, topic string, h Handler) Subscriber {
	reader := NewReader(cfg, topic)

	if cfg.Reader.Workers < 1 {
		cfg.Reader.Workers = 1
	}

	return &subscriber{
		reader:  reader,
		handler: h,
		workers: cfg.Reader.Workers,
		topic:   topic,
	}
}

func NewReader(cfg Config, topic string) *kafka.Reader {
	readerCfg := kafka.ReaderConfig{
		Brokers:                cfg.Addr,
		GroupID:                cfg.Reader.Group,
		Topic:                  topic,
		QueueCapacity:          cfg.Reader.QueueCapacity,
		MinBytes:               cfg.Reader.MinBytes,
		MaxBytes:               cfg.Reader.MaxBytes,
		MaxWait:                cfg.Reader.MaxWait.Duration,
		ReadLagInterval:        cfg.Reader.ReadLagInterval.Duration,
		HeartbeatInterval:      cfg.Reader.HeartbeatInterval.Duration,
		CommitInterval:         cfg.Reader.CommitInterval.Duration,
		PartitionWatchInterval: cfg.Reader.PartitionWatchInterval.Duration,
		SessionTimeout:         cfg.Reader.SessionTimeout.Duration,
		RebalanceTimeout:       cfg.Reader.RebalanceTimeout.Duration,
		JoinGroupBackoff:       cfg.Reader.JoinGroupBackoff.Duration,
		RetentionTime:          cfg.Reader.RetentionTime.Duration,
		StartOffset:            cfg.Reader.StartOffset,
		ReadBackoffMin:         cfg.Reader.ReadBackoffMin.Duration,
		ReadBackoffMax:         cfg.Reader.ReadBackoffMax.Duration,
		MaxAttempts:            cfg.Reader.MaxAttempts,
	}

	conn, err := parseAuth(cfg.Username, cfg.Password, cfg.TypeAuth)
	if err != nil {
		logger.DefaultLogger.Fatal(context.Background(), err.Error())
	}

	readerCfg.Dialer = &kafka.Dialer{
		SASLMechanism: conn,
	}

	return kafka.NewReader(readerCfg)
}

func (s *subscriber) Subscribe(ctx context.Context) error {
	log := logger.FromContextLogger(ctx)
	task := make(chan kafka.Message)

	for i := 0; i < s.workers; i++ {
		s.Add(1)
		go func() {
			s.worker(ctx, task)
		}()
	}

	for {
		select {
		case <-ctx.Done():
			close(task)
			return s.Unsubscribe()
		default:
		}
		msg, err := s.reader.FetchMessage(ctx)
		if err != nil {
			log.Error(ctx, fmt.Sprintf("error fetch message: %+v", err))
			continue
		}
		task <- msg
	}
}

func (s *subscriber) Unsubscribe() error {
	s.Wait()
	return s.reader.Close()
}

func (s *subscriber) worker(ctx context.Context, taskQueue <-chan kafka.Message) {
	log := logger.FromContextLogger(ctx)
	defer s.Done()

	for msg := range taskQueue {
		select {
		case <-ctx.Done():
			return
		default:
		}

		frame := &Frame{msg.Value}
		if len(msg.Headers) > 0 {
			ctx = NewContextMetadata(ctx, msg.Headers)
		}

		log.Info(
			ctx,
			string(msg.Value),
			logger.Operation(fmt.Sprintf("Got message from %s", s.topic)),
		)

		err := s.handler.Handler(ctx, frame)
		if err != nil {
			log.Error(ctx, fmt.Sprintf("error handling message: %+v", err))
			if s.HandlerErr != nil {
				s.HandlerErr.HandlerError(ctx, frame)
			}
			continue
		}
		if err = s.reader.CommitMessages(ctx, msg); err != nil {
			log.Error(ctx, fmt.Sprintf("error committing message: %+v", err))
		}
	}
}
