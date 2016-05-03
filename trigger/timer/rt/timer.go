package timer

import (

	"github.com/TIBCOSoftware/flogo-lib/core/ext/trigger"
    "github.com/carlescere/scheduler"
	"github.com/op/go-logging"
	"github.com/TIBCOSoftware/flogo-lib/core/flowinst"
	"time"
	"strconv"
	"fmt"
	"math"
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

	seconds := getInitialStartInSeconds(endpoint)
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

	timerJob, err := timerJob.Seconds().NotImmediately().Run(fn)
	if err != nil {
		log.Error("Error scheduleOnce flo err: ", err.Error())
	}

	t.timers[endpoint.FlowURI] = timerJob
}

func (t *TimerTrigger) scheduleRepeating(endpoint *trigger.EndpointConfig) {
	log.Debug("Scheduling a repeating job")

	seconds := getInitialStartInSeconds(endpoint)

	fn2 := func() {
		log.Debug("-- Starting \"Repeating\" (repeat) timer process")
		id, err := t.flowStarter.StartFlowInstance(endpoint.FlowURI, nil, nil, nil)
		if err != nil {
			log.Error("Error starting flow: ", err.Error())
		}

		log.Debug("Flow ID: " +id)
	}

	if(endpoint.Settings["notImmediate"] == "false") {
		t.scheduleJobEverySecond(endpoint, fn2)
	} else {

		log.Debug("Seconds till trigger fires: ", seconds)
		timerJob := scheduler.Every(seconds)
		if timerJob == nil {
			log.Error("timerJob is nil")
		}

		t.scheduleJobEverySecond(endpoint, fn2)

		timerJob, err := timerJob.Seconds().NotImmediately().Run(fn2)
		if err != nil {
			log.Error("Error scheduleRepeating (first) flo err: ", err.Error())
		}
		if timerJob == nil {
			log.Error("timerJob is nil")
		}

		t.timers[endpoint.FlowURI] = timerJob
	}
}

func getInitialStartInSeconds(endpoint *trigger.EndpointConfig) (int) {

	if(endpoint.Settings["startDate"] == "") {
		return 0;
	}

	//layout := "2006-01-02T15:04:05Z07:00"
	layout := "Jan 2, 2006 at 3:04pm (MST)"
	log.Debug("startDate: ",  endpoint.Settings["startDate"])
	triggerDate, err := time.Parse(layout, endpoint.Settings["startDate"])
	if err != nil {
		log.Error("Error parsing time err: ", err.Error())
	}

	log.Debug("Current time: ", time.Now())
	log.Debug("Setting start time: ", triggerDate)
	duration := time.Since(triggerDate)

	return int(math.Abs(duration.Seconds()))
}

type PrintJob struct {
	Msg string
}

func (j *PrintJob) Run() error {
	log.Debug(j.Msg)
	return nil
}

func (t *TimerTrigger) scheduleJobEverySecond (endpoint *trigger.EndpointConfig, fn func()) {

	var interval int = 0;
	if(endpoint.Settings["seconds"] != "") {
		seconds, _ := strconv.Atoi(endpoint.Settings["seconds"])
		interval = interval + seconds
	}
	if(endpoint.Settings["minutes"] != "") {
		minutes, _ := strconv.Atoi(endpoint.Settings["minutes"])
		interval = interval + minutes * 60
	}
	if(endpoint.Settings["hours"] != "") {
		hours, _ := strconv.Atoi(endpoint.Settings["hours"])
		interval = interval + hours * 3600
	}

	log.Debug("Repeating seconds: ", interval)
	// schedule repeating
	timerJob, err := scheduler.Every(interval).Seconds().Run(fn)
	if err != nil {
		log.Error("Error scheduleRepeating (repeat seconds) flo err: ", err.Error())
	}
	if timerJob == nil {
		log.Error("timerJob is nil")
	}
	t.timers[endpoint.FlowURI + "_r"] = timerJob
}