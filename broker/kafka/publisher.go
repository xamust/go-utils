package kafka

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/xamust/go-utils/logger"
	"github.com/xamust/go-utils/metadata"
	"github.com/xamust/go-utils/metadata/request_id"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type Publisher interface {
	Publish(ctx context.Context, msg any, topic string, opts ...PublishOption) error
}

type publisher struct {
	writer *kafka.Writer

	opts PublishOptions
}

func NewPublisher(cfg Config) Publisher {

	conn, err := parseAuth(cfg.Username, cfg.Password, cfg.TypeAuth)
	if err != nil {
		logger.DefaultLogger.Fatal(context.Background(), err.Error())
	}
	writer := NewWriter(cfg.Writer)

	writer.Transport = &kafka.Transport{
		SASL: conn,
	}
	writer.Addr = kafka.TCP(cfg.Addr...)
	writer.WriteTimeout = cfg.Timeout.Duration

	return &publisher{
		writer: writer,
	}
}

func NewWriter(cfg Writer) *kafka.Writer {
	if cfg.MaxAttempts < 1 {
		cfg.MaxAttempts = defaultMaxAttempts
	}
	if cfg.BatchSize < 1 {
		cfg.BatchSize = defaultBatchSize
	}
	if cfg.BatchBytes < 1 {
		cfg.BatchBytes = defaultBatchBytes
	}

	return &kafka.Writer{
		MaxAttempts: cfg.MaxAttempts,
		BatchSize:   cfg.BatchSize,
		BatchBytes:  cfg.BatchBytes,
	}
}

func (b *publisher) Publish(ctx context.Context, msg any, topic string, opts ...PublishOption) error {
	pOpts := NewOptions(opts...)

	data, err := pOpts.codec.Marshal(msg)
	if err != nil {
		return err
	}

	return b.publish(ctx, data, topic)
}

func (b *publisher) publish(ctx context.Context, msg []byte, topic string) (err error) {
	log := logger.FromContextLogger(ctx)
	key := uuid.New().String()
	if ruid, ok := request_id.FromContext(ctx); ok {
		key = ruid
	}

	log.Info(
		ctx,
		string(msg),
		logger.Operation(fmt.Sprintf("Send to topic %s", topic)),
	)

	kmsg := kafka.Message{
		Topic: topic,
		Value: msg,
		Key:   []byte(key),
	}

	if err = b.writer.WriteMessages(ctx, kmsg); err != nil {
		log.Error(ctx, fmt.Sprintf("error send message: %v", err))
		return err
	}

	metadata.SetRsUidContextHeader(ctx, key)
	log.Debug(ctx, fmt.Sprintf("successful publication of the message to topic: [%s]", topic))
	return err
}

func parseAuth(name, pass, typeConnect string) (sasl.Mechanism, error) {
	var err error
	switch strings.ToUpper(typeConnect) {
	case "PLAIN":
		return &plain.Mechanism{
			Username: name,
			Password: pass,
		}, nil
	case "SCRAM-SHA-512":
		return scram.Mechanism(scram.SHA512, name, pass)
	case "SCRAM-SHA-256":
		return scram.Mechanism(scram.SHA256, name, pass)
	default:
		err = errors.New("unknown type connection to broker")
	}
	return nil, err
}
