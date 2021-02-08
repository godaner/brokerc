package main

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/mqttv1"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/signal"
)

var MQTTSubscribeCommand = cli.Command{
	Name:  "mqttsub",
	Usage: "mqtt subscribe message",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "t",
			Usage:    "topic.",
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
			Name:     "c",
			Usage:    "disable 'clean session' (store subscription and pending messages when client disconnects).",
			Required: false,
		},
		cli.StringFlag{
			Name:     "cafile",
			Usage:    "path to a file containing trusted CA certificates to enable encrypted communication.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "cert",
			Usage:    "client certificate for authentication, if required by server.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "key",
			Usage:    "client private key for authentication, if required by server.",
			Required: false,
		},
		cli.BoolFlag{
			Name:     "insecure",
			Usage:    "do not check that the server certificate hostname matches the remote hostname. Using this option means that you cannot be sure that the remote host is the server you wish to connect to and so is insecure. Do not use this option in a production environment.",
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
		h, p, u, P, i, t, d, q, c, wt, wp, wr, wq, cafile, cert, key, insecure := context.String("h"), context.String("p"), context.String("u"), context.String("P"), context.String("i"), context.String("t"), context.Bool("d"), context.Int("q"), context.Bool("c"), context.String("will-topic"), context.String("will-payload"), context.Bool("will-retain"), context.Int("will-qos"), context.String("cafile"), context.String("cert"), context.String("key"), context.Bool("insecure")
		logger.SetDebug(d)
		if d {
			mqtt.CRITICAL = log.New(os.Stdout, "MQTT_CRITICAL ", 0)
			mqtt.ERROR = log.New(os.Stdout, "MQTT_ERROR ", 0)
			mqtt.WARN = log.New(os.Stdout, "MQTT_WARN ", 0)
			mqtt.DEBUG = log.New(os.Stdout, "MQTT_DEBUG ", 0)
		}
		b := mqttv1.MQTTBrokerV1{
			Host:           h,
			Port:           p,
			Username:       u,
			Password:       P,
			CID:            i,
			WT:             wt,
			WP:             wp,
			WR:             wr,
			WQ:             byte(wq),
			C:              c,
			CACertFile:     cafile,
			ClientCertFile: cert,
			ClientKeyFile:  key,
			Insecure:       insecure,
			Logger:         logger,
		}
		err := b.Connect()
		if err != nil {
			return err
		}
		defer b.Disconnect()
		s, err := b.Subscribe([]string{t}, func(event broker.Event) error {
			logger.Infof("SUBSCRIBE=> h:%v, p:%v, u:%v, P:%v, i:%v, t:%v, d:%v, q:%v, c:%v, m:%v !", h, p, u, P, i, t, d, q, c, string(event.Message().Body))
			return nil
		}, broker.SetSubQOS(q))
		defer s.Unsubscribe()
		if err != nil {
			return err
		}
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig
		return nil
	},
}
