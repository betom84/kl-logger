package models_test

import (
	"testing"
	"time"

	"github.com/betom84/kl-logger/klimalogg/models"
	"github.com/stretchr/testify/assert"
)

func TestMarshalSendTimeData(t *testing.T) {
	now, err := time.Parse("Mon 2006-01-02 15:04:05", "Thu 2014-10-30 21:58:25")
	assert.NoError(t, err)

	model := models.SendTimeData{Now: now}
	expected := []byte{0x0, 0x0, 0x25, 0x58, 0x21, 0x04, 0x03, 0x41, 0x01}

	data, err := model.MarshalBinary()
	assert.NoError(t, err)

	assert.Equal(t, expected, data)
}
