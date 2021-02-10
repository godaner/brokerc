# Brokerc
brokerc is a cross-platform publish and subscribe command line client tool, including mqtt client, amqp client, http client, kafka client.
# Install
To install the library, follow the classical:

    $ go get github.com/godaner/brokerc
    
Or get it from the released version: 

    https://github.com/godaner/brokerc/releases
    
> Note: wget -O brokerc https://github.com/godaner/brokerc/releases/download/1.0.1/brokerclinux

# Supported platforms

This library works (and is tested) on the following platforms:

<table>
  <thead>
    <tr>
      <th>Platform</th>
      <th>Architecture</th>
      <th>Status</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td rowspan="2">Linux</td>
      <td><code>amd64</code></td>
      <td>✅</td>
    </tr>
    <tr>
      <td><code>386</code></td>
      <td>✅</td>
    </tr>
    <tr>
      <td rowspan="2">Windows</td>
      <td><code>amd64</code></td>
      <td>✅</td>
    </tr>
    <tr>
      <td><code>386</code></td>
      <td>✅</td>
    </tr>
    <tr>
      <td>Others</td>
      <td><code>Others</code></td>
      <td>⏳</td>
    </tr>
  </tbody>
</table>

# Usage
```
$ ./brokerc --help
NAME:
   brokerc - brokerc is a cross-platform publish and subscribe command line client tool, including mqtt client, amqp client, http client, kafka client.

USAGE:
   brokerc [global options] command [command options] [arguments...]

VERSION:
   1.0.2

COMMANDS:
   mqttpub  publish mqtt message
   mqttsub  subscribe mqtt message
   amqpsub  subscribe amqp message
   amqppub  publish amqp message
   httppub  publish http message
   httpsub  subscribe http message
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help     show help
   --version  print the version
```
## MQTT
#### Publish
    ./brokerc mqttpub \
    tcp://system:manager@192.168.2.60:1883 \
    -t "/a/b" \
    -i "mqttpub" \
    -m 'cas' \
    --will-payload 'pub bye' \
    --will-topic 'will'
#### Publish with tls
    ./brokerc mqttpub \
    ssl://system:manager@localhost:1883 \
    -t "/a/b" \
    -i "mqttpub" \
    -m 'cas' \
    --will-payload 'pub bye' \
    --will-topic 'will' \
    -insecure \
    -cafile '/opt/OmniVista_2500_NMS/data/cert/wma/ca.cer' \
    -cert /opt/OmniVista_2500_NMS/data/cert/wma/wma.pem \
    -key /opt/OmniVista_2500_NMS/data/cert/wma/wma.key
#### Subscribe
    ./brokerc mqttsub \
    tcp://system:manager@192.168.2.60:1883 \
    -t "/a/b" \
    -i "mqttsub" \
    --will-payload 'sub bye' \
    --will-topic 'will'
#### Subscribe with tls
    ./brokerc mqttsub \
    ssl://system:manager@localhost:1883 \
    -t "/a/b" \
    -i "mqttsub" \
    --will-payload 'sub bye' \
    --will-topic 'will' \
    -insecure \
    -cafile '/opt/OmniVista_2500_NMS/data/cert/wma/ca.cer' \
    -cert /opt/OmniVista_2500_NMS/data/cert/wma/wma.pem \
    -key /opt/OmniVista_2500_NMS/data/cert/wma/wma.key
## AMQP
#### Publish
    ./brokerc amqppub \
    amqp://system:manager@192.168.2.60:5672 \
    -t "/a/b" \
    -i "amqpsubclient" \
    --exchange "amqpexchange" \
    -m 'hey man!'
#### Subscribe
    ./brokerc amqpsub \
    amqp://system:manager@192.168.2.60:5672 \
    -t "/a/b" \
    -i "amqpsubclient" \
    --queue "amqpqueue" \
    --exchange "amqpexchange" \
    --exchange-type "direct" \
    --queue-ad \
    --exchange-ad
## HTTP
#### Publish
    ./brokerc httppub \
    http://127.0.0.1:2222/apiv1/do \
    -H "K1:A=C;K2:B=D;K1:E=F;" \
    -m 'hey man!'
#### Subscribe
    ./brokerc httpsub \
    -h :2222
