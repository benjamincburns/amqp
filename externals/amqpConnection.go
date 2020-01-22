package externals

import (
	"github.com/streadway/amqp"
)

// AMQPConnection represents the needed functionality from a amqp.Connection
type AMQPConnection interface {
	Channel() (*amqp.Channel, error)
	Close() error
	NotifyBlocked(receiver chan amqp.Blocking) chan amqp.Blocking
	NotifyClose(receiver chan *amqp.Error) chan *amqp.Error
}
