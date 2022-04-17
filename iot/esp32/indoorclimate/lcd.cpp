
#include <M5Core2.h>
#include "lcd.h"

Lcd::Lcd() {
    M5.Lcd.clear();
    M5.Lcd.fillScreen(BLACK);
    M5.Lcd.setTextColor(WHITE);
    M5.Lcd.setTextSize(2);
}

void Lcd::sleep() {
  M5.Lcd.sleep();
}

void Lcd::wakeup() {
  M5.Lcd.wakeup();
}

void Lcd::initWifiStatus() {
    initLine(line_number_wifi, "WiFi");
}

void Lcd::updatetWifiStatus(const char* status) {
    writeValue(line_number_wifi, status); 
}

void Lcd::initAwsIotStatus() {
    initLine(line_number_awsiot, "AWS IOT");
}

void Lcd::updatetAwsIotStatus(const char* status) {
    writeValue(line_number_awsiot, status); 
}

void Lcd::initBleDeviceCount() {
    initLine(line_number_blecnt, "Devices");
}

void Lcd::updateBleDeviceCount(uint8_t device_count, uint8_t device_max) {
    char buf[10];
    sprintf(buf, "%d/%d", device_count, device_max);
    writeValue(line_number_blecnt, buf);
}

void Lcd::initBleDevice(const char* ble_address) {
    initLine(line_number_bledevice, ble_address);
    initLine(line_number_blestatus, "Status");
    
}

void Lcd::updateBleDeviceStatus(const char* status) {
    writeValue(line_number_blestatus, status);
}

void Lcd::initBleCharacteristics() {
    initLine(line_number_charc_temp, "Temp");
    initLine(line_number_charc_hum, "Hum");
    initLine(line_number_charc_bat, "Battery");
}

void Lcd::updateTemperatureStatus(const char* status) {
    writeValue(line_number_charc_temp, status);
}


void Lcd::updateHumidityStatus(const char* status) {
    writeValue(line_number_charc_hum, status);
}

void Lcd::updateBatteryStatus(const char* status) {
    writeValue(line_number_charc_bat, status);
}

void Lcd::initLine(uint8_t line_number, const char* title) {
    M5.Lcd.fillRect(0, top_margin + line_height * line_number, M5.Lcd.width(), line_height, BLACK);
    writeLineTitle(line_number, title);
    writeSeparator(line_number);
}

void Lcd::writeLineTitle(uint8_t line_number, const char* title) {
    M5.Lcd.setCursor(left_margin, top_margin + line_height * line_number);
    M5.Lcd.print(title);  
}

void Lcd::writeSeparator(uint8_t line_number) {
    M5.Lcd.setCursor(separator_pos_x, top_margin + line_height * line_number);
    M5.Lcd.print(":");
}

void Lcd::writeValue(uint8_t line_number, const char* value) {
    M5.Lcd.fillRect(status_pos_x, top_margin + line_height * line_number, M5.Lcd.width() - status_pos_x - left_margin, line_height, BLACK);
    M5.Lcd.setCursor(status_pos_x, top_margin + line_height * line_number);
    M5.Lcd.print(value);
} 

void Lcd::updateBatteryLevel(uint8_t battery_level) {
    uint8_t width = 50;
    M5.Lcd.fillRect(M5.Lcd.width() - width, top_margin + line_height * line_number_bat_level, width, line_height, BLACK);
    M5.Lcd.setCursor(M5.Lcd.width() - width, top_margin + line_height * line_number_bat_level);
    M5.Lcd.printf("%d%%", battery_level);
}
    
    
