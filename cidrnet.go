package main

import "net"

// CIDRNet is a net.IPNet wrapper that implements the Marshal/UnMarshal
// interface that go-flags wants.  This allows the flag parser to produce the
// error message directly, instead of us post-processing the results.
type CIDRNet struct {
	net.IPNet
}

// UnmarshalFlag allows go-flags package to parse CIDR blocks directly
func (n *CIDRNet) UnmarshalFlag(value string) error {
	_, ipnet, err := net.ParseCIDR(value)
	if err != nil {
		return err
	}

	n.IPNet = *ipnet
	return nil
}

// MarshalFlag allows go-flags package to print CIDR blocks directly.
func (n CIDRNet) MarshalFlag() (string, error) {
	return n.IPNet.String(), nil
}
