package kafka

type Frame struct {
	Data []byte
}

const (
	defaultMaxAttempts = 5
	defaultBatchSize   = 50
	defaultBatchBytes  = 1 * 1024 * 1024
)
