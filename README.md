brokerc is a cross platform publish subscribe client, including mqtt protocol.
# Install
To install the library, follow the classical:

    $ go get github.com/godaner/brokerc
    
Or get it from the released version: 

    https://github.com/godaner/brokerc/releases, 
> Note: wget -O brokerclinux https://github.com/godaner/brokerc/releases/download/1.0.0/brokerclinux
# Examples
## MQTT
#### Publish
    ./brokerclinux mqttpub \
    -t "/a/b" \
    -h "192.168.2.60" \
    -p "1883" \
    -u "system" \
    -P "manager" \
    -i "mqttpub" \
    -m 'cas' \
    --will-payload 'pub bye' \
    --will-topic 'will'
#### Publish with tls
    ./brokerclinux mqttpub \
    -t "/a/b" \
    -h "ssl://localhost" \
    -p "1883" \
    -u "system" \
    -P "manager" \
    -i "mqttpub" \
    -m 'cas' \
    --will-payload 'pub bye' \
    --will-topic 'will' \
    -insecure \
    -cafile '/opt/OmniVista_2500_NMS/data/cert/wma/ca.cer' \
    -cert /opt/OmniVista_2500_NMS/data/cert/wma/wma.pem \
    -key /opt/OmniVista_2500_NMS/data/cert/wma/wma.key
#### Subscribe
    ./brokerclinux mqttsub \
    -t "/a/b" \
    -h "192.168.2.62" \
    -p "1883" \
    -u "system" \
    -P "manager" \
    -i "mqttsub" \
    --will-payload 'sub bye' \
    --will-topic 'will'
#### Subscribe with tls
    ./brokerclinux mqttsub \
    -t "/a/b" \
    -h "ssl://localhost" \
    -p "1883" \
    -u "system" \
    -P "manager" \
    -i "mqttsub" \
    --will-payload 'sub bye' \
    --will-topic 'will' \
    -insecure \
    -cafile '/opt/OmniVista_2500_NMS/data/cert/wma/ca.cer' \
    -cert /opt/OmniVista_2500_NMS/data/cert/wma/wma.pem \
    -key /opt/OmniVista_2500_NMS/data/cert/wma/wma.key
## AMQP
#### Publish
#### Subscribe
    ./brokerclinux amqpsub \
    -t "/a/b" \
    -h "192.168.2.60" \
    -p "5672" \
    -u "system" \
    -P "manager" \
    -i "amqpsub" \
    --queue "amqpqueue" \
    --queue-ad
