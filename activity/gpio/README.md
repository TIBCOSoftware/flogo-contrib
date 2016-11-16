# tibco-gpio
This activity provides your flogo application the ability to control raspberry pi GPIO

## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/gpio
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "method",
      "type": "string",
      "required": true,
      "allowed" : ["Direction", "Set State", "Read State", "Pull"]
    },
    {
      "name": "pinNumber",
      "type": "integer",
      "required": true
    },
    {
      "name": "direction",
      "type": "string",
      "allowed" : ["Input", "Output"]
    },
    {
      "name": "state",
      "type": "string",
      "allowed" : ["High", "Low"]
    },

    {
      "name": "Pull",
      "type": "string",
      "allowed" : ["Up", "Down", "Off"]
    }
  ],
  "outputs": [
    {
      "name": "result",
      "type": "integer"
    }
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| method      | The method to take action for GPIO|         
| pinNumber   | The pin number   |
| direction   | The direction of pin number, either Input or Output |
| state       | The state of pin number, either high or low |
| Pull        | Pull the pin number to Up, Down and Off |


## Configuration Examples
### Get pin state
Get specific pin 23's state
```json
  "attributes": [
          {
            "name": "method",
            "value": "Read State",
            "type": "string"
          },
          {
            "name": "pinNumber",
            "value": "23",
            "type": "integer"
          }
        ]
```
### Set pin state
Set pin state to High
```json
  "attributes": [
          {
            "name": "method",
            "value": "Set State",
            "type": "string"
          },
          {
            "name": "pinNumber",
            "value": "23",
            "type": "integer"
          },
          {
            "name": "state",
            "value": "High",
            "type": "string"
          }
        ]
```
### Change pin's direction
Change pin's direction to Output
```json
  "attributes": [
          {
            "name": "method",
            "value": "Direction",
            "type": "string"
          },
          {
            "name": "pinNumber",
            "value": "23",
            "type": "integer"
          },
          {
            "name": "direction",
            "value": "Output",
            "type": "string"
          }
        ]
```
