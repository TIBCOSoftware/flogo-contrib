package coap

var jsonMetadata = `{
  "name": "tibco-coap",
  "version": "0.0.1",
  "titile": "CoAP",
  "description": "Simple CoAP Trigger",
  "settings": [
    {
      "name": "port",
      "type": "integer",
      "required": true
    }
  ],
  "outputs": [
    {
      "name": "queryParams",
      "type": "params"
    },
    {
      "name": "payload",
      "type": "string"
    }
  ],
  "endpoint": {
    "settings": [
      {
        "name": "method",
        "type": "string",
        "required" : true
      },
      {
        "name": "path",
        "type": "string",
        "required" : true
      },
      {
        "name": "autoIdReply",
        "type": "boolean"
      }
    ]
  }
}
`
