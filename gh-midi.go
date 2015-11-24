package main

import (
    "fmt"
    "github.com/gvalkov/golang-evdev"
)

func ReadGuitar(c chan *evdev.InputEvent){
  guit, _ := evdev.Open("/dev/input/by-id/usb-0810_Twin_USB_Joystick-event-joystick")
  for {
    ev, _ := guit.ReadOne()
    if ev.Type == evdev.EV_KEY || ev.Type == evdev.EV_ABS {
      c <- ev
    }
  }
}

func main() {
  c := make(chan *evdev.InputEvent)
  go ReadGuitar(c)
  for {
    but := <- c
    fmt.Println(but.Code, but.Value)
  }
}
