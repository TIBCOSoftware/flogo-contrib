package rest

var jsonMetadata = `{
  "name": "tibco-rest",
  "version": "0.0.1",
  "title": "REST Activity",
  "description": "Simple REST Activity",
  "inputs":[
    {
      "name": "method",
      "type": "string",
      "required": true
    },
    {
      "name": "uri",
      "type": "string",
      "required": true
    },
    {
      "name": "pathParams",
      "type": "params"
    },
    {
      "name": "queryParams",
      "type": "params"
    },
    {
      "name": "content",
      "type": "object"
    }
  ],
  "outputs": [
    {
      "name": "result",
      "type": "object"
    }
  ]
}
`
