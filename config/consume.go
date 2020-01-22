package config

import (
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

// Consume is the configuration of the consumer of messages from the queue
type Consume struct {
	// Consumer is the name of the consumer
	Consumer string `mapstructure:"consumer"`
	// AutoAck causes the server to acknowledge deliveries to this consumer prior to writing the delivery to the network
	AutoAck bool
	// Exclusive: when true, the server will ensure that this is the sole consumer from this queue
	// This should always be false.
	Exclusive bool
	// NoLocal is not supported by RabbitMQ
	NoLocal bool
	// NoWait: do not wait for the server to confirm the request and immediately begin deliveries
	NoWait bool `mapstructure:"consumerNoWait"`
	// Args contains addition arguments to be provided
	Args amqp.Table
}

// NewConsume gets Consume from the values in viper
func NewConsume(v *viper.Viper) (Consume, error) {
	out := Consume{
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
	}
	return out, v.Unmarshal(&out)
}

func SetupConsume(v *viper.Viper) {
	v.BindEnv("consumer", "CONSUMER")
	v.BindEnv("consumerNoWait", "CONSUMER_NO_WAIT")

	v.SetDefault("consumer", "genesis")
	v.SetDefault("consumerNoWait", false)
}
