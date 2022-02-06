package indoorclimate

import (
	"time"

	btdevice "github.com/muka/go-bluetooth/bluez/profile/device"
	config "github.com/tommzn/go-config"
	log "github.com/tommzn/go-log"
	metrics "github.com/tommzn/go-metrics"
	secrets "github.com/tommzn/go-secrets"
	utils "github.com/tommzn/go-utils"
	core "github.com/tommzn/hdb-datasource-core"
	events "github.com/tommzn/hdb-events-go"
)

// MqttClient is used to subsribe to a MQTT broker and process messages with indoor climate data.
type MqttClient struct {
	conf            config.Config
	secretsManager  secrets.SecretsManager
	logger          log.Logger
	targets         []MessageTarget
	metricPublisher metrics.Publisher
}

// SqsTarget sends passed indoor climate data to a AWS SQS queue.
type SqsTarget struct {

	// Publisher is a SQS client to publish messages.
	publisher core.Publisher
}

// logTarget writes given indoor climate data to an internal logger
type logTarget struct {
	logger log.Logger
}

// collectorTarget collects passed infoor climate data in local storage.
type collectorTarget struct {
	messages []events.IndoorClimate
}

// messageHandler is used to process messages received from a MQTT topic.
type messageHandler struct {
	logger log.Logger
}

// IndoorClimateSensor is used to fetch tem eprature, humidiy and bettery status
// from a Xiaomi Mijia (LYWSD03MMC) indoor climate sensor.
type IndoorClimateSensor struct {
	adapterId string
	deviceId  string
	device    *btdevice.Device1
}

// IndoorClimateMeasurement is a metric read from a sensor device.
type IndoorClimateMeasurement struct {
	DeviceId  string
	Timestamp time.Time
	Type      events.MeasurementType
	Value     string
}

// LogPublisher will log indoor climate measuremnts.
type LogPublisher struct {
	logger log.Logger
}

// Characteristic is a songle sensor value.
type Characteristic struct {
	uuid            string
	measurementType events.MeasurementType
}

// SensorDataCollector will try to fetch temperature, humidity and bettery status
// from a given list of sensors.
type SensorDataCollector struct {
	logger          log.Logger
	devices         []SensorDevice
	characteristics []Characteristic
	publisher       []Publisher
	retryCount      int
	schedule        *time.Duration
	errorStack      *utils.ErrorStack
	done            chan struct{}
}
