---
title: GPIO
weight: 4611
---

# GPIO
This activity allows you to control the GPIO pins on a Raspberry Pi

## Installation
### Flogo Web
This activity comes out of the box with the Flogo Web UI
### Flogo CLI
```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/gpio
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
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
  "output": [
    {
      "name": "result",
      "type": "integer"
    }
  ]
}
```
## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| method      | True     | The method to take action for specified pin (Allowed values are Direction, Set State, Read State, and Pull) |         
| pinNumber   | True     | The pin number of the GPIO |
| direction   | False    | Set the direction of the pin (Allowed values are Input and Output) |
| state       | False    | Set the state of the pin (Allowed values are High and Low) |
| Pull        | False    | Pull the pin to the specified value (Allowed values are Up, Down, and Off) |
| result      | False    | The result of the operation |

## Examples
### Get pin state
The below example retrieves the state of pin 23:
```json
"input": {
  "method": "Read State",
  "npinNumberame": 23
}
```

### Set pin state
The below example sets the state of pin 23 to High:
```json
"input": {
  "method": "Set State",
  "npinNumberame": 23,
  "state": "High"
}
```

### Change pin's direction
The below example changes the direction of the pin to Output:
```json
"input": {
  "method": "Direction",
  "npinNumberame": 23,
  "direction": "Output"
}
```