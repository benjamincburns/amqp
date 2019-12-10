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
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

//AMQPService acts as a simple interface to the command queue
type AMQPService interface {
	//Consume immediately starts delivering queued messages.
	Consume() (<-chan amqp.Delivery, error)
	//send places a message into the queue
	Send(pub amqp.Publishing) error
	//Requeue rejects the oldMsg and queues the newMsg in a transaction
	Requeue(oldMsg amqp.Delivery, newMsg amqp.Publishing) error
	//CreateQueue attempts to publish a queue
	CreateQueue() error
}

type qpService struct {
	repo AMQPRepository
	conf AMQPConfig
	log  logrus.Ext1FieldLogger
}

// NewAMQPService creates a new AMQPService
func NewAMQPService(
	conf AMQPConfig,
	repo AMQPRepository,
	log logrus.Ext1FieldLogger) AMQPService {

	return &qpService{repo: repo, conf: conf, log: log}
}

func (as qpService) Send(pub amqp.Publishing) error {
	ch, err := as.repo.GetChannel()
	if err != nil {
		return err
	}
	defer ch.Close()
	as.log.WithFields(logrus.Fields{
		"exchange": as.conf.Publish.Exchange,
		"queue":    as.conf.QueueName,
	}).Trace("publishing a message")
	return ch.Publish(as.conf.Publish.Exchange, as.conf.QueueName, as.conf.Publish.Mandatory, as.conf.Publish.Immediate, pub)
}

// Consume immediately starts delivering queued messages.
func (as qpService) Consume() (<-chan amqp.Delivery, error) {
	ch, err := as.repo.GetChannel()
	if err != nil {
		return nil, err
	}
	as.log.WithFields(logrus.Fields{
		"queue":    as.conf.QueueName,
		"consumer": as.conf.Consume.Consumer,
	}).Trace("consuming")
	return ch.Consume(as.conf.QueueName, as.conf.Consume.Consumer, as.conf.Consume.AutoAck,
		as.conf.Consume.Exclusive, as.conf.Consume.NoLocal, as.conf.Consume.NoWait,
		as.conf.Consume.Args)
}

// Requeue rejects the oldMsg and queues the newMsg in a transaction
func (as qpService) Requeue(oldMsg amqp.Delivery, newMsg amqp.Publishing) error {
	ch, err := as.repo.GetChannel()
	if err != nil {
		return err
	}
	defer ch.Close()
	err = ch.Tx()
	if err != nil {
		return err
	}

	err = ch.Publish(oldMsg.Exchange, oldMsg.RoutingKey, as.conf.Publish.Mandatory, as.conf.Publish.Immediate, newMsg)
	if err != nil {
		ch.TxRollback()
		return err
	}

	err = as.repo.RejectDelivery(oldMsg, false)
	if err != nil {
		ch.TxRollback()
		return err
	}
	return ch.TxCommit()
}

// CreateQueue attempts to publish a queue
func (as qpService) CreateQueue() error {
	ch, err := as.repo.GetChannel()
	if err != nil {
		return err
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(as.conf.QueueName, as.conf.Queue.Durable, as.conf.Queue.AutoDelete,
		as.conf.Queue.Exclusive, as.conf.Queue.NoWait, as.conf.Queue.Args)
	return err
}
