package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

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
		return &hostEntry{parts[1], ip, parts[2:]}, nil
	}

	return nil, fmt.Errorf("%s is not a valid IP", parts[0])
}

func readHosts(filename string) (hosts []hostEntry) {
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
			log.Printf("WARN %v on line %v, skipping\n", err, i)
			continue
		}
		if host != nil {
			hosts = append(hosts, *host)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return
}

func filterHosts(hosts []hostEntry, networks []CIDRNet) []hostEntry {
	output := []hostEntry{}
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
