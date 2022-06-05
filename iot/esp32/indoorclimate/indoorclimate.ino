/**
 *  Indoor Climate Data Collector
 *  
 *    Collects temperature, humidity and battery level from sensor devices, 
 *    e.g Xiaomi Mi Temperature & Humidity Sensor 2 and publishes this data to 
 *    a MQTT topic on AWS IOT.
 *    It scans sensor data in a defined schedule and uses deep sleep for idle phases.
 *
 *    Author: tommzn <tommzn@gmx.de>
 */
#include <M5Core2.h>

// Contains AWS IOT devices certificates
#include "certs.h"

// Contains WiFi SSID and passwird
#include "wifi_credentials.h"

// Contains sleep settings and AWS IOT connection settings
#include "settings.h"

#include "lcd.h"
#include "WiFiConnect.h"
#include "WiFiClientSecure.h"
#include "BLEDevice.h"
#include "BLEIndoorClimate.h"
#include "MQTTClient.h"
#include "ArduinoJson.h"
#include "base64.hpp"

// Time Client
#include <NTPClient.h>
#include <WiFiUdp.h>

// LCD handler for status updates
static Lcd lcd;

// Secure connection client, used to connect to AWS IOT
WiFiClientSecure secureClient = WiFiClientSecure();  

// Indoor climate data collector, uses Bluetooth connect/scan to get indoor climate data
static BLEIndoorClimate indoorClimateCollector = BLEIndoorClimate();

// List of BLE sensir devices
static BLEAddress deviceAddresses[] = {BLEAddress("A4:C1:38:0A:26:41")};

// WiFi connection handler, handles connect and disonnect for WiFi networks
static WiFiConnect wifi = WiFiConnect(WIFI_SSID, WIFI_PASSWORD, MAX_WIFI_CONNECT_ATTEMPS);

// Indoor climate data publisher
MQTTClient iotClient = MQTTClient(2048);
  
// NTP Setup
WiFiUDP ntpUDP;
NTPClient ntp(ntpUDP, "europe.pool.ntp.org", 0);

/**
 *  Loops all defined sensor devices, connects to, collect all measurements
 *  and publishes this data as JSON to a topic on AWS IOT.
 *  In some cases Xiaomi Mi sensors doesn't advertise it's data, so a direct connection is required.
 */
void collectIndoorClimate() {
  
  uint8_t device_count = sizeof(deviceAddresses) / sizeof(deviceAddresses[0]);
  uint8_t devices_ok = 0;
  lcd.initBleDeviceCount();
  lcd.updateBleDeviceCount(devices_ok, device_count);

  // Publish own battery level
  const char batLevel = (char) uint8_t(M5.Axp.GetBatteryLevel());
  publishMeasurement(wifi.getMacAddress().c_str(), &batLevel, "battery", ntp.getEpochTime());  

  
  for (BLEAddress deviceAddress : deviceAddresses) {

    lcd.initBleDevice(deviceAddress.toString().data());
    lcd.updateBleDeviceStatus("Connecting");
    lcd.initBleCharacteristics();

    Serial.println("Connecting...");
    if (indoorClimateCollector.connect(deviceAddress)) {

      Serial.println("Connected!");
      lcd.updateBleDeviceStatus("Connected");

      unsigned long timestamp = ntp.getEpochTime();
      lcd.updateBatteryStatus("Fetching");
      const char* battery_level = indoorClimateCollector.getBatteryLevel().c_str();
      lcd.updateBatteryStatus("OK");
      publishMeasurement(deviceAddress.toString().data(), battery_level, "battery", timestamp);  
      lcd.updateBatteryStatus("Published");
      
      lcd.updateTemperatureStatus("Fetching");
      const char* temperature = indoorClimateCollector.getTemperature().c_str();
      lcd.updateTemperatureStatus("OK");
      publishMeasurement(deviceAddress.toString().data(), temperature, "temperature", timestamp);  
      lcd.updateTemperatureStatus("Published");
      
      lcd.updateHumidityStatus("Fetching");
      const char* humidity = indoorClimateCollector.getHumidity().c_str();
      lcd.updateHumidityStatus("OK");
      publishMeasurement(deviceAddress.toString().data(), humidity, "humidity", timestamp);  
      lcd.updateHumidityStatus("Published");
      
      indoorClimateCollector.disconnect();  
      lcd.updateBleDeviceStatus("Disonnected");
      devices_ok++;
      lcd.updateBleDeviceCount(devices_ok, device_count);

    } else {
      lcd.updateBleDeviceStatus("Failed");
      lcd.updateTemperatureStatus("Failed");
      lcd.updateBatteryStatus("Failed");
      lcd.updateHumidityStatus("Failed");
    }
  }
}

/**
 *  Convert passed measurment to a JSON object, measurement values will be base64 encoded, and publish
 *  this data to a MQTT topic on AWS IOT.
 */
void publishMeasurement(const char* address, std::string value, const char* characteristic, unsigned long timestamp) {

  unsigned char base64[10];
  unsigned int base64_length = encode_base64((unsigned char *) value.c_str(), strlen(value.c_str()), base64);
  
  StaticJsonDocument<200> doc;
  doc["device_id"]      = address;
  doc["characteristic"] = characteristic;
  doc["value"]          = base64;
  doc["timestamp"]      = timestamp;
  
  char jsonBuffer[512];
  serializeJson(doc, jsonBuffer);

  Serial.println(jsonBuffer);
  iotClient.publish(AWS_IOT_TOPIC, jsonBuffer);
}

/**
 *  Connect to local WiFi networks and update connection status on LCD.
 */
bool connectToWifi() {
  lcd.updatetWifiStatus("Connecting");
  if (wifi.connect()) {
    lcd.updatetWifiStatus("Connected");
    return true;
  } else {
    lcd.updatetWifiStatus("Failed");
    return false;
  }
}

/**
 *  Establishes a connection to AWS IOT to publish measurements
 *  and update connection status in LCD.
 */
void connectToAwsIot() {
  
  lcd.updatetAwsIotStatus("Connecting");
  int retries = 0;
  while (!iotClient.connect(AWS_IOT_THING_NAME) && retries < MAX_AWSIOT_CONNECT_ATTEMPS) {
    Serial.print(".");
    retries++;
    delay(500);
  }
  if (iotClient.connected()) {
    lcd.updatetAwsIotStatus("Connected");
  } else {
    lcd.updatetAwsIotStatus("Failed");
  }
}

void setup() {

  Serial.begin(115200);
  
  M5.begin();
  BLEDevice::init("");
    
  // Init WiFi and AWS IOT connection status on LCD.
  Lcd.initWifiStatus();
  lcd.updateBatteryLevel(uint8_t(M5.Axp.GetBatteryLevel()));
  lcd.updatetWifiStatus("Disconnected");
  lcd.initAwsIotStatus();
  lcd.updatetAwsIotStatus("Disconnected");

  // Configure WiFiClientSecure to use the AWS IoT device credentials
  secureClient.setCACert(AWS_CERT_CA);
  secureClient.setCertificate(AWS_CERT_CRT);
  secureClient.setPrivateKey(AWS_CERT_PRIVATE);

  // Connect to WiFi
  if (connectToWifi()) {

    // Connect to AWS IOT
    iotClient.begin(AWS_IOT_ENDPOINT, 8883, secureClient);
    // Extend default timeout because data collection may take some seconds.
    iotClient.setKeepAlive(60);
    connectToAwsIot();

    // Update current time for measurment timestamp
    ntp.begin();
    ntp.update();
  }
}

void loop() {

  // Climate data collection requires active WiFi and AWS IOT connection
  if (wifi.connected() && iotClient.connected()) { 

      // Collect indoor climate date from all defined sensors
      collectIndoorClimate();
  }  
  // run MQTT client to handle send/receive packages
  iotClient.loop();

  // Shutdown NTP client
  ntp.end();  

  // Disconnect from AWs IOT and update connection status on LCD
  iotClient.disconnect();
  lcd.updatetAwsIotStatus("Disconnected");
  
  // Disconnect from WiFi and update connection status on LCD
  wifi.disconnect();
  lcd.updatetWifiStatus("Disconnected");
  
  // Some delay, to provide to opportunity to read all this information in LCD
  delay(DISPLAY_TIMEOUT * mS_TO_S_FACTOR);

  // Going to deep sleep until next interation
  Serial.flush(); 
  M5.Axp.DeepSleep(SECONDS_TO_SLEEP * uS_TO_S_FACTOR);
  
}
