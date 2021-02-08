package main

import (
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/amqpv1"
	"github.com/urfave/cli"
	"os"
	"os/signal"
)

var AMQPSubscribeCommand = cli.Command{
	Name:  "amqpsub",
	Usage: "amqp subscribe message",
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
			Value:    "5672",
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
			Name:     "exchange",
			Usage:    "exchange name.",
			Required: false,
			Value:    "",
		},
		cli.StringFlag{
			Name:     "exchange-type",
			Usage:    "exchange type.",
			Required: false,
			Value:    "direct",
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
		h, p, u, P, i, t, d, cafile, cert, key, insecure, exchange, exchangeType, queue, queueAD, queueDuration :=
			context.String("h"),
			context.String("p"),
			context.String("u"),
			context.String("P"),
			context.String("i"),
			context.String("t"),
			context.Bool("d"),
			context.String("cafile"),
			context.String("cert"),
			context.String("key"),
			context.Bool("insecure"),
			context.String("exchange"),
			context.String("exchange-type"),
			context.String("queue"),
			context.Bool("queue-ad"),
			context.Bool("queue-duration")
		logger.SetDebug(d)
		// if d {
		// 	amqp.CRITICAL = log.New(os.Stdout, "AMQP_CRITICAL ", 0)
		// 	amqp.ERROR = log.New(os.Stdout, "AMQP_ERROR ", 0)
		// 	amqp.WARN = log.New(os.Stdout, "AMQP_WARN ", 0)
		// 	amqp.DEBUG = log.New(os.Stdout, "AMQP_DEBUG ", 0)
		// }
		b := amqpv1.AMQPBrokerV1{
			Host:           h,
			Port:           p,
			Username:       u,
			Password:       P,
			CID:            i,
			CACertFile:     cafile,
			ClientCertFile: cert,
			ClientKeyFile:  key,
			Insecure:       insecure,
			Exchange:       exchange,
			ExchangeType:   exchangeType,
			Logger:         logger,
		}
		err := b.Connect()
		if err != nil {
			return err
		}
		defer b.Disconnect()
		s, err := b.Subscribe([]string{t}, func(event broker.Event) error {
			logger.Infof("SUBSCRIBE=> h:%v, p:%v, u:%v, P:%v, i:%v, t:%v, exchange:%v, exchangeType:%v, queue:%v, m:%v !", h, p, u, P, i, t, exchange, exchangeType, queue, string(event.Message().Body))
			return nil
		}, broker.SetSubQueue(queue), broker.SetSubAutoAck(true), broker.SetSubAutoDel(queueAD), broker.SetSubDuration(queueDuration))
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
