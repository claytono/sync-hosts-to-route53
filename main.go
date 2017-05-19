package main

import (
	"fmt"
	"log"
	"net"
	"os"

	flags "github.com/jessevdk/go-flags"
)

type hostEntry struct {
	hostname string
	ip       net.IP
	// Aliases are read from the /etc/hosts file, but not mapped to Route53
	aliases []string
}

var opts struct {
	File     string    `short:"f" long:"file" description:"Input file in /etc/hosts format" default:"/etc/hosts" value-name:"HOSTFILE"`
	Networks []CIDRNet `long:"network" description:"Filter by CIDR network" value-name:"x.x.x.x/len"`
	Domain   string    `short:"d" long:"domain" description:"Domain to update records in" required:"true"`
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

func compareHosts(hosts []hostEntry, r53hosts []hostEntry) {

}

func main() {
	parseOpts()
	hosts := readHosts(opts.File)
	hosts = filterHosts(hosts, opts.Networks)

	r53 := newRoute53()
	r53_hosts, err := r53.getRoute53Hosts(opts.Domain)
	if err != nil {
		log.Fatal(fmt.Errorf("error when retrieving zones: %v", err))
	}

	compareHosts(hosts, r53_hosts)
}
