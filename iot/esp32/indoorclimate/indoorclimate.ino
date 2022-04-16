
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
#include "timer.h"

// Time Client
#include <NTPClient.h>
#include <WiFiUdp.h>

int counter = 0;

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

Timer timer = Timer(&ntp, SECONDS_TO_SLEEP, DISPLAY_TIMEOUT);

void collectIndoorClimate() {
  
  uint8_t device_count = sizeof(deviceAddresses) / sizeof(deviceAddresses[0]);
  uint8_t devices_ok = 0;
  lcd.initBleDeviceCount();
  lcd.updateBleDeviceCount(devices_ok, device_count);
  
  for (BLEAddress deviceAddress : deviceAddresses) {

    lcd.initBleDevice(deviceAddress.toString().data());
    lcd.updateBleDeviceStatus("Connecting");
    lcd.initBleCharacteristics();
    
    if (indoorClimateCollector.connect(deviceAddress)) {

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

  iotClient.publish(AWS_IOT_TOPIC, jsonBuffer);
}


bool connectToWifi() {
  lcd.updatetWifiStatus("Connecting");
  bool connected = wifi.connect();
  showWiFiStatus();
  return connected;
}

void connectToAwsIot() {
  
  lcd.updatetAwsIotStatus("Connecting");
  int retries = 0;
  while (!iotClient.connect(AWS_IOT_THING_NAME) && retries < MAX_AWSIOT_CONNECT_ATTEMPS) {
    Serial.print(".");
    retries++;
    delay(500);
  }
  showAwsIotStatus();
}

void showWiFiStatus() {
  if (wifi.connect()) {
    lcd.updatetWifiStatus("Connected");
  } else {
    lcd.updatetWifiStatus("Failed");
  }
}

void showAwsIotStatus() {
  if (iotClient.connected()) {
    lcd.updatetAwsIotStatus("Connected");
  } else {
    lcd.updatetAwsIotStatus("Failed");
  }
}

void setup() {

  M5.begin();
  BLEDevice::init("");
  Serial.begin(115200);
  
  counter = 0;
  Lcd.initLoopCounter();
  Lcd.updateLoopCounter(counter);
  
  Lcd.initWifiStatus();
  lcd.updatetWifiStatus("Disconnected");

  lcd.initAwsIotStatus();
  lcd.updatetAwsIotStatus("Disconnected");

  // Configure WiFiClientSecure to use the AWS IoT device credentials
  secureClient.setCACert(AWS_CERT_CA);
  secureClient.setCertificate(AWS_CERT_CRT);
  secureClient.setPrivateKey(AWS_CERT_PRIVATE);

  if (connectToWifi()) {
    iotClient.begin(AWS_IOT_ENDPOINT, 8883, secureClient);
    iotClient.setKeepAlive(60);
    connectToAwsIot();

    // Update current time for measurment timestamp
    ntp.begin();
    ntp.update();
  }
}

void loop() {

  counter++;
  Lcd.updateLoopCounter(counter);

  // In case any connection get lost restart
  if (!wifi.connected() || !iotClient.connected()) { 
      showWiFiStatus();
      showAwsIotStatus();
      delay(5000);
      M5.shutdown(3);  
  }  
    
  if (timer.isExecTimerExpired()) {

    timer.initExecTimer();

    lcd.wakeup();    
    delay(500);
    
    // Collect indoor climate date from all defined sensors
    collectIndoorClimate();

    timer.initDisplayTimer();
  }

  if (timer.isDisplayTimerActive() && timer.isDisplayTimerExpired()) {
    timer.disableDisplayTimer();
    lcd.sleep();
  }
  
  iotClient.loop();
  Serial.flush(); 
  delay(1000);
  
}
