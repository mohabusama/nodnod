# NodNod

NodNod is a websocket server that can stream node stats to any websocket client. NodNod can be deployed as single or multinode cluster.

In case of multinode cluster deployment, the NodNod cluster is *masterless*, so each NodNod server should stream stats for all connected/live NodNod peers in the cluster straight to the client.

## How does it work?

Assuming you already have a cluster composed of mutliple servers. Start by installing and running **NodNod** on each server with configuration file specifying all other **NodNod** peers ([see Config file](#config-file)).

Once a **NodNod** starts, it will start connecting to other peers via a websocket connection. The reason behind establishing those connections; is that NodNod aims at constructing a *masterless* stats collector cluster. So, **any** NodNod server can stream stats for the whole cluster instead of just a single node. 

This setup also eliminates the need of a master/controller node, allowing loadbalanced setup while serving large number of clients.

Please check the [Tutorial](#tutorial) section to see **NodNod** in action.

## Installation

    go get github.com/mohabusama/nodnod

## Stats

NodNod is using [Gorilla Websocket](https://github.com/gorilla/websocket) for websocket server implementation and [gopsutil](https://github.com/shirou/gopsutil) for stats gathering. Current gathered stats are:

* CPUUsage
* Total Disk
* Used Disk
* Used Disk percentage
* Total Memory
* Used Memory
* Used Memory percentage

## Usage

A NodNod server requires a config file to start.

### Config file

The config file is a `json` file which describes the existing **NodNod** nodes/peers in the cluster.

    {
        "nodes": [
            "192.168.20.15",
            "192.168.20.16",
            "192.168.20.17"
        ],
        "mode": "pull"
    }

### Startup

Start NodNod server by passing `config` file path and `listen` address in the form `"<ip>:<port>"`.

    Usage: nodnod [OPTIONS]
    
    OPTIONS:
    
      -config string
            Path to configuration path. (default "./conf/conf.json")
      -debug
            Set logging level to DEBUG!
      -help
            Print help!
      -listen string
            Websocket service listen address (default "127.0.0.1:7070")
      -version
            Show version!

## Client

The client package can be used to connect and stat a **NodNod** cluster.

    import "github.com/mohabusama/nodnod/client"

Please check [nodnod_client.go](https://github.com/mohabusama/nodnod/blob/master/_examples/nodnod_client.go) example to see the client usage.

### Examples

The `_examples` directory includes scripts that could be used to illustrate interaction with a NodNod server.

- **nodnod_client.go**: Uses the client package to stat a NodNod cluster.
- **nodnod_dial.go**: Connects to NodNod server, and continously requests server stats with the specified duration. Uses raw Gorilla websocket.
- **nodnod_concurrent.go**: Runs a concurrency test against a NodNod server/cluster. Uses raw Gorilla websocket.

## Tutorial

**First**, create a sample `config.json` file. Here, we will run a cluster of two nodes.

    {
        "nodes": [
            "127.0.0.1:7070",
            "127.0.0.1:7071"
        ],
        "mode": "pull"
    }

**Second**, start the first server

    $ nodnod --listen 127.0.0.1:7070 --config <path-to-config.json>
    
    INFO[0000] Starting NodNod websocket server:127.0.0.1:7070 
    INFO[0000] Starting discovery service                   
    WARN[0000] Failed to connect to node: 127.0.0.1:7071. Error: dial tcp 127.0.0.1:7071: getsockopt: connection refused 
    INFO[0000] Accepted connection with client: 127.0.0.1:63161 
    INFO[0005] Established connection with node: 127.0.0.1:7071 

In another terminal, start the second server

    $ nodnod --listen 127.0.0.1:7071 --config <path-to-config.json>
    
    INFO[0000] Starting NodNod websocket server: 127.0.0.1:7071 
    INFO[0000] Starting discovery service
    INFO[0000] Established connection with node: 127.0.0.1:7070 
    INFO[0004] Accepted connection with client: 127.0.0.1:63162 

The next step is to run one of the `_examples` scripts

    $ go run nodnod_dial.go
    
    INFO[0000] Received response from node:127.0.0.1:7070   
    INFO[0000] Duration:27.499909ms                        
    {
        "host": "127.0.0.1:7070",
        "nodes": {
            "127.0.0.1:7070": {
                "address": "127.0.0.1:7070",
                "cpuUsed": 5,
                "diskTotal": 249769419776,
                "diskUsed": 66792268800,
                "diskUsedPercent": 26.741571830491146,
                "memoryTotal": 8589934592,
                "memoryUsed": 5699850240,
                "memoryUsedPercent": 66.35499000549316,
                "error": ""
            },
            "127.0.0.1:7071": {
                "address": "127.0.0.1:7071",
                "cpuUsed": 9.090909090909092,
                "diskTotal": 249769419776,
                "diskUsed": 66792268800,
                "diskUsedPercent": 26.741571830491146,
                "memoryTotal": 8589934592,
                "memoryUsed": 5700075520,
                "memoryUsedPercent": 66.35761260986328,
                "error": ""
            }
        },
        "error": ""
    }

Or run `nodnod_concurrent.go`. Here we will make 200 stats requests, with concurrency of 20 requests, and we are validating returned stats against 2 nodes. 

    $ go run nodnod_concurrent.go --concurrent 20 --number 200 --nodes 2

    INFO[0000] Launched all goroutines!                     
    INFO[0000] Launched all requests!                       
    INFO[0000] In progress ...                              
    INFO[0015] ====SUMMARY====                              
    INFO[0015] Total Duration:3.632180036s                 
    INFO[0015] Total number of requests:200                 
    INFO[0015] Total number of responses:200                
    INFO[0015] Total number of missing responses:0          
    INFO[0015] Total number of failed responses:0           
    INFO[0015] Total number of invalid responses:0           
    INFO[0015] ========                                     

## TODO

- Tests
- PUSH mode
- More realistic benchmarking
- Serve static HTML5 demo

## License

The MIT License (MIT)
