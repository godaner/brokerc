# brokerc
    brokerc is a cross platform publish subscribe client, including mqtt protocol.
## mqtt
#### publish
    ./brokerclinux mqttpub -t "/a/b" -h "192.168.2.60" -p "1883" -u "system" -P "manager" -i "mqttpub" -m 'cas' --will-payload 'pub bye' --will-topic 'will'
#### publish with tls
    ./brokerclinux mqttpub -t "/a/b" -h "ssl://localhost" -p "1883" -u "system" -P "manager" -i "mqttpub" -m 'cas' --will-payload 'pub bye' --will-topic 'will' -cafile '/opt/OmniVista_2500_NMS/data/cert/wma/ca.cer' -cert /opt/OmniVista_2500_NMS/data/cert/wma/wma.pem -key /opt/OmniVista_2500_NMS/data/cert/wma/wma.key -insecure
#### subscribe
    ./brokerclinux mqttsub -t "/a/b" -h "192.168.2.60" -p "1883" -u "system" -P "manager" -i "mqttsub" --will-payload 'sub bye' --will-topic 'will'
#### subscribe with tls
    ./brokerclinux mqttsub -t "/a/b" -h "ssl://localhost" -p "1883" -u "system" -P "manager" -i "mqttsub" --will-payload 'sub bye' --will-topic 'will' -cafile '/opt/OmniVista_2500_NMS/data/cert/wma/ca.cer' -cert /opt/OmniVista_2500_NMS/data/cert/wma/wma.pem -key /opt/OmniVista_2500_NMS/data/cert/wma/wma.key -insecure