package main

import (
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/amqpv1"
	"github.com/urfave/cli"
)

var AMQPPublishCommand = cli.Command{
	Name:      "amqppub",
	Usage:     "publish amqp message",
	UsageText: "Usage: brokerc amqppub [options...] <uri>, uri arg format: amqp[s]://[username][:password]@host.domain[:port][vhost]",
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
			Name:     "i",
			Usage:    "client id.",
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
		cli.StringFlag{
			Name:     "exchange",
			Usage:    "exchange name.",
			Required: false,
		},
		cli.StringFlag{
			Name:     "exchange-type",
			Usage:    "exchange type.",
			Required: false,
		},
		cli.BoolFlag{
			Name:     "exchange-ad",
			Usage:    "exchange ad.",
			Required: false,
		},
		cli.BoolFlag{
			Name:     "exchange-duration",
			Usage:    "exchange duration.",
			Required: false,
		},
	},
	Action: func(context *cli.Context) error {
		uri := context.Args().Get(0)
		i, m, t, d, cafile, cert, key, insecure, exchange, exchangeType, exchangeAD, exchangeDuration :=
			context.String("i"),
			context.String("m"),
			context.String("t"),
			context.Bool("d"),
			context.String("cafile"),
			context.String("cert"),
			context.String("key"),
			context.Bool("insecure"),
			context.String("exchange"),
			context.String("exchange-type"),
			context.Bool("exchange-ad"),
			context.Bool("exchange-duration")
		logger.SetDebug(d)
		b := amqpv1.AMQPBrokerV1{
			URI:        uri,
			CID:        i,
			CACertFile: cafile,
			CertFile:   cert,
			KeyFile:    key,
			Insecure:   insecure,
			Logger:     logger,
		}
		err := b.Connect()
		if err != nil {
			return err
		}
		defer b.Disconnect()
		err = b.Publish(t, &broker.Message{
			Header: nil,
			Body:   []byte(m),
		}, broker.SetPubExchangeName(exchange),
			broker.SetPubExchangeType(exchangeType),
			broker.SetPubExchangeDuration(exchangeDuration),
			broker.SetPubExchangeAD(exchangeAD))
		if err != nil {
			return err
		}
		logger.Infof("PUBLISH=> uri:%v, i:%v, t:%v, m:%v !", uri, i, t, m)
		return nil
	},
}
