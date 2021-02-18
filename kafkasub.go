package main

import (
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/kafkav1"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"strings"
)

var KafkaSubscribeCommand = cli.Command{
	Name:      "kafkasub",
	Usage:     "subscribe kafka message",
	UsageText: "Usage: brokerc kafkasub [options...] <uri>, uri arg format: host.domain[:port],host.domain[:port],...",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "t",
			Usage:    "topic , this topic will be created when both parameters p and r are not 0.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "q",
			Usage:    "queue name.",
			Required: false,
		},
		cli.IntFlag{
			Name:     "o",
			Usage:    "offset , new : -1 , old : -2.",
			Required: false,
			Value:    -1,
		},
		cli.IntFlag{
			Name:     "p",
			Usage:    "part number.",
			Required: false,
		},
		cli.IntFlag{
			Name:     "r",
			Usage:    "replica number.",
			Required: false,
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
		cli.BoolFlag{
			Name:     "insecure",
			Usage:    "do not check that the server certificate hostname matches the remote hostname.",
			Required: false,
		},
	},
	Action: func(context *cli.Context) error {
		uri := context.Args().Get(0)
		t, o, d, cafile, cert, key, insecure, q, p, r :=
			context.String("t"),
			context.Int("o"),
			context.Bool("d"),
			context.String("cafile"),
			context.String("cert"),
			context.String("key"),
			context.Bool("insecure"),
			context.String("q"),
			context.Int("p"),
			context.Int("r")
		logger.SetDebug(d)
		var kt *kafkav1.TLS
		if insecure || cafile != "" || cert != "" || key != "" {
			kt = &kafkav1.TLS{
				Insecure:   insecure,
				CertFile:   cafile,
				KeyFile:    cert,
				CaCertFile: key,
			}
		}
		b := kafkav1.KafkaBrokerV1{
			Logger:    logger,
			URIs:      strings.Split(uri, ","),
			TLS:       kt,
			Insecure:  insecure,
			SubSarama: &kafkav1.Sarama{ConsumerOffsetsInitial: int64(o)},
		}
		err := b.Connect()
		if err != nil {
			return err
		}
		defer b.Disconnect()
		s, err := b.Subscribe([]string{t}, func(event broker.Event) error {
			logger.Infof("SUBSCRIBE=> uri:%v, t:%v, q:%v, m:%v !", uri, t, q, string(event.Message().Body))
			return nil
		}, broker.SetSubQueue(q),
			broker.SetSubPart(p),
			broker.SetSubReplica(r))
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
