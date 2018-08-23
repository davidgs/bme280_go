# bme20_go

## A Native Golang driver for the Adafruit BME280 Breakout board.

## Usage

Create a BME Object:

    bme := BME280{}

Initialize it with the proper device (for the Raspberry Pi 3, this works):

    dev := "/dev/i2c-1"
	r := bme.bme280Init(dev)

Any return other than -256 is success.

To read values:

    readings := bme.bme280ReadValues()
	

readings will be a 3-value array of ints. Temperature, Pressure, Humidity, in that order. Temperature is the Temperature in C * 100 so to get the actual temperature:

    f := float64(float64(readings[0]) / 100.00)
	
Pressure is currently not properly implemented, and should not be iused. Humidity is returned as realtive humidity * 1024, so to get the humidity value:

    f = float64(rets[2]) / 1024.00

It's a good idea to call

    defer bme.Dev.Close()

after the call to init so that the device will be closed after you're done with it. 

## Error Handling

There is none at this point. You have to handle them. but if the return array has a -256 for a value in a spot, it's likely that an error occured. 