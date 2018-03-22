
# MySQL Fetch
This activity provides your flogo application the ability to make a select query and returns the table selected as a JSON string.


## Installation

```bash
flogo add activity github.com/Ganitagya/flogo-contrib/activity/mySQLFetch
```

## Schema
Inputs and Outputs:

```json
"input":[
   {
      "name": "host",
      "type": "string"
    },
    {
      "name": "username",
      "type": "string"
    },
    {
      "name": "password",
      "type": "string"
    },
    {
      "name": "database",
      "type": "string"
    },
    {
      "name": "query",
      "type": "string"
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "string"
    }
  ]
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| host       | The IP address (along with port ) where mysql is hosted (eg: localhost:3306 )|    
| username   | mySQL username |  
| password   | mySQL password for above username |  
| database   | the mySQL db from where you want to fetch data from |  
| query      | The complete select statement |  

## Configuration


### Flow Configuration
Configure a task in flow to perform Fast Fourier Transform

```json
{
  "name": "MySql Select",
  "type": "flogo:activity",
  "ref": "github.com/Ganitagya/flogo-contrib/activity/mySQLFetch",
  "version": "0.0.1",
  "title": "MySQL Select",
  "description": "Reads MySQL DB",
  "author": "Akash Mahapatra <amahapat@tibco.com>",
  "input":[
   {
      "name": "host",
      "type": "string"
    },
    {
      "name": "username",
      "type": "string"
    },
    {
      "name": "password",
      "type": "string"
    },
    {
      "name": "database",
      "type": "string"
    },
    {
      "name": "query",
      "type": "string"
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
