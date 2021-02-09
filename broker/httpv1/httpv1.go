package httpv1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/log"
	"github.com/godaner/brokerc/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HTTPBrokerV1 struct {
	c               *http.Client
	s               *http.Server
	PublishCallBack func(topic string, resp *http.Response) error `json:"-"`
	CACertFile      string
	CertFile        string
	KeyFile         string
	Insecure        bool
	Logger          log.Logger
}

func (h *HTTPBrokerV1) String() string {
	return h.Marshal()
}

func (h *HTTPBrokerV1) Marshal() string {
	bs, err := json.Marshal(h)
	if err != nil {
		return err.Error()
	}
	return string(bs)
}

func (h *HTTPBrokerV1) Publish(topic string, msg *broker.Message, opt ...broker.PublishOption) error {
	h.Logger.Debugf("HTTPBrokerV1#Publish : info is : %v , topic is : %v , msg is : %v !", h, topic, msg)
	ss := strings.SplitN(topic, "#", 2)
	if len(ss) != 2 {
		return broker.ErrPublish
	}
	url, err := url.Parse(ss[1])
	if err != nil {
		return err
	}

	// client
	ct, err := tls.GetClientTLSConfig(h.Insecure, h.CACertFile, h.CertFile, h.KeyFile)
	if err != nil {
		return err
	}
	h.c = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: ct,
		},
		Timeout: 1 * time.Hour,
	}
	// do
	body := ioutil.NopCloser(bytes.NewReader(msg.Body))
	resp, err := h.c.Do(&http.Request{
		Method: ss[0],
		URL:    url,
		Header: msg.Header,
		Body:   body,
	})
	if err != nil {
		return err
	}
	if h.PublishCallBack != nil {
		err = h.PublishCallBack(topic, resp)
		if err != nil {
			return err
		}
	}
	return nil
}
func (h *HTTPBrokerV1) Subscribe(topics []string, callBack broker.CallBack, opt ...broker.SubscribeOption) (broker.Subscriber, error) {
	h.Logger.Debugf("HTTPBrokerV1#Subscribe : info is : %v , topics is : %v !", h, topics)
	addr := ""
	if len(topics) >= 1 {
		addr = topics[0]
	}
	st, err := tls.GetServerTLSConfig(h.CACertFile)
	if err != nil {
		return nil, err
	}
	h.s = &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			bs, _ := ioutil.ReadAll(r.Body)
			_ = callBack(&httpEvent{
				uri: r.URL.RequestURI(),
				m: &broker.Message{
					Header: r.Header,
					Body:   bs,
				},
			})
		}),
		TLSConfig: st,
	}
	if h.CertFile != "" && h.KeyFile != "" {
		err := h.s.ListenAndServeTLS(h.CertFile, h.KeyFile)
		if err != nil {
			return nil, err
		}
		return &noSubscriber{}, nil
	}
	err = h.s.ListenAndServe()
	if err != nil {
		return nil, err
	}
	return &noSubscriber{}, nil
}

type httpEvent struct {
	uri string
	m   *broker.Message
}

func (h *httpEvent) Topic() string {
	return h.uri
}

func (h *httpEvent) Ack() error {
	return nil
}

func (h *httpEvent) Message() *broker.Message {
	return h.m
}

func (h *httpEvent) Context() context.Context {
	return context.Background()
}

type noSubscriber struct {
}

func (n *noSubscriber) Unsubscribe() error {
	return nil
}

func (n *noSubscriber) String() string {
	return ""
}

func (h *HTTPBrokerV1) Connect() error {
	return nil
}

func (h *HTTPBrokerV1) Disconnect() error {
	return nil
}
