package config

import (
	"time"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

// Queue is the configuration of the queue to be created if the queue does not yet exist
type Queue struct {
	// Durable sets the queue to be persistent
	Durable bool `mapstructure:"queueDurable"`
	// AutoDelete tells the queue to drop messages if there are not any consumers
	AutoDelete bool `mapstructure:"queueAutoDelete"`
	// Exclusive queues are only accessible by the connection that declares them
	Exclusive bool
	// NoWait is true, the queue will assume to be declared on the server
	NoWait bool
	// Args contains addition arguments to be provided
	Args amqp.Table

	// ReConnRetries is the number of times to attempt to re-connect to a queue
	// if a connection error is found
	ReconnRetries int `mapstructure:"queueReconnRetries`

	// Retries is the number of times to retry an action before giving up
	Retries int `mapstructure:"queueRetries"`

	// RetryDelay is the delay between retries
	RetryDelay time.Duration `mapstructure:"queueRetryDelay"`
}

// NewQueue gets Queue from the values in viper
func NewQueue(v *viper.Viper) (Queue, error) {
	out := Queue{
		Exclusive: false,
		NoWait:    false,
	}
	return out, v.Unmarshal(&out)
}

func SetupQueue(v *viper.Viper) {
	v.BindEnv("queueReconnRetries", "QUEUE_RECONN_RETRIES")
	v.BindEnv("queueRetries", "QUEUE_RETRIES")
	v.BindEnv("queueRetryDelay", "QUEUE_RETRY_DELAY")
	v.BindEnv("queueDurable", "QUEUE_DURABLE")
	v.BindEnv("queueAutoDelete", "QUEUE_AUTO_DELETE")

	v.SetDefault("queueDurable", true)
	v.SetDefault("queueAutoDelete", false)
	v.SetDefault("queueReconnRetries", 10)
	v.SetDefault("queueRetries", 10)
	v.SetDefault("queueRetryDelay", 2*time.Second)
}
