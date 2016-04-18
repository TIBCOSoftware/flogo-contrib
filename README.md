# flogo-contrib
Collection of Flogo activities, triggers and models.

## Contributions

### Activities
* [rest](activity/rest): Simple REST invoker 
* [log](activity/log): Simple flow Logger 

### Triggers
* [rest](trigger/rest): Start flow via REST
* [mqtt](trigger/mqtt): Start flow via MQTT
 
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

### Support
For Q&A you can post your questions on [Slack](https://tibco-cloud.slack.com/messages/flogo-general/)

