# Sync /etc/hosts to AWS Route 53

[![Build Status](https://travis-ci.org/claytononeill/sync-hosts-to-route53.svg?branch=master)](https://travis-ci.org/claytononeill/sync-hosts-to-route53) [![Go Report Card](https://goreportcard.com/badge/github.com/claytononeill/sync-hosts-to-route53)](https://goreportcard.com/report/github.com/claytononeill/sync-hosts-to-route53)

## What does this do?

This is a tool that will read a local file in `/etc/hosts` format and
synchronize the contents of that file to a specific domain in AWS Route 53.
The contents of the /etc/hosts file can be filtered such that only records in
specific networks will be synchronized.  By default the program will run
continuously, detecting changes in the hosts file automatically via inotify and
running periodic syncs to clean up any rogue Route 53 changes.

This was written to be used on Ubiquiti EdgeOS devices like the EdgeRouter
Lite.  This can be useful if you have public IP ranges hosted behind your
router, or if you have a VPN connection into the a private network.  This will
allow you to use any DNS caching resolver, instead of forcing DNS resolution
through a private DNSMasq instance.

There is nothing inherently specific to EdgeOS in this project, so if you need
something to sync an `/etc/hosts` file to Route53 once, or on a regular basis,
this might work well for you.

## What do I need to know?

* The local hosts file given is considered to be the source of truth for the
  networks  you specify to be managed.  Changes in Route 53 that don't match
  what is in the hosts file will be removed or overwritten.

* Host aliases in the input file are ignored.  Entries in `/etc/hosts` on
  EdgeOS devices that are added via DHCP never have alias entries.

* This app is written in Go, mostly because it generates binaries with no
  dependencies and it makes cross compiles for other architectures easy.

* Binaries can be found on the [GitHub releases
  page](https://github.com/claytononeill/sync-hosts-to-route53/releases)

* MIPS64 binaries are built and available on the [GitHub releases
  page](https://github.com/claytononeill/sync-hosts-to-route53/releases), but
  aren't recommended.  When tested on EdgeOS 1.9.1, inotify does not appear to
  work with the MIPS64 binary, but does work with the MIPS (32 bit) build.  In
  addition, when testing this in QEMU w/MIPS64, the binary just crashes on
  startup in the Golang GC code.

## AWS Authentication

This program uses the [official AWS golang
libary](https://github.com/aws/aws-sdk-go).  Because of this, any of the normal
methods of storing credentials for AWS that work with the CLI tools should work
for this program.  For example, you can set and export the `AWS_ACCESS_KEY_ID`
and `AWS_SECRET_ACCESS_KEY` environment variables.  More information can be
found in the "[SDK Configuration](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials)"
section of the library documentation.

It is recommended that you set up an AWS IAM user specifically to be used with
this program.  That IAM user should be limited to only be able to add and
remove records in Route 53.  If you have multiple domains, you may wish to
limit that user to just the specific domain you wish to synchronize against. 

## Usage

Use `sync-hosts-to-route53 --help` for a complete list of available options.

The options available include:

### -m|--mode [oneshot|daemon]

This options must be either `oneshot` or `daemon`.  The default is `daemon`.

When run with `--mode daemon` or no `--mode` argument, the program will
synchronize the host file with Route 53 once, then setup inotify watches for
the host file, and synchronize automatically by default every 15 minutes.

When run with `--mode oneshot` the program will synchronize the host file
given with Route 53 once, then exit.

### -f|--file=HOSTFILE

This specifies the local hosts file to keep in sync with Route 53.  This should
be in the format of UNIX style `/etc/hosts` file.  This defaults to
`/etc/hosts`.

### --network=x.x.x.x/len

This option will direct the program to ignore all host entries in the local
file, or in the Route 53 domain that are not inside of the network blocks
given.  This option is required and can be specified more than once.  If you
wish to affect **all** entries in the domain, then you can specify `0.0.0.0/0`
to match all IP addresses.

### -d|--domain=

This specifies the Route 53 domain to synchronize with the local hosts file.
This option is required and has no default.

### -i|--interval=

How often to run the synchronization, even if no changes have been detected in
the local hosts file.  Local changes should be detected automatically, so this
mostly serves to correct any changes made in Route 53 that don't match the
local file.  This defaults to 15 minutes.  This is ignored in oneshot mode.

### --ttl=

This is the DNS record TTL in seconds to set on new Route 53 records.  This
defaults to 3600 seconds, or one hour.

### --no-qualify-hosts

By default the Route53 domain will be appended to the end of any host file
entries that appear to be lacking it.  To disable this behavior, specify
`--no-qualifiy-hosts`.

### --exclude-hosts

Exclude specific hosts from being synced to Route53.  This can be used to
prevent manually created items from being deleted during the sync process.
This can be specified multiple times.

### --no-wait

By default the program will wait for Route 53 updates to propagate after
submitting them.  To disable this behavior, specify `--no-wait`.

### --syslog

Enable logging to syslog in addition to stdout.

### --syslog-facility=

Syslog facility to log to.  Defaults to `user`.  See `syslog(3)` for list of
options.

### --syslog-only

Log just to syslog, skipping stdout entirely.  If this is specified then
`--syslog` is implied.  Note a single message will be logged to stdout
indicating all further output will go to syslog.

### --debug

Enable debug level logging.

### --version

Print version to stdout and exit.

### -h|--help

Print summary of these options

## Building From Source

Dependencies for this project are managed using
[Glide](https://github.com/Masterminds/glide).  To build you will want to run
`glide install` to pull in the dependencies, then run `make build` to build the
binary with version information.

To build binaries for all platforms (used for releases), run `make
build-all-arch`.

## Author

Clayton O'Neill
<clayton@oneill.net>
