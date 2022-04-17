package main

import (
	"context"
	b64 "encoding/base64"
	"encoding/binary"
	"fmt"
	"time"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	indoorclimate "github.com/tommzn/hdb-datasource-indoorclimate"
)

func New(logger log.Logger, conf config.Config) *IotMessageHandler {
	return &IotMessageHandler{
		logger:    logger,
		conf:      conf,
		publisher: []indoorclimate.Publisher{},
	}
}

func (handler *IotMessageHandler) HandleEvent(ctx context.Context, indoorClimateData IndoorClimateData) error {

	defer handler.logger.Flush()

	decodedValue, err := b64.StdEncoding.DecodeString(indoorClimateData.Value)
	if err != nil {
		handler.logger.Error("Unable to decode measurement value, reason: ", err)
		return err
	}

	measurementType := indoorclimate.MeasurementType(indoorClimateData.Characteristic)
	measurementValue := convertMeasurementValue(decodedValue, measurementType)
	timeStamp := time.Unix(indoorClimateData.TimeStamp, 0)
	handler.logger.Infof("Measurement received: Device: %s, Characteristic: %s, TimeStamp: %s, Value: %s",
		indoorClimateData.DeviceId, indoorClimateData.Characteristic, timeStamp.String(), measurementValue)
	handler.publish(indoorClimateData.DeviceId, measurementType, timeStamp, measurementValue)
	return nil
}

// AooendTarget will append passed target to internal publisher list.
func (handler *IotMessageHandler) appendTarget(newTarget indoorclimate.Publisher) {
	handler.publisher = append(handler.publisher, newTarget)
}

// ConvertMeasurementValue will decode passed value from byte array to a formatted string.
// Battery level is formatted as integer, temperature and humidity as float with one decimal place.
func convertMeasurementValue(value []byte, measurementType indoorclimate.MeasurementType) string {

	if measurementType == indoorclimate.MEASUREMENTTYPE_BATTERY {
		if len(value) < 2 {
			value = append(value, byte(0))
		}
		return fmt.Sprintf("%d", int64(binary.LittleEndian.Uint16(value)))

	} else {
		i := int64(binary.LittleEndian.Uint16(value))
		return fmt.Sprintf("%.1f", float64(i)*0.01)
	}

}

// Publish will send passed measurement to all available publishers.
func (handler *IotMessageHandler) publish(deviceId string, measurementType indoorclimate.MeasurementType, timeStamp time.Time, value string) {
	measurement := indoorclimate.IndoorClimateMeasurement{
		DeviceId:  deviceId,
		Timestamp: timeStamp,
		Type:      measurementType,
		Value:     value,
	}
	for _, publisher := range handler.publisher {
		if err := publisher.SendMeasurement(measurement); err != nil {
			handler.logger.Error(err)
		}
	}
}
