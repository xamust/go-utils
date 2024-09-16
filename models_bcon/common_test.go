package models_bcon

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_DurationParsing(t *testing.T) {

	testCases := []struct {
		name       string
		body       []byte
		errorCheck assert.ValueAssertionFunc
		expected   time.Duration
	}{
		{
			"ValidString",
			[]byte(`{"timeout":"30s"}`),
			assert.Nil,
			30 * time.Second,
		},
		{
			"ValidInt",
			[]byte(`{"timeout":60}`),
			assert.Nil,
			60,
		},
		{
			name:       "InvalidString",
			body:       []byte(`{"timeout":"30"}`),
			errorCheck: assert.NotNil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			dest := struct {
				Timeout Duration `json:"timeout"`
			}{}

			err := json.Unmarshal(tt.body, &dest)
			tt.errorCheck(t, err)
			assert.Equal(t, tt.expected, dest.Timeout.Duration)
		})
	}
}
