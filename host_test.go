package main

import (
	"net"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

var hostCases = []struct {
	name     string
	unsorted hostList
	sorted   hostList
}{
	{name: "already sorted",
		unsorted: hostList{
			{hostname: "test1.com", ip: net.ParseIP("1.2.3.4")},
			{hostname: "test2.com", ip: net.ParseIP("1.2.3.4")}},
		sorted: hostList{
			{hostname: "test1.com", ip: net.ParseIP("1.2.3.4")},
			{hostname: "test2.com", ip: net.ParseIP("1.2.3.4")}},
	},
	{name: "reverse sorted",
		unsorted: hostList{
			{hostname: "test2.com", ip: net.ParseIP("1.2.3.4")},
			{hostname: "test1.com", ip: net.ParseIP("1.2.3.4")}},
		sorted: hostList{
			{hostname: "test1.com", ip: net.ParseIP("1.2.3.4")},
			{hostname: "test2.com", ip: net.ParseIP("1.2.3.4")}},
	},
	{name: "same name, different ip",
		unsorted: hostList{
			{hostname: "test.com", ip: net.ParseIP("1.2.3.5")},
			{hostname: "test.com", ip: net.ParseIP("1.2.3.4")}},
		sorted: hostList{
			{hostname: "test.com", ip: net.ParseIP("1.2.3.4")},
			{hostname: "test.com", ip: net.ParseIP("1.2.3.5")}},
	},
}

func TestSortHostList(t *testing.T) {
	for _, c := range hostCases {
		t.Run(c.name, func(t *testing.T) {
			sort.Sort(c.unsorted)
			assert.Equal(t, c.unsorted, c.sorted)
		})
	}
}
