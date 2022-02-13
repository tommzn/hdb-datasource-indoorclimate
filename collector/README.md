![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hdb-datasource-indoorclimate/collector)

# HomeDashboard Indoor Climate DataCollector
Package to create and run a Indoor Climate data collector, e.g. as a systemd service.

## Params
### configfile
Provide path to a config file using this param. Default is /etc/hdb/config.yml.
### secretsfile
Provide path to a credentials file. This is necessary if you e.g. want to use a SQS or Timestream target to publish indoor climate date. Default is ~/.hdb/credentials.

## Config
With following config a AWS SQS and AWS Timesteam target is appended to the internal list of indoor climate data publishers. A LogTarget is used always as default publisher.
```yaml
hdb:
  queue: sqs-queue
  
aws:
  sqs:
    region: eu-west-1
  timestream:
    region: eu-west-1
    database: timestreamdb
    table: timestreamtable
    batch_size: 10

```

# Links
- [HomeDashboard Documentation](https://github.com/tommzn/hdb-docs/wiki)
