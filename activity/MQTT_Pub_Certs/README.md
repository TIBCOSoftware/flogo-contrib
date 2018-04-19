github.com/Ganitagya/flogo-contrib/activity/MQTT_Pub


github.com/Ganitagya/flogo-contrib/activity/MQTT_noCert


# Publish MQTT Message
This activity provides your flogo application the ability to publish a message on an MQTT topic.


## Installation

```bash
flogo install github.com/Ganitagya/flogo-contrib/activity/MQTT_noCert
```
Link for flogo web:
```
https://github.com/Ganitagya/flogo-contrib/activity/MQTT_noCert
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
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
      "name": "topic",
      "type": "string"
    },
    {
      "name": "qos",
      "type": "integer",
      "allowed" : ["0", "1", "2"]
    },
    {
      "name": "message",
      "type": "string"
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| broker    | the MQTT Broker URI (tcp://[hostname]:[port])|
| id        | The MQTT Client ID |         
| user      | The UserID used when connecting to the MQTT broker |
| password  | The Password used when connecting to the MQTT broker |
| topic     | Topic on which the message is published |
| qos       | MQTT Quality of Service |
| message   | The message payload |


## Configuration Examples
### Simple
Configure a task in flow to publish a "hello world" message on MQTT topic called "flogo":

```json
{
  "name": "MQTT Publisher",
  "type": "flogo:activity",
  "ref": "github.com/Ganitagya/flogo-contrib/activity/MQTT_noCert",
  "version": "0.0.1",
  "title": "Publisher MQTT Message",
  "description": "Pubishes message on MQTT topic",
  "author": "Akash Mahapatra <amahapat@tibco.com>",
  "input":[
   {
      "name": "broker",
      "type": "string",
      eg: tcp://localhost:1883
    },
    {
      "name": "id",
      "type": "string",
      eg: abc
    },
    {
      "name": "user",
      "type": "string",
      eg: abc
    },
    {
      "name": "password",
      "type": "string",
      eg: abc
    },
    {
      "name": "topic",
      "type": "string",
      eg: any_topic
    },
    {
      "name": "qos",
      "type": "integer",
      "allowed" : ["0", "1", "2"]
      eg: 0
    },
    {
      "name": "message",
      "type": "string"
      eg: Hello from Flogo
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "string"
    }
  ]
}
```
