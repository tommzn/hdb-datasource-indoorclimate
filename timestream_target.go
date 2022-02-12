package indoorclimate

import (
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	metrics "github.com/tommzn/go-metrics"
)

// NewTimestreamTarget returns a new target which writes measurements to AWS Timestream.
func newTimestreamTarget(conf config.Config, logger log.Logger) Publisher {
	return &timestreamTarget{
		metricPublisher: metrics.NewTimestreamPublisher(conf, logger),
	}
}

// Send given indoor climate data to AWS Timestream.
func (target *timestreamTarget) SendMeasurement(measurement IndoorClimateMeasurement) error {
	timesteamMetrics := metrics.Measurement{
		MetricName: "hdb-datasource-indoorclimate",
		Tags: []metrics.MeasurementTag{
			metrics.MeasurementTag{
				Name:  "deviceid",
				Value: measurement.DeviceId,
			},
			metrics.MeasurementTag{
				Name:  "type",
				Value: measurement.Type.String(),
			},
		},
		Values: []metrics.MeasurementValue{
			metrics.MeasurementValue{
				Name:  "count",
				Value: 1,
			},
			metrics.MeasurementValue{
				Name:  measurement.Type.String(),
				Value: measurement.Value,
			},
		},
	}
	target.metricPublisher.Send(timesteamMetrics)
	target.metricPublisher.Flush()
	return target.metricPublisher.Error()
}
