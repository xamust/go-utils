package json

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ReadBody(t *testing.T) {
	codec := NewCodec()
	bodyTest := []byte(`{"test": "string"}`)

	var destStr string
	err := codec.ReadBody(bytes.NewReader(bodyTest), &destStr)
	assert.Nil(t, err)
	assert.Equal(t, string(bodyTest), destStr)

	var destByte []byte
	err = codec.ReadBody(bytes.NewReader(bodyTest), &destByte)
	assert.Nil(t, err)
	assert.Equal(t, bodyTest, destByte)

	type testStruct struct {
		Test string `json:"test"`
	}

	var destStruct testStruct
	err = codec.ReadBody(bytes.NewReader(bodyTest), &destStruct)
	assert.Nil(t, err)
	assert.Equal(t, testStruct{Test: "string"}, destStruct)

}

func Test_Write(t *testing.T) {
	codec := NewCodec()
	b := bytes.NewBuffer(nil)

	testCases := []struct {
		name   string
		expect string
		msg    any
	}{
		{
			`string`,
			`write big string`,
			`write big string`,
		},
		{
			`json`,
			`{"test":"bigdata"}`,
			struct {
				Test string `json:"test"`
			}{
				Test: "bigdata",
			},
		},
		{
			"bytes",
			`write big bytes`,
			[]byte(`write big bytes`),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := codec.Write(b, tt.msg)
			assert.Nil(t, err)

			assert.Contains(t, b.String(), tt.expect)
			b.Reset()
		})
	}
}
