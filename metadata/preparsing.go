package metadata

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/xamust/go-utils/validate"
)

const (
	keyHeader = `header`

	keyService        = `service`
	keySourceSystem   = `sourceSystem`
	keyRqUID          = `rqUid`
	keyOperUID        = `operUid`
	keyRqTM           = `rqTm`
	keyReceiverSystem = `receiverSystem`
)

type ParserHeaderSource interface {
	ExtractHeaderFromBytes(b []byte) (HeaderSource, error)
	Validate(HeaderSource) error
}

func (h *HeaderSource) ExtractHeaderFromBytes(b []byte) (HeaderSource, error) {
	var result HeaderSource

	err := json.Unmarshal(b, &result)

	return result, err
}

func (h *HeaderSource) RequestID() string {
	result := h.Header.RqUid
	if len(result) == 0 {
		result = uuid.New().String()
	}
	return result
}

func (h *HeaderSource) Validate(source HeaderSource) error {
	return validate.Validate(source)
}
