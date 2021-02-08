# brokerc
    brokerc is a cross platform publish subscribe client, including mqtt protocol.
## mqtt pub
    ./brokerclinux mqttpub -t "/a/b" -h "192.168.2.60" -p "1883" -u "system" -P "manager" -i "mqttpub" -m 'cas' --will-payload 'pub bye' --will-topic 'will'
## mqtt sub
    ./brokerclinux mqttsub -t "/a/b" -h "192.168.2.60" -p "1883" -u "system" -P "manager" -i "mqttsub" --will-payload 'sub bye' --will-topic 'will'