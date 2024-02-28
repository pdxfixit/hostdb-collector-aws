# hostdb-collector-aws

Gather useful information from AWS accounts and put it into HostDB.

## Build it and Run it

The container image is built [here](https://builds.pdxfixit.com/gh/hostdb-collector-aws).

The image is stored in [the registry](https://registry.pdxfixit.com/hostdb-collector-aws).

The container is run as a Kubernetes cronjob, in the [hostdb-server helm chart](https://github.com/pdxfixit/hostdb-server-chart/blob/master/hostdb-server/templates/collector-aws.yaml).

The data is published in [HostDB](https://hostdb.pdxfixit.com/?type=/aws/).
