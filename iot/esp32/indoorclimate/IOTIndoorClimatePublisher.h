
#include <WiFiClientSecure.h>
#include "iot_config.h"
//#include "MQTTClient.h"

class IOTIndoorClimatePublisher {
public:

  IOTIndoorClimatePublisher(WiFiClientSecure secureClient, IOTConfig iot_config);
  
  bool connect();
  void disconnect();
  
  void publishBatteryLevel(const char* address, std::string value, unsigned long timestamp);
  void publishTemperature(const char* address, std::string value, unsigned long timestamp);
  void publishHumidity(const char* address, std::string value, unsigned long timestamp);
  
private:

  IOTConfig m_iot_config;
  
  const char* m_battery_charc     = "battery";
  const char* m_temperature_charc = "temperature";
  const char* m_humidity_charc    = "humidity";
  
  void publish(const char* address, std::string value, const char* characteristic, unsigned long timestamp);

  //MQTTClient m_iotClient = MQTTClient(2048);
  WiFiClientSecure m_net;
  int m_connect_retries;
};
