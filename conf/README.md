# NodNod Config file

`conf.json` holds sample configuration file loaded by NodNod server.

## Sections

### nodes

Holds all the **NodNod** peers adresses that run on the cluster. It is ok to include the address of the NodNod server loading this config file.

Example:
    
    "nodes": [
        "127.0.0.1:7070",
        "127.0.0.1:7071",
        "127.0.0.1:7072"
    ]

### mode

Current implementation only supports `pull` mode. Pull mode means that the client should request the stats in order to be collected.

    "mode": "pull"
