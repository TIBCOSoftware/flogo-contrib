---
title: Return
weight: 4602
---
# flogo-return
This activity provides your flogo action/flow the ability to return immediately and set output values.

## Installation

```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/actreturn
```

## Schema
Input and Output:

```json
{
  "input":[
    {
      "name": "mappings",
      "type": "array",
      "required": true
    }
  ],
  "output": [
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| mappings    | The mappings to the action/flow ouputs |         


## Configuration Examples
### Simple
Configure a activity to return and set the output values to literals "1" and 2.

```json
{
  "id": "reply",
  "type": 1,
  "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/actreturn",
  "name": "Reply",
  "input": { 
  	"mappings":[
      { "type": 2, "value": "1", "mapTo": "Output1" },
      { "type": 2, "value": 2, "mapTo": "Output2" }
    ]
  }
}
```
