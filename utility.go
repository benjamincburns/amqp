package queue

import (
	"encoding/json"
	"fmt"

	"github.com/whiteblock/amqp/config"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// RetryCountHeader is the header for the retry count of the amqp message
const RetryCountHeader = "retryCount"

// AMQPMessage contains utilities for manipulating AMQP messages
type AMQPMessage interface {
	CreateMessage(body interface{}) (amqp.Publishing, error)
	// GetKickbackMessage takes the delivery and creates a message from it
	// for requeuing on non-fatal error
	GetKickbackMessage(msg amqp.Delivery) (amqp.Publishing, error)
	GetNextMessage(msg amqp.Delivery, body interface{}) (amqp.Publishing, error)
}

type amqpMessage struct {
	maxRetries int64
}

// NewAMQPMessage creates a new AMQPMessage
func NewAMQPMessage(maxRetries int64) AMQPMessage {
	return &amqpMessage{maxRetries: maxRetries}
}

// CreateMessage creates a message from the given body
func (am amqpMessage) CreateMessage(body interface{}) (amqp.Publishing, error) {
	return CreateMessage(body)
}

// GetNextMessage is similar to GetKickbackMessage but takes in a new body, and does not increment the
// retry count
func (am amqpMessage) GetNextMessage(msg amqp.Delivery, body interface{}) (amqp.Publishing, error) {
	return GetNextMessage(msg, body)
}

// GetKickbackMessage takes the delivery and creates a message from it
// for requeuing on non-fatal error. It returns an error if the number of retries is
// exceeded
func (am amqpMessage) GetKickbackMessage(msg amqp.Delivery) (amqp.Publishing, error) {
	return GetKickbackMessage(am.maxRetries, msg)
}

// OpenAMQPConnection attempts to dial a new AMQP connection
func OpenAMQPConnection(conf config.Endpoint) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		conf.QueueProtocol,
		conf.QueueUser,
		conf.QueuePassword,
		conf.QueueHost,
		conf.QueuePort,
		conf.QueueVHost))
}

// TryCreateQueues atempts to create the given queues, but doesn't error out if it fails
func TryCreateQueues(log logrus.Ext1FieldLogger, queues ...AMQPService) {
	errChan := make(chan error)
	for i := range queues {
		go func(i int) {
			errChan <- queues[i].CreateQueue()
		}(i)

		go func(i int) {
			errChan <- queues[i].CreateExchange()
		}(i)
	}

	for i := 0; i < len(queues)*2; i++ {
		err := <-errChan
		if err != nil {
			log.WithFields(logrus.Fields{"err": err}).Debug("failed to create a queue or exchange")
		}
	}
}

// CreateMessage creates a message from the given body
func CreateMessage(body interface{}) (amqp.Publishing, error) {
	rawBody, err := json.Marshal(body)
	if err != nil {
		return amqp.Publishing{}, err
	}

	pub := amqp.Publishing{
		Headers: map[string]interface{}{
			RetryCountHeader: int64(0),
		},
		Body: rawBody,
	}
	return pub, nil
}

// GetNextMessage is similar to GetKickbackMessage but takes in a new body, and does not increment the
// retry count
func GetNextMessage(msg amqp.Delivery, body interface{}) (amqp.Publishing, error) {
	rawBody, err := json.Marshal(body)
	if err != nil {
		return amqp.Publishing{}, err
	}
	pub := amqp.Publishing{
		Headers: msg.Headers,
		// Properties
		ContentType:     msg.ContentType,
		ContentEncoding: msg.ContentEncoding,
		DeliveryMode:    msg.DeliveryMode,
		Type:            msg.Type,
		Body:            rawBody,
	}
	if pub.Headers == nil {
		pub.Headers = map[string]interface{}{}
	}
	pub.Headers[RetryCountHeader] = int64(0) //reset retry count

	return pub, nil
}

// GetKickbackMessage takes the delivery and creates a message from it
// for requeuing on non-fatal error. It returns an error if the number of retries is
// exceeded
func GetKickbackMessage(maxRetries int64, msg amqp.Delivery) (amqp.Publishing, error) {
	pub := amqp.Publishing{
		Headers: msg.Headers,
		// Properties
		ContentType:     msg.ContentType,
		ContentEncoding: msg.ContentEncoding,
		DeliveryMode:    msg.DeliveryMode,
		Priority:        msg.Priority,
		CorrelationId:   msg.CorrelationId,
		ReplyTo:         msg.ReplyTo,
		Expiration:      msg.Expiration,
		MessageId:       msg.MessageId,
		Timestamp:       msg.Timestamp,
		Type:            msg.Type,
		Body:            msg.Body,
	}
	if pub.Headers == nil {
		pub.Headers = map[string]interface{}{}
	}
	_, exists := pub.Headers[RetryCountHeader]
	if !exists {
		pub.Headers[RetryCountHeader] = int64(0)
	}
	if pub.Headers[RetryCountHeader].(int64) > maxRetries {
		return amqp.Publishing{}, errors.New("too many retries")
	}
	pub.Headers[RetryCountHeader] = pub.Headers[RetryCountHeader].(int64) + 1
	return pub, nil
}
