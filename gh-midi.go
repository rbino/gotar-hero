package main

import (
    "fmt"
    "flag"
    "github.com/gvalkov/golang-evdev"
)

func ReadGuitar(c chan *evdev.InputEvent, dev string){
  guit, _ := evdev.Open(dev)
  for {
    ev, _ := guit.ReadOne()
    if ev.Type == evdev.EV_KEY || ev.Type == evdev.EV_ABS {
      c <- ev
    }
  }
}

func main() {
  dev := flag.String("d", "/dev/input/by-id/usb-0810_Twin_USB_Joystick-event-joystick", "The event device associated with the GH controller")
  flag.Parse()

  c := make(chan *evdev.InputEvent)
  go ReadGuitar(c, *dev)
  for {
    but := <- c
    fmt.Println(but.Code, but.Value)
  }
}
