package main

import (
	"fmt"
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/httpv1"
	"github.com/urfave/cli"
	"net/http"
	"strings"
)

var HTTPPublishCommand = cli.Command{
	Name:      "httppub",
	Usage:     "publish http message",
	UsageText: "Usage: brokerc httppub [options...] <uri>, uri arg format: http[s]://[username][:password]@host.domain[:port][suburi]",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:     "X",
			Usage:    "method.",
			Required: true,
			Value:    "GET",
		},
		cli.StringFlag{
			Name:     "H",
			Usage:    "header.",
			Required: true,
		},
		cli.StringFlag{
			Name:     "m",
			Usage:    "message.",
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
		X, H, d, m, cafile, cert, key, insecure :=
			context.String("X"),
			context.String("H"),
			context.Bool("d"),
			context.String("m"),
			context.String("cafile"),
			context.String("cert"),
			context.String("key"),
			context.Bool("insecure")
		logger.SetDebug(d)

		b := httpv1.HTTPBrokerV1{
			PublishCallBack: func(topic string, resp *http.Response) error {
				// todo
				return nil
			},
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
		err = b.Publish(fmt.Sprintf("%v#%v", X, uri), &broker.Message{
			Header: parseH(H),
			Body:   []byte(m),
		})
		if err != nil {
			return err
		}
		logger.Infof("PUBLISH=> uri:%v, X:%v, m:%v !", uri, X, m)
		logger.Debugf("PUBLISH=> uri:%v, H:%v !", uri, H)
		return nil
	},
}

func parseH(H string) (hr map[string][]string) {
	hr = map[string][]string{}
	kvStrs := strings.Split(H, ";")
	for _, kvStr := range kvStrs {
		kvs := strings.Split(kvStr, ":")
		if len(kvs) != 2 {
			continue
		}
		k, v := kvs[0], kvs[1]
		vs := hr[k]
		hr[k] = append(vs, v)
	}
	return hr
}
