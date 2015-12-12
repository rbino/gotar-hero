package main

import (
    "fmt"
    "flag"
    "github.com/gvalkov/golang-evdev"
    "github.com/rakyll/portmidi"
)

const (
  butGreen = 293
  butRed = 289
  butYellow = 288
  butBlue = 290
  butOrange = 291
  butTilt = 292
  butStrumbar = 1
  butStart = 297
  butSelect = 296
)

var (
  butMap map[uint16]uint32 = map [uint16]uint32 {
    butGreen : 0,
    butRed : 1,
    butYellow : 2,
    butBlue : 3,
    butOrange : 4,
  }
  playing int64 = -1
  octave int64 = 0
  hold bool = false
  butState int64 = 0
  baseNote int64
  vel int64
  midiChan int64
)

func ReadGuitar(c chan *evdev.InputEvent, errs chan error, dev string){
  guit, err := evdev.Open(dev)
  if err != nil{
    errs <- err
    return
  }
  for {
    ev, err := guit.ReadOne()
    if err != nil{
      errs <- err
      return
    }
    if ev.Type == evdev.EV_KEY || ev.Type == evdev.EV_ABS {
      c <- ev
    }
  }
}

func SetNote(shift uint32, value int32) {
  if value != 0 {
    butState |= 1 << shift
  } else {
    butState &= ^(1 << shift)
  }
}

func NoteOn(stream *portmidi.Stream, note int64, velocity int64){
  playing = note
  stream.WriteShort(0x90|(midiChan-1), note, velocity)
}

func NoteOff(stream *portmidi.Stream, note int64, velocity int64){
  stream.WriteShort(0x80|(midiChan-1), note, velocity)
  playing = -1
}

func SwapNote(stream *portmidi.Stream, newNote int64, velocity int64){
  oldNote := playing
  NoteOff(stream, oldNote, velocity)
  NoteOn(stream, newNote, velocity)
  playing = newNote
}

func main() {
  dev := flag.String("d", "/dev/input/by-id/usb-0810_Twin_USB_Joystick-event-joystick", "The GH controller event device")
  flag.Int64Var(&baseNote, "b", 48, "The base midi note with no button pressed")
  flag.Int64Var(&vel, "v", 100, "Midi note velocity")
  flag.Int64Var(&midiChan, "c", 1, "Midi channel")
  flag.Parse()
  portmidi.Initialize()
  out, err := portmidi.NewOutputStream(portmidi.GetDefaultOutputDeviceId(), 32, 0)
  if err != nil{
    fmt.Println(err)
    return
  }

  c := make(chan *evdev.InputEvent)
  e := make(chan error)
  go ReadGuitar(c, e, *dev)

  for {
    select{
    case but := <- c:
      switch but.Code {
      case butGreen, butRed, butYellow, butBlue, butOrange:
        SetNote(butMap[but.Code], but.Value)
        if playing != -1 {
          SwapNote(out, baseNote + butState + octave, vel)
        }
      case butStrumbar:
        if but.Value == 255 || but.Value == 0 {
          NoteOn(out, baseNote + butState + octave, vel)
        } else if !hold {
          NoteOff(out, playing, vel)
        }
      case butSelect:
        if but.Value != 0 {
          hold = !hold
          if (!hold){
            NoteOff(out, playing, vel)
          }
        }
      case butTilt:
        if but.Value == 1 {
          octave += 12
        } else {
          octave -= 12
        }
        if playing != -1 {
          SwapNote(out, baseNote + butState + octave, vel)
        }
      }
    case err := <- e:
      fmt.Println(err)
      close(c)
      close(e)
      return
    default:
    }
  }
}
