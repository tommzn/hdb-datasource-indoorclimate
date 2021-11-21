[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/hdb-datasource-indoorclimate.svg)](https://pkg.go.dev/github.com/tommzn/hdb-datasource-indoorclimate)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hdb-datasource-indoorclimate)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tommzn/hdb-datasource-indoorclimate)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/hdb-datasource-indoorclimate)](https://goreportcard.com/report/github.com/tommzn/hdb-datasource-indoorclimate)
[![Actions Status](https://github.com/tommzn/hdb-datasource-indoorclimate/actions/workflows/go.pkg.auto-ci.yml/badge.svg)](https://github.com/tommzn/hdb-datasource-indoorclimate/actions)

# HomeDashboard Indoor Climate DataSource
Fetches indoor climate data from a MQTT broker and publishes to HomebDaskboard backend.

## Topics
This indoor climate consumer subscribes to MQTT topics for Bluetooth data send from ioBroker to consume indoor climate data like temperature, humidity and battery status of a sensor.
You can specific a prefix for this three topics by config.

### Consumed Topics
| Measurement Type      | Topic                |
| --------------------- | -------------------- |
| Temperature           | /ble/+/+/temperature |
| Humidity              | /ble/+/+/humidity    |
| Device Battery Status | /ble/+/+/battery     |

## Config
Config can be used to specific MQTT broker, port and a prefix for topics. If nothing has been defined it tries to connect to a MQTT broker at localhost:1883 and subscribes to ioBroker topics with indoor climate data send for Bluetooth devices.
More details about loading config at https://github.com/tommzn/go-config

### Config example
```yaml
mqtt:
  topic_prefix: iobroker
  broker: mqtt-broker-01
  port: 1883
```

## Targets
[MessgeTarget](https://github.com/tommzn/hdb-datasource-indoorclimate/blob/main/interfaces.go) interface is used for destinations indoor climate data are send to. By default a consumer contains a log target, only. You can register additional targets using 

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
