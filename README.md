# bme20_go

## A Native Golang driver for the Adafruit BME280 Breakout board.

## Usage

    func main() {
	    dev := "/dev/i2c-1"
	    bme := BME280{}
	    r := bme.bme280Init(dev)
	    if r < 0 {
		    fmt.Println("Error")
	    }
	    rets := bme.bme280ReadValues()
	    f := float64(float64(rets[0]) / 100.00)
	    fmt.Println("Temp: ", f)
	    f = float64(rets[2]) / 1024.00
	    fmt.Println("Humidity: ", f)
	    fmt.Println("Pressure: ", rets[1])
	    bme.dev.Close()
    }

