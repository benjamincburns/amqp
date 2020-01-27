package queue

import (
	"testing"

	"github.com/whiteblock/amqp/config"
	"github.com/whiteblock/amqp/mocks"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewAMQPService(t *testing.T) {
	conf := config.Config{
		QueueName: "test queue",
	}
	repo := new(mocks.AMQPRepository)
	repo.On("AddListeners", mock.Anything, mock.Anything).Once()
	serv := NewAMQPService(conf, repo, nil)
	assert.NotNil(t, serv)
}

func TestAMQPService_Consume(t *testing.T) {
	conf := config.Config{
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
	ch := new(mocks.AMQPChannel)

	ch.On("Consume", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Run(
		func(args mock.Arguments) {
			assert.True(t, args[0:len(args)-2].Is(conf.QueueName, conf.Consume.Consumer, conf.Consume.AutoAck,
				conf.Consume.Exclusive, conf.Consume.NoLocal, conf.Consume.NoWait))
		}).Once()
	repo := new(mocks.AMQPRepository)
	repo.On("GetChannel").Return(ch, nil).Once()
	repo.On("AddListeners", mock.Anything, mock.Anything).Once()

	serv := NewAMQPService(conf, repo, logrus.New())

	_, err := serv.Consume()
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	ch.AssertExpectations(t)
}

func TestAMQPService_Requeue_Success(t *testing.T) {
	conf := config.Config{
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

	ch := new(mocks.AMQPChannel)
	ch.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			require.Len(t, args, 5)
			assert.Equal(t, conf.Exchange.Name, args.Get(0))
			assert.Equal(t, conf.QueueName, args.Get(1))
			assert.Equal(t, conf.Publish.Mandatory, args.Get(2))
			assert.Equal(t, conf.Publish.Immediate, args.Get(3))
			assert.Equal(t, newMsg, args.Get(4))
		}).Once()
	ch.On("Tx").Return(nil).Once()
	ch.On("Close").Return(nil).Once()
	ch.On("TxCommit").Return(nil).Once()

	repo := new(mocks.AMQPRepository)
	repo.On("GetChannel").Return(ch, nil).Once()
	repo.On("RejectDelivery", mock.Anything, mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			require.Len(t, args, 2)
			assert.Equal(t, oldMsg, args.Get(0))
			assert.False(t, args.Bool(1))
		}).Once()
	repo.On("AddListeners", mock.Anything, mock.Anything).Once()
	serv := NewAMQPService(conf, repo, logrus.New())

	err := serv.Requeue(oldMsg, newMsg)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	ch.AssertExpectations(t)
}

func TestAmqpService_CreateQueue(t *testing.T) {
	conf := config.Config{
		QueueName: "test queue",
		Queue: config.Queue{
			Durable:    true,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       map[string]interface{}{},
		},
	}

	ch := new(mocks.AMQPChannel)
	ch.On("QueueDeclare", mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything).Return(amqp.Queue{}, nil).Run(
		func(args mock.Arguments) {
			require.Len(t, args, 6)
			assert.True(t, args[0:5].Is(conf.QueueName, conf.Queue.Durable, conf.Queue.AutoDelete,
				conf.Queue.Exclusive, conf.Queue.NoWait, conf.Queue.Args))
		}).Once()
	ch.On("Close").Return(nil).Once()
	repo := new(mocks.AMQPRepository)
	repo.On("GetChannel").Return(ch, nil).Once()
	repo.On("AddListeners", mock.Anything, mock.Anything).Once()
	serv := NewAMQPService(conf, repo, logrus.New())

	err := serv.CreateQueue()
	assert.NoError(t, err)

	repo.AssertExpectations(t)
	ch.AssertExpectations(t)
}

func TestAMQPService_Send_Success(t *testing.T) {
	ch := new(mocks.AMQPChannel)
	ch.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
	ch.On("Close").Return(nil).Once()

	repo := new(mocks.AMQPRepository)
	repo.On("GetChannel").Return(ch, nil).Once()
	repo.On("AddListeners", mock.Anything, mock.Anything).Once()
	serv := NewAMQPService(config.Config{}, repo, logrus.New())

	err := serv.Send(amqp.Publishing{})
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}
