
#include "WiFiConnect.h"
#include "WiFi.h"
#include "Arduino.h"

WiFiConnect::WiFiConnect(char* ssid, char* password, int maxConnectRetries) {
  m_connect_retries = maxConnectRetries;
  m_ssid = ssid;
  m_password = password;
}


bool WiFiConnect::connect() {

  if (WiFi.status() == WL_CONNECTED) {
    return true;
  }
  
  WiFi.mode(WIFI_STA);
  WiFi.begin(m_ssid, m_password);

  Serial.print("Try to connect to WiFi ");
  Serial.print(m_ssid);
  int retries = 0;
  while (WiFi.status() != WL_CONNECTED && retries < m_connect_retries){
    delay(500);
    Serial.print(".");
    retries++;
  }

  if (WiFi.status() == WL_CONNECTED) {
    Serial.println("Success");
    return true;
  } else {
    Serial.println("Failed");
    return false;
  }
}

bool WiFiConnect::connected() {
  return WiFi.status() == WL_CONNECTED;
}

void WiFiConnect::disconnect() {
  WiFi.disconnect(true);
  WiFi.mode(WIFI_OFF);
  Serial.print("WiFi disconnected!");
}

String WiFiConnect::getMacAddress() {
  return WiFi.macAddress();
}
