brokerc is a cross platform publish subscribe client, including mqtt client, amqp client, http client.
# Install
To install the library, follow the classical:

    $ go get github.com/godaner/brokerc
    
Or get it from the released version: 

    https://github.com/godaner/brokerc/releases
    
> Note: wget -O brokerc https://github.com/godaner/brokerc/releases/download/1.0.1/brokerclinux

# Examples
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
