package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"
)

type hostEntry struct {
	hostname string
	ip       net.IP
	// Aliases are read from the /etc/hosts file, but not mapped to Route53
	aliases []string
	// rrset only exists for imported Route 53 records
	rrset *route53.ResourceRecordSet
}

type hostList []hostEntry

func (h hostList) Len() int {
	return len(h)
}

func (h hostList) Less(i, j int) bool {
	if h[i].hostname != h[j].hostname {
		return h[i].hostname < h[j].hostname
	}
	return bytes.Compare(h[i].ip, h[j].ip) < 0
}

func (h hostList) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func parseLine(line string) (*hostEntry, error) {
	if i := strings.Index(line, "#"); i >= 0 {
		line = line[0:i]
	}

	parts := strings.Fields(line)
	if len(parts) == 0 {
		return nil, nil
	}

	if len(parts) < 2 {
		return nil, fmt.Errorf("should contain at least two fields")
	}

	if ip := net.ParseIP(parts[0]); ip != nil {
		return &hostEntry{
			hostname: parts[1],
			ip:       ip,
			aliases:  parts[2:],
		}, nil
	}

	return nil, fmt.Errorf("%s is not a valid IP", parts[0])
}

func readHosts(filename string) (hosts hostList) {
	file, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		i++
		host, err := parseLine(scanner.Text())
		if err != nil {
			log.Warnf("%v on line %v, skipping\n", err, i)
			continue
		}
		if host != nil {
			host.hostname = canonifyHostname(host.hostname)
			hosts = append(hosts, *host)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

func filterHosts(hosts hostList, networks []CIDRNet) hostList {
	output := hostList{}
	for _, host := range hosts {
		for _, net := range networks {
			if net.Contains(host.ip) {
				output = append(output, host)
				break
			}
		}
	}
	return output
}

func qualifyHosts(hosts hostList, domain string) hostList {
	result := make(hostList, len(hosts))
	for i, h := range hosts {
		result[i] = h
		if !strings.HasSuffix(h.hostname, domain) {
			result[i].hostname = h.hostname + "." + domain
		}
	}

	return result
}
