package adxl345streamrpi

import (
	"github.com/kidoman/embd"
	"fmt"
	"errors"
	"encoding/binary"
	"bytes"
)

const (
	regDevid         = 0x00
	regThreshTap     = 0x1d
	regOfsX          = 0x1e
	regOfsY          = 0x1f
	regOfsZ          = 0x20
	regDur           = 0x21
	regLatent        = 0x22
	regWindow        = 0x23
	regThreshAct     = 0x24
	regThreshInact   = 0x25
	regTimeInact     = 0x26
	regActInact_Ctl  = 0x27
	regThreshFF      = 0x28
	regTimeFF        = 0x29
	regTapAxes       = 0x2a
	regActTap_Status = 0x2b
	regBWRate        = 0x2c
	regPowerCtl      = 0x2d
	regIntEnable     = 0x2e
	regIntMap        = 0x2f
	regIntSource     = 0x30
	regDataFormat    = 0x31
	regDataX0        = 0x32
	regDataX1        = 0x33
	regDataY0        = 0x34
	regDataY1        = 0x35
	regDataZ0        = 0x36
	regDataZ1        = 0x37
	regFifoCtl       = 0x38
	regFifoStatus    = 0x39
)

const (
	powerCtl8Hz       byte = 0x00
	powerCtl4Hz       byte = 0x01
	powerCtl2Hz       byte = 0x02
	powerCtl1Hz       byte = 0x03
	powerCtlSleep     byte = 0x04
	powerCtlMeasure   byte = 0x08
	powerCtlAutoSleep byte = 0x10
	powerCtlLink      byte = 0x20
)

const (
	dataFormatRange2g   byte = 0x00
	dataFormatRange4g   byte = 0x01
	dataFormatRange8g   byte = 0x02
	dataFormatRange16g  byte = 0x03
	dataFormatJustify   byte = 0x04
	dataFormatFullRes   byte = 0x08
	dataFormatIntInvert byte = 0x20
	dataFormatSpi       byte = 0x40
	dataFormatSelfTest  byte = 0x80
)

const (
	bwRate6_25 byte = 0x06
	bwRate12_5 byte = 0x07
	bwRate25   byte = 0x08
	bwRate50   byte = 0x09
	bwRate100  byte = 0x0a
	bwRate200  byte = 0x0b
	bwRate400  byte = 0x0c
	bwRate800  byte = 0x0d
	bwRate1600 byte = 0x0e
	bwRate3200 byte = 0x0f
)


const deviceID byte = 0xE5
const DeviceAddr byte = 0x53
const fullResolutionScaleFactor float64 = 3.9

type Adxl345 struct {
	Bus      embd.I2CBus
	Opt      *Opt
	device  int
	address uint8
}

type Opt struct {

}

// NewOpt : initialize opts
func NewOpt() *Opt {
	return &Opt{
		
	}
}

type Acceleration struct {
	data [3]float64 /* mg */
}


// New : initialize ADXL345
func New(bus embd.I2CBus, opt *Opt) (*Adxl345, error) {
	adxl := &Adxl345{Bus: bus, Opt: opt}
	if err := adxl.setup(); err != nil {
		return nil, err
	}
	return adxl, nil
}

func (adxl *Adxl345) setup() error {

	if err := adxl.checkDevID(); err != nil {
		fmt.Sprintf(err.Error())
	}
	adxl.Bus.WriteByteToReg(DeviceAddr, regDataFormat, dataFormatRange16g|dataFormatFullRes)
	adxl.Bus.WriteByteToReg(DeviceAddr, regBWRate, bwRate400)
	adxl.Bus.WriteByteToReg(DeviceAddr, regPowerCtl, powerCtlMeasure)
	return nil

}

func (adxl *Adxl345) checkDevID() error {

	log.Info("Writing byte %v", regDevid)
	err := adxl.Bus.WriteByte(DeviceAddr, regDevid)
	if err != nil {
		log.Errorf("Error while writing byte to device !", err.Error())
		return err
	}
	data, err := adxl.Bus.ReadByte(DeviceAddr)
	if err != nil {
		log.Errorf("Error while reading byte from device !", err.Error())
		return err
	}

	if data != deviceID {
		errors.New(fmt.Sprintf("ADXL345 at %x on bus %d returned wrong device id: %x\n", adxl.address, adxl.device, data))
	} else
	{
		log.Debug("Device ID is correct.")
	}

	return nil
}

func (adxl *Adxl345) Destroy() {
}


func (adxl *Adxl345) Read() (*Acceleration, error) {

	log.Debug("Start reading.....")
	ret := &Acceleration{}
	var xReg int16
	var yReg int16
	var zReg int16

	log.Debug("Start adxl.Bus.WriteByte(DeviceAddr, regDataX0).....")
	err := adxl.Bus.WriteByte(DeviceAddr, regDataX0)

	if err != nil {
		return ret, err
	}

	log.Debug("adxl.Bus.ReadByte(DeviceAddr)")
	reading, err := adxl.Bus.ReadBytes(DeviceAddr, 6)
	if err != nil {
		return ret, err
	}

	buf := bytes.NewBuffer(reading)

	binary.Read(buf, binary.LittleEndian, &xReg)
	binary.Read(buf, binary.LittleEndian, &yReg)
	binary.Read(buf, binary.LittleEndian, &zReg)
	if err != nil {
		return ret, err
	}
	
	ret.data[0] = float64(xReg) * fullResolutionScaleFactor
	ret.data[1] = float64(yReg) * fullResolutionScaleFactor
	ret.data[2] = float64(zReg) * fullResolutionScaleFactor

	return ret, nil
}

