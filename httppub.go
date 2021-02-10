package main

import (
	"encoding/json"
	"fmt"
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/httpv1"
	"github.com/urfave/cli"
	"io/ioutil"
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
		cli.StringFlag{
			Name:     "o",
			Usage:    "write to file instead of stdout.",
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
			Usage:    "do not check that the server certificate hostname matches the remote hostname.",
			Required: false,
		},
	},
	Action: func(context *cli.Context) error {
		uri := context.Args().Get(0)
		X, H, d, m, o, cafile, cert, key, insecure :=
			context.String("X"),
			context.String("H"),
			context.Bool("d"),
			context.String("m"),
			context.String("o"),
			context.String("cafile"),
			context.String("cert"),
			context.String("key"),
			context.Bool("insecure")
		logger.SetDebug(d)

		b := httpv1.HTTPBrokerV1{
			PublishCallBack: func(topic string, resp *http.Response) error {
				bs, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					logger.Errorf("PUBLISH RECV ERR=> err:%v !", err)
					return nil
				}
				defer resp.Body.Close()
				if o != "" {
					err = ioutil.WriteFile(o, bs, 0755)
					if err != nil {
						logger.Errorf("PUBLISH RECV ERR=> err:%v !", err)
						return nil
					}
				} else {
					hs, _ := json.Marshal(resp.Header)
					logger.Infof("PUBLISH RECV=> uri:%v, X:%v, resp:%v !", uri, X, string(bs))
					logger.Debugf("PUBLISH RECV=> uri:%v, H:%v !", uri, string(hs))
				}
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
		logger.Infof("PUBLISH=> uri:%v, X:%v, m:%v !", uri, X, m)
		logger.Debugf("PUBLISH=> uri:%v, H:%v !", uri, H)
		err = b.Publish(fmt.Sprintf("%v#%v", X, uri), &broker.Message{
			Header: parseH(H),
			Body:   []byte(m),
		})
		if err != nil {
			return err
		}
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
