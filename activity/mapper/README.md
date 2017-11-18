# flogo-mapper
This activity provides your flogo application the ability to map values on to the action/flow working attribute set.

## Installation

```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/mapper
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
| mappings    | The mappings to the action/flow working data |         


## Configuration Examples
### Simple
Configure a activity to set the flow attributes to literals "1" and 2.

```json
{
  "id": "mapper",
  "type": 1,
  "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/mapper",
  "name": "Mapper",
  "input": { 
  	"mappings":[
      { "type": 2, "value": "1", "mapTo": "FlowAttr1" },
      { "type": 2, "value": 2, "mapTo": "FlowAttr2" }
    ]
  }
}
```
