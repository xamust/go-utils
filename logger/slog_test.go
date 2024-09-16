package logger

import (
	"bytes"
	"context"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/xamust/go-utils/kafka_logger"
	kafka_mock "github.com/xamust/go-utils/kafka_logger/mock"
	"github.com/xamust/go-utils/util/dual_writer"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tFields = map[string]any{
	"request_": "qwery123",
}

func TestFields(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewSlogLogger(
		WithLevel(DebugLevel),
		WithOutput(buf),
		WithFields(tFields),
		WithSource(),
	)

	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	f := map[string]any{
		"key": "val",
	}
	nl := l.Fields(f)

	nl.Info(ctx, "test")

	assert.Contains(t, buf.String(), `"message":"test"`)
	assert.Contains(t, buf.String(), `"key":"val"`)
}

func Test_Error(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewSlogLogger(
		WithLevel(DebugLevel),
		WithOutput(buf),
		WithSource(),
	)

	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l.Error(ctx, "important error")

	assert.Contains(t, buf.String(), `"ErrorText":"important error"`)
}

func Test_CheckSkipLevel(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewSlogLogger(
		WithLevel(InfoLevel),
		WithOutput(buf),
		WithSource(),
	)
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l.Debug(ctx, "need skip")

	assert.Nil(t, buf.Bytes())
}

func Test_Events(t *testing.T) {
	s, r := "source", "receiver"
	ctx := NewContextEvent(context.TODO(), s, r)
	buf := bytes.NewBuffer(nil)
	l := NewSlogLogger(
		WithLevel(InfoLevel),
		WithOutput(buf),
	)
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	l.Info(ctx, "test msg")

	assert.NotNil(t, buf.Bytes())
	assert.Contains(t, buf.String(), fmt.Sprintf(`"%s":"%s"`, keyEventRec, r))
	assert.Contains(t, buf.String(), fmt.Sprintf(`"%s":"%s"`, keyEventSource, s))
}

func Test_Ctx(t *testing.T) {
	ctx := NewContextEvent(context.Background(), "test", "test2")
	ctx = NewContextEvent(ctx, "test2", "test3")
	ctx = NewContextEvent(ctx, "test3", "test4")
	_ = getEvent(ctx)
}

func Test_LongMsg(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewSlogLogger(
		WithLevel(InfoLevel),
		WithOutput(buf),
		WithSource(),
		MaxBytesMessage(50),
	)

	if err := l.Init(); err != nil {
		t.Fatal(err)
	}
	msg := "MESSAGE_MESSAGE_MESSAGE_MESSAGE_MESSAGE_MESSAGE_MESSAGE_MESSAGE_MESSAGE_MESSAGE"
	l.Info(ctx, msg)

	assert.Contains(t, buf.String(), fmt.Sprintf("too much content: have [%d], expected [<%d]", len(msg), 50))
}

func Test_RqUID(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewSlogLogger(
		WithLevel(InfoLevel),
		WithOutput(buf),
		WithSource(),
	)

	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	var fields = make(map[string]any)

	fields[keyUID] = "a1-s2-d3"

	l.Fields(fields).Info(ctx, "test", Operation("testing something"))

	fmt.Printf("buf: %s\n", buf.String())

	assert.Contains(t, buf.String(), `"RqUID":"a1-s2-d3"`)
}

func Test_KafkaLogger(t *testing.T) {
	ctx := context.TODO()
	buf := bytes.NewBuffer(nil)
	l := NewSlogLogger(
		WithLevel(InfoLevel),
		WithOutput(buf),
		WithKafka("test_service", kafka_logger.TestStage),
		//	WithSource(),
	)

	kafkaMock := &kafka_mock.KafkaPublisherMock{}

	kafkaMock.
		On("Write", mock.Anything, mock.Anything).
		Return(0, nil).
		Run(func(args mock.Arguments) {
			in := args.Get(0).([]byte)
			fmt.Printf("kafka: %s", string(in))
			assert.Contains(t, string(in), "{\"k8s.pod_name\":\"test_service\",\"k8s.container_name\":\"test_service\",\"ProcStatus\":\"SUCCESS\",\"Operation\":\"testing something\",\"EventReceiver\":\"\",\"EventSource\":\"\"")
		})

	l.(*slogLogger).opts.out = dual_writer.NewDuplicator(buf, kafkaMock)
	if err := l.Init(); err != nil {
		t.Fatal(err)
	}

	l.Info(ctx, "test", Operation("testing something"))
	fmt.Printf("buf: %s\n", buf.String())
	assert.Contains(t, buf.String(), "{\"k8s.pod_name\":\"test_service\",\"k8s.container_name\":\"test_service\",\"ProcStatus\":\"SUCCESS\",\"Operation\":\"testing something\",\"EventReceiver\":\"\",\"EventSource\":\"\"")

}
