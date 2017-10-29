package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/hyperboloide/pipe/piped/service"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	version = "0.1.1"
)

var (
	_ = kingpin.New("piped", "Piped Server")

	port = kingpin.Flag("port", "Port number for of the HTTP service.").
		Default("7890").
		OverrideDefaultFromEnvar("PIPED_PORT").
		Short('p').
		Int()

	silent = kingpin.Flag("silent", "Do not log requests.").
		Short('s').
		Bool()

	configPath = kingpin.Arg("config", "Path to the configuration file.").
			ExistingFile()
)

func readConfig() json.RawMessage {
	if *configPath == "" {
		if _, err := os.Stat("./piped.json"); !os.IsNotExist(err) {
			*configPath = "./piped.json"
		} else if _, err := os.Stat("/etc/piped/piped.json"); !os.IsNotExist(err) {
			*configPath = "/etc/piped/piped.json"
		} else if _, err := os.Stat(os.Getenv("HOME") + "/.piped.json"); !os.IsNotExist(err) {
			*configPath = os.Getenv("HOME") + "/.piped.json"
		} else {
			log.Fatal("no configuration file found!")
		}
	}
	cfg, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func main() {
	kingpin.Version(version)
	kingpin.Parse()
	log.SetFlags(0)

	config := readConfig()
	r := service.RouterFromConfig(config, *silent)
	log.Printf("piped listenning for http connections on port %d", *port)
	host := fmt.Sprintf(":%d", *port)
	if err := http.ListenAndServe(host, r); err != nil {
		log.Fatal(err)
	}
}
