[![Actions Status](https://github.com/tommzn/hdb-datasource-exchangerate/actions/workflows/go.image.build.yml/badge.svg)](https://github.com/tommzn/hdb-datasource-exchangerate/actions)

# Indoor Climate Data Collector
This package composes a [contimuous collector](https://github.com/tommzn/hdb-datasource-core/collector.go) and [indoor climate data source](https://github.com/tommzn/hdb-datasource-indoorclimate) to subscribe to MQTT topics, from ioBroker, which delivers indoor climate data. 

## Config
This collector requires a config to get settings for MQTT broker and logging.

### Example 
```yaml
log:
  loglevel: error
  shipper: logzio  

mqtt:
  topic_prefix: iobroker
  broker: mqtt-broker-01
  port: 1883
