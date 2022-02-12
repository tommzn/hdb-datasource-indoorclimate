package indoorclimate

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	utils "github.com/tommzn/go-utils"
	core "github.com/tommzn/hdb-datasource-core"
	events "github.com/tommzn/hdb-events-go"
)

func NewSensorDataCollector(conf config.Config, logger log.Logger) core.Collector {

	schedule := conf.GetAsDuration("indoorclimate.schedule", nil)
	retryCount := conf.GetAsInt("indoorclimate.retry", config.AsIntPtr(3))
	adapterId := conf.Get("indoorclimate.adapter", config.AsStringPtr("hci0"))

	deviceIds := deviceIdsFromConfig(conf)
	devices := []SensorDevice{}
	for _, deviceId := range deviceIds {
		devices = append(devices, NewIndoorClimateSensor(*adapterId, deviceId))
	}
	logger.Debugf("Number of observed devices: %d", len(devices))

	characteristics := characteristicsFromConfig(conf)
	logger.Debugf("Number of observed characteristics: %d", len(characteristics))

	publisher := []Publisher{newLogPublisher(logger)}
	if queue := conf.Get("hdb.queue", nil); queue != nil {
		publisher = append(publisher, NewSqsTarget(conf, logger))
		logger.Debug("SQS Publisher added.")
	}
	if timestreamTable := conf.Get("aws.timestream.table", nil); timestreamTable != nil {
		publisher = append(publisher, newTimestreamTarget(conf, logger))
		logger.Debug("Timestream Publisher added.")
	}

	return &SensorDataCollector{
		schedule:        schedule,
		logger:          logger,
		devices:         devices,
		characteristics: characteristics,
		publisher:       publisher,
		retryCount:      *retryCount,
	}
}

// Run will start collecting sensor data from all defined devices.
func (collector *SensorDataCollector) Run(ctx context.Context) error {

	collector.errorStack = utils.NewErrorStack()
	if collector.schedule == nil {
		collector.RunSingle(ctx)
	} else {
		collector.RunContinouous(ctx)
	}
	return collector.errorStack.AsError()
}

func (collector *SensorDataCollector) RunSingle(ctx context.Context) {

	defer collector.logger.Flush()

	collector.done = make(chan struct{})
	for _, device := range collector.devices {

		go collector.readDevciceData(device)
		select {
		case <-collector.done:
			collector.logger.Debug("Data collection finished for device: ", device.Id())
		case <-ctx.Done():
			collector.errorStack.Append(errors.New("Sensor data collection has been canceled."))
		}
	}

}

func (collector *SensorDataCollector) RunContinouous(ctx context.Context) {

	collector.logger.Debugf("Run continuous collection with schedule of: %s", *collector.schedule)
	for {
		collector.RunSingle(ctx)
		select {
		case <-time.After(*collector.schedule):
			collector.logger.Debug("Restart Data collection")
		case <-ctx.Done():
			collector.logger.Debug("Stop sensor data collection: ", ctx.Err())
			return
		}
	}
}

func (collector *SensorDataCollector) readDevciceData(device SensorDevice) {

	defer collector.deviceDataCollectd()

	err := device.Connect()
	if err != nil {
		collector.errorStack.Append(err)
		return
	}
	defer device.Disconnect()

	collector.logger.Debug("Read data from device: ", device.Id())
	for _, characteristic := range collector.characteristics {

		collector.logger.Debug("Read characteristic ", characteristic.uuid)
		if measurementValue, err := collector.readDevciceCharacteristic(device, characteristic); err == nil {
			collector.publish(device.Id(), characteristic.measurementType, measurementValue)
		} else {
			collector.errorStack.Append(err)
		}
	}
}

func (collector *SensorDataCollector) readDevciceCharacteristic(device SensorDevice, characteristic Characteristic) (string, error) {

	for attemps := 0; attemps < collector.retryCount; attemps++ {
		if val, err := device.ReadValue(characteristic.uuid); err == nil {
			if characteristic.measurementType == events.MeasurementType_BATTERY {
				if len(val) < 2 {
					val = append(val, byte(0))
				}
				return fmt.Sprintf("%d", int64(binary.LittleEndian.Uint16(val))), nil

			} else {
				i := int64(binary.LittleEndian.Uint16(val))
				return fmt.Sprintf("%.1f", float64(i)*0.01), nil
			}
		} else {
			collector.logger.Errorf("Unable to fetch characteristic %s from device %s, reason: %s",
				characteristic.uuid, device.Id(), err)
		}
	}
	return "", fmt.Errorf("Unable to fetch characteristic %s from device %s",
		characteristic.uuid, device.Id())
}

func (collector *SensorDataCollector) publish(deviceId string, measurementType events.MeasurementType, value string) {
	measurement := IndoorClimateMeasurement{
		DeviceId:  deviceId,
		Timestamp: time.Now(),
		Type:      measurementType,
		Value:     value,
	}
	for _, publisher := range collector.publisher {
		if err := publisher.SendMeasurement(measurement); err != nil {
			collector.logger.Error(err)
		}
	}
}

func (collector *SensorDataCollector) deviceDataCollectd() {
	collector.done <- struct{}{}
}
