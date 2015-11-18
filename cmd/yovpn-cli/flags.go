package main

import (
	"flag"
	"log"
)

var (
	region     string
	token      string
	configFile string
	destroy    bool
)

func init() {
	flag.StringVar(&region, "region", "", "Region to spawn endpoint in.")
	flag.StringVar(&token, "token", "", "Token for using DigitalOcean API.")
	flag.StringVar(&configFile, "output", "yovpn.ovpn", "File to write configuration to.")
	flag.BoolVar(&destroy, "destroy", false, "If true, will destroy endpoints.")
}

func checkFlags() {
	flag.Parse()
	if len(token) == 0 {
		log.Fatal("You need to provide a token!")
	}
}
