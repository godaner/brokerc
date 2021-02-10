# Brokerc
brokerc is a cross-platform publish and subscribe command line client tool, including mqtt client, amqp client, http client.
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
$ ./brokerclinux --help
NAME:
   brokerc - brokerc is a cross platform publish subscribe client, including mqtt client, amqp client, http client.

USAGE:
   brokerc [global options] command [command options] [arguments...]

VERSION:
   1.0.1

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
```
$ ./brokerclinux mqttpub --help
NAME:
   brokerc mqttpub - publish mqtt message

USAGE:
   Usage: brokerc mqttpub [options...] <uri>, uri arg format: mqtt[s]://[username][:password]@host.domain[:port]

OPTIONS:
   -t value              topic.
   -m value              message.
   -i value              client id.
   -d                    debug.
   -q value              quality of service level to use for all messages. Defaults to 0. (default: 0)
   -r                    message should be retained.
   --cafile value        path to a file containing trusted CA certificates to enable encrypted communication.
   --cert value          client certificate for authentication, if required by server.
   --key value           client private key for authentication, if required by server.
   --insecure            do not check that the server certificate hostname matches the remote hostname. Using this option means that you cannot be sure that the remote host is the server you wish to connect to and so is insecure. Do not use this option in a production environment.
   --will-payload value  payload for the client Will, which is sent by the broker in case of unexpected disconnection. If not given and will-topic is set, a zero length message will be sent.
   --will-topic value    the topic on which to publish the client Will.
   --will-retain         if given, make the client Will retained.
   --will-qos value      QoS level for the client Will.
```
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
```
$ ./brokerclinux mqttsub --help
NAME:
   brokerc mqttsub - subscribe mqtt message

USAGE:
   Usage: brokerc mqttsub [options...] <uri>, uri arg format: mqtt[s]://[username][:password]@host.domain[:port]

OPTIONS:
   -t value              topic.
   -i value              client id.
   -d                    debug.
   -q value              quality of service level to use for all messages. Defaults to 0. (default: 0)
   -c                    disable 'clean session' (store subscription and pending messages when client disconnects).
   --cafile value        path to a file containing trusted CA certificates to enable encrypted communication.
   --cert value          client certificate for authentication, if required by server.
   --key value           client private key for authentication, if required by server.
   --insecure            do not check that the server certificate hostname matches the remote hostname. Using this option means that you cannot be sure that the remote host is the server you wish to connect to and so is insecure. Do not use this option in a production environment.
   --will-payload value  payload for the client Will, which is sent by the broker in case of unexpected disconnection. If not given and will-topic is set, a zero length message will be sent.
   --will-topic value    the topic on which to publish the client Will.
   --will-retain         if given, make the client Will retained.
   --will-qos value      QoS level for the client Will.
```
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
```
$ ./brokerclinux amqppub --help
NAME:
   brokerc amqppub - publish amqp message

USAGE:
   Usage: brokerc amqppub [options...] <uri>, uri arg format: amqp[s]://[username][:password]@host.domain[:port][vhost]

OPTIONS:
   -t value               topic.
   -m value               message.
   -i value               client id.
   -d                     debug.
   --cafile value         path to a file containing trusted CA certificates to enable encrypted communication.
   --cert value           client certificate for authentication, if required by server.
   --key value            client private key for authentication, if required by server.
   --insecure             do not check that the server certificate hostname matches the remote hostname. Using this option means that you cannot be sure that the remote host is the server you wish to connect to and so is insecure. Do not use this option in a production environment.
   --exchange value       exchange name.
   --exchange-type value  exchange type.
   --exchange-ad          exchange ad.
   --exchange-duration    exchange duration.
```
    ./brokerc amqppub \
    amqp://system:manager@192.168.2.60:5672 \
    -t "/a/b" \
    -i "amqpsubclient" \
    --exchange "amqpexchange" \
    -m 'hey man!'
#### Subscribe
```
$ ./brokerclinux amqpsub --help
NAME:
   brokerc amqpsub - subscribe amqp message

USAGE:
   Usage: brokerc amqpsub [options...] <uri>, uri arg format: amqp[s]://[username][:password]@host.domain[:port][vhost]

OPTIONS:
   -t value               topic.
   -i value               client id.
   -d                     debug.
   --cafile value         path to a file containing trusted CA certificates to enable encrypted communication.
   --cert value           client certificate for authentication, if required by server.
   --key value            client private key for authentication, if required by server.
   --insecure             do not check that the server certificate hostname matches the remote hostname. Using this option means that you cannot be sure that the remote host is the server you wish to connect to and so is insecure. Do not use this option in a production environment.
   --exchange value       exchange name.
   --exchange-type value  exchange type.
   --exchange-ad          exchange ad.
   --exchange-duration    exchange duration.
   --queue value          queue name.
   --queue-ad             queue auto delete.
   --queue-duration       queue duration.
```
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
```
$ ./brokerclinux httppub --help
NAME:
   brokerc httppub - publish http message

USAGE:
   Usage: brokerc httppub [options...] <uri>, uri arg format: http[s]://[username][:password]@host.domain[:port][suburi]

OPTIONS:
   -X value              method. (default: "GET")
   -H value              header.
   -m value              message.
   -o value              write to file instead of stdout.
   -d                    debug.
   --cafile value        path to a file containing trusted CA certificates to enable encrypted communication.
   --cert value          client certificate for authentication, if required by server.
   --key value           client private key for authentication, if required by server.
   --insecure            do not check that the server certificate hostname matches the remote hostname. Using this option means that you cannot be sure that the remote host is the server you wish to connect to and so is insecure. Do not use this option in a production environment.
   --will-payload value  payload for the client Will, which is sent by the broker in case of unexpected disconnection. If not given and will-topic is set, a zero length message will be sent.
   --will-topic value    the topic on which to publish the client Will.
   --will-retain         if given, make the client Will retained.
   --will-qos value      QoS level for the client Will.
```
    ./brokerc httppub \
    http://127.0.0.1:2222/apiv1/do \
    -H "K1:A=C;K2:B=D;K1:E=F;" \
    -m 'hey man!'
#### Subscribe
```
$ ./brokerclinux httpsub --help
NAME:
   brokerc httpsub - subscribe http message

USAGE:
   Usage: brokerc httpsub [options...]

OPTIONS:
   -h value              host.
   -d                    debug.
   --cafile value        path to a file containing trusted CA certificates to enable encrypted communication.
   --cert value          client certificate for authentication, if required by server.
   --key value           client private key for authentication, if required by server.
   --will-payload value  payload for the client Will, which is sent by the broker in case of unexpected disconnection. If not given and will-topic is set, a zero length message will be sent.
   --will-topic value    the topic on which to subscribe the client Will.
   --will-retain         if given, make the client Will retained.
   --will-qos value      QoS level for the client Will.
```
    ./brokerc httpsub \
    -h :2222
