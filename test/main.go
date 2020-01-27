package main

import (
	"time"

	"github.com/whiteblock/amqp"
	"github.com/whiteblock/amqp/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	config.Setup(viper.GetViper())
	conf := config.Must(viper.GetViper())
	conf.Exchange = conf.Exchange.AsXDelay()
	conf.QueueName = "test"
	conf = conf.SetExchangeName("xdelay")
	serv := queue.NewService(conf, log.New())
	queue.TryCreateQueues(log.New(), serv)
	queue.BindQueuesToExchange(log.New(), serv)
	go listen(serv)
	send(queue.NewService(conf, log.New())) //give it a different service incase the library has magic
}

func send(serv queue.AMQPService) {
	cnt := 0
	delay := time.Duration(10 * time.Second)
	delayed := false
	for {
		pub, err := queue.CreateMessage(cnt)
		if err != nil {
			log.Panic(err)
		}
		if cnt%4 != 0 {
			delayed = true
			pub.Headers["x-delay"] = int32(delay.Milliseconds())
		} else {
			delayed = false
		}
		cnt++
		log.WithFields(log.Fields{
			"delay":   delay,
			"delayed": delayed,
			"num":     cnt}).Info("publishing message")

		err = serv.Send(pub)
		if err != nil {
			log.Panic(err)
		}
		time.Sleep(3 * time.Second)

	}
}

func listen(serv queue.AMQPService) {
	msgs, err := serv.Consume()
	if err != nil {
		log.Panic(err)
	}
	for msg := range msgs {
		log.WithField("body", string(msg.Body)).Info("received a message")
	}
	log.Panic("channel was closed")
}
