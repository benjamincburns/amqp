package main

import (
	"time"

	"github.com/whiteblock/amqp"
	"github.com/whiteblock/amqp/config"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getServ(name string) queue.AMQPService {
	conf := config.Must(viper.GetViper())
	conf.Exchange = conf.Exchange.AsXDelay()
	conf.QueueName = name
	conf = conf.SetExchangeName("xdelay")
	serv := queue.NewService(conf, log.New())
	queue.TryCreateQueues(log.New(), serv)
	queue.BindQueuesToExchange(log.New(), serv)
	return serv
}

func main() {
	config.Setup(viper.GetViper())
	serv := getServ("test")
	go listen(serv)

	go send(getServ("test")) //give it a different service incase the library has magic

	serv2 := getServ("test2")

	go listen(serv2)

	send(serv2)
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
		msg.Ack(false)
	}
	log.Panic("channel was closed")
}
