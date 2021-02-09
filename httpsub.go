package main

import (
	"encoding/json"
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/httpv1"
	"github.com/urfave/cli"
	"os"
	"os/signal"
)

var HTTPSubscribeCommand = cli.Command{
	Name:      "httpsub",
	Usage:     "subscribe http message",
	UsageText: "Usage: brokerc httpsub [options...]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "h",
			Usage:    "host.",
			Required: true,
		},
		cli.BoolFlag{
			Name:     "d",
			Usage:    "debug.",
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
		cli.StringFlag{
			Name:     "will-payload",
			Usage:    "payload for the client Will, which is sent by the broker in case of unexpected disconnection. If not given and will-topic is set, a zero length message will be sent.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "will-topic",
			Usage:    "the topic on which to subscribe the client Will.",
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
		h, d, cafile, cert, key :=
			context.String("h"),
			context.Bool("d"),
			context.String("cafile"),
			context.String("cert"),
			context.String("key")
		logger.SetDebug(d)

		b := httpv1.HTTPBrokerV1{
			CACertFile: cafile,
			CertFile:   cert,
			KeyFile:    key,
			Logger:     logger,
		}
		err := b.Connect()
		if err != nil {
			return err
		}
		defer b.Disconnect()
		s, err := b.Subscribe([]string{h}, func(event broker.Event) error {
			hs, _ := json.Marshal(event.Message().Header)
			logger.Infof("SUBSCRIBE=> uri:%v, m:%v !", event.Topic(), string(event.Message().Body))
			logger.Debugf("SUBSCRIBE=> uri:%v, H:%v !", event.Topic(), string(hs))
			return nil
		})
		if err != nil {
			return err
		}
		defer s.Unsubscribe()
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig
		return nil
	},
}
