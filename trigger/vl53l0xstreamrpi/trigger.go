package vl53l0xstreamrpi

import (
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"context"
	"encoding/binary"
	"errors"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/rpi"
	"strconv"
	"time"
)

// Constants used for VL53L0X driver
const (
	VL53L0XAddress = 0x29
)

var interval = 500

// log is the default package logger
var log = logger.GetLogger("trigger-VL53L0XStream-RPI")

// Errors and things
var (
	ErrMeasureTimeout = errors.New("vl53l0x: measure timeout")
	ErrOutOfBounds    = errors.New("vl53l0x: measurement out of bounds")
)

// VL530LXDriver represents the I2C driver for the VL530LX proximity chip.
type VL530LXDriver struct {
	bus     embd.I2CBus
	address byte
}

// VL53L0XFactory My Trigger factory
type VL53L0XFactory struct {
	metadata *trigger.Metadata
}

// VL53L0XTrigger is a stub for your Trigger implementation
type VL53L0XTrigger struct {
	metadata *trigger.Metadata
	runner   action.Runner
	config   *trigger.Config
	driver *VL530LXDriver
}

//NewFactory create a new Trigger factory
func NewFactory(md *trigger.Metadata) trigger.Factory {
	return &VL53L0XFactory{metadata: md}
}

//New Creates a new trigger instance for a given id
func (t *VL53L0XFactory) New(config *trigger.Config) trigger.Trigger {
	return &VL53L0XTrigger{metadata: t.metadata, config: config}
}

func doEvery(d time.Duration, f func()) {
	for _ = range time.Tick(d) {
		f()
	}
}

// Init implements trigger.Trigger.Init
func (t *VL53L0XTrigger) Init(runner action.Runner) {
	t.runner = runner
	
	if t.config.Settings == nil {
		log.Infof("No configuration set for the trigger... Using default configuration...")
	} else {
		if t.config.Settings["delay_ms"] != nil && t.config.Settings["delay_ms"] != "" {
			interval, _ = strconv.Atoi(t.config.GetSetting("delay_ms"))
		} else {
			log.Infof("No delay has been set. Using default value (", interval, "ms)")
		}
	}

	t.driver = NewDriver(embd.NewI2CBus(1))
	
	log.Infof("In init, id: '%s', Metadata: '%+v', Config: '%+v'", t.config.Id, t.metadata, t.config)
}

// Metadata implements trigger.Trigger.Metadata
func (t *VL53L0XTrigger) Metadata() *trigger.Metadata {
	return t.metadata
}

// Start implements trigger.Trigger.Start
func (t *VL53L0XTrigger) Start() error {
	// start the trigger
	log.Debug("Start Trigger VL53L0XStream for Raspberry PI")
	handlers := t.config.Handlers
	//t.timers = make(map[string]*scheduler.Job)

	log.Debug("Processing handlers")
	for _, handler := range handlers {
		t.scheduleRepeating(handler)
		log.Debugf("Processing Handler: %s", handler.ActionId)
	}
	return nil
}

// Stop implements trigger.Trigger.Start
func (t *VL53L0XTrigger) Stop() error {
	// stop the trigger
	return nil
}

func NewDriver(bus embd.I2CBus) *VL530LXDriver {
	d := &VL530LXDriver{
		bus:     bus,
		address: VL53L0XAddress,
	}
	//d.bus.WriteByteToReg(d.address, 0x00, 0x01)
	
	return d
}

// Measure measures the distance detected by the driver.
func (d *VL530LXDriver) Measure() (distance int , err error) {

	byteA, err := d.bus.ReadByteFromReg(d.address, 0x1E)
	if err != nil {
		log.Error(err)
	}
	byteB, err := d.bus.ReadByteFromReg(d.address, 0x1F)
	if err != nil {
		log.Error(err)
	}
	d.bus.WriteByteToReg(d.address, 0x02, 0x01)
	distance = int(binary.BigEndian.Uint16([]byte{byteA, byteB}))

	return distance, err
}

func (t *VL53L0XTrigger) getDataFromSensor(endpoint *trigger.HandlerConfig) (distance int, err error) {
	distance, err = t.driver.Measure()
	return distance, err
}

func (t *VL53L0XTrigger) scheduleRepeating(endpoint *trigger.HandlerConfig) {

	log.Debug("Repeating every ", interval, "ms")
	fn2 := func() {
		act := action.Get(endpoint.ActionId)

		data := make(map[string]interface{})

		distance, err := t.getDataFromSensor(endpoint)
		if err != nil {
			log.Error("Error while reading sensor data. Err: ", err.Error())
		}
		data["Distance"] = distance
		log.Debug("Distance: ", distance, "mm")
		startAttrs, err := t.metadata.OutputsToAttrs(data, true)

		if err != nil {
			log.Errorf("After run error' %s'\n", err)
		}

		ctx := trigger.NewContext(context.Background(), startAttrs)
		results, err := t.runner.RunAction(ctx, act, nil)

		if err != nil {
			log.Errorf("An error occured while starting the flow. Err:", err)
		}
		log.Info("Exec: ", results)
	}

	// schedule repeating
	doEvery(time.Duration(interval)*time.Millisecond, fn2)
}