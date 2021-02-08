brokerc is a cross platform publish subscribe client, including mqtt protocol.
# Install
To install the library, follow the classical:

    $ go get github.com/godaner/brokerc
    
Or get it from the released version: 

    https://github.com/godaner/brokerc/releases
    
# Examples
## MQTT
#### Publish
    ./brokerclinux mqttpub -t "/a/b" -h "192.168.2.60" -p "1883" -u "system" -P "manager" -i "mqttpub" -m 'cas' --will-payload 'pub bye' --will-topic 'will'
#### Publish with tls
    ./brokerclinux mqttpub -t "/a/b" -h "ssl://localhost" -p "1883" -u "system" -P "manager" -i "mqttpub" -m 'cas' --will-payload 'pub bye' --will-topic 'will' -cafile '/opt/OmniVista_2500_NMS/data/cert/wma/ca.cer' -cert /opt/OmniVista_2500_NMS/data/cert/wma/wma.pem -key /opt/OmniVista_2500_NMS/data/cert/wma/wma.key -insecure
#### Subscribe
    ./brokerclinux mqttsub -t "/a/b" -h "192.168.2.60" -p "1883" -u "system" -P "manager" -i "mqttsub" --will-payload 'sub bye' --will-topic 'will'
#### Subscribe with tls
    ./brokerclinux mqttsub -t "/a/b" -h "ssl://localhost" -p "1883" -u "system" -P "manager" -i "mqttsub" --will-payload 'sub bye' --will-topic 'will' -cafile '/opt/OmniVista_2500_NMS/data/cert/wma/ca.cer' -cert /opt/OmniVista_2500_NMS/data/cert/wma/wma.pem -key /opt/OmniVista_2500_NMS/data/cert/wma/wma.key -insecure