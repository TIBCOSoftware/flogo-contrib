# tibco-kafkapub
This activity provides your flogo application the ability to send a Kafka message


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/kafkapub
```

## Schema
Inputs and Outputs:

```json
{
 "inputs":[
    {
      "name": "BrokerUrls",
      "type": "string"
    },
    {
      "name": "Topic",
      "type": "string"
    },
    {
      "name": "Message",
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
      "name": "truststore",
      "type": "string"
    }
  ],
  "outputs": [
    {
      "name": "partition",
      "type": "int"
    },
    {
      "name": "offset",
      "type": "long"
    }
  ]
}
```
## Settings
| Setting     | Description    | Cardinality |
|:------------|:---------------|--------------|
| BrokerUrls | The Kafka cluster to connect to |Required|         
| Token  | The Kafka topic on which to place the message  |Required|
| Message       | The text message to send |Required|
| user  | If connectiong to a SASL enabled port, the userid to use for authentication | Optional|
| password  | If connectiong to a SASL enabled port, the password to use for authentication | Optional|
| truststore  | If connectiong to a TLS secured port, the directory containing the certificates representing the trust chain for the connection.  This is usually just the CACert used to sign the server's certificate | Optional|

## Outputs
| Value     | Description    |
|:------------|:---------------|
| partition | Documents the partition that the message was placed on |
| offset | Documents the offset for the message |

## Configuration Examples
### Simple
Configure a task to send the time of day to the 'syslog' topic
//TODO
```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-kafkapub",
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

### Advanced
Configure a task in flow to send 'my text message' to a number from a REST trigger's query parameter:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-kafkapub",
  "name": "Send Text Message",
  "attributes": [
    { "name": "accountSID", "value": "A...9" },
    { "name": "authToken", "value": "A...9" },
    { "name": "from", "value": "+12016901385" },
    { "name": "message", "value": "my text message" }
  ],
  "inputMappings": [
    { "type": 1, "value": "[T.queryParams].From", "mapTo": "to" }
  ]
}
```
