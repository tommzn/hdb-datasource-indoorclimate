#define uS_TO_S_FACTOR 1000000
#define mS_TO_S_FACTOR 1000

// Interval for BLE sensor scan, current 10m
#define SECONDS_TO_SLEEP 600

// Switch display to sleep mode after 5s
#define DISPLAY_TIMEOUT 5

// Max Wifi connect attemps
#define MAX_WIFI_CONNECT_ATTEMPS 10

// Max AWS IOT connect attemps
#define MAX_AWSIOT_CONNECT_ATTEMPS 10

// AWS IOT settings
const char AWS_IOT_THING_NAME[] = "ThingName";
const char AWS_IOT_ENDPOINT[]   = "xyz-ats.iot.us-west-1.amazonaws.com";
const char AWS_IOT_TOPIC[]      = "iot/topic";

  
