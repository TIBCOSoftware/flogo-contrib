
# MySQL DML
This activity provides your flogo application the ability to fire a select, update, Delete and insert queries to SQL database.

# Third Party Drivers Used
https://github.com/golang/go/wiki/SQLDrivers


## Schema
Inputs and Outputs:

```json
"inputs":[
    {
      "name": "driverName",
      "type": "string",
      "required": true,
      "allowed": [
        "mysql"
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


```json
{
  "type": 1,
  "activityType": "MYSQLDML",
  "description": "Fire Select,Update,Insert and delete queries to MYSQL Database",
  "attributes": [
    { "name": "driverName", "value": "mysql" },
    { "name": "datasourceName", "value": "username:password@tcp(hostserver:port)/dbName" },
    { "name": "query", "value": "select * from table_name" }
  ]
}
```

