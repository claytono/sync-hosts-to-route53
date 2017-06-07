package main

import (
	"log/syslog"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	logrus_syslog "github.com/Sirupsen/logrus/hooks/syslog"
	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

var log = logrus.New()

var opts struct {
	Mode           string        `short:"m" long:"mode" description:"Operating mode" default:"daemon" choice:"daemon" choice:"oneshot"`
	File           string        `short:"f" long:"file" description:"Input file in /etc/hosts format" default:"/etc/hosts" value-name:"HOSTFILE"`
	Networks       []CIDRNet     `long:"network" description:"Filter by CIDR network" value-name:"x.x.x.x/len"`
	Domain         string        `short:"d" long:"domain" description:"Domain to update records in" required:"true"`
	Interval       time.Duration `short:"i" long:"interval" description:"Seconds between scheduled resync times." default:"15m"`
	TTL            int64         `long:"ttl" description:"TTL to use for Route 53 records" default:"3600"`
	NoQualifyHosts bool          `long:"no-qualify-hosts" description:"Don't force domain to be added to end of hosts"`
	NoWait         bool          `long:"no-wait" description:"Don't wait for Route 53 to finish update"`
	Syslog         bool          `long:"syslog" description:"Send logging to syslog in addition to stdout"`
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

	// Accept trailing dot, but ignore it for consistency sake
	if strings.HasSuffix(opts.Domain, ".") {
		opts.Domain = opts.Domain[:len(opts.Domain)-1]
	}
}

func configureLogging() {
	// logrus defaults to stderr, but stdout is more conventional
	log.Out = os.Stdout

	if opts.Syslog {
		log.Info("Disabling color for syslog")
		tf := &logrus.TextFormatter{DisableColors: true}
		log.Formatter = tf

		hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_USER, "sync-hosts-to-route53")
		if err == nil {
			log.Hooks.Add(hook)
		} else {
			log.Fatal("Cannot initialize syslog: ", err)
		}
	}

}

func canonifyHostname(hostname string) string {
	hostname = strings.ToLower(hostname)
	// Accept trailing dot, but ignore it for consistency sake
	if strings.HasSuffix(hostname, ".") {
		hostname = hostname[:len(hostname)-1]
	}

	return hostname
}

// compareHosts takes the contents of the local /etc/hosts file and the
// Route 53 hosts that should be compared and produces two arrays.  The first
// array is a list of Route 53 records that need to be updated, and the second
// array is a list of Route 53 records that should be deleted.
func compareHosts(hosts hostList, r53hosts hostList) (hostList, hostList) {
	// Build index on name, we'll delete entries out of here a we match them
	// against /etc/hosts entries.  The remaining entries aren't present
	// locally anymore and will need to be deleted.
	rhByName := map[string]hostEntry{}
	for _, rh := range r53hosts {
		rhByName[rh.hostname] = rh
	}

	toUpdate := hostList{}
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

	toDelete := make(hostList, 0, len(rhByName))
	for _, rh := range rhByName {
		toDelete = append(toDelete, rh)
	}

	return toUpdate, toDelete
}

func removeDupes(hosts hostList) hostList {
	found := make(map[string]bool, len(hosts))

	// Sort hostlist to ensure stable duplication suppression.  We don't want
	// to ping pont between choosing different options because of parse order.
	sort.Sort(hosts)

	dupCount := 0
	result := make(hostList, 0, len(hosts))
	for _, h := range hosts {
		if _, ok := found[h.hostname]; ok {
			log.Warnf("Duplicate hostname found in /etc/hosts, ignoring (%v/%v)",
				h.hostname, h.ip.String())
			dupCount++
		} else {
			found[h.hostname] = true
			result = append(result, h)
		}
	}

	return result
}

func runOnce() {
	hosts := readHosts(opts.File)
	hosts = filterHosts(hosts, opts.Networks)
	if !opts.NoQualifyHosts {
		hosts = qualifyHosts(hosts, opts.Domain)
	}
	hosts = removeDupes(hosts)

	r53 := newRoute53()
	r53Hosts, err := r53.getHosts(opts.Domain)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error when retrieving zones"))
	}
	r53Hosts = filterHosts(r53Hosts, opts.Networks)

	toUpdate, toDelete := compareHosts(hosts, r53Hosts)
	if len(toUpdate) > 0 || len(toDelete) > 0 {
		if err := r53.sync(opts.Domain, opts.TTL, !opts.NoWait, toUpdate, toDelete); err != nil {
			log.Fatal(errors.Wrap(err, "Could not sync records to Route 53"))
		}
	} else {
		log.Info("No changes needed.  Everything in sync.")
	}

}

func main() {
	parseOpts()
	configureLogging()
	if opts.Mode == "oneshot" {
		runOnce()
	} else {
		daemon(opts.Interval, opts.File)
	}
}
