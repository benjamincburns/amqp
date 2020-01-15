package externals

import (
	"github.com/streadway/amqp"
)

//AMQPChannel represents the needed functionality from a amqp.Channel
type AMQPChannel interface {
	Close() error
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
	QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error)
	Tx() error
	TxCommit() error
	TxRollback() error
}
