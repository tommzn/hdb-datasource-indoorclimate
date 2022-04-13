
#include "WiFi.h"
#include "WiFiClientSecure.h"

class WiFiConnect {
public:

  WiFiConnect(char* ssid, char* password, int maxConnectRetries);
  
  bool connect();
  void disconnect();
  
private:

  int m_connect_retries;
  char* m_ssid;
  char* m_password;
};
