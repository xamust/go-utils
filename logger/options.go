package logger

import (
	"context"
	"github.com/xamust/go-utils/kafka_logger"
	"github.com/xamust/go-utils/util/dual_writer"
	"github.com/xamust/go-utils/util/env"
	"io"
	"os"
)

type Option func(*Options)

const (
	defaultCallerSkipCount = 3
	defaultMaxBytesMessage = 1 << 10
)

type customKey struct {
	// TimeKey is the key used for the time of the log call
	timeKey string
	// LevelKey is the key used for the level of the log call
	levelKey string
	// MessageKey is the key used for the message of the log call
	messageKey string
	// SourceKey is the key used for the source file and line of the log call
	sourceKey string
}

type Options struct {
	// out holds the output writer
	out io.Writer
	// context holds external options
	context context.Context
	// fields to always be logged
	fields map[string]any
	// The logging lvl the logger should log
	lvl Level
	// custom keys value for logging
	keys customKey
	// addSource the flag that shows the culler
	addSource bool
	// callerSkipCount skipping the call queue
	callerSkipCount int
	// maxBytesMessage check maximum symbol in msg
	maxBytesMessage int
}

func NewOptions(opts ...Option) Options {
	options := Options{
		lvl:             DefaultLevel,
		fields:          make(map[string]interface{}),
		out:             os.Stdout,
		context:         context.Background(),
		addSource:       false,
		callerSkipCount: defaultCallerSkipCount,
		maxBytesMessage: defaultMaxBytesMessage,
	}
	customKeys()(&options)

	for _, o := range opts {
		o(&options)
	}

	return options
}

func WithOutput(out io.Writer) Option {
	return func(o *Options) {
		o.out = out
	}
}

func WithContext(ctx context.Context) Option {
	return func(o *Options) {
		o.context = ctx
	}
}

func WithFields(fields map[string]interface{}) Option {
	return func(o *Options) {
		for k, v := range fields {
			o.fields[k] = v
		}
	}
}

func WithLevel(level Level) Option {
	return func(o *Options) {
		o.lvl = level
	}
}

func WithSource() Option {
	return func(options *Options) {
		options.addSource = true
	}
}

func WithCallerSkipCount(c int) Option {
	return func(o *Options) {
		o.callerSkipCount = c
	}
}

func MaxBytesMessage(m int) Option {
	return func(options *Options) {
		if m < 1 {
			m = defaultMaxBytesMessage
		}
		options.maxBytesMessage = m
	}
}

// WithKafka add kafka logger (if set kafkaConfig, stage is ignored)
// stage get default config data
func WithKafka(serviceName string, stage kafka_logger.Stage, kafkaConfig ...kafka_logger.KafkaLogConfig) Option {
	return func(options *Options) {
		options.out = dual_writer.NewDuplicator(options.out, kafka_logger.NewKafkaLogger(stage, kafkaConfig...))
		options.fields[kafka_logger.K8sPodName] = env.GetHostName(serviceName)
		options.fields[kafka_logger.K8sContainerName] = serviceName
	}
}

func customKeys() Option {
	return func(o *Options) {
		o.keys.timeKey = "@timestamp"
		o.keys.levelKey = "level"
		o.keys.messageKey = "message"
		o.keys.sourceKey = "caller"
	}
}
