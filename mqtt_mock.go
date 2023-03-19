package indoorclimate

type mqttMessage struct {
	payload string
	topic   string
}

func newMqttMessage(topic, payload string) *mqttMessage {
	return &mqttMessage{payload: payload, topic: topic}
}

func (m *mqttMessage) Duplicate() bool {
	return false
}

func (m *mqttMessage) Qos() byte {
	return 1
}

func (m *mqttMessage) Retained() bool {
	return true
}

func (m *mqttMessage) Topic() string {
	return m.topic
}

func (m *mqttMessage) MessageID() uint16 {
	return 1
}

func (m *mqttMessage) Payload() []byte {
	return []byte(m.payload)
}

func (m *mqttMessage) Ack() {

}
