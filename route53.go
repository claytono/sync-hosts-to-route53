package main

import (
	"fmt"
	"log"
	"net"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/pkg/errors"
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
		return nil, errors.Wrap(err, "Cannot list zones")
	}

	if len(resp.HostedZones) == 0 {
		return nil, fmt.Errorf("could not find domain '%v'", domain)
	}

	if *resp.HostedZones[0].Name != (domain + ".") {
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
		return nil, errors.Wrap(err, "Cannot get records")
	}

	return resp.ResourceRecordSets, nil
}

func (r53 route53Client) getHosts(domain string) (hostList, error) {
	zone, err := r53.getZone(domain)
	if err != nil {
		return hostList{}, errors.Wrap(err, "Cannot get zone")
	}

	rawHosts, err := r53.getRecords(*zone.Id)
	if err != nil {
		return hostList{}, errors.Wrap(err, "Cannot get hosts")
	}

	return convertR53RecordsToHosts(rawHosts), nil
}

func (r53 route53Client) sync(domain string, ttl int64, toUpdate []hostEntry, toDelete []hostEntry) error {
	zone, err := r53.getZone(domain)
	if err != nil {
		return errors.Wrap(err, "Cannot get zone")
	}

	changes := make([]*route53.Change, 0, len(toUpdate)+len(toDelete))
	for _, h := range toUpdate {
		change := route53.Change{
			Action: aws.String("UPSERT"),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name: aws.String(h.hostname),
				Type: aws.String("A"),
				TTL:  aws.Int64(ttl),
				ResourceRecords: []*route53.ResourceRecord{
					{Value: aws.String(h.ip.String())},
				},
			},
		}
		changes = append(changes, &change)
	}

	for _, h := range toDelete {
		change := route53.Change{
			Action: aws.String("DELETE"),
			ResourceRecordSet: &route53.ResourceRecordSet{
				Name: &h.hostname,
				Type: aws.String("A"),
				TTL:  &ttl,
				ResourceRecords: []*route53.ResourceRecord{
					{Value: aws.String(h.ip.String())},
				},
			},
		}
		changes = append(changes, &change)
	}

	input := route53.ChangeResourceRecordSetsInput{
		HostedZoneId: zone.Id,
		ChangeBatch: &route53.ChangeBatch{
			Changes: changes,
		},
	}

	fmt.Printf("Adding/updating %v records, deleting %v out of date records\n",
		len(toUpdate), len(toDelete))

	_, err = r53.svc.ChangeResourceRecordSets(&input)
	if err != nil {
		return errors.Wrapf(err, "Could not update Route 53 records")
	}

	return nil
}

func convertR53RecordsToHosts(rawHosts []*route53.ResourceRecordSet) hostList {
	hosts := hostList{}
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
		host := hostEntry{
			hostname: canonifyHostname(*rh.Name),
			ip:       ip,
		}

		hosts = append(hosts, host)
	}

	return hosts
}
