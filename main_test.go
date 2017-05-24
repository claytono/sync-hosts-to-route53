package main

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cases = []struct {
	name     string
	hosts    []hostEntry
	r53hosts []hostEntry
	toUpdate []hostEntry
	toDelete []hostEntry
}{
	{"noop",
		[]hostEntry{
			{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
		},
		[]hostEntry{
			{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
		},
		[]hostEntry{},
		[]hostEntry{},
	},
	{"need-update",
		[]hostEntry{
			{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
		},
		[]hostEntry{
			{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.5")},
		},
		[]hostEntry{
			{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
		},
		[]hostEntry{},
	},
	{"add-new",
		[]hostEntry{
			{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
		},
		[]hostEntry{},
		[]hostEntry{
			{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
		},
		[]hostEntry{},
	},
	{"remove-stale",
		[]hostEntry{
			{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
		},
		[]hostEntry{
			{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
			{hostname: "test2.test.com", ip: net.ParseIP("1.2.3.5")},
		},
		[]hostEntry{},
		[]hostEntry{
			{hostname: "test2.test.com", ip: net.ParseIP("1.2.3.5")},
		},
	},
}

func TestCompareHosts(t *testing.T) {
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			toUpdate, toDelete := compareHosts(c.hosts, c.r53hosts)
			assert.Equal(t, c.toUpdate, toUpdate)
			assert.Equal(t, c.toDelete, toDelete)
		})
	}
}
