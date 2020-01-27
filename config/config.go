package config

import (
	"github.com/spf13/viper"
)

// Config is the configuration for AMQP
type Config struct {
	// QueueName the name of the queue to connect to
	QueueName string
	// Queue is the configuration for the queue
	Queue Queue
	// Consume is the configuration of the consumer
	Consume Consume
	// Publish is the configuration for the publishing of messages
	Publish Publish
	// Endpoint is the configuration for the connection endpoint
	Endpoint Endpoint
	// Exchange is the configuration for the exchange, if it is created
	Exchange Exchange
}

func (c Config) SetExchangeName(name string) Config {
	c.Exchange.Name = name
	c.Publish.Exchange = name
	return c
}

// Must is like New but panics on error
func Must(v *viper.Viper) Config {
	conf, err := New(v)
	if err != nil {
		panic(err)
	}
	return conf
}

// New creates a new instance of Config from viper
func New(v *viper.Viper) (out Config, err error) {
	out.Queue, err = NewQueue(v)
	if err != nil {
		return
	}

	out.Consume, err = NewConsume(v)
	if err != nil {
		return
	}

	out.Publish, err = NewPublish(v)
	if err != nil {
		return
	}

	out.Endpoint, err = NewEndpoint(v)
	if err != nil {
		return
	}

	out.Exchange, err = NewExchange(v)
	return
}

// Setup sets up the given viper with the defaults and env bindings
func Setup(v *viper.Viper) {
	SetupQueue(v)
	SetupEndpoint(v)
	SetupConsume(v)
	SetupPublish(v)
	SetupExchange(v)
}
