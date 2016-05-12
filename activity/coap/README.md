# tibco-rest
This activity provides your flogo application the ability to send a CoAP message.


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/coap
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "address",
      "type": "string",
      "required": true
    },
    {
      "name": "method",
      "type": "string",
      "required": true
    },
    {
      "name": "type",
      "type": "string"
    },
    {
      "name": "messageId",
      "type": "integer"
    },
    {
      "name": "options",
      "type": "params"
    },
    {
      "name": "payload",
      "type": "string"
    }

  ],
  "outputs": [
    {
      "name": "response",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| address   | The CoAP address to dial |         
| method    | The CoAP method (POST,GET,PUT,DELETE)|
| type      | Message Type (Confirmable, NonConfirmable, Acknowledgement, Reset) |
| messageId | ID used to detect duplicates and for optional reliability |
| options   | CoAP options |
| payload   | The message payload |


## Configuration Examples
### Simple
Configure a task in flow to get pet '1234' from the [swagger petstore](http://petstore.swagger.io):

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-rest",
  "name": "Query for pet 1234",
  "attributes": [
    { "name": "address", "value": "localhost:7777" },
    { "name": "method", "value": "GET" },
    { "name": "type", "value": "Confirmable" },
    { "name": "messageId", "value": 12345 },
    { "name": "payload", "value": "hello world" },
    { "name": "options", "value": {"ETag":"tag", "MaxAge":2, "URIPath":"/mypath" }
  ]
}
```
