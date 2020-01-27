package queue

import (
	"sync"
	"time"

	"github.com/whiteblock/amqp/config"
	"github.com/whiteblock/amqp/externals"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// AMQPService acts as a simple interface to the command queue
type AMQPService interface {
	//Consume immediately starts delivering queued messages.
	Consume() (<-chan amqp.Delivery, error)
	//send places a message into the queue
	Send(pub amqp.Publishing) error
	//Requeue rejects the oldMsg and queues the newMsg in a transaction
	Requeue(oldMsg amqp.Delivery, newMsg amqp.Publishing) error
	//CreateQueue attempts to publish a queue
	CreateQueue() error

	// CreateExchange will attempt to create an exchange
	CreateExchange() error

	// Channel gets a channel
	Channel() (externals.AMQPChannel, error)

	Config() config.Config
}

type qpService struct {
	repo AMQPRepository
	conf config.Config
	log  logrus.Ext1FieldLogger

	closeChan  chan *amqp.Error
	failedChan chan *amqp.Error
	blockChan  chan amqp.Blocking
	mux        *sync.Mutex
}

// NewService is a more convienent form of NewAMQPService. It calls log.Panic on error
func NewService(conf config.Config, log logrus.Ext1FieldLogger) AMQPService {
	conn, err := OpenAMQPConnection(conf.Endpoint)
	if err != nil {
		log.Panic(err)
	}
	log.WithFields(logrus.Fields{
		"amqpHost": conf.Endpoint.QueueHost,
		"amqpUser": conf.Endpoint.QueueUser,
		"amqpPort": conf.Endpoint.QueuePort,
	}).Info("opened an amqp connection")
	return NewAMQPService(conf, NewAMQPRepository(conn), log)
}

// NewAMQPService creates a new AMQPService
func NewAMQPService(
	conf config.Config,
	repo AMQPRepository,
	log logrus.Ext1FieldLogger) AMQPService {

	out := &qpService{repo: repo, conf: conf, log: log, mux: &sync.Mutex{}}
	out.registerCallbacks()
	return out
}

func (as *qpService) registerCallbacks() {
	if as.closeChan == nil {
		as.closeChan = make(chan *amqp.Error)
	}
	if as.blockChan == nil {
		as.blockChan = make(chan amqp.Blocking)
	}

	if as.failedChan == nil {
		as.failedChan = make(chan *amqp.Error)
	}

	as.repo.AddListeners(as.closeChan, as.blockChan)

	go as.handleClose()
	go as.handleBlocked()

}

func (as *qpService) handleClose() {
	for {
		var (
			amqpErr *amqp.Error
			ok      bool
		)
		select {
		case amqpErr, ok = <-as.closeChan:
			if !ok {
				as.mux.Lock()
				as.closeChan = make(chan *amqp.Error)
				as.mux.Unlock()

				as.log.WithField("queue", as.conf.QueueName).Error("close notify channel closed")
			}
		case amqpErr = <-as.failedChan:
		}
		if amqpErr != nil {
			as.log.WithFields(logrus.Fields{
				"error": amqpErr,
			}).Warn("detected a closed connection")
		}

		as.log.Info("attempting to reconnect")

		i := 0
		for ; i < as.conf.Queue.ReconnRetries; i++ {
			err := as.attemptReconnect()
			if err == nil {
				break
			}
			as.log.WithFields(logrus.Fields{
				"error":   err,
				"queue":   as.conf.QueueName,
				"attempt": i,
			}).Error("failed to reconnect")
		}
		if i == as.conf.Queue.ReconnRetries {
			as.log.WithFields(logrus.Fields{
				"queue": as.conf.QueueName}).Fatal("could not connect to rabbitmq")
		}
	}
}

func (as *qpService) handleBlocked() {
	for {
		blocking, ok := <-as.blockChan
		if !ok {
			as.mux.Lock()
			as.blockChan = make(chan amqp.Blocking)
			as.mux.Unlock()

			as.log.WithField("queue", as.conf.QueueName).Error("block notify channel closed")
		}
		as.log.WithFields(logrus.Fields{
			"blocking": blocking.Active,
			"reason":   blocking.Reason,
		}).Warn("detected a change in connection blocking")
	}
}

func (as *qpService) attemptReconnect() (out error) {
	for i := 0; i < as.conf.Queue.ReconnRetries; i++ {
		conn, err := OpenAMQPConnection(as.conf.Endpoint)
		if err != nil {
			as.log.WithFields(logrus.Fields{
				"error": err,
				"host":  as.conf.Endpoint.QueueHost,
				"try":   i,
			}).Warn("could not connect to Queue")
			out = err
			continue
		}
		as.repo.SwapConn(conn)
		as.repo.AddListeners(as.closeChan, as.blockChan)
		return nil
	}
	return errors.Wrap(out, "could not reconnect to queue")
}

func (as *qpService) send(pub amqp.Publishing) error {
	ch, err := as.repo.GetChannel()
	if err != nil {
		as.log.WithFields(logrus.Fields{
			"queue": as.conf.QueueName,
			"error": err}).Error("had an unrecoverable error")
		return err
	}
	defer ch.Close()
	as.log.WithFields(logrus.Fields{
		"exchange": as.conf.Publish.Exchange,
		"queue":    as.conf.QueueName,
	}).Trace("publishing a message")
	return ch.Publish(as.conf.Publish.Exchange, as.conf.QueueName,
		as.conf.Publish.Mandatory, as.conf.Publish.Immediate, pub)
}

func (as *qpService) Send(pub amqp.Publishing) (err error) {
	for i := 0; i <= as.conf.Queue.Retries; i++ {
		err = as.send(pub)
		if err == nil {
			return
		}
		if i == 0 {
			as.mux.Lock()
			as.closeChan <- new(amqp.Error)
			as.mux.Unlock()
		}
		as.log.WithFields(logrus.Fields{
			"queue":   as.conf.QueueName,
			"attempt": i,
			"error":   err.Error(),
		}).Warn("unable to send a message to the queue")
		time.Sleep(as.conf.Queue.RetryDelay)
	}
	return
}

func (as *qpService) consume() (<-chan amqp.Delivery, error) {
	ch, err := as.repo.GetChannel()
	if err != nil {
		return nil, err
	}
	as.log.WithFields(logrus.Fields{
		"queue":    as.conf.QueueName,
		"consumer": as.conf.Consume.Consumer,
	}).Trace("consuming")
	return ch.Consume(as.conf.QueueName, as.conf.Consume.Consumer, as.conf.Consume.AutoAck,
		as.conf.Consume.Exclusive, as.conf.Consume.NoLocal, as.conf.Consume.NoWait,
		as.conf.Consume.Args)
}

// Consume immediately starts delivering queued messages.
func (as *qpService) Consume() (del <-chan amqp.Delivery, err error) {
	for i := 0; i <= as.conf.Queue.Retries; i++ {
		del, err = as.consume()
		if err == nil {
			return
		}
		if i == 0 {
			as.mux.Lock()
			as.closeChan <- new(amqp.Error)
			as.mux.Unlock()
		}
		as.log.WithFields(logrus.Fields{
			"queue":   as.conf.QueueName,
			"attempt": i,
			"error":   err.Error(),
		}).Warn("unable to start consuming")
		time.Sleep(as.conf.Queue.RetryDelay)
	}
	return
}

func (as *qpService) requeue(oldMsg amqp.Delivery, newMsg amqp.Publishing) error {
	ch, err := as.repo.GetChannel()
	if err != nil {
		return err
	}
	defer ch.Close()
	err = ch.Tx()
	if err != nil {
		return err
	}

	err = ch.Publish(oldMsg.Exchange, oldMsg.RoutingKey, as.conf.Publish.Mandatory, as.conf.Publish.Immediate, newMsg)
	if err != nil {
		ch.TxRollback()
		return err
	}

	err = as.repo.RejectDelivery(oldMsg, false)
	if err != nil {
		ch.TxRollback()
		return err
	}
	return ch.TxCommit()
}

// Requeue rejects the oldMsg and queues the newMsg in a transaction
func (as *qpService) Requeue(oldMsg amqp.Delivery, newMsg amqp.Publishing) (err error) {
	for i := 0; i <= as.conf.Queue.Retries; i++ {
		err = as.requeue(oldMsg, newMsg)
		if err == nil {
			return
		}
		if i == 0 {
			as.mux.Lock()
			as.closeChan <- new(amqp.Error)
			as.mux.Unlock()
		}
		as.log.WithFields(logrus.Fields{
			"queue":   as.conf.QueueName,
			"attempt": i,
			"error":   err.Error(),
		}).Warn("unable to requeue a message")
		time.Sleep(as.conf.Queue.RetryDelay)
	}
	return
}

// CreateQueue attempts to publish a queue
func (as *qpService) CreateQueue() error {
	ch, err := as.repo.GetChannel()
	if err != nil {
		return err
	}
	defer ch.Close()
	_, err = ch.QueueDeclare(as.conf.QueueName, as.conf.Queue.Durable, as.conf.Queue.AutoDelete,
		as.conf.Queue.Exclusive, as.conf.Queue.NoWait, as.conf.Queue.Args)
	return err
}

// CreateExchange will attempt to create an exchange
func (as qpService) CreateExchange() error {
	if as.conf.Exchange.Name == "" {
		return nil
	}
	ch, err := as.repo.GetChannel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.ExchangeDeclare(as.conf.Exchange.Name, as.conf.Exchange.Kind, as.conf.Exchange.Durable,
		as.conf.Exchange.AutoDelete, as.conf.Exchange.Internal, as.conf.Exchange.NoWait, as.conf.Exchange.Args)
}


// GetChannel gets a channel
func (as qpService) Channel() (externals.AMQPChannel, error) {
	return as.repo.GetChannel()
}

func (as qpService) Config() config.Config {
	return as.conf
}