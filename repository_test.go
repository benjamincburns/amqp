package queue

import (
	"testing"

	"github.com/whiteblock/amqp/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewAMQPRepository(t *testing.T) {
	repo := NewAMQPRepository(new(mocks.AMQPConnection))
	assert.NotNil(t, repo)
}

func TestAMQPRepository_GetChannel(t *testing.T) {
	conn := new(mocks.AMQPConnection)
	conn.On("Channel").Return(nil, nil).Once()
	repo := NewAMQPRepository(conn)

	ch, err := repo.GetChannel()
	assert.NoError(t, err)
	assert.Nil(t, ch)
}

func TestAMQPRepository_RejectDelivery(t *testing.T) {
	msg := new(mocks.AMQPDelivery)
	msg.On("Reject", mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			require.Len(t, args, 1)
			assert.Equal(t, true, args.Get(0))
		}).Once()
	conn := new(mocks.AMQPConnection)
	repo := NewAMQPRepository(conn)
	require.NotNil(t, repo)

	err := repo.RejectDelivery(msg, true)
	assert.NoError(t, err)
	msg.AssertExpectations(t)

}
