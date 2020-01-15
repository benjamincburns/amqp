package externals

import (
//"github.com/streadway/amqp"
)

//AMQPDelivery represents the needed functionality from a amqp.Delivery
type AMQPDelivery interface {
	Reject(requeue bool) error
}
