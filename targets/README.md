[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/hdb-datasource-indoorclimate/targets.svg)](https://pkg.go.dev/github.com/tommzn/hdb-datasource-indoorclimate/targets)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hdb-datasource-indoorclimate/targets)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tommzn/hdb-datasource-indoorclimate/targets)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/hdb-datasource-indoorclimate/targets)](https://goreportcard.com/report/github.com/tommzn/hdb-datasource-indoorclimate/targets)

# Indoor Climate Data Targets
This package contains different targets indoor climate data can be send to. You can assign such a target to the SensorDataCollector using it's AppendTarget method-

## Stdout Publisher
Writes indoor climate data to stout using fmt.

## Log Publisher
Target, which writes indoor climate data to given logger with log level Info. See [Log](https://github.com/tommzn/go-log) for more details about used logger.

## AWS SQS Publisher
This target will send indoor climate measurements to a AWS SQS queue if one has been defined in config as followed.
```yaml
hdb:
  queue: sqs-queue
```
See [Metrics](https://github.com/tommzn/go-metrics) for ore details about timesteam integration.

## AWS Timestream Publisher
To collect indoor climate data in a timestream database (AWS Timestream) you can provide a timestream config and a correcponding publisher will be added.
```yaml
aws:
  timestream:
    region: eu-west-1
    database: timestreamdb
    table: timestreamtable
    batch_size: 10
```
<strong>Note: In addition to this config you've to provide AWS access keys with correct permissions.</strong>

# Links
- [HomeDashboard Documentation](https://github.com/tommzn/hdb-docs/wiki)
