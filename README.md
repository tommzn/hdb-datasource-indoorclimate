[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/hdb-datasource-indoorclimate.svg)](https://pkg.go.dev/github.com/tommzn/hdb-datasource-indoorclimate)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hdb-datasource-indoorclimate)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tommzn/hdb-datasource-indoorclimate)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/hdb-datasource-indoorclimate)](https://goreportcard.com/report/github.com/tommzn/hdb-datasource-indoorclimate)
[![Actions Status](https://github.com/tommzn/hdb-datasource-indoorclimate/actions/workflows/go.pkg.auto-ci.yml/badge.svg)](https://github.com/tommzn/hdb-datasource-indoorclimate/actions)

# HomeDashboard Indoor Climate DataSource
Fetches indoor climate data from Bluetooth sensor devices, e.g. Xiaomi Mi Temperature and Humidity Monitor 2, and publishes this data to a AWS SQS queue for futher processing.

## Define Devices and Characteristics
You've to add a list of devices and characteristics you want to observe in config.
```yaml
indoorclimate:
  schedule: 10m
  devices:
    - id: "A4:XX:XX:XX:26:41"
  characteristics:
    - uuid: "00002a6e-0000-1000-8000-00805f9b34fb"
      type: "temperature"
    - uuid: "00002a6f-0000-1000-8000-00805f9b34fb"
      type: "humidity"
    - uuid: "00002a19-0000-1000-8000-00805f9b34fb"
      type: "battery"
```
### Schedule
This value defines how often sensor data should be read. See [Config](https://github.com/tommzn/go-config) for more details about supported values.

### Devices
Provide a list of MAC addesses for Bluetooth environment sensors.

### Characteristics
Define a list of observed characteristics by their UUID and specifiy a indoor cliamte data type.

## Targets
In addition to a default log target, a SQS publisher target or a AWS Timestream target will be added automatically if correcponding config is set.

### Log Publisher
Default target, will be added all time and writes indoor climate measurements with log level debug. See [Log](https://github.com/tommzn/go-log) for more details about used logger.

### AWS SQS Publisher
This target will send indoor climate measurements to a AWS SQS queue if one has been defined in config as followed.
```yaml
hdb:
  queue: sqs-queue
```
See [Metrics](https://github.com/tommzn/go-metrics) for ore details about timesteam integration.

### AWS Timestream Publisher
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

## Usage
After creating a new collector you can call it's Run method to start consuming new indoor climate data from MQTT broker. By default all received indoor climate data are send
to default target which is a logger, only. Collector will run until you cancel passed context.
```golang

    import (
       indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"  
       config "github.com/tommzn/go-config"
	     log "github.com/tommzn/go-log"
	     secrets "github.com/tommzn/go-secrets"
    )
    
    collector, err := indoorclimate.New(conf, logger, secretsmanager)
    if err != nil {
        panic(err)
    }

    ctx, cancelFunc := context.WithCancel(context.Background())
    err := collector.Run(ctx)
    if err != nil {
        panic(err)
    }
```

# Links
- [HomeDashboard Documentation](https://github.com/tommzn/hdb-docs/wiki)
