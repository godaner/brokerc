package kafkav1

import (
	"fmt"
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/log/logv1"
	"os"
	"sync"
	"testing"
	"time"
)

func TestKafkaEvent_Ack(t *testing.T) {
	var once sync.Once
	go func() {
		f(once)
	}()
	go func() {
		f(once)
	}()
	select {}
}
func f(once sync.Once) {
	once.Do(func() {
		<-time.After(100 * time.Second)
	})
	fmt.Println("finish")
}
func TestKafka_Publish(t *testing.T) {

	b := &KafkaBrokerV1{
		Logger: &logv1.LoggerV1{
			DebugWriter: os.Stdout,
			InfoWriter:  os.Stdout,
			WarnWriter:  os.Stdout,
			ErrorWriter: os.Stdout,
		},
		URIs: []string{"192.168.2.10:9092"},
		// TLS: &TLS{
		// 	// openssl x509 -in cert.crt -out cert.der -outform DER
		// 	// openssl x509 -in cert.der -inform DER -out cert.pem -outform PEM
		// 	ClientCertFile: "/home/godaner/Desktop/kafka/cert.pem",
		// 	ClientKeyFile:  "/home/godaner/Desktop/kafka/key.key",
		// 	CaCertFile:     "/home/godaner/Desktop/kafka/ca-cert",
		// },
	}
	err := b.Connect()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			time.Sleep(time.Second * 1)
			err := b.Publish("gotopic5", &broker.Message{
				Header: nil,
				Body:   []byte("111111111111111"),
			}, broker.SetPubPart(3), broker.SetPubReplica(1))
			if err != nil {
				panic(err)
			}
		}
	}()
	go func() {
		_, err := b.Subscribe([]string{"gotopic5"}, func(event broker.Event) error {
			fmt.Println("topic2 : " + event.Topic() + " !\n")
			fmt.Println("msg2 : " + string(event.Message().Body) + " !\n")
			return nil
		}, broker.SetSubQueue("kafkatest"), broker.SetSubPart(3), broker.SetSubReplica(1))
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		for {
			time.Sleep(time.Second * 1)
			err := b.Publish("gotopic5", &broker.Message{
				Header: nil,
				Body:   []byte("222222222222222"),
			}, broker.SetPubPart(3), broker.SetPubReplica(1))
			if err != nil {
				panic(err)
			}
		}
	}()
	go func() {
		_, err := b.Subscribe([]string{"gotopic5"}, func(event broker.Event) error {
			fmt.Println("topic : " + event.Topic() + " !\n")
			fmt.Println("msg : " + string(event.Message().Body) + " !\n")
			return nil
		}, broker.SetSubQueue("kafkatest"), broker.SetSubQueue("kafkatest"), broker.SetSubPart(3), broker.SetSubReplica(1))
		if err != nil {
			panic(err)
		}
	}()

	f := make(chan int, 1)
	<-f
}

func TestKafka_Subscribe(t *testing.T) {
	b := &KafkaBrokerV1{
		Logger: &logv1.LoggerV1{
			DebugWriter: os.Stdout,
			InfoWriter:  os.Stdout,
			WarnWriter:  os.Stdout,
			ErrorWriter: os.Stdout,
		},
		URIs: []string{"192.168.2.10:9092"},
		// TLS: &TLS{
		// 	// openssl x509 -in cert.crt -out cert.der -outform DER
		// 	// openssl x509 -in cert.der -inform DER -out cert.pem -outform PEM
		// 	ClientCertFile: "/home/godaner/Desktop/kafka/cert.pem",
		// 	ClientKeyFile:  "/home/godaner/Desktop/kafka/key.key",
		// 	CaCertFile:     "/home/godaner/Desktop/kafka/ca-cert",
		// },
		SubSarama: &Sarama{
			ConsumerOffsetsInitial: ConsumerOffsetOldest,
		},
	}
	err := b.Connect()
	if err != nil {
		panic(err)
	}
	// go func() {
	err = b.Publish("gotopic51", &broker.Message{
		Header: nil,
		Body:   []byte("111111111111111"),
	})
	if err != nil {
		panic(err)
	}
	// }()
	go func() {
		_, err := b.Subscribe([]string{"gotopic51"}, func(event broker.Event) error {
			event.Topic()
			fmt.Println("topic : " + event.Topic() + " !\n")
			fmt.Println("msg : " + string(event.Message().Body) + " !\n")
			event.Ack()
			return nil
		}, broker.SetSubQueue("kafkatest1"), broker.SetSubAutoAck(false))
		if err != nil {
			panic(err)
		}
	}()

	f := make(chan int, 1)
	<-f
}
