github.com/Ganitagya/flogo-contrib/activity/Database_Query


# databasequery 
## Query data from mysql/ postgres/ sqlite3
This activity provides your flogo application the ability to fire a select query to SQL database and fetch the required data returning it as a JSON string.

# Third Party Drivers Used
https://github.com/golang/go/wiki/SQLDrivers


## Installation

```bash
flogo add activity github.com/Ganitagya/flogo-contrib/activity/databasequery
```

## Schema
Inputs and Outputs:

```json
"inputs":[
    {
      "name": "driverName",
      "type": "string",
      "required": true,
      "allowed": [
        "mysql",
        "postgres",
        "sqlite3"
      ]
    },
    {
      "name": "datasourceName",
      "type": "string",
      "required": true
    },
    {
      "name": "query",
      "type": "string",
      "required": true
    }
  ],
  "outputs": [
    {
      "name": "result",
      "type": "string"
    }
  ]
  ```
  
  ## Settings
| Setting     | Description    |
|:------------|:---------------|
| driverName  | The database driver to be used |    
| datasourceName | The connection string for a particular db driver |  
| query      | The complete select statement |  

## Configuration Examples
### MYSQL
Configure a task to increment a 'messages' counter:

```json
{
  "type": 1,
  "activityType": "MYSQL Query",
  "description": "Fetch data from MYSQL Database",
  "attributes": [
    { "name": "driverName", "value": "mysql" },
    { "name": "datasourceName", "value": "username:password@tcp(hostserver:port)/dbName" },
    { "name": "query", "value": "select * from table_name" }
  ]
}
```

### POSTGRES
Configure a task to increment a 'messages' counter:

```json
{
  "type": 2,
  "activityType": "POSTGRES Query",
  "description": "Fetch data from POSTGRES Database",
  "attributes": [
    { "name": "driverName", "value": "postres" },
    { "name": "datasourceName", "value": "host=localhost port=5432 user=postgres password=mysecretpassword dbname=postgres sslmode=disable" },
    { "name": "query", "value": "select * from table_name" }
  ]
}
```

### SQLITE3
Configure a task to increment a 'messages' counter:

```json
{
  "type": 3,
  "activityType": "SQLITE3 Query",
  "description": "Fetch data from SQLITE3 Database",
  "attributes": [
    { "name": "driverName", "value": "sqlite3" },
    { "name": "datasourceName", "value": "sqlite3db" },
    { "name": "query", "value": "select * from table_name" }
  ]
}
```

