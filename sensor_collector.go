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
)

// NewSensorDataCollector creates a new collector.
// Please note: No publisher tartget will be added. You've to do it by yourself using AppenTarget method.
func NewSensorDataCollector(conf config.Config, logger log.Logger) core.Collector {

	schedule := conf.GetAsDuration("indoorclimate.schedule", nil)
	retryCount := conf.GetAsInt("indoorclimate.retry", config.AsIntPtr(3))
	adapterId := conf.Get("indoorclimate.adapter", config.AsStringPtr("hci0"))

	deviceIds := deviceIdsFromConfig(conf)
	devices := []SensorDevice{}
	for _, deviceId := range deviceIds {
		devices = append(devices, NewIndoorClimateSensor(*adapterId, deviceId))
	}
	characteristics := characteristicsFromConfig(conf)

	return &SensorDataCollector{
		schedule:        schedule,
		logger:          logger,
		devices:         devices,
		characteristics: characteristics,
		publisher:       []Publisher{},
		retryCount:      *retryCount,
	}
}

// AooendTarget will append passed target to internal publisher list.
func (collector *SensorDataCollector) AppendTarget(newTarget Publisher) {
	collector.publisher = append(collector.publisher, newTarget)
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

// RunSingle run indoor climate data fetch once for all devices.
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

// RunContinouous will run in an endless loop and fetches device data in a defines schedule.
func (collector *SensorDataCollector) RunContinouous(ctx context.Context) {

	collector.logger.Debugf("Run continuous collection with schedule of: %s", *collector.schedule)
	for {
		collector.errorStack = utils.NewErrorStack()
		collector.RunSingle(ctx)
		if err := collector.errorStack.AsError(); err != nil {
			collector.logger.Error(err)
		}
		collector.logger.Flush()
		select {
		case <-time.After(*collector.schedule):
			collector.logger.Debug("Restart Data collection")
		case <-ctx.Done():
			collector.logger.Debug("Stop sensor data collection: ", ctx.Err())
			return
		}
	}
}

// ReadDevciceData reads all defines characteristics from passed device.
func (collector *SensorDataCollector) readDevciceData(device SensorDevice) {

	defer collector.deviceDataCollected()

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

// ReadDevciceCharacteristic will read single characteristic from passed device.
func (collector *SensorDataCollector) readDevciceCharacteristic(device SensorDevice, characteristic Characteristic) (string, error) {

	for attemps := 0; attemps < collector.retryCount; attemps++ {
		if val, err := device.ReadValue(characteristic.uuid); err == nil {
			if characteristic.measurementType == MEASUREMENTTYPE_BATTERY {
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

// Publish will send passed measurement to all available publishers.
func (collector *SensorDataCollector) publish(deviceId string, measurementType MeasurementType, value string) {
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

// DeviceDataCollected will finish indoor climate data collection by writing to internal channel.
func (collector *SensorDataCollector) deviceDataCollected() {
	collector.done <- struct{}{}
}
