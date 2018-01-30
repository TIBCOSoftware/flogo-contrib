package BME280StreamRPI

// Original file coming from https://github.com/taiyoh/go-embd-bme280/blob/master/bme280.go
import (
	"github.com/kidoman/embd"
)

// DeviceAddr is address for BME280
// ReadDataAddr is address for fetching data
const (
	DeviceAddr   = 0x77
	ReadDataAddr = 0xf7
)

// Opt : options for BME280
type Opt struct {
	TemperatureOverSampling uint
	PressureOverSampling    uint
	HumidityOverSampling    uint
	Mode                    uint
	TStandby                uint
	Filter                  uint
	SPI3WEnable             bool
}

// NewOpt : initialize opts
func NewOpt() *Opt {
	return &Opt{
		TemperatureOverSampling: 1, // Temperature oversampling x 1
		PressureOverSampling:    1, // Pressure oversampling x 1
		HumidityOverSampling:    1, // Humidity oversampling x 1
		Mode:                    3, // Normal mode
		TStandby:                5, // 1000 msec
		Filter:                  0, // off
		SPI3WEnable:             false,
	}
}

// MeasureReg : return byte for measuring setting
func (o *Opt) MeasureReg() byte {
	return byte((o.TemperatureOverSampling << 5) | (o.PressureOverSampling << 2) | o.Mode)
}

// ConfigReg : return byte for config
func (o *Opt) ConfigReg() byte {
	var spi3wEnable uint
	if o.SPI3WEnable {
		spi3wEnable = 1
	} else {
		spi3wEnable = 0
	}
	return byte((o.TStandby << 5) | (o.Filter << 2) | spi3wEnable)
}

// BME280 : client for device handling
type BME280 struct {
	Bus      embd.I2CBus
	Opt      *Opt
	calibval []byte
	tfine    int32
	digT1    uint16
	digT2    int16
	digT3    int16
	digP1    uint16
	digP2    int16
	digP3    int16
	digP4    int16
	digP5    int16
	digP6    int16
	digP7    int16
	digP8    int16
	digP9    int16
	digH1    uint8
	digH2    int16
	digH3    uint8
	digH4    int16
	digH5    int16
	digH6    int8
}

// New : initialize BME280
func New(bus embd.I2CBus, opt *Opt) (*BME280, error) {
	bme280 := &BME280{Bus: bus, Opt: opt}
	if err := bme280.setup(); err != nil {
		return nil, err
	}
	if err := bme280.calibrate(); err != nil {
		return nil, err
	}
	return bme280, nil
}

func (d *BME280) setup() error {
	regs := []byte{byte(d.Opt.HumidityOverSampling), d.Opt.MeasureReg(), d.Opt.ConfigReg()}
	for i, addr := range []byte{0xf2, 0xf4, 0xf5} {
		if err := d.Bus.WriteByteToReg(DeviceAddr, addr, regs[i]); err != nil {
			return err
		}
	}
	return nil
}

func (d *BME280) calibrateTemp() {
	d.digT1 = uint16(d.calibval[1])<<8 | uint16(d.calibval[0])
	d.digT2 = int16(d.calibval[3])<<8 | int16(d.calibval[2])
	d.digT3 = int16(d.calibval[5])<<8 | int16(d.calibval[4])
}

func (d *BME280) calibratePres() {
	d.digP1 = uint16(d.calibval[7])<<8 | uint16(d.calibval[6])
	d.digP2 = int16(d.calibval[9])<<8 | int16(d.calibval[8])
	d.digP3 = int16(d.calibval[11])<<8 | int16(d.calibval[10])
	d.digP4 = int16(d.calibval[13])<<8 | int16(d.calibval[12])
	d.digP5 = int16(d.calibval[15])<<8 | int16(d.calibval[14])
	d.digP6 = int16(d.calibval[17])<<8 | int16(d.calibval[16])
	d.digP7 = int16(d.calibval[19])<<8 | int16(d.calibval[18])
	d.digP8 = int16(d.calibval[21])<<8 | int16(d.calibval[20])
	d.digP9 = int16(d.calibval[23])<<8 | int16(d.calibval[22])
}

func (d *BME280) calibrateHum() {
	d.digH2 = int16(d.calibval[1])<<8 | int16(d.calibval[0])
	d.digH3 = uint8(d.calibval[2])
	d.digH4 = int16(d.calibval[3])<<4 | (0x0f & int16(d.calibval[4]))
	d.digH5 = int16(d.calibval[5])<<4 | (int16(d.calibval[4]) >> 4)
	d.digH6 = int8(d.calibval[6])
}

func (d *BME280) calibrate() error {
	d.calibval = make([]byte, 26)
	if err := d.Bus.ReadFromReg(DeviceAddr, byte(0x88), d.calibval); err != nil {
		return err
	}

	d.calibrateTemp()
	d.calibratePres()

	d.digH1 = uint8(d.calibval[25])
	d.calibval = make([]byte, 7)
	if err := d.Bus.ReadFromReg(DeviceAddr, byte(0xe1), d.calibval); err != nil {
		return err
	}

	d.calibrateHum()

	return nil
}

func (d *BME280) compensateTemp(raw int32) float64 {
	t1 := float64(d.digT1)
	t2 := float64(d.digT2)
	t3 := float64(d.digT3)
	raw64 := float64(raw)

	v1 := (raw64/16384.0 - t1/1024.0) * t2
	v2 := (raw64/131072.0 - t1/8192.0) * (raw64/131072.0 - t1/8192.0) * t3
	tfine := v1 + v2
	d.tfine = int32(tfine)

	return tfine / 5120.0
}

func (d *BME280) compensatePres(raw int32) float64 {
	p1 := float64(d.digP1)
	p2 := float64(d.digP2)
	p3 := float64(d.digP3)
	p4 := float64(d.digP4)
	p5 := float64(d.digP5)
	p6 := float64(d.digP6)
	p7 := float64(d.digP7)
	p8 := float64(d.digP8)
	p9 := float64(d.digP9)
	raw64 := float64(raw)

	pres := 1048576.0 - raw64

	v1 := float64(d.tfine)/2.0 - 64000.0
	v2 := v1 * v1 * p6 / 32768.0
	v2 = v2 + v1*p5*2.0
	v2 = (v2 / 4.0) + (p4 * 65536.0)
	v1 = (p3*v1*v1/524288.0 + p2*v1) / 524288.0
	v1 = (1.0 + v1/32768.0) * p1

	if v1 != 0.0 {
		pres = (pres - (v2 / 4096.0)) * 6250.0 / v1
	} else {
		return 0.0 // invalid
	}
	v1 = p9 * pres * pres / 2147483648.0
	v2 = pres * p8 / 32768.0
	pres = pres + (v1+v2+p7)/16.0

	return pres
}

func (d *BME280) compensateHum(raw int32) float64 {
	hum := float64(d.tfine) - 76800.0

	h1 := float64(d.digH1)
	h2 := float64(d.digH2)
	h3 := float64(d.digH3)
	h4 := float64(d.digH4)
	h5 := float64(d.digH5)
	h6 := float64(d.digH6)

	raw64 := float64(raw)

	if hum != 0.0 {
		hum = (raw64 - (h4*64.0 + h5/16384.0*hum)) * (h2 / 65536.0 * (1.0 + h6/67108864.0*hum*(1.0+h3/67108864.0*hum)))
	} else {
		return 0.0 // invalid
	}

	hum = hum * (1.0 - h1*hum/524288.0)
	if hum > 100.0 {
		hum = 100.0
	} else if hum < 0.0 {
		hum = 0.0
	}

	return hum
}

// Read : fetching data from BME280
func (d *BME280) Read() ([]float64, error) {
	data := make([]byte, 8)
	if err := d.Bus.ReadFromReg(DeviceAddr, ReadDataAddr, data); err != nil {
		return nil, err
	}

	presRaw := int32(uint32(data[0])<<12 | uint32(data[1])<<4 | uint32(data[2])>>4)
	tempRaw := int32(uint32(data[3])<<12 | uint32(data[4])<<4 | uint32(data[5])>>4)
	humRaw := int32(uint32(data[6])<<8 | uint32(data[7]))

	d.tfine = 0 // initialize
	temp := d.compensateTemp(tempRaw)
	pres := d.compensatePres(presRaw)
	hum := d.compensateHum(humRaw)

	return []float64{temp, pres, hum}, nil
}
