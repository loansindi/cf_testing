package main

import (
	"flag"
	"github.com/loansindi/cf_lcd"
	"github.com/tarm/serial"
	"log"
	"time"
)

func main() {
	const KEYUP = 0x01
	const KEYENTER = 0x02
	const KEYCANCEL = 0x04
	const KEYLEFT = 0x08
	const KEYRIGHT = 0x10
	const KEYDOWN = 0x20
	port := flag.String("serialport", "/dev/ttyUSB0", "Serial port to use")
	brightness := flag.Int("brightness", 100, "brightness?")
	clear := flag.Int("clear", 0, "clear?")
	message := flag.String("message", "", "etc")
	row := flag.Int("row", 0, "etc")
	col := flag.Int("col", 0, "etc")
	flag.Parse()
	sc := &serial.Config{Name: *port, Baud: 19200}
	s, err := serial.OpenPort(sc)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(300 * time.Millisecond)
	mask := make([]byte, 2)
	_ = cf_lcd.KeyReporting(s, mask)
	if *brightness != 100 {
		cf_lcd.Backlight(s, *brightness)
	}
	if *clear == 1 {
		cf_lcd.Clear(s)
	}
	err = cf_lcd.Write(s, *row, *col, *message)
	if err != nil {
		log.Fatal(err)
	}
	for {
		buttons := cf_lcd.GetKeys(s)
		log.Println(buttons)
		if buttons[4] == KEYUP {
			*brightness = *brightness + 10
			log.Println(*brightness)
			cf_lcd.Backlight(s, *brightness)
		}
		if buttons[4] == KEYDOWN {
			*brightness = *brightness - 10
			log.Println(*brightness)
			cf_lcd.Backlight(s, *brightness)
		}
		if buttons[4] == KEYLEFT {
			cf_lcd.Write(s, 0, 0, "hello")
		}
		if buttons[4] == KEYRIGHT {
			cf_lcd.Write(s, 1, 0, "how are you")
		}
		if buttons[4] == KEYENTER {
			cf_lcd.Clear(s)
		}
		time.Sleep(100 * time.Millisecond)
	}
}
