package amqpv1

import (
	"fmt"
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/log/logv1"
	"os"
	"testing"
)

func TestAMQPBrokerV1_Connect(t *testing.T) {
	a := &AMQPBrokerV1{
		URI: "amqp://system:manager@192.168.2.62:5672/",
		CID: "test1",
		Logger: &logv1.LoggerV1{
			DebugWriter: os.Stdout,
			InfoWriter:  os.Stdout,
			WarnWriter:  os.Stdout,
			ErrorWriter: os.Stdout,
		},
	}
	a.Connect()
	defer a.Disconnect()
	fmt.Println(a.Subscribe([]string{"test1"}, func(event broker.Event) error {
		fmt.Println(event)
		return nil
	}, broker.SetSubQueue("testqueue1")))
}
