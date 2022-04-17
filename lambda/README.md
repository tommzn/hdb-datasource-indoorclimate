[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/hdb-datasource-indoorclimate.svg)](https://pkg.go.dev/github.com/tommzn/hdb-datasource-indoorclimate/lambda)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hdb-datasource-indoorclimate/lambda)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/hdb-datasource-indoorclimate)](https://goreportcard.com/report/github.com/tommzn/hdb-datasource-indoorclimate/lambda)

# Indoor Climate Measurement Processor
Lambda function to process indoor climate measurements send from a device to AWS IOT. 

## Config
You've to provice a config file to enable forwarding for measurements. By default no forwarding target is set. This means your measurement will get lost.
A config can be provided as a local file, config.yml, or via S3 bucket. See [go-config](https://github.com/tommzn/go-config) for more details.
```yaml
hdb:
  queue: AWS SQS Queue
  archive: AWS SQS Queue

aws:
  sqs:
    region: us-west-1
  timestream:
    region: us-west-1
    database: database
    table: table
    batch_size: 10
```
### AWS SQS Target
To enable event publoshing to an AWS SQS queue, provide queue name and SQS region.

### AWS Timestream
To publish metrics to AWS Timestream provide region, database and a table.

# Links
- [HomeDashboard Documentation](https://github.com/tommzn/hdb-docs/wiki)
