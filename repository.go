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
	"sync"

	"github.com/whiteblock/amqp/externals"

	"github.com/streadway/amqp"
)

//AMQPRepository represents functions for connecting to a AMQP provider
type AMQPRepository interface {
	GetChannel() (externals.AMQPChannel, error)
	RejectDelivery(msg externals.AMQPDelivery, requeue bool) error
	SwapConn(conn externals.AMQPConnection)
	AddListeners(closeChan chan *amqp.Error, blockChan chan amqp.Blocking)
}

type amqpRepository struct {
	conn externals.AMQPConnection
	mux  *sync.Mutex
}

//NewAMQPRepository creates a new AMQPRepository
func NewAMQPRepository(conn externals.AMQPConnection) AMQPRepository {
	return &amqpRepository{conn: conn, mux: &sync.Mutex{}}
}

func (ar amqpRepository) GetChannel() (externals.AMQPChannel, error) {
	ar.mux.Lock()
	defer ar.mux.Unlock()
	return ar.conn.Channel()
}

func (ar amqpRepository) RejectDelivery(msg externals.AMQPDelivery, requeue bool) error {
	return msg.Reject(requeue)
}

func (ar *amqpRepository) SwapConn(conn externals.AMQPConnection) {
	ar.mux.Lock()
	defer ar.mux.Unlock()

	ar.conn.Close()

	ar.conn = conn
}

func (ar *amqpRepository) GetConn() externals.AMQPConnection {
	return ar.conn
}

func (ar *amqpRepository) AddListeners(closeChan chan *amqp.Error,
	blockChan chan amqp.Blocking) {
	ar.mux.Lock()
	defer ar.mux.Unlock()

	ar.conn.NotifyClose(closeChan)
	ar.conn.NotifyBlocked(blockChan)
}
