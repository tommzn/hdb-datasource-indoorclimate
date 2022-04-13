
#include <NTPClient.h>
#include "timeclient.h"
  
bool RTC_TimeSource::begin() {

  // In case RTC is initial or lost power
  if (rtc.lostPower()) {

    // Try to get current time from NTP
    NTP_TimeSource ntp = NTP_TimeSource();
    if (ntp.begin()) {
      rtc.adjust(DateTime(ntp.unixtime()));
      ntp.end();
    } else {
      // Fallback to compile time
      rtc.adjust(DateTime(F(__DATE__), F(__TIME__)));
    }
  }
  return true;
}

uint32_t RTC_TimeSource::unixtime() {
  return rtc.now().unixtime();
}

void RTC_TimeSource::end() {
  // Nothing to do here
}

bool NTP_TimeSource::begin() {
  ntpClient.begin();
  ntpClient.update();
  return true;
}

NTP_TimeSource::NTP_TimeSource() {
  ntpClient = NTPClient(ntpUDP, "europe.pool.ntp.org", 0);
}

uint32_t NTP_TimeSource::unixtime() {
  return ntpClient.getEpochTime();
}

void NTP_TimeSource::end() {
  ntpClient.end();
}

void TimeClient::init() {
  
  RTC_TimeSource rtc = RTC_TimeSource();
  if (rtc.begin()) {
    m_timesource = &rtc;
  } else {
    NTP_TimeSource ntp = NTP_TimeSource();
    m_timesource = &ntp;
  }
}

bool TimeClient::begin() {
  return m_timesource->begin();
}

uint32_t TimeClient::unixtime() {
  return m_timesource->unixtime();
}

void TimeClient::end() {
  m_timesource->end();
}
