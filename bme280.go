package bme280_go

import (
	"fmt"

	"time"

	"golang.org/x/exp/io/i2c"
)

// BME280 holdsw the device and all configuration data for the sensor
type BME280 struct {
	Dev            *i2c.Device
	tempressConfig [3]int
	pressConfig    [9]int
	hConfig        [6]int
}

// BMEData contains the read-values from the sensor
type BMEData struct {
	Temperature int
	Humidity    int
	Pressure    int
}

//
// BME280Init Opens a file system handle to the I2C device
// reads the calibration data and sets the device
// into auto sensing mode
//
func (bme280 *BME280) BME280Init(channel string) int {
	//device := BME280{}
	bme280.tempressConfig = make([]int, 3)
	bme280.pressConfig = make([]int, 9)
	bme280.hConfig = make([]int, 6)
	ucCal := make([]byte, 36)
	var err error
	bme280.Dev, err = i2c.Open(&i2c.Devfs{Dev: channel}, 0x77)
	if err != nil {
		panic(err)
	}
	//defer device.dev.Close()
	// get ID
	b := []byte{0x00}
	err = bme280.Dev.ReadReg(0xD0, b)
	if err != nil {
		fmt.Println("Chip ID Err ", err)
		return -256
	}
	if int(b[0]) != 0x60 {
		fmt.Println("Invalid chip ID! ", b)
		return -256
	}
	cB := []byte{0xE0, 0xB6}

	err = bme280.Dev.Write(cB)
	if err != nil {
		fmt.Println("Device Reset Failed ", err)
		return -256
	}
	time.Sleep(300 * time.Millisecond)

	calib1 := make([]byte, 24)
	// Read 24 bytes of calibration data
	err = bme280.Dev.ReadReg(0x88, calib1)

	if err != nil {
		fmt.Println("calibration data not read correctly", err)
		return -256
	}
	x := 0
	for x < 24 {
		ucCal[x] = calib1[x]
		x++
	}
	b = []byte{0x00}
	// get humidity calibration byte
	err = bme280.Dev.ReadReg(0xA1, b)
	if err != nil {
		fmt.Println("Failed to read humidity calibration byte", err)
		return -256
	}
	ucCal[x] = b[0]
	x++
	calib2 := make([]byte, 7)
	err = bme280.Dev.ReadReg(0xE1, calib2)
	if err != nil {
		fmt.Println("Failed to read humiduty calibration byte 2: ", err)
		return -256
	}
	y := 0
	for x < 32 {
		ucCal[x] = calib2[y]
		x++
		y++
	}
	// time to set up the calibration
	bme280.tempressConfig[0] = int(ucCal[0]) + (int(ucCal[1]) << 8)
	bme280.tempressConfig[1] = int(ucCal[2]) + (int(ucCal[3]) << 8)
	if bme280.tempressConfig[1] > 32767 {
		bme280.tempressConfig[1] -= 65536
	}
	bme280.tempressConfig[2] = int(ucCal[4]) + (int(ucCal[5]) << 8)
	if bme280.tempressConfig[2] > 32767 {
		bme280.tempressConfig[2] -= 65536
	}
	// Prepare pressure calibration data
	bme280.pressConfig[0] = int(ucCal[6]) + (int(ucCal[7]) << 8)
	bme280.pressConfig[1] = int(ucCal[8]) + (int(ucCal[9]) << 8)
	if bme280.pressConfig[1] > 32767 {
		bme280.pressConfig[1] -= 65536
	}
	bme280.pressConfig[2] = int(ucCal[10]) + (int(ucCal[11]) << 8)
	if bme280.pressConfig[2] > 32767 {
		bme280.pressConfig[2] -= 65536
	}
	bme280.pressConfig[3] = int(ucCal[12]) + (int(ucCal[13]) << 8)
	if bme280.pressConfig[3] > 32767 {
		bme280.pressConfig[3] -= 65536
	}
	bme280.pressConfig[4] = int(ucCal[14]) + (int(ucCal[15]) << 8)
	if bme280.pressConfig[4] > 32767 {
		bme280.pressConfig[4] -= 65536
	}
	bme280.pressConfig[5] = int(ucCal[16]) + (int(ucCal[17]) << 8)
	if bme280.pressConfig[5] > 32767 {
		bme280.pressConfig[5] -= 65536
	}
	bme280.pressConfig[6] = int(ucCal[18]) + (int(ucCal[19]) << 8)
	if bme280.pressConfig[6] > 32767 {
		bme280.pressConfig[6] -= 65536
	}
	bme280.pressConfig[7] = int(ucCal[20]) + (int(ucCal[21]) << 8)
	if bme280.pressConfig[7] > 32767 {
		bme280.pressConfig[7] -= 65536
	}
	bme280.pressConfig[8] = int(ucCal[22]) + (int(ucCal[23]) << 8)
	if bme280.pressConfig[8] > 32767 {
		bme280.pressConfig[8] -= 65536
	}
	// Prepare humidity calibration data
	bme280.hConfig[0] = int(ucCal[24])
	bme280.hConfig[1] = int(ucCal[25]) + (int(ucCal[26]) << 8)
	if bme280.hConfig[1] > 32767 {
		bme280.hConfig[1] -= 65536
	}
	bme280.hConfig[2] = int(ucCal[27])
	bme280.hConfig[3] = (int(ucCal[28]) << 4) + (int(ucCal[29]) & 0xf)
	if bme280.hConfig[3] > 2047 {
		bme280.hConfig[3] -= 4096
	}
	bme280.hConfig[4] = (int(ucCal[30]) << 4) + (int(ucCal[29]) >> 4)
	if bme280.hConfig[4] > 2047 {
		bme280.hConfig[4] -= 4096
	}
	bme280.hConfig[5] = int(ucCal[31])
	if bme280.hConfig[5] > 127 {
		bme280.hConfig[5] -= 256
	}
	tB := []byte{0xF2, 0x01}
	err = bme280.Dev.Write(tB)
	if err != nil {
		fmt.Println("Humidity control error: ", err)
		return -256
	}
	tB[0] = 0xF4
	tB[1] = 0x27
	err = bme280.Dev.Write(tB)
	if err != nil {
		fmt.Println("Measurement mode set error: ", err)
		return -256
	}
	tB[0] = 0xF5
	tB[1] = 0xA0
	err = bme280.Dev.Write(tB)
	if err != nil {
		fmt.Println("Configuration write error: ", err)
		return -256
	}
	return 0
} /* bme280Init() */

//
// BME280ReadValues Reads the sensor register values
// and translate them into calibrated readings
// using the previously loaded calibration data
// Temperature is expressed in Celsius degrees as T * 100 (for 2 decimal places)
// Pressure is <future>
// Humidity is express as H * 1024 (10 bit fraction)
//
func (bme280 *BME280) BME280ReadValues() BMEData {
	data := BMEData{}
	ret := make([]byte, 8)
	// r := BMEData{}

	err := bme280.Dev.ReadReg(0xF7, ret)
	if err != nil {
		fmt.Println("Failed to read Data: ", err)
		return data
	}
	p := (int(ret[0]) << 12) + (int(ret[1]) << 4) + (int(ret[2]) >> 4)
	t := (int(ret[3]) << 12) + (int(ret[4]) << 4) + (int(ret[5]) >> 4)
	h := (int(ret[6]) << 8) + int(ret[7])
	//fmt.Println("Raw values: ", p, t, h)
	// Calculate calibrated temperature value
	// the value is 100x C (e.g. 2601 = 26.01C)
	var1 := (((t >> 3) - (bme280.tempressConfig[0] << 1)) * (bme280.tempressConfig[1])) >> 11
	var2 := (((((t >> 4) - (bme280.tempressConfig[1])) * ((t >> 4) - (bme280.tempressConfig[0]))) >> 12) * (bme280.tempressConfig[2])) >> 14
	tFine := var1 + var2
	T := (tFine*5 + 128) >> 8
	fmt.Println("Calibrated Temp: ", T)
	// Calculate calibrated humidity value
	var1 = (tFine - 76800)
	var1 = (((((h << 14) - ((bme280.hConfig[3]) << 20) - ((bme280.hConfig[4]) * var1)) +
		(16384)) >> 15) * (((((((var1*(bme280.hConfig[5]))>>10)*(((var1*(bme280.hConfig[2]))>>11)+(32768)))>>10)+(2097152))*(bme280.hConfig[1]) + 8192) >> 14))
	var1 = (var1 - (((((var1 >> 15) * (var1 >> 15)) >> 7) * (bme280.hConfig[0])) >> 4))
	if var1 < 0 {
		var1 = 0
	}
	if var1 > 419430400 {
		var1 = 419430400
	}
	H := var1 >> 12
	P := 0
	// Calculate calibrated pressure value
	var1_64 := uint64(tFine - 128000)
	var2_64 := uint64(var1_64 * var1_64 * uint64(bme280.pressConfig[5]))
	var2_64 = var2_64 + ((var1_64 * uint64(bme280.pressConfig[4])) << 17)
	var2_64 = var2_64 + ((uint64(bme280.pressConfig[3])) << 35)
	var1_64 = ((var1_64 * var1_64 * uint64(bme280.pressConfig[2])) >> 8) + ((var1_64 * uint64(bme280.pressConfig[1])) << 12)
	var1_64 = (((1) << 47) + var1_64) * (uint64(bme280.pressConfig[0])) >> 33
	if var1_64 == 0 {
		P = 0
	} else {
		p64 := uint64(1048576 - p)
		p64 = (((p64 << 31) - var2_64) * 3125) / var1_64
		var1_64 = ((uint64(bme280.pressConfig[8])) * (p64 >> 13) * (p64 >> 13)) >> 25
		var2_64 = ((uint64(bme280.pressConfig[7])) * p64) >> 19
		p64 = ((p64 + var1_64 + var2_64) >> 8) + ((uint64(bme280.pressConfig[6])) << 4)
		P = int(p64 / 100)
	}

	data.Temperature = T
	data.Pressure = P
	data.Humidity = H
	return data
}
