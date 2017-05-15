package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	File     string    `short:"f" long:"file" description:"Input file in /etc/hosts format" default:"/etc/hosts" value-name:"HOSTFILE"`
	Networks []CIDRNet `long:"network" description:"Filter by CIDR network" value-name:"x.x.x.x/len"`
}

func parseOpts() {
	parser := flags.NewParser(&opts, flags.Default)

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

func main() {
	parseOpts()
	hosts := readHosts(opts.File)
	hosts = filterHosts(hosts, opts.Networks)
	fmt.Println(hosts)
}
