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
