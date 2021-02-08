# brokerc
    brokerc is a cross platform publish subscribe client, including mqtt protocol.
## mqtt pub
    ./brokerclinux mqttpub -t "/a/b" -h "192.168.2.60" -p "1883" -u "system" -P "manager" -i "cc" -m 'cas' --will-payload 'bye' --will-topic 'will'