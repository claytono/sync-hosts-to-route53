package main

import (
	"fmt"
	"log"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type route53Client struct {
	sess *session.Session
	svc  *route53.Route53
}

func newRoute53() route53Client {
	var r53 route53Client

	// awsSession is global so that we only read the config once and reuse it
	r53.sess = session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	r53.svc = route53.New(r53.sess)

	return r53
}

func (r53 route53Client) getZone(domain string) (*route53.HostedZone, error) {
	params := &route53.ListHostedZonesByNameInput{
		DNSName: aws.String(domain),
	}
	resp, err := r53.svc.ListHostedZonesByName(params)
	if err != nil {
		return nil, err
	}

	if len(resp.HostedZones) == 0 {
		return nil, fmt.Errorf("could not find domain '%v'", domain)
	}

	return resp.HostedZones[0], nil
}

func (r53 route53Client) getRecords(zid string) ([]*route53.ResourceRecordSet, error) {
	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: &zid,
	}
	resp, err := r53.svc.ListResourceRecordSets(params)
	if err != nil {
		return nil, err
	}

	return resp.ResourceRecordSets, nil
}

func (r53 route53Client) getHosts(domain string) ([]hostEntry, error) {
	zone, err := r53.getZone(domain)
	if err != nil {
		return []hostEntry{}, err
	}

	rawHosts, err := r53.getRecords(*zone.Id)
	if err != nil {
		return []hostEntry{}, err
	}

	return convertR53RecordsToHosts(rawHosts), nil
}

func convertR53RecordsToHosts(rawHosts []*route53.ResourceRecordSet) []hostEntry {
	hosts := []hostEntry{}
	for _, rh := range rawHosts {
		if *rh.Type != "A" {
			continue
		}

		if len(rh.ResourceRecords) > 1 {
			log.Printf("WARN %v has too many resource records (%d), ignoring record",
				*rh.Name, len(rh.ResourceRecords))
			continue
		}
		ip := net.ParseIP(*rh.ResourceRecords[0].Value)
		if ip == nil {
			log.Printf("WARN cannot parse IP %v for %v, ignoring record",
				*rh.ResourceRecords[0].Value, *rh.Name)
			continue
		}
		host := hostEntry{hostname: *rh.Name, ip: ip}
		hosts = append(hosts, host)
	}

	return hosts
}
