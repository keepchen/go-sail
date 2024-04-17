package dto

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type demoDto struct {
	Base
	Data struct {
		Sentence string `json:"word"`
	} `json:"data"`
}

func TestMarshalAndUnmarshal(t *testing.T) {
	var demo demoDto

	demo.RequestID = "5686efa5-c747-4f63-8657-e6052f8181a9"
	demo.Success = true
	demo.Code = 200
	demo.Message = "SUCCESS"
	demo.Timestamp = time.Now().UnixMilli()
	demo.Data.Sentence = "hello world!"

	js, err := json.Marshal(&demo)
	t.Log(string(js))

	assert.NoError(t, err)

	var newDemo demoDto
	err = json.Unmarshal(js, &newDemo)
	assert.NoError(t, err)
	assert.Equal(t, newDemo.RequestID, newDemo.RequestID)
	assert.Equal(t, newDemo.Success, newDemo.Success)
	assert.Equal(t, newDemo.Code, newDemo.Code)
	assert.Equal(t, newDemo.Message, newDemo.Message)
	assert.Equal(t, newDemo.Timestamp, newDemo.Timestamp)
	assert.Equal(t, newDemo.Data.Sentence, newDemo.Data.Sentence)
}
