package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/mqttv1"
	"github.com/urfave/cli"
	"log"
	"os"
)

var MQTTPublishCommand = cli.Command{
	Name:  "mqttpub",
	Usage: "mqtt publish message",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "t",
			Usage:    "topic.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "m",
			Usage:    "message.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "h",
			Usage:    "host",
			Value:    "localhost",
			Required: true,
		},
		cli.StringFlag{
			Name:     "p",
			Usage:    "port.",
			Value:    "1883",
			Required: true,
		},
		cli.StringFlag{
			Name:     "u",
			Usage:    "username.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "P",
			Usage:    "password.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "i",
			Usage:    "client id.",
			Required: false,
		},
		cli.BoolFlag{
			Name:     "d",
			Usage:    "debug.",
			Required: false,
		},
		cli.IntFlag{
			Name:     "q",
			Usage:    "quality of service level to use for all messages. Defaults to 0.",
			Required: false,
		},
		cli.BoolFlag{
			Name:     "r",
			Usage:    "message should be retained.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "will-payload",
			Usage:    "payload for the client Will, which is sent by the broker in case of unexpected disconnection. If not given and will-topic is set, a zero length message will be sent.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "will-topic",
			Usage:    "the topic on which to publish the client Will.",
			Required: false,
		},
		cli.BoolFlag{
			Name:     "will-retain",
			Usage:    "if given, make the client Will retained.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "will-qos",
			Usage:    "QoS level for the client Will.",
			Required: false,
		},
	},
	Action: func(context *cli.Context) error {
		h, p, u, P, i, t, d, q, r, m, wt, wp, wr, wq := context.String("h"), context.String("p"), context.String("u"), context.String("P"), context.String("i"), context.String("t"), context.Bool("d"), context.Int("q"), context.Bool("r"), context.String("m"), context.String("will-topic"), context.String("will-payload"), context.Bool("will-retain"), context.Int("will-qos")
		logger.SetDebug(d)
		if d {
			mqtt.CRITICAL = log.New(os.Stdout, "MQTT_CRITICAL ", 0)
			mqtt.ERROR = log.New(os.Stdout, "MQTT_ERROR ", 0)
			mqtt.WARN = log.New(os.Stdout, "MQTT_WARN ", 0)
			mqtt.DEBUG = log.New(os.Stdout, "MQTT_DEBUG ", 0)
		}
		b := mqttv1.MQTTBrokerV1{
			IP:       h,
			Port:     p,
			Username: u,
			Password: P,
			CID:      i,
			WT:       wt,
			WP:       wp,
			WR:       wr,
			WQ:       byte(wq),
			Logger:   logger,
		}
		err := b.Connect()
		if err != nil {
			return err
		}
		defer b.Disconnect()
		err = b.Publish(t, &broker.Message{
			Header: nil,
			Body:   []byte(m),
		}, broker.SetPubQOS(q), broker.SetPubRetained(r))
		if err != nil {
			return err
		}
		logger.Infof("PUBLISH=> h:%v, p:%v, u:%v, P:%v, i:%v, t:%v, d:%v, q:%v, r:%v, m:%v !", h, p, u, P, i, t, d, q, r, m)
		return nil
	},
}
