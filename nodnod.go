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
	VERSION     = "0.1"
	CONFIG_PATH = "./conf/conf.json"
)

var (
	flHelp    = flag.Bool("help", false, "Print help!")
	flAddress = flag.String("listen", "127.0.0.1:7070",
		"Websocket service listen address")
	flConfigPath = flag.String("config", CONFIG_PATH, "Path to configuration path.")
	flName       = flag.String("name", "",
		"Name of this NodNod server. Default is Host name.")

	flVersion = flag.Bool("version", false, "Show version!")

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

	if *flName == "" {
		*flName, _ = os.Hostname()
		log.Warn("Using hostname as NodNod name:", *flName)
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

	log.Infof("Starting NodNod websocket server: [%s : %s]", *flName, *flAddress)

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
