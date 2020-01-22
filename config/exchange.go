package config

import (
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

// Exchange is the configuration for the exchange to be created/used
type Exchange struct {
	Name string `mapstructure:"exchangeName"`

	Kind string `mapstructure:"exchangeKind"`

	Durable bool `mapstructure:"exchangeDurable"`
	// AutoDelete tells the queue to drop messages if there are not any consumers
	AutoDelete bool `mapstructure:"exchangeAutoDelete"`
	// Exclusive queues are only accessible by the connection that declares them
	Internal bool `mapstructure:"exchangeInternal"`
	// NoWait is true, the queue will assume to be declared on the server
	NoWait bool `mapstructure:"exchangeNoWait"`

	Args amqp.Table
}

// NewExchange gets Exchange from the values in viper
func NewExchange(v *viper.Viper) (out Exchange, err error) {
	return out, v.Unmarshal(&out)
}

func SetupExchange(v *viper.Viper) {
	v.BindEnv("exchangeName", "EXCHANGE_NAME")
	v.BindEnv("exchangeKind", "EXCHANGE_KIND")
	v.BindEnv("exchangeDurable", "EXCHANGE_DURABLE")
	v.BindEnv("exchangeAutoDelete", "EXCHANGE_AUTO_DELETE")
	v.BindEnv("exchangeInternal", "EXCHANGE_INTERNAL")
	v.BindEnv("exchangeNoWait", "EXCHANGE_NO_WAIT")

	v.SetDefault("exchangeName", "")
	v.SetDefault("exchangeKind", "")
	v.SetDefault("exchangeDurable", true)
	v.SetDefault("exchangeAutoDelete", false)
	v.SetDefault("exchangeInternal", false)
	v.SetDefault("exchangeNoWait", false)
}
