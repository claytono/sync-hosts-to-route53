# Sync /etc/hosts to AWS Route 53

[![Build Status](https://travis-ci.org/claytononeill/sync-hosts-to-route53.svg?branch=master)](https://travis-ci.org/claytononeill/sync-hosts-to-route53) [![Go Report Card](https://goreportcard.com/badge/github.com/claytononeill/sync-hosts-to-route53)](https://goreportcard.com/report/github.com/claytononeill/sync-hosts-to-route53)

*NOTE: THIS IS A WORK IN PROGRESS*

This is a tool that will read a local file in `/etc/hosts` format any synchronize the contents of that file to a specific domain in AWS Route 53.  The contents of the /etc/hosts file can be filtered such that only records in specific networks will be synchronized.

This is intended to be used on Ubiquiti EdgeOS devices like the EdgeRouter Lite.  This can be useful if you have public IP ranges hosted behind your router, or if you have a VPN connection into the a private network.  This will allow you to use any DNS caching resolver, instead of forcing DNS resolution through a private DNSMasq instance.