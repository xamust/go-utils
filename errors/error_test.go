package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Error_Cast(t *testing.T) {
	err := NewInternalErrorRsp("internalation EERRRRORRR")
	assert.NotNil(t, err)
}
