package metadata

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_SetContext(t *testing.T) {

	ctx := context.Background()

	header := HeaderSource{Header: Header{
		Service:        "",
		SourceSystem:   "",
		RqUid:          "",
		OperUID:        "",
		RqTm:           time.Now().Local().Format(time.RFC3339Nano),
		ReceiverSystem: "",
		Platform:       "source",
		RsUid:          "",
	}}

	ctx = NewContextHeader(ctx, header)

	SetTimeContextHeader(ctx)

	h, _ := FromContextHeader(ctx)
	assert.NotEmpty(t, h.Header.RsTm)

	SetRsUidContextHeader(ctx, uuid.New().String())

	SetServiceCaller(ctx)

	h, _ = FromContextHeader(ctx)
	assert.NotEmpty(t, h.Header.RsUid)
	assert.Equal(t, "Test_SetContext", h.Header.Service)
}

func Test_MockTime(t *testing.T) {
	header := HeaderSource{Header: Header{
		RqTm:     time.Now().Local().Format(time.RFC3339Nano),
		Platform: "source",
	}}

	ctx := context.Background()
	ctx = NewContextHeader(ctx, header)

	SetTimeContextHeader(ctx)
	h, _ := FromContextHeader(ctx)
	assert.NotEmpty(t, h.Header.RsTm)
	assert.NotEqual(t, time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC).Local().Format(time.RFC3339), h.Header.RsTm)

	ctx = NewContextMockTime(ctx, func() time.Time {
		return time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC)
	})
	SetTimeContextHeader(ctx)
	h, _ = FromContextHeader(ctx)
	assert.NotEmpty(t, h.Header.RsTm)
	assert.Equal(t, time.Date(2020, 11, 01, 00, 00, 00, 0, time.UTC).Local().Format(time.RFC3339), h.Header.RsTm)
}

func Test_CopyRqUid(t *testing.T) {
	header := HeaderSource{Header: Header{
		RqTm:     time.Now().Local().Format(time.RFC3339Nano),
		Platform: "source",
		RqUid:    uuid.NewString(),
	}}

	ctx := context.Background()
	ctx = NewContextHeader(ctx, header)

	CopyRqToRsUid(ctx)
	h, _ := FromContextHeader(ctx)
	assert.NotEqual(t, 0, len(h.Header.RsUid))
}
