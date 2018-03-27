---
title: Reply (Legacy)
weight: 4617
---
# Reply (Legacy)
This activity allows you to reply to a trigger invocation and map output values. While the activity is still here, you probably want to use the `actreply`

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/reply
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "code",
      "type": "integer",
      "required": true
    },
    {
      "name": "data",
      "type": "any"
    }
  ],
  "output": [
  ]
}
```
## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| code        | True     | The response code to send back to the trigger |         
| data        | False    | The response data to send back to the trigger |

## Examples
The below example responds with an HTTP success code.

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-reply",
  "name": "Respond OK",
  "attributes": [
    { "name": "code", "value": 200 }
  ]
}
```
