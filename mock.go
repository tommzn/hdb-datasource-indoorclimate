package indoorclimate

type mqttMessage struct {
	topic, payload string
	ack            bool
}

func (message *mqttMessage) Duplicate() bool {
	return false
}

func (message *mqttMessage) Qos() byte {
	return 1
}
func (message *mqttMessage) Retained() bool {
	return false
}

func (message *mqttMessage) Topic() string {
	return message.topic
}

func (message *mqttMessage) MessageID() uint16 {
	return 1
}

func (message *mqttMessage) Payload() []byte {
	return []byte(message.payload)
}

func (message *mqttMessage) Ack() {
	message.ack = true
}
