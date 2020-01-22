package config

import (
	"github.com/spf13/viper"
)

// Endpoint is the configuration needed to connect to a rabbitmq vhost
type Endpoint struct {
	QueueProtocol string `mapstructure:"queueProtocol"`
	QueueUser     string `mapstructure:"queueUser"`
	QueuePassword string `mapstructure:"queuePassword"`
	QueueHost     string `mapstructure:"queueHost"`
	QueuePort     int    `mapstructure:"queuePort"`
	QueueVHost    string `mapstructure:"queueVHost"`
}

// NewEndpoint gets Endpoint from the values in viper
func NewEndpoint(v *viper.Viper) (out Endpoint, err error) {
	return out, v.Unmarshal(&out)
}

// SetupEndpoint configures viper for Endpoint
func SetupEndpoint(v *viper.Viper) {
	v.BindEnv("queueProtocol", "QUEUE_PROTOCOL")
	v.BindEnv("queueUser", "QUEUE_USER")
	v.BindEnv("queuePassword", "QUEUE_PASSWORD")
	v.BindEnv("queueHost", "QUEUE_HOST")
	v.BindEnv("queuePort", "QUEUE_PORT")
	v.BindEnv("queueVHost", "QUEUE_VHOST")

	v.SetDefault("queueProtocol", "amqp")
	v.SetDefault("queueUser", "user")
	v.SetDefault("queuePassword", "password")
	v.SetDefault("queueHost", "localhost")
	v.SetDefault("queuePort", 5672)
	v.SetDefault("queueVHost", "/test")
}
