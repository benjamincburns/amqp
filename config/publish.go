package config

import (
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

// Publish is the configuration for the publishing of messages
type Publish struct {
	Mandatory bool   `mapstructure:"publishMandatory"`
	Immediate bool   `mapstructure:"publishImmediate"`
	Exchange  string `mapstructure:"exchangeName"`
}

// NewPublish gets Publish from the values in viper
func NewPublish(v *viper.Viper) (out Publish, err error) {
	return out, v.Unmarshal(&out)
}

func SetupPublish(v *viper.Viper) {
	v.BindEnv("publishMandatory", "PUBLISH_MANDATORY")
	v.BindEnv("publishImmediate", "PUBLISH_IMMEDIATE")
	v.SetDefault("publishMandatory", false)
	v.SetDefault("publishImmediate", false)
}
