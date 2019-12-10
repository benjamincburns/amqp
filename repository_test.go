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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	externalsMock "github.com/whiteblock/genesis/mocks/pkg/externals"
	"testing"
)

func TestNewAMQPRepository(t *testing.T) {
	repo := NewAMQPRepository(new(externalsMock.AMQPConnection))
	assert.NotNil(t, repo)
}

func TestAMQPRepository_GetChannel(t *testing.T) {
	conn := new(externalsMock.AMQPConnection)
	conn.On("Channel").Return(nil, nil).Once()
	repo := NewAMQPRepository(conn)

	ch, err := repo.GetChannel()
	assert.NoError(t, err)
	assert.Nil(t, ch)
}

func TestAMQPRepository_RejectDelivery(t *testing.T) {
	msg := new(externalsMock.AMQPDelivery)
	msg.On("Reject", mock.Anything).Return(nil).Run(
		func(args mock.Arguments) {
			require.Len(t, args, 1)
			assert.Equal(t, true, args.Get(0))
		}).Once()
	conn := new(externalsMock.AMQPConnection)
	repo := NewAMQPRepository(conn)
	require.NotNil(t, repo)

	err := repo.RejectDelivery(msg, true)
	assert.NoError(t, err)
	msg.AssertExpectations(t)

}
