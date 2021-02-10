package main

import (
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/amqpv1"
	"github.com/urfave/cli"
	"os"
	"os/signal"
)

var AMQPSubscribeCommand = cli.Command{
	Name:      "amqpsub",
	Usage:     "subscribe amqp message",
	UsageText: "Usage: brokerc amqpsub [options...] <uri>, uri arg format: amqp[s]://[username][:password]@host.domain[:port][vhost]",
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
		cli.StringFlag{
			Name:     "queue",
			Usage:    "queue name.",
			Required: false,
		},
		cli.BoolFlag{
			Name:     "queue-ad",
			Usage:    "queue auto delete.",
			Required: false,
		},
		cli.BoolFlag{
			Name:     "queue-duration",
			Usage:    "queue duration.",
			Required: false,
		},
	},
	Action: func(context *cli.Context) error {
		uri := context.Args().Get(0)
		i, t, d, cafile, cert, key, insecure, exchange, exchangeType, exchangeAD, exchangeDuration, queue, queueAD, queueDuration :=
			context.String("i"),
			context.String("t"),
			context.Bool("d"),
			context.String("cafile"),
			context.String("cert"),
			context.String("key"),
			context.Bool("insecure"),
			context.String("exchange"),
			context.String("exchange-type"),
			context.Bool("exchange-ad"),
			context.Bool("exchange-duration"),
			context.String("queue"),
			context.Bool("queue-ad"),
			context.Bool("queue-duration")
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
		s, err := b.Subscribe([]string{t}, func(event broker.Event) error {
			logger.Infof("SUBSCRIBE=> uri:%v, i:%v, t:%v, exchange:%v, exchangeType:%v, queue:%v, m:%v !", uri, i, t, exchange, exchangeType, queue, string(event.Message().Body))
			return nil
		}, broker.SetSubQueue(queue),
			broker.SetSubAutoAck(true),
			broker.SetSubAutoDel(queueAD),
			broker.SetSubDuration(queueDuration),
			broker.SetSubExchangeName(exchange),
			broker.SetSubExchangeType(exchangeType),
			broker.SetSubExchangeDuration(exchangeDuration),
			broker.SetSubExchangeAD(exchangeAD))
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
