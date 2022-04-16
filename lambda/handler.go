package main

import (
	"context"
	b64 "encoding/base64"
	"encoding/binary"
	"fmt"
	"time"

	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
)

func New(logger log.Logger, conf config.Config) MessageHandler {
	return &IotMessageHandler{logger: logger, conf: conf}
}

func (handler *IotMessageHandler) HandleEvent(ctx context.Context, indoorClimateDate IndoorClimateDate) error {

	decodedValue, err := b64.StdEncoding.DecodeString(indoorClimateDate.Value)
	if err != nil {
		handler.logger.Error("Unable to decode measurement value, reason: ", err)
		return err
	}

	measurementValue := convertMeasurementValue(decodedValue, indoorClimateDate.Characteristic)
	timeStamp := time.Unix(indoorClimateDate.TimeStamp, 0)
	handler.logger.Infof("Measurement received: Device: %s, Characteristic: %s, TimeStamp: %s, Value: %s",
		indoorClimateDate.DeviceId, indoorClimateDate.Characteristic, timeStamp.String(), measurementValue)

	return nil
}

func convertMeasurementValue(value []byte, characteristic string) string {

	if characteristic == "battery" {
		if len(value) < 2 {
			value = append(value, byte(0))
		}
		return fmt.Sprintf("%d", int64(binary.LittleEndian.Uint16(value)))

	} else {
		i := int64(binary.LittleEndian.Uint16(value))
		return fmt.Sprintf("%.1f", float64(i)*0.01)
	}

}
