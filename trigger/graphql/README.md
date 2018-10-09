---
title: GraphQL
weight: 4706
---
# tibco-graphql
This trigger serves as a GraphQL HTTP endpoint. You can pass in GraphQL queries via `GET` and `POST` requests.

## Installation

```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/trigger/graphql
```

## Schema
Settings, Outputs and Endpoint:

```json
    "settings": [
      {
        "name": "port",
        "type": "integer",
        "required": true
      },
      {
        "name": "types",
        "type": "array",
        "required": true
      },
      {
        "name": "schema",
        "type": "object",
        "required": true
      },
      {
        "name": "operation",
        "type": "string",
        "required": false,
        "value": "QUERY",
        "allowed" : ["QUERY"]
      },
      {
        "name": "path",
        "type": "string",
        "required" : true
      }
    ],
    "output": [
      {
        "name": "args",
        "type": "any"
      }
    ],
    "reply": [
      {
        "name": "data",
        "type": "any"
      }
    ],
    "handler": {
      "settings": [
        {
          "name": "resolverFor",
          "type": "string",
          "required" : true
        }
      ]
    }
```
## Settings
### Trigger:
| Setting     | Description    |
|:------------|:---------------|
| port | The port to listen on |         
| types | The GraphQL object types |
| schema | The GraphQL schema |
| operation | The GraphQL operation to support, QUERY is the only valid option |
| path | The HTTP resource path |
### Output:
| Setting     | Description    |
|:------------|:---------------|
| args      | The GraphQL query arguments |
### Handler:
| Setting     | Description    |
|:------------|:---------------|
| resolverFor      | Indicates that this handler can resolve the specified GraphQL field. The value here must match a field from the schema. |

## Example GraphQL Types

```json
        "types": [
          {
            "Name": "user",
            "Fields": {
              "id": {
                "Type": "graphql.String"
              },
              "name": {
                "Type": "graphql.String"
              }
            }
          },
          {
            "Name": "address",
            "Fields": {
              "street": {
                "Type": "graphql.String"
              },
              "number": {
                "Type": "graphql.String"
              }
            }
          }
        ]
```

## Example GraphQL Schemas

```json
        "schema": {
          "Query": {
            "Name": "Query",
            "Fields": {
              "user": {
                "Type": "user",
                "Args": {
                  "id": {
                    "Type": "graphql.String"
                  },
                  "name": {
                    "Type": "graphql.String"
                  }
                }
              },
              "address": {
                "Type": "address",
                "Args": {
                  "street": {
                    "Type": "graphql.String"
                  },
                  "number": {
                    "Type": "graphql.String"
                  }
                }
              }
            }
          }
        }
```

Note that if `user` and `address` are both to be resolvable, then a handler, which specifies `address` and `user` in the `resolverFor` field is required. Currently one Flogo action can be used to resolve a single GraphQL field, you may resolve as many fields as required with multiple handlers.

## Example Application

To build the example application, follow the steps below:

```bash
flogo create -flv github.com/TIBCOSoftware/flogo-contrib/action/flow@master,github.com/TIBCOSoftware/flogo-lib/engine@master -f ~/Downloads/example.json
```

Note that the above command assumes that you've downloaded the example.json and placed it in your Downloads dir.

```bash
cd Example
flogo build
```

Now, run the application:

```bash
cd bin
./Example
```

To test the application, send a `GET` request:

```bash
curl -g 'http://localhost:7777/graphql?query={user(name:"Matt"){name,id},address{street,number}}'
```

The following response will be returned:

```json
{"data":{"address":{"number":"123","street":"Main St."},"user":{"id":"123","name":"Matt"}}}
```