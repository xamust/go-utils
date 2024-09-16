package kafka_logger

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/xamust/go-utils/encoder/json"
	"strings"

	"github.com/xamust/go-utils/metadata"
	"github.com/xamust/go-utils/metadata/request_id"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type Publisher interface {
	Publish(ctx context.Context, in any, topic string) error

	Write(p []byte) (n int, err error)
}

type publisher struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaLogger(stage Stage, kafkaConfig ...KafkaLogConfig) Publisher {
	cfg := GetConfig(stage)
	if len(kafkaConfig) != 0 {
		cfg = kafkaConfig[0]
	}
	conn, err := parseAuth(cfg.Config.Username, cfg.Config.Password, cfg.Config.TypeAuth)
	if err != nil {
		log.Fatal(err.Error())
	}
	writer := NewWriter(cfg.Config.Writer)

	writer.Transport = &kafka.Transport{
		SASL: conn,
	}
	writer.Addr = kafka.TCP(cfg.Config.Addr...)
	writer.WriteTimeout = cfg.Config.Timeout.Duration

	return &publisher{
		writer: writer,
		topic:  cfg.Topic,
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

func (b *publisher) Publish(ctx context.Context, msg any, topic string) error {
	data, err := json.NewCodec().Marshal(msg)
	if err != nil {
		return err
	}
	return b.publish(ctx, data, topic)
}

func (b *publisher) publish(ctx context.Context, msg []byte, topic string) (err error) {
	key := uuid.New().String()
	if ruid, ok := request_id.FromContext(ctx); ok {
		key = ruid
	}
	if err = b.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Value: msg,
		Key:   []byte(key),
	}); err != nil {
		return err
	}
	metadata.SetRsUidContextHeader(ctx, key)
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

func (b *publisher) Write(p []byte) (n int, err error) {
	if err = b.Publish(context.Background(), string(p), b.topic); err != nil {
		// skip err
		log.Error(context.Background(), fmt.Errorf("failed to publish message to kafEL: %w", err).Error(), err.Error())
		return 0, nil
	}
	return len(p), nil
}
