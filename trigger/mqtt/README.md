# tibco-mqtt
This trigger provides your flogo application the ability to start a flow via MQTT


## Installation

```bash
flogo add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/mqtt
```

## Schema
Settings, Outputs and Endpoint:

```json
{
  "settings":[
      {
        "name": "broker",
        "type": "string"
      },
      {
        "name": "id",
        "type": "string"
      },
      {
        "name": "user",
        "type": "string"
      },
      {
        "name": "password",
        "type": "string"
      },
      {
        "name": "store",
        "type": "string"
      },
      {
        "name": "topic",
        "type": "string"
      },
      {
        "name": "qos",
        "type": "number"
      },
      {
        "name": "cleansess",
        "type": "boolean"
      }
    ],
    "endpoint": {
      "settings": [
        {
          "name": "topic",
          "type": "string"
        }
      ]
    }
}
```

## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the MQTT Trigger.

### Start flow
Configure the Trigger to handle a start flow. "settings" "topic" is the topic it uses to listen for incoming messages. So in this case the "endpoints" "settings" "topic" is "test_start" will start the "flowURI". The incoming message payload has to define "replyTo" which is the the topic used to return the flow ID.

```json
{
  "triggers": [
    {
      "name": "tibco-mqtt",
      "settings": {
        "topic": "flogo/#",
        "broker": "tcp://192.168.1.12:1883",
        "id": "flogoEngine",
        "user": "",
        "password": "",
        "store": "",
        "qos": "0",
        "cleansess": "false"
      },
      "endpoints": [
        {
          "flowURI": "local://testFlow",
          "settings": {
            "topic": "test_start"
          }
        }
      ]
    }
  ]
}
```
