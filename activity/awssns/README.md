# Amazon SNS - SMS
This activity sends SMS using Amazon Simple Notification Services (SNS).


## Installation

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/awssns
```

## Schema
Inputs and Outputs:

```json
{
  "input":[
    {
      "name": "accessKey",
      "type": "string",
      "required": "true"
    },
    {
      "name": "secretKey",
      "type": "string",
      "required": "true"
    },
    {
      "name": "region",
      "type": "string",
      "required": "true",
      "allowed" : ["us-east-1","us-west-2","eu-west-1","ap-northeast-1","ap-southeast-1","ap-southeast-2"]
    },
    {
      "name": "smsType",
      "type": "string",
      "allowed" : ["Promotional", "Transactional"]
    },
    {
      "name": "from",
      "type": "string",
      "required": "true"
    },
    {
      "name": "to",
      "type": "string",
      "required": "true"
    },
    {
      "name": "message",
      "type": "string",
      "required": "true"
    }
  ],
  "output": [
  	{
      "name": "messageId",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| accessKey        | Amazon Access Key ID |         
| secretKey        | Amazon Secret Access Key |
| region        | Default region to use to send the SMS. [See Here](http://docs.aws.amazon.com/sns/latest/dg/sms_supported-countries.html) |
| smsType        | Type of SMS to be sent (Promotional or Transactional) |
| from        | Sender ID for the SMS |
| to        | Phone number (International format) to which send the SMS |
| message        | The message itself |
| messageId        | The unique message ID returned by AWS SNS |


## More details
Please find more details about SNS regions and countries here: http://docs.aws.amazon.com/sns/latest/dg/sms_supported-countries.html
