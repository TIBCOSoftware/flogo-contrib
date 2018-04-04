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
flogo install github.com/TIBCOSoftware/flogo-contrib/activity/log
```
#### Install Trigger
Example: install **rest** trigger

```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/trigger/rest
```
#### Install Model
Example: install **simple** model

```bash
flogo install github.com/TIBCOSoftware/flogo-contrib/model/simple
```

## Contributing and support

### Contributing

New activites, triggers and models are welcomed. If you would like to contribute, please following the [contribution guidelines](https://github.com/TIBCOSoftware/flogo/blob/master/CONTRIBUTING.md). If you have any questions, issues, etc feel free to chat with us on [Gitter](https://gitter.im/project-flogo/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link).

## License
flogo-contrib is licensed under a BSD-type license. See TIBCO LICENSE.txt for license text.

### Support
For Q&A you can post your questions on [Gitter](https://gitter.im/project-flogo/Lobby?utm_source=share-link&utm_medium=link&utm_campaign=share-link)

