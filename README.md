# flogo-contrib

[![Build Status](https://travis-ci.org/TIBCOSoftware/flogo-contrib.svg?branch=master)](https://travis-ci.org/TIBCOSoftware/flogo-contrib.svg?branch=master)

Collection of Flogo activities, triggers and models.

## Contributions

### Activities
* [awsiot](activity/awsiot): Aws IOT shadow update
* [coap](activity/coap): CoAP messaging 
* [counter](activity/counter): Global counter  
* [log](activity/log): Simple flow Logger 
* [rest](activity/rest): Simple REST invoker
* [twilio](activity/twilio): Simple Twilio SMS sender
* [websocket] (activity/wsmessage): Simple Websocket Message

### Triggers
* [coap](trigger/coap): Start flow via CoAP
* [mqtt](trigger/mqtt): Start flow via MQTT
* [rest](trigger/rest): Start flow via REST
* [timer](trigger/timer): Start flow via Timer
 
### Models
* [simple](model/simple): Basic flow model

## Installation

#### Install Activity
Example: install **log** activity

```bash
flogo add activity github.com/TIBCOSoftware/flogo-contrib/activity/log
```
#### Install Trigger
Example: install **rest** trigger

```bash
flogo add trigger github.com/TIBCOSoftware/flogo-contrib/trigger/rest
```
#### Install Model
Example: install **simple** model

```bash
flogo add model github.com/TIBCOSoftware/flogo-contrib/model/simple
```

## Contributing and support

### Contributing

New activites, triggers and models are welcome. If you would like to submit one, contact us via [Slack](https://tibco-cloud.slack.com/messages/flogo-general/).  Contributions should follow these guidelines:

* Code must be [gofmt](https://golang.org/cmd/gofmt/) compliant.
* Execute [golint](https://github.com/golang/lint) on your code.
* Document all funcs, structs and types.
* Ensure that 'go test' succeeds.

## License
flogo-contrib is licensed under a BSD-type license. See TIBCO LICENSE.txt for license text.

### Support
For Q&A you can post your questions on [Slack](https://tibco-cloud.slack.com/messages/flogo-general/)

