package main

import (
	"log"
	"net"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
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

// compareHosts takes the contents of the local /etc/hosts file and the
// route53 hosts that should be compared and produces two arrays.  The first
// array is a list of Route53 records that need to be updated, and the second
// array is a list of Route53 records that should be deleted.
func compareHosts(hosts []hostEntry, r53hosts []hostEntry) ([]hostEntry, []hostEntry) {
	// Build index on name, we'll delete entries out of here a we match them
	// against /etc/hosts entries.  The remaining entries aren't present
	// locally anymore and will need to be deleted.
	rhByName := map[string]hostEntry{}
	for _, rh := range r53hosts {
		rhByName[rh.hostname] = rh
	}

	toUpdate := []hostEntry{}
	// Find existing hosts
	for _, h := range hosts {
		rh, ok := rhByName[h.hostname]
		if ok {
			delete(rhByName, h.hostname)
			if !h.ip.Equal(rh.ip) {
				toUpdate = append(toUpdate, h)
			}
		} else {
			toUpdate = append(toUpdate, h)
		}
	}

	toDelete := make([]hostEntry, 0, len(rhByName))
	for _, rh := range rhByName {
		toDelete = append(toDelete, rh)
	}

	return toUpdate, toDelete
}

func main() {
	parseOpts()
	hosts := readHosts(opts.File)
	hosts = filterHosts(hosts, opts.Networks)

	r53 := newRoute53()
	r53Hosts, err := r53.getHosts(opts.Domain)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "error when retrieving zones: %v"))
	}
	r53Hosts = filterHosts(r53Hosts, opts.Networks)
}
