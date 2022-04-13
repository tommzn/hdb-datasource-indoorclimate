
#include "iot_config.h"

IOTConfig::IOTConfig() {
  m_max_onnect_retries = 3;
  m_topic = "topic";
  m_aws_iot_endpoint = "";
  m_thing_name = "thing_name";
}

IOTConfig::IOTConfig(const char* thing_name, const char* aws_iot_endpoint, const char* topic, int max_onnect_retries) {
  m_max_onnect_retries = max_onnect_retries;
  m_topic = topic;
  m_aws_iot_endpoint = aws_iot_endpoint;
  m_thing_name = thing_name;
}

int IOTConfig::getMaxCnnectRetries() {
  return m_max_onnect_retries;
}

const char* IOTConfig::getTopic() {
  return m_topic;
}

const char* IOTConfig::getAwsIotEndpoint() {
  return m_aws_iot_endpoint;
}

const char* IOTConfig::getThingName() {
  return m_thing_name;
}
