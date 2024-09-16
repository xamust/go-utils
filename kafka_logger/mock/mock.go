package mock

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type KafkaPublisherMock struct {
	mock.Mock
}

func (k *KafkaPublisherMock) Publish(ctx context.Context, in any, topic string) error {
	args := k.Called(in, topic)
	return args.Error(0)
}

func (k *KafkaPublisherMock) Write(p []byte) (n int, err error) {
	args := k.Called(p)
	return args.Get(0).(int), args.Error(1)
}
