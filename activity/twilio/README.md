---
title: Twilio
weight: 4620
---
# Twilio
This activity allows you to send a SMS via Twilio.

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/twilio
```

## Schema
Inputs and Outputs:
```json
{
  "input":[
    {
      "name": "accountSID",
      "type": "string"
    },
    {
      "name": "authToken",
      "type": "string"
    },
    {
      "name": "from",
      "type": "string"
    },
    {
      "name": "to",
      "type": "string"
    },
    {
      "name": "message",
      "type": "string"
    }
  ],
  "output": []
}
```
## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| accountSID  | False    | The Twilio account SID |         
| authToken   | False    | The Twilio auth token  |
| from        | False    | The Twilio number you are sending the SMS from |
| to          | False    | The number you are sending the SMS to. This field should be in the format '+15555555555' |
| message     | False    | The SMS message |

## Examples
The below example sends 'my text message' to '617-555-5555' via Twilio:
```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-twilio",
  "name": "Send Text Message",
  "attributes": [
    { "name": "accountSID", "value": "A...9" },
    { "name": "authToken", "value": "A...9" },
    { "name": "from", "value": "+12016901385" },
    { "name": "to", "value": "+16175555555" },
    { "name": "message", "value": "my text message" }
  ]
}
```