# Amazon SNS - SMS
This activity sends SMS using Amazon Simple Notification Services (SNS).


## Installation

```bash
flogo add activity github.com/philippegabert/flogo-contrib/activity/awssns
```

## Schema
Inputs and Outputs:

```json
{
  "inputs":[
    {
      "name": "AWS_ACCESS_KEY_ID",
      "type": "string",
      "required": "true"
    },
    {
      "name": "AWS_SECRET_ACCESS_KEY",
      "type": "string",
      "required": "true"
    },
    {
      "name": "AWS_DEFAULT_REGION",
      "type": "string",
      "required": "true",
      "allowed" : ["us-east-1","us-west-2","eu-west-1","ap-northeast-1","ap-southeast-1","ap-southeast-2"]
    },
    {
      "name": "SMS_TYPE",
      "type": "string",
      "allowed" : ["Promotional", "Transactional"]
    },
    {
      "name": "SMS_FROM",
      "type": "string",
      "required": "true"
    },
    {
      "name": "SMS_TO",
      "type": "string",
      "required": "true"
    },
    {
      "name": "SMS_MESSAGE",
      "type": "string",
      "required": "true"
    }
  ],
  "outputs": [
  	{
      "name": "MESSAGE_ID",
      "type": "string"
    }
  ]
}
```
## Settings
| Setting     | Description    |
|:------------|:---------------|
| AWS_ACCESS_KEY_ID        | Amazon Access Key ID |         
| AWS_SECRET_ACCESS_KEY        | Amazon Secret Access Key |
| AWS_DEFAULT_REGION        | Default region to use to send the SMS. [See Here](http://docs.aws.amazon.com/sns/latest/dg/sms_supported-countries.html) |
| SMS_TYPE        | Type of SMS to be sent (Promotional or Transactional) |
| SMS_FROM        | Sender ID for the SMS |
| SMS_TO        | Phone number (International format) to which send the SMS |
| SMS_MESSAGE        | The message itself |
| MESSAGE_ID        | The unique message ID returned by AWS SNS |


## More details
Please find more details about SNS regions and countries here: http://docs.aws.amazon.com/sns/latest/dg/sms_supported-countries.html
