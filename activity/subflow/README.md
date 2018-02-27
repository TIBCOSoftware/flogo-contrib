# flogo-subflow
This activity provides your flogo flow the ability to start a sub-flow.

## Installation

```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/subflow
```

## Schema
The Input/Output schema is determined from the Input/Output metadata
of the sub-flow that is being executed

## Settings

```json
{
  "settings":[
    {
      "name": "flowURI",
      "type": "string",
      "required": true
    }
  ]

}
```

| Setting     | Description    |
|:------------|:---------------|
| flowURI    | The URI of the flow to execute |         


## Configuration Examples
### Simple
Configure a activity to execute "mysubflow" and set its input values to literals "1" and "2".

```json
{
  "id": "RunSubFlow",
  "activity": {
    "ref": "github.com/TIBCOSoftware/flogo-contrib/activity/subflow",
    "settings" : {
      "flowURI" : "res://flow:mysubflow"
    },
    "input": { 
  	  "mappings":[
        { "type": "literal", "value": "1", "mapTo": "FlowIn1" },
        { "type": "literal", "value": "2", "mapTo": "FlowIn2" }
      ]
    }
  }
}
```
