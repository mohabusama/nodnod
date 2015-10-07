package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/mohabusama/nodnod/stats"
	"net/http"
	"os"
)

const (
	VERSION    = "0.1"
	configPath = "./conf/conf.json"
)

var (
	flHelp    = flag.Bool("help", false, "Print help!")
	flAddress = flag.String("listen", "127.0.0.1:7070",
		"Websocket service listen address")
	flConfigPath = flag.String("config", configPath, "Path to configuration path.")
	flVersion    = flag.Bool("version", false, "Show version!")

	flDebug = flag.Bool("debug", false, "Set logging level to DEBUG!")

	// flStatic     = flag.Bool([]string{"s", "-static"}, false, "Serve static html demo")

	globalConfig = new(Configuration)
	globalNodes  []*stats.Node
)

func main() {
	// 1. Handle cmd line args
	flag.Usage = showHelp

	flag.Parse()

	if *flHelp {
		flag.Usage()
		return
	}

	log.SetLevel(log.InfoLevel)
	if *flDebug {
		log.SetLevel(log.DebugLevel)
	}

	if *flVersion {
		showVersion()
		return
	}

	// 2. Load and validate config
	err := loadConfig()
	if err != nil {
		panic(err)
	}

	if err = loadNodes(); err != nil {
		panic(err)
	}

	// 3. Start discovery
	go discover()

	// 4. Launch server
	http.HandleFunc("/", serveWebsocket)

	log.Info("Starting NodNod websocket server:", *flAddress)

	err = http.ListenAndServe(*flAddress, nil)
	if err != nil {
		log.Fatal("Failed to start server!")
	}

}

func showVersion() {
	fmt.Printf("NodNod version\t%s\n", VERSION)
}

func showHelp() {
	fmt.Fprint(os.Stdout, "Usage: nodnod [OPTIONS]\n\n")
	fmt.Fprint(os.Stdout, "OPTIONS:\n\n")
	flag.CommandLine.SetOutput(os.Stdout)
	flag.PrintDefaults()
}
