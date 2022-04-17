# Setup IOT device to collect indoor climate date
Step by step guid to setup AWS IOT devices, policy and certificates to enable Arduino/EPS32 scetch at () to publish indoor climate data to MQTT topics.
## Create device and policy
## Create device and policy
Use iot-device.yml as CloudFormation template to create a new IOT device and corresponding policy.
## Generate device certificates
Run following command to generate public/private key for a device.
```
aws iot create-keys-and-certificate \
    --set-as-active \
    --certificate-pem-outfile "IndoorClimateHub01.cert.pem" \
    --public-key-outfile "IndoorClimateHub01.public.key" \
    --private-key-outfile "IndoorClimateHub01.private.key" \
    --region "eu-west-1"
```
## Get Root CA certificate
Visit https://docs.aws.amazon.com/iot/latest/developerguide/server-authentication.html and download the Root CA certificate.

## Attach a thing to the certificate
Run following command to assign the indoor climate data collector thing to the previous generated certificate.
```
aws iot attach-thing-principal \
    --principal <CertificateArnFromCreateCertRequest> \
    --thing-name <YourThingName> \
    --region "eu-west-1"
```
## Attach policy to the certificate
Run following command to attach iot policy created via CloudFormation to the certificate.
```
aws iot attach-policy \
    --target <CertificateArnFromCreateCertRequest> \
    --policy-name <PolicyNameFromCfnStack>
```