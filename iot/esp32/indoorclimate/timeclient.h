
#include <NTPClient.h>
#include <WiFiUdp.h>

// RTC support for M5Stack
#include "RTClib.h"

class TimeSource {
public:

  TimeSource() {}
  virtual ~TimeSource() {}
  
  // Init time source
  virtual bool begin();

  // Get current unix time, seconds since 1970-01-01
  virtual uint32_t unixtime() = 0;

  // Stop time source, maybe clean up
  virtual void end();
};

class RTC_TimeSource: public TimeSource {
public:

  RTC_TimeSource() {};
  virtual ~RTC_TimeSource() {};
  
  virtual bool begin() override;
  virtual uint32_t unixtime() override;
  virtual void end() override;
  
private:
  
  RTC_DS3231 rtc;
};

class NTP_TimeSource: public TimeSource {
public:

  NTP_TimeSource();
  virtual ~NTP_TimeSource() {};
  
  virtual bool begin() override;
  virtual uint32_t unixtime() override;
  virtual void end() override;
  
private:
  
  WiFiUDP ntpUDP;
  NTPClient ntpClient;
};


class TimeClient {
public:

  void init();
  
  // Init time source
  bool begin();

  // Get current unix time, seconds since 1970-01-01
  uint32_t unixtime();

  // Stop time source, maybe clean up
  void end();

private:

  TimeSource* m_timesource;
};
