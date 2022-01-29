package indoorclimate

import (
	"math/rand"
	"regexp"
	"strings"
	"time"

	config "github.com/tommzn/go-config"
)

// extractDeviceId try to extract a device id, a mac address, from given topic.
func extractDeviceId(topic string) *string {

	macAddressExp := regexp.MustCompile("(?:[0-9A-Fa-f]{2}[:]){5}(?:[0-9A-Fa-f]{2})")
	match := macAddressExp.FindStringSubmatch(topic)
	if len(match) == 1 {
		firstMatch := match[0]
		return &firstMatch
	}
	return nil
}

// extractMeasurementType returns topic suffix which is used as measurement type.
func extractMeasurementType(topic string) *string {

	if !strings.Contains(topic, "/") || strings.HasSuffix(topic, "/") {
		return nil
	}
	suffix := topic[strings.LastIndex(topic, "/")+1:]
	return &suffix
}

// topicsToSubsrcibe generates a list of topic to listen for.
func topicsToSubsrcibe(prefix *string) []string {
	topics := []string{"/ble/+/+/temperature", "/ble/+/+/humidity", "/ble/+/+/battery"}
	if prefix != nil {
		*prefix = strings.TrimSuffix(*prefix, "/")
		for idx, topic := range topics {
			topics[idx] = *prefix + topic
		}
	}
	return topics
}

const letterBytes = "abcdefghijklmnopqrstuvwxyz1234567890"

// randStringBytes returns random bytes of given length from letterBytes.
func randStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// devicesFromConfig will load devices and their characteristics from given config.
func devicesFromConfig(conf config.Config) []Device {

	devices := []Device{}
	characteristics := []Characteristic{}

	characteristicsConfig := conf.GetAsSliceOfMaps("indoorclimate.characteristics")
	for _, characteristicConfig := range characteristicsConfig {
		if uuid, ok := characteristicConfig["uuid"]; ok {
			if measurementType, ok := characteristicConfig["type"]; ok {
				characteristics = append(characteristics,
					Characteristic{
						uuid:            uuid,
						measurementType: measurementType,
					})
			}
		}
	}

	devicesConfig := conf.GetAsSliceOfMaps("indoorclimate.devices")
	for _, deviceConfig := range devicesConfig {
		if deviceId, ok := deviceConfig["id"]; ok {
			device := Device{
				deviceId:        deviceId,
				characteristics: characteristics,
			}
			devices = append(devices, device)
		}
	}
	return devices
}
