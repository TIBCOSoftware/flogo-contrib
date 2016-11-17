# tibco-device
This trigger provides your flogo application the ability to start a flow from a device 


## Installation

```bash
flogo add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/device
```

## Schema
Settings, Outputs and Endpoint:

```json
  "settings":[
    {
      "name": "mqtt_server",
      "type": "string",
      "required" : true
    },
    {
      "name": "mqtt_user",
      "type": "string"
    },
    {
      "name": "mqtt_password",
      "type": "string"
    },
    {
      "name": "device:name",
      "type": "string",
      "required" : true
    },
    {
      "name": "device:ssid",
      "type": "string",
      "required" : true
    },
    {
      "name": "device:wifi_password",
      "type": "string",
      "required" : true
    },
    {
      "name": "device:board",
      "type": "string",
      "required" : true,
      "allowed" : ["feather_m0_wifi"]
    }
  ],
  "endpoint": {
    "settings": [
      {
        "name": "device:pin",
        "type": "string",
        "required" : true
      },
      {
        "name": "device:condition",
        "type": "string",
        "required" : true
      },
      {
        "name": "device:response_pin",
        "type": "string"
      }
    ]
  }
```

## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the Device Trigger.

### Start a flow
Configure the Trigger to start "test".

```json
{
  "triggers": [
    {
      "name": "tibco-device",
      "type": "device",
      "settings": {
        "mqtt_server":"192.168.1.50",
        "mqtt_user":"",
        "mqtt_pass":"",
        "device:name":"myarduino",
        "device:type":"feather_m0_wifi",
        "device:ssid":"mynetwork",
        "device:wifi_password": "mypass"
      },
      "endpoints": [
        {
          "actionType": "flow",
          "actionURI": "embedded://test",
          "settings": {
            "device:pin": "D:A3",
            "device:condition": "== HIGH",
            "device:response_pin": "D:A4"
          }
        }
      ]
    }
  ]
}

```
