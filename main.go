package main

import (
	"flag"
	"fmt"
	"github.com/loansindi/cf_lcd"
	"github.com/tarm/serial"
	"log"
	"sync"
	"time"
)

func pollKeys(p *serial.Port, c chan byte) {
	for {
		buttons, err := cf_lcd.GetKeys(p)
		if err != nil {
			log.Print("Broke in pollkeys")
			log.Fatal(err)
		}
		if buttons != nil {
			c <- buttons[1]
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func drawScreen(p *serial.Port, screen chan int, bright chan int) {
	for {
		select {
		case menu := <-screen:
			log.Print("Menu case")
			if menu == 1 {
				log.Print("Screen 1")
				brightness := <-bright
				cf_lcd.Clear(p)
				cf_lcd.Write(p, 0, 0, "Brightness")
				cf_lcd.Write(p, 1, 0, fmt.Sprintf("%d", brightness))
			}
			if menu == 2 {
				log.Print("Screen 2")
				_ = <-bright
				cf_lcd.Clear(p)
				cf_lcd.Write(p, 1, 0, "Second Screen")
			}
			if menu == 3 {
				log.Print("Screen 3")
				_ = <-bright
				cf_lcd.Clear(p)
				cf_lcd.Write(p, 0, 0, "Third")
				cf_lcd.Write(p, 1, 0, "Screen")
			}
		}
	}
}

func handleButtons(s *serial.Port, c chan byte, screen chan int, bright chan int) {
	const KEYUP = 0x01
	const KEYENTER = 0x02
	const KEYCANCEL = 0x04
	const KEYLEFT = 0x08
	const KEYRIGHT = 0x10
	const KEYDOWN = 0x20
	brightness := 100
	page := 1
	for {
		select {
		case but := <-c:
			if but != 0 {
				log.Print("Button: ", but)
			}
			if but == KEYUP {
				if page > 1 {
					page--
				}
				log.Print("KEYUP")
				bright <- brightness
				screen <- page
			}

			if but == KEYDOWN {
				if page < 3 {
					page++
				}
				bright <- brightness
				screen <- page
			}

			if but == KEYLEFT && page == 1 {
				if brightness > 0 {
					brightness = brightness - 10
					cf_lcd.Backlight(s, brightness)
					screen <- page
					bright <- brightness
				}
			}
			if but == KEYRIGHT && page == 1 {
				if brightness < 100 {
					brightness = brightness + 10
					cf_lcd.Backlight(s, brightness)
					screen <- page
					bright <- brightness
				}
			}

		}
	}
}
func main() {
	port := flag.String("serialport", "/dev/ttyUSB0", "Serial port to use")
	flag.Parse()
	sc := &serial.Config{Name: *port, Baud: 19200}
	s, err := serial.OpenPort(sc)
	if err != nil {

		log.Fatal(err)
	}
	time.Sleep(300 * time.Millisecond)
	mask := make([]byte, 2)
	_, err = cf_lcd.KeyReporting(s, mask)
	c := make(chan byte)
	go pollKeys(s, c)
	screen := make(chan int)
	bright := make(chan int, 1)
	go drawScreen(s, screen, bright)
	go handleButtons(s, c, screen, bright)
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

}
