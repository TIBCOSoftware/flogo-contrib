# tibco-aggregate
This activity provides your flogo application with rudimentary aggregation.


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/aggregate
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "function",
      "type": "string",
      "allowed" : ["block_avg", "moving_avg", "timeblockavg"]
    },
    {
      "name": "windowSize",
      "type": "integer",
    },
    {
      "name": "autoReset",
      "type": "boolean"
    },
    {
      "name": "value",
      "type": "number"
    }
  ],
  "outputs": [
    {
      "name": "result",
      "type": "number"
    },
    {
      "name": "report",
      "type": "boolean"
    }
  ]
}
```
## Settings
| Setting   | Description    |
|:----------|:---------------|
| function   | The aggregate fuction, currently only average is supported |
| windowSize  | The window size of the values to aggregate |
| autoReset | Flag indicating if the window should be reset after it reports |
| value | The value to aggregate |


## Configuration Examples

Configure a task to aggregate a 'temperature' attribute with a moving window of size 5:

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-aggregate",
  "name": "Aggregate Temperature",
  "attributes": [
      { "name": "function", "value": "average" }
      { "name": "windowSize", "value": "5" }
  ]
  "inputMappings": [
    { "type": 1, "value": "temperature", "mapTo": "value" }
  ]
}
```
