package main

import (
	"github.com/urfave/cli"
)

var MQTTSubscribeCommand = cli.Command{
	Name:  "mqttsub",
	Usage: "mqtt subscribe message",
	Flags: []cli.Flag{

	},
	Action: func(context *cli.Context) error {
		return nil
	},
}
