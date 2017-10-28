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
	Version = "0.0.1"
)

var (
	_ = kingpin.New("piped", "Piped Server")

	Port = kingpin.Flag("port", "Port number for of the HTTP service.").
		Default("7890").
		OverrideDefaultFromEnvar("PIPED_PORT").
		Short('p').
		Int()

	Silent = kingpin.Flag("silent", "Do not log requests.").
		Short('s').
		Bool()

	ConfigPath = kingpin.Arg("config", "Path to the configuration file.").
			ExistingFile()
)

func ReadConfig() json.RawMessage {
	if *ConfigPath == "" {
		if _, err := os.Stat("./piped.json"); os.IsExist(err) {
			*ConfigPath = "./piped.json"
		} else if _, err := os.Stat("/etc/piped/piped.json"); os.IsExist(err) {
			*ConfigPath = "/etc/piped/piped.json"
		} else if _, err := os.Stat(os.Getenv("HOME") + "/.piped.json"); os.IsExist(err) {
			*ConfigPath = os.Getenv("HOME") + "/.piped.json"
		} else {
			log.Fatal("no configuration file found!")
		}
	}

	fd, err := os.Open(*ConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	var conf json.RawMessage
	if cfg, err := ioutil.ReadFile(*ConfigPath); err != nil {
		log.Fatal(err)
	} else if err := json.Unmarshal(cfg, &conf); err != nil {
		log.Fatal(err)
	}
	return conf
}

func main() {
	kingpin.Version(Version)
	kingpin.Parse()
	log.SetFlags(0)

	config := ReadConfig()
	r := service.RouterFromConfig(config, *Silent)
	log.Printf("piped listenning for http connections on port %d", *Port)
	host := fmt.Sprintf(":%d", *Port)
	http.ListenAndServe(host, r)
}
