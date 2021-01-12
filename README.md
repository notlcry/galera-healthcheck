galera-healthcheck
==================

[![Build Status](https://travis-ci.org/cloudfoundry-incubator/galera-healthcheck.svg?branch=master)](https://travis-ci.org/cloudfoundry-incubator/galera-healthcheck)


This go-based process is designed to run on a MariaDB Galera node and monitor the health of the node.
An http endpoint is opened, by default at '/' on port 9200.
A healthy node will return HTTP status 200, and a node that should not be accessed returns a 503.

Several commandline flags are supported (see the code.)

- build

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build
```

- run cmd

```bash
galera-healthcheck -host=169.254.254.1 -user=system -password=System@123 -port=8081 -availWhenDonor=true -availWhenReadOnly=true -pidfile=/var/run/galera-healthcheck-2.pid
```
