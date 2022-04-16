# AWS IoT Media Device State Publisher

When video or camera is on,
this program will publish event on json data which is inputed by config yaml to AWS.

On the other hand, when both video and camera are off,
this program will publish event off json data to AWS.


## Config example

Please set config like this.

```
QoS: 1
topic: $aws/things/xxxxx/shadow/update
endpoint: xxxxxxxxxxx.iot.ap-northeast-1.amazonaws.com
port: 443
root_ca: ./AmazonRootCA1.pem
private_key: ./private.pem.key
certificate: ./certificate.pem.crt
event:
  on:  '{"state": {"desired": {"device_state": true}}}'
  off: '{"state": {"desired": {"device_state": false}}}'
```
