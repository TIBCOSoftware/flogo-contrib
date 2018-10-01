---
title: Lambda
weight: 4704
---

# tibco-lambda

The Lambda trigger provides your Flogo application the ability to start flows running as an AWS Lambda function. 

## How it works

In Lambda, functions can be triggered by a variety of event sources (like Amazon S3, or the Amazon API Gateway) and in Lambda terms, each of them would need a specific method to serve as starting point. To overcome that, and make it possible to create just one flow that can be triggered by multiple events, the Lambda trigger abstracts the incoming event, making sure that no change is needed when you want to hook up other events to the same trigger.

## Installation

### Flogo Web

This activity comes out of the box with the Flogo Web UI

### Flogo CLI

```bash
flogo install trigger github.com/TIBCOSoftware/flogo-contrib/trigger/lambda
```

## Schema
Settings, Outputs:

```json
{
  "settings": [
  ],
  "output": [
    {
      "name": "context",
      "type": "object"
    },
    {
      "name": "evt",
      "type": "object"
    }
  ]
}
```

## Context

The context object contains details about the Flogo app and the runtime

```json
{
    "awsRequestId":"6301513f-c5c5-11e8-9997-c9e96bfbdd33", <-- The request ID that AWS associated with the Lambda execution
    "functionName":"TestApp", <-- The name of the app that was triggered
    "functionVersion":"$LATEST", <-- The version of the function that is executed
    "logGroupName":"/aws/lambda/TestApp", <-- The name of the log group in AWS CloudWatch where logs are sent to
    "logStreamName":"2018/10/01/[$LATEST]ec9b2c81da4b4307ab50738bb1eeb28f", <-- The logstream to which events are sent
    "memoryLimitInMB":128 <-- The amount of RAM your function has available
}
```

## Example events

The events themselves are not changed by the Lambda trigger, though without changing your Lambda code a single app can handle multiple types of events. Each of the properties of the event are accessible using the dot notation in the mappings.

### Amazon S3

When Amazon S3 emits a "PUT" event, the event details that are passed to the flow will look like

```json
{
    "Records": [
        {
            "awsRegion": "us-west-2",
            "eventName": "ObjectCreated:Put",
            "eventSource": "aws:s3",
            "eventTime": "1970-01-01T00:00:00.000Z",
            "eventVersion": "2.0",
            "requestParameters": {
                "sourceIPAddress": "127.0.0.1"
            },
            "responseElements": {
                "x-amz-id-2": "EXAMPLE123/5678abcdefghijklambdaisawesome/mnopqrstuvwxyzABCDEFGH",
                "x-amz-request-id": "EXAMPLE123456789"
            },
            "s3": {
                "bucket": {
                    "arn": "arn:aws:s3:::example-bucket",
                    "name": "example-bucket",
                    "ownerIdentity": {
                        "principalId": "EXAMPLE"
                    }
                },
                "configurationId": "testConfigRule",
                "object": {
                    "eTag": "0123456789abcdef0123456789abcdef",
                    "key": "test/key",
                    "sequencer": "0A1B2C3D4E5F678901",
                    "size": 1024
                },
                "s3SchemaVersion": "1.0"
            },
            "userIdentity": {
                "principalId": "EXAMPLE"
            }
        }
    ]
}
```

### Amazon DynamoDB

When Amazon DynamoDB emits an "Update" event, the event details that are passed to the flow will look like

```json
{
    "Records": [
        {
            "awsRegion": "us-west-2",
            "dynamodb": {
                "Keys": {
                    "Id": {
                        "N": "101"
                    }
                },
                "NewImage": {
                    "Id": {
                        "N": "101"
                    },
                    "Message": {
                        "S": "New item!"
                    }
                },
                "SequenceNumber": "111",
                "SizeBytes": 26,
                "StreamViewType": "NEW_AND_OLD_IMAGES"
            },
            "eventID": "1",
            "eventName": "INSERT",
            "eventSource": "aws:dynamodb",
            "eventSourceARN": "arn:aws:dynamodb:us-west-2:account-id:table/ExampleTableWithStream/stream/2015-06-27T00:48:05.899",
            "eventVersion": "1.0"
        }
    ]
}
```

### AWS CloudWatch

When AWS CloudWatch emits a scheduled event, the event details that are passed to the flow will look like

```json
{
    "account": "123456789012",
    "detail": {},
    "detail-type": "Scheduled Event",
    "id": "cdc73f9d-aea9-11e3-9d5a-835b769c0d9c",
    "region": "us-east-1",
    "resources": [
        "arn:aws:events:us-east-1:123456789012:rule/my-schedule"
    ],
    "source": "aws.events",
    "time": "1970-01-01T00:00:00Z"
}
```

### Amazon API Gateway

When Amazon API Gateway forwards an API call, the event details that are passed to the flow will look like

```json
{
    "body": "eyJ0ZXN0IjoiYm9keSJ9",
    "headers": {
        "Accept": "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
        "Accept-Encoding": "gzip, deflate, sdch",
        "Accept-Language": "en-US,en;q=0.8",
        "Cache-Control": "max-age=0",
        "CloudFront-Forwarded-Proto": "https",
        "CloudFront-Is-Desktop-Viewer": "true",
        "CloudFront-Is-Mobile-Viewer": "false",
        "CloudFront-Is-SmartTV-Viewer": "false",
        "CloudFront-Is-Tablet-Viewer": "false",
        "CloudFront-Viewer-Country": "US",
        "Host": "1234567890.execute-api.us-west-2.amazonaws.com",
        "Upgrade-Insecure-Requests": "1",
        "User-Agent": "Custom User Agent String",
        "Via": "1.1 08f323deadbeefa7af34d5feb414ce27.cloudfront.net (CloudFront)",
        "X-Amz-Cf-Id": "cDehVQoZnx43VYQb9j2-nvCh-9z396Uhbp027Y2JvkCPNLmGJHqlaA==",
        "X-Forwarded-For": "127.0.0.1, 127.0.0.2",
        "X-Forwarded-Port": "443",
        "X-Forwarded-Proto": "https"
    },
    "httpMethod": "POST",
    "isBase64Encoded": "false",
    "path": "/path/to/resource",
    "pathParameters": {
        "proxy": "/path/to/resource"
    },
    "queryStringParameters": {
        "foo": "bar"
    },
    "requestContext": {
        "accountId": "123456789012",
        "apiId": "1234567890",
        "httpMethod": "POST",
        "identity": {
            "accessKey": null,
            "accountId": null,
            "caller": null,
            "cognitoAuthenticationProvider": null,
            "cognitoAuthenticationType": null,
            "cognitoIdentityId": null,
            "cognitoIdentityPoolId": null,
            "sourceIp": "127.0.0.1",
            "user": null,
            "userAgent": "Custom User Agent String",
            "userArn": null
        },
        "path": "/prod/path/to/resource",
        "protocol": "HTTP/1.1",
        "requestId": "c6af9ac6-7b61-11e6-9a41-93e8deadbeef",
        "requestTime": "09/Apr/2015:12:34:56 +0000",
        "requestTimeEpoch": 1428582896000,
        "resourceId": "123456",
        "resourcePath": "/{proxy+}",
        "stage": "prod"
    },
    "resource": "/{proxy+}",
    "stageVariables": {
        "baz": "qux"
    }
}
```
