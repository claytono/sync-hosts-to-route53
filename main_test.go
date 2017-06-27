package main

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompareHosts(t *testing.T) {
	cases := []struct {
		name     string
		hosts    hostList
		r53hosts hostList
		toUpdate hostList
		toDelete hostList
	}{
		{"noop",
			hostList{
				{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
			},
			hostList{
				{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
			},
			hostList{},
			hostList{},
		},
		{"need-update",
			hostList{
				{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
			},
			hostList{
				{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.5")},
			},
			hostList{
				{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
			},
			hostList{},
		},
		{"add-new",
			hostList{
				{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
			},
			hostList{},
			hostList{
				{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
			},
			hostList{},
		},
		{"remove-stale",
			hostList{
				{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
			},
			hostList{
				{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
				{hostname: "test2.test.com", ip: net.ParseIP("1.2.3.5")},
			},
			hostList{},
			hostList{
				{hostname: "test2.test.com", ip: net.ParseIP("1.2.3.5")},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			toUpdate, toDelete := compareHosts(c.hosts, c.r53hosts)
			assert.Equal(t, c.toUpdate, toUpdate)
			assert.Equal(t, c.toDelete, toDelete)
		})
	}
}

func TestRemoveExcludedHosts(t *testing.T) {
	cases := []struct {
		name     string
		hosts    hostList
		excludes []string
		result   hostList
	}{
		{"noop",
			hostList{
				{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
			},
			[]string{"test2.test.com"},
			hostList{
				{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
			},
		},
		{"remove-only-one",
			hostList{
				{hostname: "test1.test.com", ip: net.ParseIP("1.2.3.4")},
			},
			[]string{"test1.test.com"},
			hostList{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := removeExcludedHosts(c.hosts, c.excludes)
			assert.Equal(t, result, c.result)
		})
	}
}
