package timer

import (

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
    "github.com/carlescere/scheduler"
	"github.com/op/go-logging"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"time"
	"strconv"
	"fmt"
)

// log is the default package logger
var log = logging.MustGetLogger("trigger-tibco-mqtt")

type TimerTrigger struct {
	metadata    *trigger.Metadata
	flowStarter flowinst.Starter
	settings    map[string]string
	config      *trigger.Config
	timers      map[string]*scheduler.Job
}

func init() {
	md := trigger.NewMetadata(jsonMetadata)
	trigger.Register(&TimerTrigger{metadata: md})
}

// Metadata implements trigger.Trigger.Metadata
func (t *TimerTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Init implements ext.Trigger.Init
func (t *TimerTrigger) Init(flowStarter flowinst.Starter, config *trigger.Config) {

	t.flowStarter = flowStarter
	t.settings = config.Settings
	t.config = config
}

// Start implements ext.Trigger.Start
func (t *TimerTrigger) Start() error {

	log.Debug("Start")
	t.timers = make(map[string]*scheduler.Job)
	endpoints := t.config.Endpoints

	log.Debug("Processing endpoints")
	for _, endpoint := range endpoints {

		repeating := endpoint.Settings["repeating"]
		log.Debug("Repeating: ", repeating);
		if(repeating == "false") {
			t.scheduleOnce(endpoint)
		} else if(repeating == "true") {
			t.scheduleRepeating(endpoint)
		} else {
			log.Error("No match for repeating: ", repeating)
		}
		log.Debug("Settings repeating: ", endpoint.Settings["repeating"])
		log.Debugf("Processing endpoint: %s", endpoint.FlowURI)
	}

	return nil
}

// Stop implements ext.Trigger.Stop
func (t *TimerTrigger) Stop() {

	log.Debug("Stopping endpoints")
	for k, v := range t.timers {
		fmt.Println("k:", k, "v:", v)
		if(t.timers[k].IsRunning()) {
			log.Debug("Stopping timer for : ", k)
			v.Quit <- true;
		} else {
			log.Debugf("Timer: %s is not running", k)
		}
	}
}

func (t *TimerTrigger) scheduleOnce(endpoint *trigger.EndpointConfig) {
	log.Debug("Scheduling a run one time job")

	seconds := getSeconds(endpoint)
	log.Debug("Seconds till trigger fires: ", seconds)
	timerJob := scheduler.Every(int(seconds))

	if timerJob == nil {
		log.Error("timerJob is nil")
	}

	fn := func() {
		log.Debug("-- Starting \"Once\" timer process")
		id, err := t.flowStarter.StartFlowInstance(endpoint.FlowURI, nil, nil, nil)
		if err != nil {
			log.Error("Error starting flow: ", err.Error())
		}
		log.Debug("Flow ID: " +id)
		timerJob.Quit <- true
	}

	timerJob, err := timerJob.Seconds().Run(fn)
	if err != nil {
		log.Error("Error scheduleOnce flo err: ", err.Error())
	}

	t.timers[endpoint.FlowURI] = timerJob
}

func (t *TimerTrigger) scheduleRepeating(endpoint *trigger.EndpointConfig) {
	log.Debug("Scheduling a repeating job")

	fn2 := func() {
		log.Debug("-- Starting \"Repeating\" (repeat) timer process")
		id, err := t.flowStarter.StartFlowInstance(endpoint.FlowURI, nil, nil, nil)
		if err != nil {
			log.Error("Error starting flow: ", err.Error())
		}

		log.Debug("Flow ID: " +id)
	}

	seconds := getSeconds(endpoint)
	log.Debug("Seconds till trigger fires: ", seconds)
	timerJob := scheduler.Every(int(seconds))
	if timerJob == nil {
		log.Error("timerJob is nil")
	}

	fn := func() {
		log.Debug("-- Starting \"Repeating\" (first) timer process")

		id, err := t.flowStarter.StartFlowInstance(endpoint.FlowURI, nil, nil, nil)
		if err != nil {
			log.Error("Error starting flow: ", err.Error())
		}

		log.Debug("Flow ID: " +id)
		timerJob.Quit <- true

		// schedule repeating
		if(endpoint.Settings["hours"] != "") {

			log.Debug("repeatHours: ", endpoint.Settings["hours"])
			hours, _:= strconv.Atoi(endpoint.Settings["hours"])
			timerJob, err := scheduler.Every(hours).Hours().Run(fn2)
			if err != nil {
				log.Error("Error scheduleRepeating (repeat hours) flo err: ", err.Error())
			}
			if timerJob == nil {
				log.Error("timerJob is nil")
			}
			t.timers[endpoint.FlowURI + "_r"] = timerJob
		} else if(endpoint.Settings["minutes"] != "") {

			log.Debug("minutes: ", endpoint.Settings["minutes"])
			minutes, _:= strconv.Atoi(endpoint.Settings["minutes"])
			timerJob, err := scheduler.Every(minutes).Minutes().Run(fn2)
			if err != nil {
				log.Error("Error scheduleRepeating (repeat minutes) flo err: ", err.Error())
			}
			if timerJob == nil {
				log.Error("timerJob is nil")
			}
			t.timers[endpoint.FlowURI + "_r"] = timerJob
		} else if(endpoint.Settings["seconds"] != "") {

			log.Debug("repeatSeconds: ", endpoint.Settings["seconds"])
			seconds, _:= strconv.Atoi(endpoint.Settings["seconds"])
			timerJob, err := scheduler.Every(seconds).Seconds().Run(fn2)
			if err != nil {
				log.Error("Error scheduleRepeating (repeat seconds) flo err: ", err.Error())
			}
			if timerJob == nil {
				log.Error("timerJob is nil")
			}
			t.timers[endpoint.FlowURI + "_r"] = timerJob
		}
	}

	timerJob, err := timerJob.Seconds().Run(fn)
	if err != nil {
		log.Error("Error scheduleRepeating (first) flo err: ", err.Error())
	}
	if timerJob == nil {
		log.Error("timerJob is nil")
	}

	t.timers[endpoint.FlowURI] = timerJob
}

func getSeconds(endpoint *trigger.EndpointConfig) int64 {

	if(endpoint.Settings["startDate"] == "") {
		return 1;
	}

	layout := "01/02/2006, 15:04:05"
	log.Debug("startDate: ",  endpoint.Settings["startDate"])
	triggerDate, err := time.Parse(layout, endpoint.Settings["startDate"])
	if err != nil {
		log.Error("Error parsing time err: ", err.Error())
	}

	duration := time.Since(triggerDate)

	return int64(duration.Seconds())
}

type PrintJob struct {
	Msg string
}

func (j *PrintJob) Run() error {
	log.Debug(j.Msg)
	return nil
}