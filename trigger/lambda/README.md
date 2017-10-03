# tibco-lambda
This trigger provides your flogo application the ability to start a flow as an AWS Lambda function

## Installation

```bash
flogo install trigger github.com/TIBCOSoftware/flogo-contrib/trigger/lambda
```

## Schema
Settings, Outputs:

```json
{
  "settings": [
  ],
  "outputs": [
    {
      "name": "logStreamName",
      "type": "string"
    },
    {
      "name": "logGroupName",
      "type": "string"
    },
    {
      "name": "awsRequestId",
      "type": "string"
    },
    {
      "name": "memoryLimitInMB",
      "type": "string"
    },
    {
      "name": "evt",
      "type": "string"
    }
  ]
}
```