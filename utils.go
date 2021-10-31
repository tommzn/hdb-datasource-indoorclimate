package indoorclimate

import (
	"regexp"
	"strings"
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
