package provisioner

import (
	"fmt"
	"io/ioutil"
	"os"
)

func readTemplate() (string, error) {
	file, err := os.Open("share/client.conf")
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func createConfigFile(dropletIP string, secret string) (string, error) {
	template, err := readTemplate()
	if err != nil {
		return "", err
	}

	config := fmt.Sprintf("%s\nremote %s\n<secret>\n%s\n</secret>\n", template, dropletIP, secret)
	return config, nil
}
