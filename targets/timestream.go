package targets

import (
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	metrics "github.com/tommzn/go-metrics"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

// NewTimestreamTarget returns a new target which writes measurements to AWS Timestream.
func newTimestreamTarget(conf config.Config, logger log.Logger) indoorclimate.Publisher {
	return &TimestreamTarget{
		metricPublisher: metrics.NewTimestreamPublisher(conf, logger),
	}
}

// Send given indoor climate data to AWS Timestream.
func (target *TimestreamTarget) SendMeasurement(measurement indoorclimate.IndoorClimateMeasurement) error {

	measurementType := toEventType(measurement.Type)
	timesteamMetrics := metrics.Measurement{
		MetricName: "hdb-datasource-indoorclimate",
		Tags: []metrics.MeasurementTag{
			metrics.MeasurementTag{
				Name:  "deviceid",
				Value: measurementType.String(),
			},
			metrics.MeasurementTag{
				Name:  "type",
				Value: measurementType.String(),
			},
		},
		Values: []metrics.MeasurementValue{
			metrics.MeasurementValue{
				Name:  "count",
				Value: 1,
			},
			metrics.MeasurementValue{
				Name:  measurementType.String(),
				Value: measurement.Value,
			},
		},
	}
	target.metricPublisher.Send(timesteamMetrics)
	target.metricPublisher.Flush()
	return target.metricPublisher.Error()
}
