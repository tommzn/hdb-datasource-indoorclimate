/**
 * Inddor climate data collector 
 * 
 *  Collects indoor climate date (temperature, humidity and sensor battery level) from sensor devices.
 *  All indoor climate measurements are published to a MQTT topic on AWS IOT.
 *  
 *  Requires a board with a RTC to get measurement timestamp, e.g. M5Stack Core 2
 *  
 *  @Author tommzn <tommzn@gmx.de>
 *  
 */

// Contains WiFi SSID and passwird
#include "wifi_credentials.h"

// Contains AWS IOT devices certificates
#include "certs.h"

// Contains sleep settings, connection settings and a list of BLE devices
#include "settings.h"

// Config for AW IOT connections
#include "iot_config.h"

#include "BLEDevice.h"
#include "BLEIndoorClimate.h"
#include "IOTIndoorClimatePublisher.h"
#include "WiFiConnect.h"
#include "WiFiClientSecure.h"
#include "timeclient.h"

// WiFi connection handler, handles connect and disonnect for WiFi networks
static WiFiConnect wifi = WiFiConnect(WIFI_SSID, WIFI_PASSWORD, MAX_RECONNECT_ATTEMPS);

// Secure connection client, used to connect to AWS IOT
WiFiClientSecure secureClient = WiFiClientSecure();  

// Indoor climate data collector, uses Bluetooth connect/scan to get indoor climate data
static BLEIndoorClimate indoorClimateCollector = BLEIndoorClimate();

// List of BLE sensir devices
static BLEAddress deviceAddresses[] = {BLEAddress("A4:C1:38:0A:26:41")};

// Config for AWS IOT
IOTConfig iotConfig = IOTConfig(AWS_IOT_THING_NAME, AWS_IOT_ENDPOINT, AWS_IOT_TOPIC, MAX_RECONNECT_ATTEMPS);

// AWS IOT publisher for collected indoor climate data
static IOTIndoorClimatePublisher indoorClimatePublisher = IOTIndoorClimatePublisher(secureClient, iotConfig);

TimeClient timeclient;

/**
 *  CollectIndoorClimate collects and published indoor climate data.
 *  
 *    - Ensure WiFi connection
 *    - Create AWs IOT indoor climate data publisher
 *    - Loop all BLE sensors to get indoor climate data and publish collected measurements to AWS IOT (MQTT Topic)
 *    - Disconnect from AWS IOT and WiFi
 */
void collectIndoorClimate() {
  
  if (!wifi.connect()) {   
    wifi.disconnect(); 
    return;
  }

  if (!indoorClimatePublisher.connect()) {   
    wifi.disconnect(); 
    return;
  }
  
  for (BLEAddress deviceAddress : deviceAddresses) {

    if (indoorClimateCollector.connect(deviceAddress)) {

      uint32_t timestamp = timeclient.unixtime();
      indoorClimatePublisher.publishBatteryLevel(deviceAddress.toString().data(), indoorClimateCollector.getBatteryLevel(), timestamp);
      indoorClimatePublisher.publishTemperature(deviceAddress.toString().data(), indoorClimateCollector.getTemperature(), timestamp);
      indoorClimatePublisher.publishHumidity(deviceAddress.toString().data(), indoorClimateCollector.getHumidity(), timestamp);  

      indoorClimateCollector.disconnect();  
    }
  }

  indoorClimatePublisher.disconnect();
  wifi.disconnect(); 
}

/**
 *  Setup
 *    - Init Serial monitor and BLE 
 *    - Assign AWS IOT certificates to secure client
 *    - Set deep sleep timer
 */
void setup() {
  
  Serial.begin(115200);

  // Init BLE device
  BLEDevice::init("");

  // Configure WiFiClientSecure to use the AWS IoT device credentials
  secureClient.setCACert(AWS_CERT_CA);
  secureClient.setCertificate(AWS_CERT_CRT);
  secureClient.setPrivateKey(AWS_CERT_PRIVATE);

  esp_sleep_enable_timer_wakeup(SECONDS_TO_SLEEP * uS_TO_S_FACTOR);
  
}

/**
 *  At each execution current time will be updated via NTP, indoor climate data 
 *  is collected from all defined devices and is published to a MQTT topic in AWs IOT.
 */
void loop() {

  timeclient.begin();
  
  collectIndoorClimate(); 

  timeclient.end();
  
  Serial.flush(); 
  esp_deep_sleep_start();
}
