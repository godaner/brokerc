package mqttv1

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/log"
	"github.com/google/uuid"
	"golang.org/x/net/context"
	"io/ioutil"
	"sync"
	"time"
)

type MQTTBrokerV1 struct {
	sync.Once
	Host           string // localhost , ssl://localhost
	Port           string
	Username       string
	Password       string
	CID            string // client id
	WT             string // will topic
	WP             string // will payload
	WR             bool   // will retain
	WQ             byte   // will qos
	C              bool   // clean session , for subscribe
	CACertFile     string
	ClientCertFile string
	ClientKeyFile  string
	Insecure       bool
	Logger         log.Logger
	subscribers    *sync.Map
	c              MQTT.Client
}

func (s *MQTTBrokerV1) Connect() error {
	s.subscribers = &sync.Map{}
	s.Logger.Debugf("MQTTBrokerV1#Connect : info is : %v !", s)
	// opts
	opts := MQTT.NewClientOptions()
	if s.Host == "" || s.Port == "" {
		return broker.ErrConnectParam
	}
	opts.AddBroker(s.Host + ":" + s.Port)
	cid := uuid.New().String()
	if s.CID != "" {
		cid = s.CID
	}
	opts.SetClientID(cid)
	opts.SetCleanSession(s.C)
	if s.Username != "" {
		opts.SetUsername(s.Username)
	}
	if s.Password != "" {
		opts.SetPassword(s.Password)
	}
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)

	opts.OnConnect = s.mqttConnectEvent
	opts.OnConnectionLost = s.mqttConnectionLostEvent
	if s.WT != "" && s.WP != "" {
		opts.WillEnabled = true
		opts.WillPayload = []byte(s.WP)
		opts.WillRetained = s.WR
		opts.WillQos = s.WQ
		opts.WillTopic = s.WT
	}
	t, err := s.getTLSConfig()
	if err != nil {
		return err
	}
	opts.TLSConfig = t
	// NewClient
	s.c = MQTT.NewClient(opts)
	if token := s.c.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (s *MQTTBrokerV1) getTLSConfig() (t *tls.Config, err error) {
	t = &tls.Config{
		InsecureSkipVerify: s.Insecure,
	}
	if s.CACertFile != "" {
		var ca *x509.CertPool
		data, err := ioutil.ReadFile(s.CACertFile)
		if err != nil {
			return nil, err
		}
		ca = x509.NewCertPool()
		if ok := ca.AppendCertsFromPEM(data); !ok {
			return nil, broker.ErrTLS
		}
		t.RootCAs = ca
	}
	if s.ClientCertFile != "" && s.ClientKeyFile != "" {
		crt, err := tls.LoadX509KeyPair(s.ClientCertFile, s.ClientKeyFile)
		if err != nil {
			return nil, err
		}
		t.Certificates = []tls.Certificate{crt}
	}
	return t, nil
}

func (s *MQTTBrokerV1) Disconnect() error {
	if s.c == nil {
		return nil
	}
	s.c.Disconnect(250)
	return nil
}

func (s *MQTTBrokerV1) String() string {
	return s.Marshal()
}

func (s *MQTTBrokerV1) Marshal() string {
	bs, err := json.Marshal(s)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

// Publish
func (s *MQTTBrokerV1) Publish(topic string, msg *broker.Message, opt ...broker.PublishOption) error {
	opts := broker.PublishOptions{
		QOS:      0,
		Retained: false,
	}
	for _, o := range opt {
		o(&opts)
	}
	if opts.QOS < 0 || opts.QOS > 2 {
		return broker.ErrQOS
	}
	for i := 0; i < 1; i++ {
		token := s.c.Publish(topic, byte(opts.QOS), opts.Retained, string(msg.Body))
		if !token.Wait() {
			return token.Error()
		} else {
			return nil
		}
	}
	return nil
}

// Subscribe
func (s *MQTTBrokerV1) Subscribe(topics []string, callBack broker.CallBack, opt ...broker.SubscribeOption) (broker.Subscriber, error) {
	subscriber := &mqttSubscriber{
		id:       uuid.NewString(),
		sub:      false,
		topics:   topics,
		callBack: callBack,
		opt:      opt,
		opts: broker.SubscribeOptions{ // default options
			QOS: 0,
		},
		broker: s,
	}
	s.subscribers.Store(subscriber.id, subscriber)
	return subscriber, subscriber.subscribe()
}
func (s *MQTTBrokerV1) rmSubscriber(id string) {

}

func (s *MQTTBrokerV1) mqttConnectEvent(client MQTT.Client) {
	s.Logger.Debug("MQTTBrokerV1#mqttConnectEvent : connect connect connect connect !")
	s.subscribers.Range(func(key, value interface{}) bool { // reconnect
		value.(*mqttSubscriber).mqttConnectEvent(client)
		return true
	})
}

func (s *MQTTBrokerV1) mqttConnectionLostEvent(client MQTT.Client, err error) {
	s.Logger.Debugf("MQTTBrokerV1#mqttConnectionLostEvent : connection lost connection lost connection lost connection lost , err is : %v !", err)

	s.subscribers.Range(func(key, value interface{}) bool { // reconnect
		value.(*mqttSubscriber).mqttConnectionLostEvent(client, err)
		return true
	})
}

// mqttEvent
type mqttEvent struct {
	topic string
	cxt   context.Context
	m     *broker.Message
}

func (e *mqttEvent) Ack() error {
	return nil
}

func (e *mqttEvent) Topic() string {
	return e.topic
}

func (e *mqttEvent) Message() *broker.Message {
	return e.m
}

func (e *mqttEvent) Context() context.Context {
	return e.cxt
}

// String
func (s *mqttSubscriber) String() string {
	return fmt.Sprintf("/mqttbroker/subscriber/%v", s.topics)
}

// mqttSubscriber
type mqttSubscriber struct {
	sync.Once
	sync.Mutex
	id       string
	sub      bool
	topics   []string
	callBack broker.CallBack
	opt      []broker.SubscribeOption
	opts     broker.SubscribeOptions
	broker   *MQTTBrokerV1
}

func (s *mqttSubscriber) mqttConnectEvent(client MQTT.Client) {
	s.subscribe()
}

func (s *mqttSubscriber) mqttConnectionLostEvent(client MQTT.Client, err error) {
	s.Lock()
	defer s.Unlock()
	s.sub = false
}

// subscribe
func (s *mqttSubscriber) subscribe() error {
	s.Lock()
	defer s.Unlock()
	c, logger := s.broker.c, s.broker.Logger
	if s.sub {
		return nil
	}
	if c == nil || !c.IsConnected() {
		return broker.ErrConnectionIsNotOK
	}
	// default opt
	if s.opts.QOS < 0 || s.opts.QOS > 2 {
		return broker.ErrQOS
	}
	for _, o := range s.opt {
		o(&s.opts)
	}
	logger.Debugf("MQTTBrokerV1#subscribe : subscribe topics is : %v , opts is : %v !", s.topics, s.opts)
	for _, topic := range s.topics {
		if token := c.Subscribe(topic, byte(s.opts.QOS), func(client MQTT.Client, message MQTT.Message) {
			if s.callBack != nil {
				e := &mqttEvent{
					topic: topic,
					cxt:   context.Background(),
					m: &broker.Message{
						Header: make(map[string]string),
						Body:   message.Payload(),
					},
				}
				err := s.callBack(e)
				if err != nil {
					logger.Errorf("MQTTBrokerV1#subscribe : callBack err , err is : %v !", err)
				}
			}
		}); token.Wait() && token.Error() != nil {
			return token.Error()
		}
	}
	s.sub = true
	return nil
}

// Unsubscribe
func (s *mqttSubscriber) Unsubscribe() error {
	s.Lock()
	defer s.Unlock()
	c := s.broker.c
	err := c.Unsubscribe(s.topics...).Error()
	if err != nil {
		return err
	}
	s.sub = false
	s.broker.rmSubscriber(s.id)
	return nil
}
