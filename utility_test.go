package queue

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestAMQPMessage_CreateMessage_Success(t *testing.T) {
	am := NewAMQPMessage(2)
	something := []byte("2544")
	msg, err := am.CreateMessage(2544)
	assert.NoError(t, err)
	assert.NotNil(t, msg.Headers)
	assert.Equal(t, int64(0), msg.Headers["retryCount"])
	assert.ElementsMatch(t, something, msg.Body)
}

func TestAMQPMessage_GetNextMessage_Success(t *testing.T) {
	am := NewAMQPMessage(2)
	msg := amqp.Delivery{
		ContentType:     "1",
		ContentEncoding: "2",
		DeliveryMode:    1,
		Type:            "4",
	}
	something := []byte("2544")
	pub, err := am.GetNextMessage(msg, 2544)
	assert.NoError(t, err)
	assert.NotNil(t, pub.Headers)
	assert.Equal(t, int64(0), pub.Headers["retryCount"])
	assert.ElementsMatch(t, something, pub.Body)
	assert.Equal(t, msg.ContentType, pub.ContentType)
	assert.Equal(t, msg.ContentEncoding, pub.ContentEncoding)
	assert.Equal(t, msg.DeliveryMode, pub.DeliveryMode)
	assert.Equal(t, msg.Type, pub.Type)
}

func TestAMQPMessage_GetKickbackMessage_Success(t *testing.T) {
	am := NewAMQPMessage(2)
	msg := amqp.Delivery{
		Body: []byte("stuff"),
	}

	res, err := am.GetKickbackMessage(msg)
	assert.NoError(t, err)
	assert.NotNil(t, res.Headers)
	val, exists := res.Headers["retryCount"]
	assert.True(t, exists)
	assert.Exactly(t, int64(1), val)
}

func TestAMQPMessage_GetKickbackMessage_Failure(t *testing.T) {
	var retries int64 = 20

	msg := amqp.Delivery{
		Headers: map[string]interface{}{
			"retryCount": int64(retries + 1),
		},
		Body: []byte("supposed to fail"),
	}
	am := NewAMQPMessage(retries)
	_, err := am.GetKickbackMessage(msg)
	assert.Error(t, err)
}
