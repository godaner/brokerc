package broker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrPublish           = errors.New("publish err")
	// ErrSubscribe         = errors.New("subscribe err")
	ErrConnectParam = errors.New("connect param err")
	// ErrConnect           = errors.New("connect err")
	ErrConnectionIsNotOK = errors.New("connection is not ok")
	ErrQOS               = errors.New("qos err")
	ErrTLS               = errors.New("tls err")
)
var CurrKafkaBroker Broker

// Message Publish or Subscribe Message
type Message struct {
	Header map[string]string
	Body   []byte
}

func (m *Message) String() string {
	return m.Marshal()
}

func (m *Message) Marshal() string {
	bs, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

func (m *Message) Unmarshal(bs []byte) {
	json.Unmarshal(bs, m)
}

// Broker
type Broker interface {
	fmt.Stringer
	Publish(topic string, msg *Message, opt ...PublishOption) error
	Subscribe(topics []string, callBack CallBack, opt ...SubscribeOption) (Subscriber, error)
	Connect() error
	Disconnect() error
}

// CallBack
type CallBack func(event Event) error

// Event
type Event interface {
	Topic() string
	Ack() error
	Message() *Message
	Context() context.Context
}

// Subscriber
type Subscriber interface {
	Unsubscribe() error
	String() string
}
