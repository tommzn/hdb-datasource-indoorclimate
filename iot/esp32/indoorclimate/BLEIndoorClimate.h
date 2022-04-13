
#include "BLEUUID.h"
#include "BLEAddress.h"
#include "BLEClient.h"

/**
 * BLEIndoorClimate collects indoor climate data from BLE sensor devices.
 * 
 * Used services and characteristics
 *  - Environment Service 0x181A
 *    - Temperature 0x2a6e
 *    - Humidity 0x2a6f
 *  - Battery Service 0x180F
 *    - Battery Level 0x2A19
 */
class BLEIndoorClimate {
public:

  bool connect(BLEAddress address);
  void disconnect();
  
  std::string getBatteryLevel();
  std::string getTemperature();
  std::string getHumidity();

private:
  BLEClient* m_bleClient;
  
  BLEUUID m_batteryService      = BLEUUID((uint16_t)0x180F);
  BLEUUID m_environmentService  = BLEUUID((uint16_t)0x181A);

  BLEUUID m_batteryLevelCharacteristic  = BLEUUID((uint16_t)0x2A19);
  BLEUUID m_temperatureCharacteristic   = BLEUUID((uint16_t)0x2a6e);
  BLEUUID m_humidityCharacteristic      = BLEUUID((uint16_t)0x2a6f);

  std::string getCharacteristic(BLEUUID serviceBLEUUID, BLEUUID characteristicBLEUUID);

};
