package main

import (
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/mqttv1"
	"github.com/urfave/cli"
	"os"
	"os/signal"
)

var MQTTSubscribeCommand = cli.Command{
	Name:      "mqttsub",
	Usage:     "subscribe mqtt message",
	UsageText: "Usage: brokerc mqttsub [options...] <uri>, uri arg format: mqtt[s]://[username][:password]@host.domain[:port]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "t",
			Usage:    "topic.",
			Required: true,
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
		uri := context.Args().Get(0)
		i, t, d, q, c, wt, wp, wr, wq, cafile, cert, key, insecure :=
			context.String("i"),
			context.String("t"),
			context.Bool("d"),
			context.Int("q"),
			context.Bool("c"),
			context.String("will-topic"),
			context.String("will-payload"),
			context.Bool("will-retain"),
			context.Int("will-qos"),
			context.String("cafile"),
			context.String("cert"),
			context.String("key"),
			context.Bool("insecure")
		logger.SetDebug(d)
		b := mqttv1.MQTTBrokerV1{
			URI:        uri,
			CID:        i,
			WT:         wt,
			WP:         wp,
			WR:         wr,
			WQ:         byte(wq),
			C:          c,
			CACertFile: cafile,
			CertFile:   cert,
			KeyFile:    key,
			Insecure:   insecure,
			Logger:     logger,
			Debug:      d,
		}
		err := b.Connect()
		if err != nil {
			return err
		}
		defer b.Disconnect()
		s, err := b.Subscribe([]string{t}, func(event broker.Event) error {
			logger.Infof("SUBSCRIBE=> uri:%v, i:%v, t:%v, q:%v, c:%v, m:%v !", uri, i, t, q, c, string(event.Message().Body))
			return nil
		}, broker.SetSubQOS(q))
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
