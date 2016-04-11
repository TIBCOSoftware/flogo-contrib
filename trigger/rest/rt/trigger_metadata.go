package rest

var jsonMetadata = `{
  "name": "rest",
  "version": "0.0.1",
  "description": "Simple REST trigger",
  "settings": [
    {
      "name": "port",
      "type": "number"
    }
  ],
  "outputs": [
    {
      "name": "params",
      "type": "string"
    },
    {
      "name": "content",
      "type": "string"
    }
  ],
  "endpoint": {
    "settings": [
      {
        "name": "method",
        "type": "string"
      },
      {
        "name": "path",
        "type": "string"
      }
    ]
  }
}`
