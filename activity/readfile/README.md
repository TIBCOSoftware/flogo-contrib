
# Read A line from a File
This activity provides your flogo application the ability to read a particular line from a file.


## Installation

```bash
flogo install github.com/Ganitagya/flogo-contrib/activity/readfile
```
Link for flogo web:
```
https://github.com/Ganitagya/flogo-contrib/activity/readfile
```

## Schema
Inputs and Outputs:

```json
{
 "input":[
    {
      "name": "filename",
      "type": "string",
      "required": true
    },
    {
      "name": "lineNumber",
      "type": "integer"
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
## Inputs
| Input   | Description    |
|:----------|:---------------|
| filename  | complete path of the file |
| lineNumber| Line of the file to be read |

## Ouputs
| Output   | Description    |
|:----------|:---------------|
| result    | Result of the operation |


## Configuration Examples
### Simple
Configure a task in flow :

```json

```