
class Lcd {
public:
    Lcd();

    void sleep();
    void wakeup();
    
    void initWifiStatus();
    void updatetWifiStatus(const char* wifi_status);

    void initAwsIotStatus();
    void updatetAwsIotStatus(const char* wifi_status);
    
    void initBleDeviceCount();
    void updateBleDeviceCount(uint8_t device_count, uint8_t device_max);

    void initBleDevice(const char* ble_address);
    void updateBleDeviceStatus(const char* device_status);

    void initBleCharacteristics();
    void updateTemperatureStatus(const char* status);
    void updateHumidityStatus(const char* status);
    void updateBatteryStatus(const char* status);

    void updateBatteryLevel(uint8_t battery_level);

private:

    void initLine(uint8_t line_number, const char* title);
    void writeLineTitle(uint8_t line_number, const char* title);
    void writeSeparator(uint8_t line_number);
    void writeValue(uint8_t line_number, const char* value);
    
    uint8_t line_height = 20;
    uint8_t left_margin = 10;
    uint8_t top_margin  = 10;

    uint8_t separator_pos_x = 100;
    uint8_t status_pos_x    = 115;

    uint8_t line_number_wifi    = 0;
    uint8_t line_number_awsiot  = 1;
    uint8_t line_number_blecnt  = 2;

    uint8_t line_number_bledevice   = 4;
    uint8_t line_number_blestatus   = 5;
    uint8_t line_number_charc_temp  = 6;
    uint8_t line_number_charc_hum   = 7;
    uint8_t line_number_charc_bat   = 8;
    uint8_t line_number_bat_level   = 9;
    
};
