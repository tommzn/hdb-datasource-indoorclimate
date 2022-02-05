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

	retryCount := conf.GetAsInt("indoorclimate.retry", config.AsIntPtr(3))
	adapterId := conf.Get("indoorclimate.adapter", config.AsStringPtr("hci0"))
	deviceIds := deviceIdsFromConfig(conf)
	devices := []SensorDevice{}
	for _, deviceId := range deviceIds {
		devices = append(devices, NewIndoorClimateSensor(*adapterId, deviceId))
	}
	characteristics := characteristicsFromConfig(conf)
	return &SensorDataCollector{
		logger:          logger,
		devices:         devices,
		characteristics: characteristics,
		publisher:       []Publisher{newLogPublisher(logger)},
		retryCount:      *retryCount,
	}
}

// Run will start collecting sensor data from all defined devices.
func (collector *SensorDataCollector) Run(ctx context.Context) error {

	defer collector.logger.Flush()

	collector.errorStack = utils.NewErrorStack()
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
	return collector.errorStack.AsError()
}

func (collector *SensorDataCollector) readDevciceData(device SensorDevice) {

	defer collector.deviceDataCollectd()

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
		publisher.Send(measurement)
	}
}

func (collector *SensorDataCollector) deviceDataCollectd() {
	collector.done <- struct{}{}
}
