
#include "WiFi.h"
#include "WiFiClientSecure.h"

class WiFiConnect {
public:

  WiFiConnect(char* ssid, char* password, int maxConnectRetries);
  
  bool connect();
  bool connected();
  void disconnect();
  
  String getMacAddress();

private:

  int m_connect_retries;
  char* m_ssid;
  char* m_password;
};
