package bme280_go

import (
	"fmt"
	"time"
	"math"
	"golang.org/x/exp/io/i2c"
)

// BME280 holds the device and all configuration data for the sensor
type BME280 struct {
	Dev            *i2c.Device
	Config         bme280Config
	tempressConfig [3]int
	pressConfig    [9]int
	hConfig        [6]int
	tFine						int32
	_addr					 int

}
const (

	bme280Address = 0x77
	bme280AltAddress = 0x76
  bme280RegisterDigT1 = 0x88
	bme280RegisterDigT2 = 0x8A
	bme280RegisterDigT3 = 0x8C

	bme280RegisterDigP1 = 0x8E
	bme280RegisterDigP2 = 0x90
	bme280RegisterDigP3 = 0x92
	bme280RegisterDigP4 = 0x94
	bme280RegisterDigP5 = 0x96
	bme280RegisterDigP6 = 0x98
	bme280RegisterDigP7 = 0x9A
	bme280RegisterDigP8 = 0x9C
	bme280RegisterDigP9 = 0x9E

	bme280RegisterDigH1 = 0xA1
	bme280RegisterDigH2 = 0xE1
	bme280RegisterDigH3 = 0xE3
	bme280RegisterDigH4 = 0xE4
	bme280RegisterDigH5 = 0xE5
	bme280RegisterDigH6 = 0xE7

	bme280RegisterChipID = 0xD0
	bme280RegisterVersion = 0xD1
	bme280RegisterSoftReset = 0xE0

	bme280RegisterCal26 = 0xE1 // R calibration stored in 0xE1-0xF0

	bme280RegisterControlHumid = 0xF2
	bme280RegisterStatus = 0xF3
	bme280RegisterControl= 0xF4
	bme280RegisterConfig = 0xF5
	bme280RegisterPressureData = 0xF7
	bme280RegisterTempData = 0xFA
	bme280RegisterHumidData = 0xFD
)

type bme280Config struct {
	digT1 uint16///< temperature compensation value
	digT2 uint16 ///< temperature compensation value
	digT3 uint16 ///< temperature compensation value

	digP1 uint16///< pressure compensation value
	digP2 uint16 ///< pressure compensation value
	digP3 uint16 ///< pressure compensation value
	digP4 uint16  ///< pressure compensation value
	digP5 uint16  ///< pressure compensation value
	digP6 uint16  ///< pressure compensation value
	digP7 uint16  ///< pressure compensation value
	digP8 uint16  ///< pressure compensation value
	digP9 uint16  ///< pressure compensation value

	digH1 uint8 ///< humidity compensation value
	digH2 uint16 ///< humidity compensation value
	digH3 uint8 ///< humidity compensation value
	digH4 uint16 ///< humidity compensation value
	digH5 uint16 ///< humidity compensation value
	digH6 uint8  ///< humidity compensation value
  }
// BMEData contains the read-values from the sensor
type BMEData struct {
	Temperature float32
	Humidity    float32
	Pressure    float32
	Altitude	float32
}

//
// BME280Init Opens a file system handle to the I2C device
// reads the calibration data and sets the device
// into auto sensing mode
//
func (bme280 *BME280) BME280Init(channel string, addr int) int {
	var err error
	bme280._addr = addr
	bme280.Dev, err = i2c.Open(&i2c.Devfs{Dev: channel}, bme280._addr)
	if err != nil {
		panic(err)
	}
	// get ID
	b := []byte{0x00}
	err = bme280.Dev.ReadReg(bme280RegisterChipID, b)
	if err != nil {
		fmt.Println("Chip ID Err ", err)
		return -256
	}
	if int(b[0]) != 0x60 {
		fmt.Println("Invalid chip ID! ", b)
		return -256
	}
	cB := []byte{bme280RegisterSoftReset, 0xB6}

	err = bme280.Dev.Write(cB)
	if err != nil {
		fmt.Println("Device Reset Failed ", err)
		return -256
	}
	time.Sleep(300 * time.Millisecond)
// if chip is still reading calibration, delay
    for {
		breaker :=  bme280.isReadingCalibration()
	 	if(breaker){
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}

    bme280.readCoefficients(); // read trimming parameters, see DS 4.2.2
	calib1 := make([]byte, 24)
	// Read 24 bytes of calibration data
	err = bme280.Dev.ReadReg(0x88, calib1)

	if err != nil {
		fmt.Println("calibration data not read correctly", err)
		return -256
	}
	return 0
} /* bme280Init() */

func  (bme280 *BME280) isReadingCalibration() bool {
	var rStatus []byte
	err := bme280.Dev.ReadReg(bme280RegisterStatus, rStatus);
	if err != nil {
		fmt.Println("Register read error: ", err)
		return false
	}
	var vstat = rStatus[0]
	const ft = uint8(1 << 0)
	return uint8(vstat & ft) != 0
}

func (bme280 *BME280) readCoefficients() int {
	var readByte []byte
	err := bme280.Dev.ReadReg(bme280RegisterDigT1, readByte)

	if err != nil {
		fmt.Println("Configuration read error: ", err)
		return -256
	}
	bme280.Config.digT1 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigT2, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digT2 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigT3, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digT3 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigP1, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digP1 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigP2, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digP2 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigP3, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digP3 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigP4, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digP4 = uint16(readByte[0] | readByte[0])
	err = bme280.Dev.ReadReg(bme280RegisterDigP5, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digP5 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigP6, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digP6 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigP7, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digP7 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigP8, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digP8 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigP9, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
  	bme280.Config.digP9 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigH1, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digH1 = uint8(readByte[0])
	err = bme280.Dev.ReadReg(bme280RegisterDigH2, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digH2 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigH3, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digH3 = uint8(readByte[0])
	err = bme280.Dev.ReadReg(bme280RegisterDigH4, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
	}
	temp := readByte[0]
	bme280.Config.digH4 = uint16(readByte[0] | readByte[1])
	err = bme280.Dev.ReadReg(bme280RegisterDigH4 + 1, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
	}
	bme280.Config.digH4 = uint16(temp << 4 | readByte[0] & 0xF)
	err = bme280.Dev.ReadReg(bme280RegisterDigH5 + 1, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
	}
	temp = readByte[0]
	err = bme280.Dev.ReadReg(bme280RegisterDigH5, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
	}
	bme280.Config.digH5 = uint16(temp << 4 | readByte[0] >> 4);
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	err = bme280.Dev.ReadReg(bme280RegisterDigH6, readByte)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	bme280.Config.digH6 = uint8(readByte[0])
	return 0
  }

//
// bme280ReadValues Reads the sensor register values
// and translate them into calibrated readings
// using the previously loaded calibration data
// Temperature is expressed in Celsius degrees as T * 100 (for 2 decimal places)
// Pressure is <future>
// Humidity is express as H * 1024 (10 bit fraction)
//
func (bme280 *BME280) BME280ReadValues() BMEData {
	data := BMEData{}

	data.Temperature = bme280.BME280ReadTemperature()
	data.Pressure = bme280.BME280ReadPressure()
	data.Humidity = bme280.BME280ReadHumidity()
	data.Altitude = bme280.BME280ReadAltitude(1013.25)

	return data
}

func (bme280 *BME280) BMESetAddress(addr byte) {
	bme280._addr = addr
}
func (bme280 *BME280) BME280ReadTemperature() float32 {
	var var1, var2, adcT int32
    var readByte []byte
	bme280.Dev.ReadReg(bme280RegisterTempData, readByte)
    adcT = int32(readByte[0])
	adcT = adcT | int32(readByte[1])
    adcT = adcT << 8;
	adcT = adcT | int32(readByte[2] )
	if adcT == 0x800000 { // value in case temp measurement was disabled
	  return -256.00;
	}
	adcT >>= 4;

	var1 = ((((adcT >> 3) - (int32(bme280.Config.digT1 << 1)))) *
			(int32(bme280.Config.digT2))) >> 11

	var2 = (((((adcT >> 4) - (int32(bme280.Config.digT1))) *
			  ((adcT >> 4) - (int32(bme280.Config.digT1)))) >> 12) *
			(int32(bme280.Config.digT3))) >> 14

	bme280.tFine = var1 + var2

	var T = float32((bme280.tFine * 5 + 128) >> 8)
	return T / 100.00
  }

  func (bme280 *BME280) BME280ReadPressure() float32 {
	var var1, var2, p int64

	_ = bme280.BME280ReadTemperature(); // must be done first to get t_fine
	var readByte []byte
	bme280.Dev.ReadReg(bme280RegisterPressureData, readByte)
    var adcP = int32(readByte[0])
	adcP = adcP | int32(readByte[1])
    adcP = adcP << 8;
	adcP = adcP | int32(readByte[2] )
	if adcP == 0x800000 { // value in case pressure measurement was disabled
	  return -256.00
	}
	adcP >>= 4

	var1 = (int64(bme280.tFine)) - 128000
	var2 = var1 * var1 * int64(bme280.Config.digP6)
	var2 = var2 + ((var1 * int64(bme280.Config.digP5)) << 17);
	var2 = var2 + ((int64(bme280.Config.digP4)) << 35)
	var1 = ((var1 * var1 * int64(bme280.Config.digP3)) >> 8) +
		   ((var1 * int64(bme280.Config.digP2)) << 12)
	var1 =
		((((int64(1)) << 47) + var1)) * (int64(bme280.Config.digP1)) >> 33

	if (var1 == 0) {
	  return 0; // avoid exception caused by division by zero
	}
	p = int64(1048576 - adcP)
	p = (((p << 31) - var2) * 3125) / var1
	var1 = ((int64(bme280.Config.digP9)) * (p >> 13) * (p >> 13)) >> 25
	var2 = ((int64(bme280.Config.digP8)) * p) >> 19

	p = ((p + var1 + var2) >> 8) + ((int64(bme280.Config.digP7)) << 4)
	return float32(p / 256)
  }
  func (bme280 *BME280) BME280ReadHumidity() float32 {
	_ = bme280.BME280ReadTemperature(); // must be done first to get t_fine
	var adcH int32
	var readByte []byte
	bme280.Dev.ReadReg(bme280RegisterHumidData, readByte)
    adcH = int32(readByte[0])
	adcH = adcH | int32(readByte[1])
    adcH = adcH << 8;
	adcH = adcH | int32(readByte[2] )
	if adcH == 0x8000{ // value in case humidity measurement was disabled
	  return -256.00
	}

	var vx1u32r int32

	vx1u32r = (bme280.tFine - (int32(76800)))

	vx1u32r = (((((adcH << 14) - ((int32(bme280.Config.digH4)) << 20) -
					((int32(bme280.Config.digH5)) * vx1u32r)) +
				   (int32(16384))) >> 15) *
				 (((((((vx1u32r * (int32(bme280.Config.digH6))) >> 10) *
					  (((vx1u32r * (int32(bme280.Config.digH3))) >> 11) +
					   (int32(32768)))) >> 10) +
					(int32(2097152))) *
					   (int32(bme280.Config.digH2)) + 8192) >> 14))

	vx1u32r = (vx1u32r - (((((vx1u32r >> 15) * (vx1u32r >> 15)) >> 7) *
							   (int32(bme280.Config.digH1))) >> 4));

	if vx1u32r < 0 {
		vx1u32r = 0
	}
	if vx1u32r > 419430400 {
		vx1u32r = 419430400
	}
	var h = (vx1u32r >> 12);
	return float32(h) / 1024.0;
  }

  func (bme280 *BME280) BME280ReadAltitude(seaLevel float32) float32 {
	// Equation taken from BMP180 datasheet (page 16):
	//  http://www.adafruit.com/datasheets/BST-BMP180-DS000-09.pdf

	// Note that using the equation from wikipedia can give bad results
	// at high altitude. See this thread for more information:
	//  http://forums.adafruit.com/viewtopic.php?f=22&t=58064

	atmospheric := bme280.BME280ReadPressure() / 100.00;
	return float32(44330.0 * (1.0 - math.Pow(float64(atmospheric / seaLevel), 0.1903)))
  }

  /*!
   *   Calculates the pressure at sea level (in hPa) from the specified
   * altitude (in meters), and atmospheric pressure (in hPa).
   *   @param  altitude      Altitude in meters
   *   @param  atmospheric   Atmospheric pressure in hPa
   *   @returns the pressure at sea level (in hPa) from the specified altitude
   */
  func (bme280 *BME280) BME280SeaLevelForAltitude( altitude float32, atmospheric float32) float32 {
	// Equation taken from BMP180 datasheet (page 17):
	//  http://www.adafruit.com/datasheets/BST-BMP180-DS000-09.pdf

	// Note that using the equation from wikipedia can give bad results
	// at high altitude. See this thread for more information:
	//  http://forums.adafruit.com/viewtopic.php?f=22&t=58064

	return float32(float64(atmospheric) / math.Pow(float64(1.0 - (altitude / 44330.0)), 5.255))
  }
