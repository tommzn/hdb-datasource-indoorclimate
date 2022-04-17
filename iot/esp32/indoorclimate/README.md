# Inddor Climate Data Collector
Sketch for a M5Stack ESP32 to collect indoor climate data from BLE sensor devices and publish all this measurements to a MQTT topic on AWS IOT.
## Hardware
Requires a 5Stack Core2.
## Setup
### WiFi 
Add your SSID and password at wifi_credentials.h
### AWS IOT
Add a thing name, your AWS IOT endpoint and a topic to settings.h and all certificates required for AWS IOT to certs.h