
#include "timer.h"

Timer::Timer(NTPClient* ntp, uint16_t exec_timer_duration, uint16_t display_timer_duration) {
  _ntp = ntp;
  _timer_duration = exec_timer_duration;
  _display_timer_duration = display_timer_duration;
}

void Timer::initExecTimer() {
  _timer = _ntp->getEpochTime() + _timer_duration;

}

void Timer::initDisplayTimer() {
  _display_timer = _ntp->getEpochTime() + _display_timer_duration;

}

bool Timer::isExecTimerExpired() {
  return _ntp->getEpochTime() > _timer;
}
    
bool Timer::isDisplayTimerActive() {
  return _display_timer > 0;
}

bool Timer::isDisplayTimerExpired() {
  return _ntp->getEpochTime() > _display_timer;
}

void Timer::disableDisplayTimer() {
  _display_timer = -1;
}
