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
	"fmt"
	"testing"

	externalsMocks "github.com/whiteblock/amqp/mocks/pkg/externals"
	"github.com/whiteblock/amqp/mocks/pkg/repository"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewAMQPService(t *testing.T) {
	conf := config.AMQP{
		QueueName: "test queue",
	}
	repo := new(mocks.AMQPRepository)

	serv := NewAMQPService(conf, repo, nil)
	assert.NotNil(t, serv)
}

func TestAMQPService_Consume(t *testing.T) {
	conf := config.AMQP{
		QueueName: "test queue",
		Consume: config.Consume{
			Consumer:  "test",
			AutoAck:   false,
			Exclusive: false,
			NoLocal:   false,
			NoWait:    false,
			Args:      nil,
		},
	}
	ch := new(externalsMocks.AMQPChannel)

	ch.On("Consume", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Run(
		func(args mock.Arguments) {
			assert.True(t, args[0:len(args)-2].Is(conf.QueueName, conf.Consume.Consumer, conf.Consume.AutoAck,
				conf.Consume.Exclusive, conf.Consume.NoLocal, conf.Consume.NoWait))
		}).Once()
	repo := new(repoMocks.AMQPRepository)
	repo.On("GetChannel").Return(ch, nil).Once()

	serv := NewAMQPService(conf, repo, logrus.New())

	_, err := serv.Consume()
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	ch.AssertExpectations(t)
}

func TestAMQPService_Requeue_Success(t *testing.T) {
	conf := config.AMQP{
		QueueName: "test queue",
		Publish: config.Publish{
			Mandatory: true,
			Immediate: true,
		},
	}

	oldMsg := amqp.Delivery{
		Exchange:   "t",
		RoutingKey: "1",
	}
	newMsg := amqp.Publishing{}

	ch := new(externalsMocks.AMQPChannel)
	ch.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			require.Len(t, args, 5)
			assert.Equal(t, oldMsg.Exchange, args.Get(0))
			assert.Equal(t, oldMsg.RoutingKey, args.Get(1))
			assert.Equal(t, conf.Publish.Mandatory, args.Get(2))
			assert.Equal(t, conf.Publish.Immediate, args.Get(3))
			assert.Equal(t, newMsg, args.Get(4))
		}).Once()
	ch.On("Tx").Return(nil).Once()
	ch.On("Close").Return(nil).Once()
	ch.On("TxCommit").Return(nil).Once()

	repo := new(repoMocks.AMQPRepository)
	repo.On("GetChannel").Return(ch, nil).Once()
	repo.On("RejectDelivery", mock.Anything, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			require.Len(t, args, 2)
			assert.Equal(t, oldMsg, args.Get(0))
			assert.False(t, args.Bool(1))
		}).Once()

	serv := NewAMQPService(conf, repo, logrus.New())

	err := serv.Requeue(oldMsg, newMsg)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	ch.AssertExpectations(t)
}

func TestAMQPService_Requeue_RejectDelivery_Failure(t *testing.T) {
	ch := new(externalsMocks.AMQPChannel)
	ch.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	ch.On("Tx").Return(nil).Once()
	ch.On("Close").Return(nil).Once()
	ch.On("TxRollback").Return(nil).Once()

	repo := new(repoMocks.AMQPRepository)
	repo.On("GetChannel").Return(ch, nil).Once()
	repo.On("RejectDelivery", mock.Anything, mock.Anything).Return(fmt.Errorf("error")).Once()

	serv := NewAMQPService(config.AMQP{}, repo, logrus.New())

	err := serv.Requeue(amqp.Delivery{}, amqp.Publishing{})
	assert.Error(t, err)

	repo.AssertExpectations(t)
	ch.AssertExpectations(t)
}

func TestAMQPService_Requeue_GetChannel_Failure(t *testing.T) {
	repo := new(repoMocks.AMQPRepository)
	repo.On("GetChannel").Return(nil, fmt.Errorf("error")).Once()

	serv := NewAMQPService(config.AMQP{}, repo, logrus.New())

	err := serv.Requeue(amqp.Delivery{}, amqp.Publishing{})
	assert.Error(t, err)

	repo.AssertExpectations(t)
}

func TestAMQPService_Requeue_Tx_Failure(t *testing.T) {
	ch := new(externalsMocks.AMQPChannel)
	ch.On("Tx").Return(fmt.Errorf("error")).Once()
	ch.On("Close").Return(nil).Once()

	repo := new(repoMocks.AMQPRepository)
	repo.On("GetChannel").Return(ch, nil).Once()

	serv := NewAMQPService(config.AMQP{}, repo, logrus.New())

	err := serv.Requeue(amqp.Delivery{}, amqp.Publishing{})
	assert.Error(t, err)

	repo.AssertExpectations(t)
	ch.AssertExpectations(t)
}

func TestAMQPService_Requeue_Publish_Failure(t *testing.T) {
	ch := new(externalsMocks.AMQPChannel)
	ch.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(
		fmt.Errorf("error")).Once()
	ch.On("Tx").Return(nil).Once()
	ch.On("Close").Return(nil).Once()
	ch.On("TxRollback").Return(nil).Once()

	repo := new(repoMocks.AMQPRepository)
	repo.On("GetChannel").Return(ch, nil).Once()

	serv := NewAMQPService(config.AMQP{}, repo, logrus.New())

	err := serv.Requeue(amqp.Delivery{}, amqp.Publishing{})
	assert.Error(t, err)

	repo.AssertExpectations(t)
	ch.AssertExpectations(t)
}

func TestAmqpService_CreateQueue(t *testing.T) {
	conf := config.AMQP{
		QueueName: "test queue",
		Queue: config.Queue{
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       map[string]interface{}{},
		},
	}

	ch := new(externalsMocks.AMQPChannel)
	ch.On("QueueDeclare", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(amqp.Queue{}, nil).Run(
		func(args mock.Arguments) {
			require.Len(t, args, 6)
			assert.True(t, args[0:5].Is(conf.QueueName, conf.Queue.Durable, conf.Queue.AutoDelete, conf.Queue.Exclusive, conf.Queue.NoWait, conf.Queue.Args))
		}).Once()
	ch.On("Close").Return(nil).Once()
	repo := new(repoMocks.AMQPRepository)
	repo.On("GetChannel").Return(ch, nil).Once()
	serv := NewAMQPService(conf, repo, logrus.New())

	err := serv.CreateQueue()
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	ch.AssertExpectations(t)
}

func TestAMQPService_Send_GetChannel_Failure(t *testing.T) {
	repo := new(repoMocks.AMQPRepository)
	repo.On("GetChannel").Return(nil, fmt.Errorf("error")).Once()

	serv := NewAMQPService(config.AMQP{}, repo, logrus.New())

	err := serv.Send(amqp.Publishing{})
	assert.Error(t, err)

	repo.AssertExpectations(t)
}

func TestAMQPService_Send_Success(t *testing.T) {
	ch := new(externalsMocks.AMQPChannel)
	ch.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	ch.On("Close").Return(nil).Once()

	repo := new(repoMocks.AMQPRepository)
	repo.On("GetChannel").Return(ch, nil).Once()

	serv := NewAMQPService(config.AMQP{}, repo, logrus.New())

	err := serv.Send(amqp.Publishing{})
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}
