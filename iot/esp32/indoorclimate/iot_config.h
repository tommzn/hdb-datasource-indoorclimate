#ifndef __IOT_CONFIG_H__
#define __IOT_CONFIG_H__

class IOTConfig {
public:
  IOTConfig();
  IOTConfig(const char* thing_name, const char* aws_iot_endpoint, const char* topic, int max_onnect_retries);

  int getMaxCnnectRetries();
  const char* getTopic();
  const char* getAwsIotEndpoint();
  const char* getThingName();

private:

  int m_max_onnect_retries;
  const char* m_topic;
  const char* m_aws_iot_endpoint;
  const char* m_thing_name;

};

#endif
