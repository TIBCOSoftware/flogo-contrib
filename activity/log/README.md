---
title: Log
weight: 4615
---
# Log
This activity allows you to write log messages.

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/log
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "message",
      "type": "string",
      "value": ""
    },
    {
      "name": "flowInfo",
      "type": "boolean",
      "value": "false"
    },
    {
      "name": "addToFlow",
      "type": "boolean",
      "value": "false"
    }
  ],
  "output": [
    {
      "name": "message",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| message     | False    | The message to log |
| flowInfo    | False    | If set to true this will append the flow information to the log message |
| addToFlow   | False    | If set to true this will add the log message to the 'message' output of the activity and make it available in further activities |
| message     | False    | The message that was logged |

## Examples
The below example logs a message 'test message':

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-log",
  "name": "Log Message",
  "attributes": [
    { "name": "message", "value": "test message" }
  ]
}
```