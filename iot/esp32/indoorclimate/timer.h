
#include <NTPClient.h>

class Timer {
public:
    Timer(NTPClient* ntp, uint16_t exec_timer_duration, uint16_t display_timer_duration);

    void initExecTimer();
    void initDisplayTimer();

    bool isExecTimerExpired();
    
    bool isDisplayTimerActive();
    bool isDisplayTimerExpired();
    void disableDisplayTimer();
    
private:
  uint16_t _timer_duration          = 0;
  uint16_t _display_timer_duration  = 0;
  
  unsigned long _timer = 0;
  long _display_timer  = -1;

  NTPClient* _ntp;

};
