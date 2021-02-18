package kafkav1

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/log"
	sync2 "github.com/godaner/brokerc/sync"
	tls2 "github.com/godaner/brokerc/tls"
	"github.com/google/uuid"
	c "golang.org/x/net/context"
	"strings"
	"sync"
	"time"
)

const (
	ConsumerOffsetNewest int64 = -1
	ConsumerOffsetOldest int64 = -2
)

type TLS struct {
	Insecure   bool
	CertFile   string // 公钥,需要转化为pem
	KeyFile    string // 私钥
	CaCertFile string // ca证书
}

func (t *TLS) getTLS() (*tls.Config, error) {
	if t == nil {
		return nil, nil
	}
	return tls2.GetClientTLSConfig(t.Insecure, t.CaCertFile, t.CertFile, t.KeyFile)
}

type Sarama struct {
	ConsumerOffsetsInitial int64
}

func (s *Sarama) toSarama(c *sarama.Config) {
	if s == nil {
		return
	}
	if s.ConsumerOffsetsInitial != 0 {
		c.Consumer.Offsets.Initial = s.ConsumerOffsetsInitial
	}
}

// KafkaBrokerV1
type KafkaBrokerV1 struct {
	sync.Once
	sync2.OnceError
	Logger    log.Logger
	URIs      []string
	TLS       *TLS
	Insecure  bool
	SubSarama *Sarama
	sc        []sarama.Client // sub client
	c         sarama.Client   // pub client
	p         sarama.SyncProducer
	rem       *sync.Map
}

func (k *KafkaBrokerV1) Connect() error {
	k.Once.Do(func() {
		k.rem = &sync.Map{}
	})
	return nil
}

func (k *KafkaBrokerV1) Disconnect() error {
	if k.p != nil {
		err := k.p.Close()
		if err != nil {
			return err
		}
	}
	if k.c != nil {
		err := k.c.Close()
		if err != nil {
			return err
		}
	}
	for _, c := range k.sc {
		if c != nil {
			err := c.Close()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (k *KafkaBrokerV1) String() string {
	return k.Marshal()
}

func (k *KafkaBrokerV1) Marshal() string {
	bs, err := json.Marshal(k)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

func (k *KafkaBrokerV1) Publish(topic string, msg *broker.Message, opt ...broker.PublishOption) error {
	opts := broker.PublishOptions{
		Context: context.Background(),
	}
	for _, o := range opt {
		o(&opts)
	}
	if opts.Part != 0 && opts.Replica != 0 {
		onceErrI, _ := k.rem.LoadOrStore("topic"+topic, sync2.OnceError{})
		onceErr := onceErrI.(sync2.OnceError)
		err := onceErr.Do(func() error {
			err := k.brokersCreateTopic(k.URIs, topic, opts.Part, opts.Replica)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	err := k.OnceError.Do(func() error {
		// For implementation reasons, the SyncProducer requires
		// `Producer.Return.Errors` and `Producer.Return.Successes`
		// to be set to true in its configuration.
		config, err := k.getPubConfig()
		c, err := sarama.NewClient(k.URIs, config)
		if err != nil {
			return err
		}
		k.c = c
		p, err := sarama.NewSyncProducerFromClient(c)
		if err != nil {
			return err
		}
		k.p = p
		// k.sc = make([]sarama.Client, 0)
		return nil
	})
	if err != nil {
		return err
	}
	// notice marshal body only
	_, _, err = k.p.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(string(msg.Body)),
	})
	return err
}

func (k *KafkaBrokerV1) Subscribe(topics []string, callBack broker.CallBack, opt ...broker.SubscribeOption) (broker.Subscriber, error) {
	opts := broker.SubscribeOptions{
		AutoAck: true,
		Queue:   uuid.New().String(),
		Context: context.Background(),
	}
	for _, o := range opt {
		o(&opts)
	}
	if opts.Part != 0 && opts.Replica != 0 {
		for _, topic := range topics {
			err := k.brokersCreateTopic(k.URIs, topic, opts.Part, opts.Replica)
			if err != nil {
				return nil, err
			}
		}
	}
	// we need to create a new client per consumer
	c, err := k.getSubClusterClient()
	if err != nil {
		return nil, err
	}
	cg, err := sarama.NewConsumerGroupFromClient(opts.Queue, c)
	if err != nil {
		return nil, err
	}
	h := &consumerGroupHandler{
		handler: callBack,
		subopts: opts,
		// kopts:   k.opts,
		cg: cg,
	}
	ctx := context.Background()
	go func() {
		for {
			select {
			case err := <-cg.Errors():
				k.Logger.Errorf("KafkaBrokerV1#Subscribe : consumer1 err , err is : %v !", err)
			default:
				err := cg.Consume(ctx, topics, h)
				// if err != nil {
				// 	k.Logger.Errorf("KafkaBrokerV1#Subscribe : consumer2 err , err is : %v !", err)
				// }
				if err == sarama.ErrClosedConsumerGroup {
					return
				}

			}
		}
	}()
	return &subscriber{cg: cg, topics: topics}, nil
}
func (k *KafkaBrokerV1) getSubClusterClient() (sarama.Client, error) {
	config, err := k.getSubConfig()
	if err != nil {
		return nil, err
	}
	cs, err := sarama.NewClient(k.URIs, config)
	if err != nil {
		return nil, err
	}
	k.sc = append(k.sc, cs)
	return cs, nil
}
func (k *KafkaBrokerV1) getPubConfig() (*sarama.Config, error) {
	// if c, ok := k.opts.Context.Value(brokerConfigKey{}).(*sarama.Config); ok {
	//	return c
	// }
	config := sarama.NewConfig()
	if k.TLS != nil {
		tlsConfig, err := k.TLS.getTLS()
		if err != nil {
			return nil, err
		}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig

	}
	config.ClientID = uuid.New().String()

	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	// 设置发送/订阅时获取连接重试次数与重试间隔时间
	config.Metadata.Retry.Max = 20
	config.Metadata.Retry.Backoff = 30 * time.Second
	config.Metadata.RefreshFrequency = 10 * time.Minute
	config.Metadata.Full = true
	return config, nil
}

func (k *KafkaBrokerV1) getSubConfig() (*sarama.Config, error) {
	// if c, ok := k.opts.Context.Value(clusterConfigKey{}).(*sarama.Config); ok {
	//	return c
	// }
	config := sarama.NewConfig()
	if k.TLS != nil {
		tlsConfig, err := k.TLS.getTLS()
		if err != nil {
			return nil, err
		}
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = tlsConfig
	}
	// the oldest supported version is V0_10_2_0
	if !config.Version.IsAtLeast(sarama.V0_10_2_0) {
		config.Version = sarama.V0_10_2_0
	}

	config.ClientID = uuid.New().String()
	// 设置发送/订阅时获取连接重试次数与重试间隔时间
	config.Consumer.Return.Errors = true
	config.Metadata.Retry.Max = 20
	config.Metadata.Retry.Backoff = 30 * time.Second
	config.Metadata.RefreshFrequency = 10 * time.Minute
	config.Metadata.Full = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	// last
	k.SubSarama.toSarama(config)
	return config, nil
}

func (k *KafkaBrokerV1) brokersCreateTopic(uris []string, name string, part, replica int) (err error) {
	for _, uri := range uris {
		err := k.brokerCreateTopic(uri, name, part, replica)
		if err != nil {
			return err
		}
	}
	return nil
}
func (k *KafkaBrokerV1) brokerCreateTopic(uri string, name string, part, replica int) (err error) {
	b := sarama.NewBroker(uri)
	defer b.Close()
	config := sarama.NewConfig()
	config.Version = sarama.V1_0_0_0
	err = b.Open(config)
	if err != nil {
		return err
	}
	ok, err := b.Connected()
	if err != nil {
		return err
	}
	if !ok {
		return broker.ErrConnect
	}
	topicDetail := &sarama.TopicDetail{}
	topicDetail.NumPartitions = int32(part)
	topicDetail.ReplicationFactor = int16(replica)
	topicDetail.ConfigEntries = make(map[string]*string)

	topicDetails := make(map[string]*sarama.TopicDetail)
	topicDetails[name] = topicDetail

	request := sarama.CreateTopicsRequest{
		Timeout:      time.Second * 15,
		TopicDetails: topicDetails,
	}
	response, err := b.CreateTopics(&request)
	if err != nil {
		return err
	}
	t := response.TopicErrors
	tErr := t[name]
	if tErr != nil {
		if tErr.Err == sarama.ErrNoError {
			return nil
		}
		if tErr.Err == sarama.ErrTopicAlreadyExists {
			k.Logger.Infof("KafkaBrokerV1#brokerCreateTopic : Topic with this name already exists , name is : %v !", name)
			return nil
		}
		return tErr.Err
	}
	return nil
}

// consumerGroupHandler
type consumerGroupHandler struct {
	handler broker.CallBack
	subopts broker.SubscribeOptions
	cg      sarama.ConsumerGroup
	sess    sarama.ConsumerGroupSession
}

func (*consumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (*consumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h *consumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// notice umMarshal body only
		m := broker.Message{
			Header: nil,
			Body:   msg.Value,
		}
		if err := h.handler(&kafkaEvent{
			autoAck: h.subopts.AutoAck,
			t:       msg.Topic,
			cg:      h.cg,
			km:      msg,
			m:       &m,
			sess:    sess,
		}); err == nil && h.subopts.AutoAck {
			sess.MarkMessage(msg, "")
		}
	}
	return nil
}

// kafkaEvent
type kafkaEvent struct {
	autoAck bool
	t       string
	cg      sarama.ConsumerGroup
	km      *sarama.ConsumerMessage
	m       *broker.Message
	sess    sarama.ConsumerGroupSession
}

func (p *kafkaEvent) Context() c.Context {
	return context.TODO() // todo
}

func (p *kafkaEvent) Topic() string {
	return p.t
}

func (p *kafkaEvent) Message() *broker.Message {
	return p.m
}

func (p *kafkaEvent) Ack() error {
	if p.autoAck {
		return nil
	}
	p.sess.MarkMessage(p.km, "")
	return nil
}

// subscriber
type subscriber struct {
	cg     sarama.ConsumerGroup
	topics []string
}

func (s *subscriber) String() string {
	return "kafka subscriber"
}
func (s *subscriber) Topic() string {
	return strings.Join(s.topics, ",")
}

func (s *subscriber) Unsubscribe() error {
	return s.cg.Close()
}
