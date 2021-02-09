package amqpv1

import (
	"encoding/json"
	"fmt"
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/log"
	"github.com/godaner/brokerc/tls"
	AMQP "github.com/streadway/amqp"
	"golang.org/x/net/context"
	"sync"
	"time"
)

const (
	reconnectInterval = time.Duration(5) * time.Second
)

type AMQPBrokerV1 struct {
	sync.Once
	URI         string // amqp[s]://[username][:password]@host.domain[:port][vhost]
	CID         string // client id
	CACertFile  string
	CertFile    string
	KeyFile     string
	Insecure    bool
	Logger      log.Logger
	conn        *AMQP.Connection
	publisherCh *AMQP.Channel // just for publisher
	once        sync.Once
	subscribers []*amqpSubscriber
}

func (a *AMQPBrokerV1) String() string {
	return a.Marshal()
}

func (a *AMQPBrokerV1) Marshal() string {
	bs, err := json.Marshal(a)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

func (a *AMQPBrokerV1) Connect() error {
	a.Do(func() {
		a.subscribers = make([]*amqpSubscriber, 0)
	})
	a.Logger.Debugf("AMQPBrokerV1#connect : info is : %v !", a)
	t, err := tls.GetClientTLSConfig(a.Insecure, a.CACertFile, a.CertFile, a.KeyFile)
	if err != nil {
		return err
	}
	// conn
	a.conn, err = AMQP.DialConfig(a.URI, AMQP.Config{
		TLSClientConfig: t,
	})
	if err != nil {
		return err
	}

	// add listener
	closeSig := make(chan *AMQP.Error, 1)
	go a.listenClose(closeSig)
	a.conn.NotifyClose(closeSig)
	return nil
}

func (a *AMQPBrokerV1) Disconnect() error {
	if a.publisherCh != nil {
		err := a.publisherCh.Close()
		if err != nil {
			return err
		}
	}
	for _, s := range a.subscribers {
		err := s.Unsubscribe()
		if err != nil {
			return err
		}
	}
	err := a.conn.Close()
	if err != nil {
		return err
	}
	return nil
}

// Publish
func (a *AMQPBrokerV1) Publish(topic string, msg *broker.Message, opt ...broker.PublishOption) (err error) {
	// a.init()
	if a.conn == nil || a.conn.IsClosed() {
		return broker.ErrConnectionIsNotOK
	}
	// default opt
	opts := broker.PublishOptions{}
	for _, o := range opt {
		o(&opts)
	}

	// publisherCh
	if a.publisherCh == nil {
		a.publisherCh, err = a.conn.Channel()
		if err != nil {
			return err
		}
	}

	// exchange
	if opts.ExchangeName != "" && opts.ExchangeType != "" {
		a.Logger.Debugf("AMQPBrokerV1#Publish : dec exchange , name is : %v , type is : %v !", opts.ExchangeName, opts.ExchangeType)
		err = a.publisherCh.ExchangeDeclare(
			opts.ExchangeName,     // name
			opts.ExchangeType,     // type
			opts.ExchangeDuration, // durable
			opts.ExchangeAD,       // auto-deleted
			false,                 // internal
			false,                 // no-wait
			nil,                   // arguments
		)
		if err != nil {
			return err
		}
	}

	// Publish
	err = a.publisherCh.Publish(
		opts.ExchangeName, // exchange
		topic,             // routing key
		false,             // mandatory
		false,             // immediate
		AMQP.Publishing{
			ContentType: "text/plain",
			Body:        msg.Body,
		})
	if err != nil {
		return err
	}

	return nil
}

// Subscribe
func (a *AMQPBrokerV1) Subscribe(topics []string, callBack broker.CallBack, opt ...broker.SubscribeOption) (broker.Subscriber, error) {
	s := &amqpSubscriber{
		topics:   topics,
		callBack: callBack,
		opt:      opt,
		broker:   a,
	}
	a.subscribers = append(a.subscribers, s)
	return s, s.subscribe()
}

// ////////////// Sub Method ///////////////

// amqpSubscriber
type amqpSubscriber struct {
	ch       *AMQP.Channel
	topics   []string
	callBack broker.CallBack
	opt      []broker.SubscribeOption
	opts     broker.SubscribeOptions
	broker   *AMQPBrokerV1
}

// String
func (s *amqpSubscriber) String() string {
	return fmt.Sprintf("/amqpbroker/subscriber/%v/%v/%v", s.topics, s.opts.ExchangeName, s.opts.Queue)
}

// Unsubscribe
func (s *amqpSubscriber) Unsubscribe() error {
	if s.ch != nil {
		s.ch.Close()
		s.ch = nil
	}
	return nil
}

// subscribe do sub
func (s *amqpSubscriber) subscribe() (err error) {
	if s.broker.conn == nil || s.broker.conn.IsClosed() {
		return broker.ErrConnectionIsNotOK
	}
	this := s.broker
	opt := s.opt
	callBack := s.callBack
	topics := s.topics
	// default opt
	opts := broker.SubscribeOptions{
		AutoAck: true,
	}
	for _, o := range opt {
		o(&opts)
	}
	s.opts = opts
	s.broker.Logger.Debugf("AMQPBrokerV1#subscribe : subscribe topics is : %v , opts is : %v !", s.topics, s.opts)

	// get self channel
	s.ch, err = this.conn.Channel()
	if err != nil {
		return err
	}
	if s.opts.ExchangeName != "" {
		err = s.ch.ExchangeDeclare(
			s.opts.ExchangeName,     // name
			s.opts.ExchangeType,     // type
			s.opts.ExchangeDuration, // durable
			s.opts.ExchangeAD,       // auto-deleted
			false,                   // internal
			false,                   // no-wait
			nil,                     // arguments
		)
		if err != nil {
			return err
		}
	}
	// queue
	q, err := s.ch.QueueDeclare(
		s.opts.Queue,    // name
		s.opts.Duration, // durable
		s.opts.AutoDel,  // delete when usused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}
	if len(topics) == 0 {
		err = s.ch.QueueBind(q.Name, "", s.opts.ExchangeName, false, nil)
		if err != nil {
			return err
		}
	} else {
		for _, topic := range topics {
			err = s.ch.QueueBind(q.Name, topic, s.opts.ExchangeName, false, nil)
			if err != nil {
				return err
			}
		}
	}
	msgs, err := s.ch.Consume(
		q.Name,         // queue
		s.broker.CID,   // consumer
		s.opts.AutoAck, // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	if err != nil {
		return err
	}
	go func() {
		for msg := range msgs {
			if callBack != nil {
				header := make(map[string][]string)
				for k, v := range msg.Headers {
					vs, _ := v.(string)
					header[k] = []string{vs}
				}
				e := &amqpEvent{
					autoAck: s.opts.AutoAck,
					d:       AMQP.Delivery{},
					topic:   msg.RoutingKey,
					cxt:     context.Background(),
					m: &broker.Message{
						Header: header,
						Body:   msg.Body,
					},
				}
				err = callBack(e)
				if err != nil {
					s.broker.Logger.Errorf("AMQPBrokerV1#subscribe : callBack err , err is : %v !", err)
				}
			}
		}
	}()
	return nil
}

func (a *AMQPBrokerV1) listenClose(c chan *AMQP.Error) {
	for err := range c {
		a.Logger.Errorf("AMQPBrokerV1#listenClose : conn close error , err is : %v , reason is :%v , code is : %v , code is : %v !", err.Error(), err.Reason, err.Code)
	}
	a.Logger.Debug("AMQPBrokerV1#listenClose : we will start reconnect !")
	for {
		err := a.Connect()
		if err == nil {
			for _, s := range a.subscribers {
				err := s.subscribe()
				if err != nil {
					a.Logger.Errorf("AMQPBrokerV1#listenClose : resubscribe err , err is : %v !", err.Error())
				}
			}
			break
		}
		<-time.After(reconnectInterval)
	}
}

// amqpEvent
type amqpEvent struct {
	autoAck bool
	d       AMQP.Delivery
	topic   string
	cxt     context.Context
	m       *broker.Message
}

func (e *amqpEvent) Ack() error {
	if e.autoAck {
		return nil
	}
	return e.d.Ack(false)
}

func (e *amqpEvent) Topic() string {
	return e.topic
}

func (e *amqpEvent) Message() *broker.Message {
	return e.m
}

func (e *amqpEvent) Context() context.Context {
	return e.cxt
}
