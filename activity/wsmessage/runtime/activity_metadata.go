package sendWSMessage

var jsonMetadata = `{
  "name": "sendWSMessage",
  "version": "0.0.1",
  "title": "Send WebSocket Message",
  "description": "This activity sends a message to a WebSocket enabled server like TIBCO eFTL",
  "inputs":[
    {
      "name": "Server",
      "type": "string",
      "value": ""
    },
    {
      "name": "Channel",
      "type": "string",
      "value": ""
    },
    {
      "name": "Destination",
      "type": "string",
      "value": ""
    },
    {
      "name": "Message",
      "type": "string",
      "value": ""
    },
    {
      "name": "Username",
      "type": "string",
      "value": ""
    },
    {
      "name": "Password",
      "type": "string",
      "value": ""
    }
  ],
  "outputs": [
    {
      "name": "output",
      "type": "string"
    }
  ]
}`
