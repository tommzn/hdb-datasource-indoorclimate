[![Go Reference](https://pkg.go.dev/badge/github.com/tommzn/hdb-datasource-indoorclimate.svg)](https://pkg.go.dev/github.com/tommzn/hdb-datasource-indoorclimate)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/tommzn/hdb-datasource-indoorclimate)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/tommzn/hdb-datasource-indoorclimate)
[![Go Report Card](https://goreportcard.com/badge/github.com/tommzn/hdb-datasource-indoorclimate)](https://goreportcard.com/report/github.com/tommzn/hdb-datasource-indoorclimate)
[![Actions Status](https://github.com/tommzn/hdb-datasource-indoorclimate/actions/workflows/go.pkg.auto-ci.yml/badge.svg)](https://github.com/tommzn/hdb-datasource-indoorclimate/actions)

# HomeDashboard Indoor Climate DataSource
Fetches indoor climate data from Bluetooth sensor devices, e.g. Xiaomi Mi Temperature and Humidity Monitor 2, and publishes this data to specified targets.

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
By default a SensorDataCollector doesn't have any target assign. This means your indoor climate date get lost. [Targets](https://github.com/tommzn/hdb-datasource-indoorclimate/tree/main/targets) package provides different publishers you can assign to a SensorDataCollector to send indoor climate data to a target.

## Usage
After creating a new collector you can call it's Run method to start consuming new indoor climate data from MQTT broker. By default all received indoor climate data are send
to default target which is a logger, only. Collector will run until you cancel passed context.
```golang

    import (
      "fmt"
      "context"

       indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"  
       targets "github.com/tommzn/hdb-datasource-indoorclimate(targets"  
       config "github.com/tommzn/go-config"
	     log "github.com/tommzn/go-log"
    )
    
    conf, _ := config.NewConfigSource().Load()
    if err != nil {
        panic(err)
    }
    logger := log.NewLoggerFromConfig(conf, nil)

    datacollector := indoorclimate.NewSensorDataCollector(conf, logger)
    if err != nil {
        panic(err)
    }

    datacollector.AppendTarget(targets.NewStdoutTarget())
    if err := datacollector.Run(context.Background()); err != nil {
      fmt.Println)err  
    }
    
```

# Links
- [HomeDashboard Documentation](https://github.com/tommzn/hdb-docs/wiki)
