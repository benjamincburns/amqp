/*
	AMQP Library
	Copyright (C) 2019 Whiteblock Inc.

	This program is free software; you can redistribute it and/or
	modify it under the terms of the GNU Lesser General Public
	License as published by the Free Software Foundation; either
	version 3 of the License, or (at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
	Lesser General Public License for more details.

	You should have received a copy of the GNU Lesser General Public License
	along with this program; if not, write to the Free Software Foundation,
	Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/

package queue

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
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
func OpenAMQPConnection(conf AMQPEndpoint) (*amqp.Connection, error) {
	return amqp.Dial(fmt.Sprintf("%s://%s:%s@%s:%d/%s",
		conf.QueueProtocol,
		conf.QueueUser,
		conf.QueuePassword,
		conf.QueueHost,
		conf.QueuePort,
		conf.QueueVHost))
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
