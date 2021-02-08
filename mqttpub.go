package main

import (
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/mqttv1"
	"github.com/urfave/cli"
)

var MQTTPublishCommand = cli.Command{
	Name:  "mqttpub",
	Usage: "mqtt publish message",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "t",
			Usage:    "topic",
			Required: true,
		},
		cli.StringFlag{
			Name:     "m",
			Usage:    "message",
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
			Usage:    "port",
			Value:    "1883",
			Required: true,
		},
		cli.StringFlag{
			Name:     "u",
			Usage:    "username",
			Required: false,
		},
		cli.StringFlag{
			Name:     "P",
			Usage:    "password",
			Required: false,
		},
		cli.StringFlag{
			Name:     "i",
			Usage:    "client id",
			Required: false,
		},
		cli.StringFlag{
			Name:     "d",
			Usage:    "debug",
			Required: false,
		},
	},
	Action: func(context *cli.Context) error {
		h, p, u, P, i, t, m, d := context.String("h"), context.String("p"), context.String("u"), context.String("P"), context.String("i"), context.String("t"), context.String("m"), context.Bool("d")
		logger.SetDebug(d)
		b := mqttv1.MQTTBrokerV1{
			IP:       h,
			Port:     p,
			Username: u,
			Password: P,
			CID:      i,
			Logger:   logger,
		}
		err := b.Connect()
		if err != nil {
			return err
		}
		err = b.Publish(t, &broker.Message{
			Header: nil,
			Body:   []byte(m),
		})
		if err != nil {
			return err
		}
		logger.Infof("publish=> h:%v, p:%v, t:%v, i:%v, u:%v, P:%v, m:%v !", h, p, t, i, u, P, m)
		return nil
	},
}
