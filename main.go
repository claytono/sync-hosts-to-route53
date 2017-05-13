package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

const HOSTS_FILE = "hosts"

func main() {
	networks := pflag.StringArray("network", []string{}, "IP ranges to manage")
	pflag.Parse()
	fmt.Println(networks)
	fmt.Println(readHosts(HOSTS_FILE))
}
