# MySQL Fetch
This activity provides your flogo application the ability to fire a select query to mySQL database and fetch the required data returning it as a JSON string.

## Installation
### Flogo Web
Use the below link for to install this activity in Flogo Web
```
https://github.com/Ganitagya/flogo-contrib/activity/mysqlfetch
```
### Flogo CLI
```bash
flogo add activity github.com/Ganitagya/flogo-contrib/activity/mysqlfetch
```

## Schema
Inputs and Outputs:

```json
"input":[
   {
      "name": "host",
      "type": "string",
      "required": true
    },
    {
      "name": "username",
      "type": "string",
      "required": true
    },
    {
      "name": "password",
      "type": "string",
      "required": true
    },
    {
      "name": "database",
      "type": "string",
      "required": true
    },
    {
      "name": "query",
      "type": "string",
      "required": true
    }
  ],
  "output": [
    {
      "name": "result",
      "type": "any"
    }
  ]
```

## Settings
| Setting     | Required | Description |
|:------------|:---------|:------------|
| host        | True     | The hostname of the MySQL server, including the portnumber of the server  (like `localhost:3306`)|    
| username    | True     | The username to connect to MySQL |  
| password    | True     | The password to connect to MySQL |  
| database    | True     | The name of the MySQL database |  
| query       | True     | The `SELECT` query you want to execute |  
