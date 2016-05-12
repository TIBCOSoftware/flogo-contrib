package coap

var jsonMetadata = `{
  "name": "tibco-coap",
  "version": "0.0.1",
  "title": "CoAP Activity",
  "description": "Simple CoAP Activity",
  "inputs":[
    {
      "name": "address",
      "type": "string",
      "required": true
    },
    {
      "name": "method",
      "type": "string",
      "required": true
    },
    {
      "name": "type",
      "type": "string"
    },
    {
      "name": "messageId",
      "type": "integer"
    },
    {
      "name": "options",
      "type": "params"
    },
    {
      "name": "payload",
      "type": "string"
    }

  ],
  "outputs": [
    {
      "name": "response",
      "type": "string"
    }
  ]
}
`
