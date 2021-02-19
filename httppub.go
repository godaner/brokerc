package main

import (
	"encoding/json"
	"fmt"
	"github.com/godaner/brokerc/broker"
	"github.com/godaner/brokerc/broker/httpv1"
	"github.com/godaner/brokerc/spinner"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
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
				hs, _ := json.Marshal(resp.Header)
				logger.Debugf("PUBLISH RECV=> uri:%v, H:%v !", uri, string(hs))
				defer resp.Body.Close()
				if o != "" {
					lens := resp.Header.Get("Content-Length")
					len, _ := strconv.ParseUint(lens, 10, 64)
					// Spinner
					stopUpdateSpinner := make(chan struct{})
					defer close(stopUpdateSpinner)
					s := spinner.Spinner{}
					s.Start()
					defer s.Stop()
					download, v, fileName := uint64(0), uint64(0), o
					s.UpdateStatus(&spinner.Status{
						Download: &download,
						Total:    &len,
						V:        &v,
						FileName: &fileName,
					})
					go func() {
						for ; ; {
							select {
							case <-time.After(time.Second):
								s.UpdateStatus(&spinner.Status{
									Download: &download,
									Total:    &len,
									V:        &v,
									FileName: &fileName,
								})
								v = 0
							case <-stopUpdateSpinner:
								return
							}
						}
					}()
					// download
					file, err := os.Create(o)
					if err != nil {
						return err
					}
					defer file.Close()
					buf := make([]byte, KB)
					for {
						n, err := resp.Body.Read(buf)
						if err != nil {
							return nil
						}
						if n <= 0 || err == io.EOF {
							break
						}
						download += uint64(n)
						v += uint64(n)
						_, err = file.Write(buf[:n])
						if err != nil {
							return nil
						}
					}
				} else {
					bs, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						return err
					}
					logger.Infof("PUBLISH RECV=> uri:%v, X:%v, resp:%v !", uri, X, string(bs))
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
