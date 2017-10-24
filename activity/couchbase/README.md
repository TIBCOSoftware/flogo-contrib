# tibco-couchbase
This activity provides your flogo application the ability to connect to a Couchbase server


## Installation

```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/couchbase
```

## Schema
Inputs and Outputs:

```json
{
   "inputs":[
      {
         "name":"key",
         "type":"string",
         "required":true
      },
      {
         "name":"data",
         "type":"string"
      },
      {
         "name":"method",
         "type":"string",
         "allowed":[
            "Insert",
            "Upsert",
            "Remove",
            "Get"
         ],
         "value":"Insert",
         "required":true
      },
      {
         "name":"expiry",
         "type":"integer",
         "value":0,
         "required":true
      },
      {
         "name":"server",
         "type":"string",
         "required":true
      },
      {
         "name":"username",
         "type":"string"
      },
      {
         "name":"password",
         "type":"string"
      },
      {
         "name":"bucket",
         "type":"string",
         "required":true
      },
      {
         "name":"bucketPassword",
         "type":"string"
      }
   ],
   "outputs":[
      {
         "name":"output",
         "type":"any"
      }
   ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| key | The document key identifier |         
| data   | The document data (raw, JSON, etc.) |
| method       | The method type (Insert, Upsert, Remove or Get) |
| expiry   | The document expiry (default: 0) |
| server   | The Couchbase server (e.g. *couchbase://127.0.0.1*) |
| username   | Cluster username |
| password   | Cluster password |
| bucket   | The bucket name |
| bucketPassword   | The bucket password if any |
Note: if method is set to Get, data is ignored
## Configuration Examples
Configure an upsert method:

```json
{  
   "key":"foo",
   "data":"bar",
   "method":"Upsert",
   "expiry":0,
   "server":"couchbase://127.0.0.1",
   "username":"Administrator",
   "password":"password",
   "bucket":"test",
   "bucketPassword":""
}
```