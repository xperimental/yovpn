package provisioner

import (
	"fmt"

	"github.com/xperimental/yovpn/share"
)

func createConfigFile(dropletIP string, secret string) (string, error) {
	template := share.ClientTemplate

	config := fmt.Sprintf("%s\nremote %s\n<secret>\n%s\n</secret>\n", template, dropletIP, secret)
	return config, nil
}
