
#include "BLEDevice.h"
#include "BLEIndoorClimate.h"
#include "Arduino.h"

/**
 *  Creates a new BLE client object an try to connect to given device.
 *  If there's an active connection to another device, it will disconnect first.
 */
bool BLEIndoorClimate::connect(BLEAddress address) {
  
  Serial.print("Connecting to device ");
  if (m_bleClient != nullptr && m_bleClient->isConnected()) {
    Serial.println("Already connected, disconnect first!");
    m_bleClient->disconnect();
  }

  m_bleClient = BLEDevice::createClient();
  if (m_bleClient->connect(address)) {
    Serial.println("Success");
    return true;
  } else {
    Serial.println("Failed");
    return false;
  }   
}

/**
 *  Disconnects from device if there's an active connection.
 */
void BLEIndoorClimate::disconnect() {
  if (m_bleClient != nullptr && m_bleClient->isConnected()) {
    m_bleClient->disconnect();
  }
  Serial.println("BLE device disconnected");  
}

/**
 *  Try to read battery level characteristics from connected device.
 *  Service:        0x180
 *  Characteristic: 0x2A19
 */
std::string BLEIndoorClimate::getBatteryLevel() {
  return getCharacteristic(m_batteryService, m_batteryLevelCharacteristic);
}

/**
 *  Try to read temperature characteristics from connected device.
 *  Service:        0x181A
 *  Characteristic: 0x2a6e
 */
std::string BLEIndoorClimate::getTemperature() {
  return getCharacteristic(m_environmentService, m_temperatureCharacteristic);
}

/**
 *  Try to read humidity characteristics from connected device.
 *  Service:        0x181A
 *  Characteristic: 0x2a6f
 */
std::string BLEIndoorClimate::getHumidity() {
  return getCharacteristic(m_environmentService, m_humidityCharacteristic);
}

/**
 *  Try to read give characteristics from passed service from connected device.
 *  Will return with an empty string if something went wrong.
 */
std::string BLEIndoorClimate::getCharacteristic(BLEUUID serviceUUID, BLEUUID characteristicUUID) {

  if (m_bleClient == nullptr || !m_bleClient->isConnected()) {
    return "";
  }

  BLERemoteService* remoteService = m_bleClient->getService(serviceUUID);
  if (remoteService == nullptr) {
      return "";
  }
  
  BLERemoteCharacteristic* remoteCharacteristic = remoteService->getCharacteristic(characteristicUUID);
  if (remoteCharacteristic == nullptr) {
    return "";
  }
  
  return (remoteCharacteristic->canRead()) ? remoteCharacteristic->readValue() : "";

}
