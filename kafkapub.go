package main

import (
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/kafkav1"
	"github.com/urfave/cli"
	"strings"
)

var KafkaPublishCommand = cli.Command{
	Name:      "kafkapub",
	Usage:     "publish kafka message",
	UsageText: "Usage: brokerc kafkapub [options...] <uri>, uri arg format: host.domain[:port],host.domain[:port],...",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "t",
			Usage:    "topic , when p and r is not zero , we will create this topic.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "m",
			Usage:    "message.",
			Required: true,
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
		m, t, d, cafile, cert, key, insecure, p, r :=
			context.String("m"),
			context.String("t"),
			context.Bool("d"),
			context.String("cafile"),
			context.String("cert"),
			context.String("key"),
			context.Bool("insecure"),
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
			SubSarama: nil,
		}
		err := b.Connect()
		if err != nil {
			return err
		}
		defer b.Disconnect()
		err = b.Publish(t, &broker.Message{
			Header: nil,
			Body:   []byte(m),
		}, broker.SetPubPart(p),
			broker.SetPubReplica(r))
		if err != nil {
			return err
		}
		logger.Infof("PUBLISH=> uri:%v, t:%v, m:%v !", uri, t, m)
		return nil
	},
}
