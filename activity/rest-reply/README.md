# tibco-rest
This activity provides your flogo application the ability to reply to a REST trigger invocation.  It is used in tandem with the REST trigger.


## Installation

```bash
flogo add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/rest
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/rest-reply
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "code",
      "type": "integer",
      "required": true
    },
    {
      "name": "data",
      "type": "any"
    }
  ],
  "outputs": [
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| code        | The HTTP response code |         
| data        | The response data |

## Configuration Examples
### Simple
Configure a activity to respone with a simple http success code.

```json
{
  "id": 3,
  "type": 1,
  "activityType": "tibco-restreply",
  "name": "Respond OK",
  "attributes": [
    { "name": "code", "value": 200 }
  ]
}
```
