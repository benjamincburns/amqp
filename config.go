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
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

//AMQPConfig is the configuration for AMQP
type AMQPConfig struct {
	//QueueName the name of the queue to connect to
	QueueName string
	//Queue is the configuration for the queue
	Queue Queue
	//Consume is the configuration of the consumer
	Consume Consume
	//Publish is the configuration for the publishing of messages
	Publish Publish

	Endpoint AMQPEndpoint
}

func NewAMQPConfig(v *viper.Viper) (out AMQPConfig, err error) {
	out.Queue, err = NewQueue(v)
	if err != nil{
		return
	}

	out.Consume, err = NewConsume(v)
	if err != nil{
		return
	}

	out.Publish, err = NewPublish(v)
	if err != nil{
		return
	}

	out.Endpoint, err = NewEndpoint(v)

	return
}

//Queue is the configuration of the queue to be created if the queue does not yet exist
type Queue struct {
	//Durable sets the queue to be persistent
	Durable bool `mapstructure:"queueDurable"`
	//AutoDelete tells the queue to drop messages if there are not any consumers
	AutoDelete bool `mapstructure:"queueAutoDelete"`
	//Exclusive queues are only accessible by the connection that declares them
	Exclusive bool
	//NoWait is true, the queue will assume to be declared on the server
	NoWait bool
	//Args contains addition arguments to be provided
	Args amqp.Table
}

//NewQueue gets Queue from the values in viper
func NewQueue(v *viper.Viper) (Queue, error) {
	out := Queue{
		Exclusive: false,
		NoWait:    false,
	}
	return out, v.Unmarshal(&out)
}

//Consume is the configuration of the consumer of messages from the queue
type Consume struct {
	//Consumer is the name of the consumer
	Consumer string `mapstructure:"consumer"`
	//AutoAck causes the server to acknowledge deliveries to this consumer prior to writing the delivery to the network
	AutoAck bool
	//Exclusive: when true, the server will ensure that this is the sole consumer from this queue
	//This should always be false.
	Exclusive bool
	//NoLocal is not supported by RabbitMQ
	NoLocal bool
	//NoWait: do not wait for the server to confirm the request and immediately begin deliveries
	NoWait bool `mapstructure:"consumerNoWait"`
	//Args contains addition arguments to be provided
	Args amqp.Table
}

//NewConsume gets Consume from the values in viper
func NewConsume(v *viper.Viper) (Consume, error) {
	out := Consume{
		AutoAck:   false,
		Exclusive: false,
		NoLocal:   false,
	}
	return out, v.Unmarshal(&out)
}

//Publish is the configuration for the publishing of messages
type Publish struct {
	Mandatory bool   `mapstructure:"publishMandatory"`
	Immediate bool   `mapstructure:"publishImmediate"`
	Exchange  string `mapstructure:"exchange"`
}

//GetPublish gets Publish from the values in viper
func NewPublish(v *viper.Viper) (out Publish,err error) {
	return out, v.Unmarshal(&out)
}

//AMQPEndpoint is the configuration needed to connect to a rabbitmq vhost
type AMQPEndpoint struct {
	QueueProtocol string `mapstructure:"queueProtocol"`
	QueueUser     string `mapstructure:"queueUser"`
	QueuePassword string `mapstructure:"queuePassword"`
	QueueHost     string `mapstructure:"queueHost"`
	QueuePort     int    `mapstructure:"queuePort"`
	QueueVHost    string `mapstructure:"queueVHost"`
}

//GetAMQPEndpoint gets AMQPEndpoint from the values in viper
func GetAMQPEndpoint(v *viper.Viper) (out AMQPEndpoint, err error) {
	return out, v.Unmarshal(&out)
}

func SetConfig(v *viper.Viper) {
	/** START Queue **/
	v.BindEnv("queueDurable", "QUEUE_DURABLE")
	v.BindEnv("queueAutoDelete", "QUEUE_AUTO_DELETE")
	v.SetDefault("queueDurable", true)
	v.SetDefault("queueAutoDelete", false)
	/** END Queue **/

	/** START Consume **/
	v.BindEnv("consumer", "CONSUMER")
	v.BindEnv("consumerNoWait", "CONSUMER_NO_WAIT")
	v.SetDefault("consumer", "genesis")
	v.SetDefault("consumerNoWait", false)
	/** END Consume **/

	/** START Publish **/
	v.BindEnv("publishMandatory", "PUBLISH_MANDATORY")
	v.BindEnv("publishImmediate", "PUBLISH_IMMEDIATE")
	v.BindEnv("exchange", "EXCHANGE")
	v.SetDefault("exchange", "")
	v.SetDefault("publishMandatory", false)
	v.SetDefault("publishImmediate", false)
	/** END Publish **/

	/** START AMQPEndpoint **/
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
	/** END AMQPEndpoint **/
}