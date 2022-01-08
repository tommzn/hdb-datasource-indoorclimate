package indoorclimate

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
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
