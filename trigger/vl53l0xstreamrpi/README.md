# VL53L0XStreamRPI


Flogo Trigger for sensor VL53L0X on Raspberry Pi (Distance in mm)

## Installation

#### Install Trigger
Example: install **VL53L0XStreamRPI** trigger

```bash
flogo install github.com/prithvimoses/flogo-contrib/trigger/devices/RaspberryPi/VL53L0XStreamRPI
```


## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "delay_ms",
      "type": "integer",
      "required": false
    }
  ],
  "outputs": [
   {
      "name": "Distance",
      "type": "number"
    }
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| delay_ms      | The delay, in milliseconds, between two measures |         
Note: |* **delay_ms**: If left blank, defaut value of 500ms



## Output
| Setting     | Description    |
|:------------|:---------------|
| Distance      | The distance, in mm |            

