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
      "name": "pin number",
      "type": "int",
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
      "type": "string"
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
### Simple
Configure a task in flow to get pet '1234' from the [swagger petstore](http://petstore.swagger.io):
