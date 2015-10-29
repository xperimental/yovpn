package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var output = flag.String("output", "yovpn.ovpn", "Output file for OpenVPN configuration")

func readTemplate() string {
	file, err := os.Open("share/client.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf)
}

func writeConfigFile(dropletIP string, secret string) string {
	template := readTemplate()
	config := fmt.Sprintf("%s\nremote %s\n<secret>\n%s\n</secret>\n", template, dropletIP, secret)
	err := ioutil.WriteFile(*output, []byte(config), 0600)
	if err != nil {
		log.Fatal(err)
	}
	return *output
}
