# timer
This trigger provides your flogo application the ability to schedule a flow via scheduling service

## Installation

```bash
flogo add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/timer
```

## Schema
Outputs and Endpoint:

```json
{
  "outputs": [
    {
      "name": "params",
      "type": "params"
    },
    {
      "name": "content",
      "type": "object"
    }
  ],
  "endpoint": {
    "settings": [
      {
        "name": "repeating",
        "type": "string"
      },
      {
        "name": "notImmediate",
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
}
```

## Example Configurations

Triggers are configured via the triggers.json of your application. The following are some example configuration of the Timer Trigger.

### repeating = false
Configure the Trigger to run a flow immediately

```json
{
  "triggers": [
    {
      "name": "timer",
      "settings": {
      },
      "endpoints": [
        {
          "flowURI": "local://new_device_flow",
          "settings": {
            "repeating": "false"
          }
        }
      ]
    }
  ]
}
```

### repeating = false
Configure the Trigger to run a flow at a certain date/time. "startDate" settings format = "mm/dd/yyyy, hours:minutes:seconds"

```json
{
  "triggers": [
    {
      "name": "rest",
      "settings": {
        "port": "8080"
      },
      "endpoints": [
        {
          "flowURI": "local://new_device_flow",
          "settings": {
            "repeating": "false",
            "startDate" : "05/01/2016, 12:25:01"
          }
        }
      ]
    }
  ]
}
```

### repeating = true
Configure the Trigger to run a flow immediately and repeating every hours|minutes|seconds. "notImmediate" set to true, the trigger will not fire immediately in this case the first execution will occur in 24 hours. If set to true the first execuction will will occur immediately.

```json
{
  "triggers": [
    {
      "name": "rest",
      "settings": {
        "port": "8080"
      },
      "endpoints": [
        {
          "flowURI": "local://new_device_flow",
          "settings": {
            "repeating": "true",
            "notImmediate": "true",
            "hours": "24"
          }
        }
      ]
    }
  ]
}
```

### repeating = true
Configure the Trigger to run a flow at a certain date/time and repeating every hours|minutes|seconds

```json
{
  "triggers": [
    {
      "name": "rest",
      "settings": {
        "port": "8080"
      },
      "endpoints": [
        {
          "flowURI": "local://new_device_flow",
          "settings": {
            "repeating": "true",
            "startDate" : "05/01/2016, 12:25:01",
            "hours": "64"
          }
        }
      ]
    }
  ]
}
```
