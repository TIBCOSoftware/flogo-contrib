# WITS0 trigger
This trigger provides your folog application the ability to connect to a local serial port and read WITS0 data. The data can be sent out as raw packets or in parsed JSON format

## Installation
```bash
flogo install github.com/swarwick/Flogo/trigger/wits0
```
## Schema
Settings, Outputs and Endpoint:

```json
  "output": [
    {
      "name": "data",
      "type": "string"
    }
  ],
  "handler": {
    "settings": [
      {
        "name": "SerialPort",
        "type": "string"
      },
      {
        "name": "BaudRate",
        "type": "number"
      },
      {
        "name": "DataBits",
        "type": "number"
      },
      {
        "name": "StopBits",
        "type": "number"
      },
      {
        "name": "Parity",
        "type": "number"
      },
      {
        "name": "ReadTimeoutSeconds",
        "type": "number"
      },    
      {
        "name": "HeartBeatInterval",
        "type": "number"
      },
      {
        "name": "HeartBeatValue",
        "type": "string"
      },
      {
        "name": "PacketHeader",
        "type": "string"
      },
      {
        "name": "PacketFooter",
        "type": "string"
      },
      {
        "name": "LineSeparator",
        "type": "string"
      },
      {
        "name": "OutputRaw",
        "type": "boolean"
      }
    ]
  }
```

Triggers are configured via the triggers.json of your application. The following is an example configuration of the WITS0 trigger.

### Start a flow

```json
{
    "id": "wits0",
    "settings": {
    },
    "handlers": [{
        "action": {
            "ref": "github.com/swarwick/flogo/action/flow",
            "data": {
              "flowURI": "res://flow:query"
            }
          },
        "settings": {
            "SerialPort": "/dev/ttyUSB0",
            "BaudRate": 9600,
            "DataBits": 8,
            "StopBits": 1,
            "Parity": 0,
            "ReadTimeoutSeconds": 1,
            "HeartBeatInterval": 30,
            "HeartBeatValue": "&&\n0111-9999\n!!",
            "PacketHeader": "&&",
            "PacketFooter": "!!",
            "LineSeparator": "\r\n",
            "OutputRaw": false
        }
    }]
}
```