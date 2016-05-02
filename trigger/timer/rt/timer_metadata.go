package timer

var jsonMetadata = `{
  "name": "tibco-timer",
  "version": "0.0.1",
  "description": "Simple Timer Trigger",
  "settings":[
  ],
  "endpoint": {
    "settings": [
      {
        "name": "repeating",
        "type": "string"
      },
      {
        "name": "startDate",
        "type": "string"
      },
      {
        "name": "hours",
        "type": "string"
      },
      {
        "name": "minutes",
        "type": "string"
      },
      {
        "name": "seconds",
        "type": "string"
      }
    ]
  }
}`