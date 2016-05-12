# tibco-coap
This trigger provides your flogo application the ability to start a flow via CoAP

## Installation

```bash
flogo add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/coap
```

## Schema
Settings, Outputs and Endpoint:

```json
"settings": [
  {
    "name": "port",
    "type": "integer",
    "required": true
  }
],
"outputs": [
  {
    "name": "payload",
    "type": "string"
  }
],
"endpoint": {
  "settings": [
    {
      "name": "method",
      "type": "string",
      "required" : true
    },
    {
      "name": "path",
      "type": "string",
      "required" : true
    },
    {
      "name": "autoIdReply",
      "type": "boolean"
    }
  ]
}
```
## Settings
### Trigger:
| Setting     | Description    |
|:------------|:---------------|
| port | The port to listen on |         
### Endpoint:
| Setting     | Description    |
|:------------|:---------------|
| method      | The CoAP method |         
| path        | The path  |
| autoIdReply | Automatically reply with the ID of the flow instance |

## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the CoAP Trigger.

### POST
Configure the Trigger to handle a CoAP POST message with path /device

```json
{
  "triggers": [
    {
      "name": "tibco-coap",
      "settings": {
        "port": "7777"
      },
      "endpoints": [
        {
          "flowURI": "embedded://coap_flow",
          "settings": {
            "method": "POST",
            "path": "/device"
          }
        }
      ]
    }
  ]
}
```