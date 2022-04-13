
#include "Arduino.h"
#include "ArduinoJson.h"
#include "IOTIndoorClimatePublisher.h"

IOTIndoorClimatePublisher::IOTIndoorClimatePublisher(WiFiClientSecure secure_client, IOTConfig iot_config) {
  m_net = secure_client;
  m_iot_config = iot_config;
}

bool IOTIndoorClimatePublisher::connect() {

  return true;
  /**
  if(m_iotClient.connected()) {
    return true;
  }
  
  m_iotClient.begin(m_aws_iot_endpoint, 8883, m_net);

  int retries = 0;
  while (!m_iotClient.connect(m_thing_name) && retries < m_connect_retries) {
    delay(500);
    retries++;
  }

  return m_iotClient.connected();
  */
}

void IOTIndoorClimatePublisher::disconnect() {
  //m_iotClient.disconnect();
}



void IOTIndoorClimatePublisher::publishBatteryLevel(const char* address, std::string value, unsigned long timestamp) {
  publish(address, value, m_battery_charc, timestamp);  
}

void IOTIndoorClimatePublisher::publishTemperature(const char* address, std::string value, unsigned long timestamp) {
  publish(address, value, m_temperature_charc, timestamp);  
}

void IOTIndoorClimatePublisher::publishHumidity(const char* address, std::string value, unsigned long timestamp) {
  publish(address, value, m_humidity_charc, timestamp);  
}

void IOTIndoorClimatePublisher::publish(const char* address, std::string value, const char* characteristic, unsigned long timestamp) {

  StaticJsonDocument<200> doc;
  doc["device_id"]      = address;
  doc["characteristic"] = characteristic;
  doc["value"]          = value;
  doc["timestamp"]      = timestamp;
  
  char jsonBuffer[512];
  serializeJson(doc, jsonBuffer);

  Serial.print(jsonBuffer);
  
  //m_iotClient.publish(m_topic, value.data());
}
