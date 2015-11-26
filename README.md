# Gotar Hero

This is a simple program that converts Guitar Hero PS2 controller events to MIDI notes.

## Dependencies

```
go get github.com/gvalkov/golang-evdev
go get github.com/rakyll/portmidi
```

Be aware that some distributions (for example Debian Testing) have an old version of libportmidi-dev, so you have to compile it to install the Go bindings.

## Usage
Run `gotar-hero -h` for the command line parameters.

The program was tested with a PS2 controller using a PS2 Port to USB adapter on Linux.

The current note mapping is defined as a basenote + offset, where offset is the number obtained taking the five Fret Buttons as bits (Green is the LSB, Orange is the MSB).  
The Select button toggles hold-mode.  
If the guitar is not in hold-mode, the Strum Bar sends a Note On when it's up or down and a Note Off when it's in neutral position.  
In hold-mode, only Note On is sent while strumming (useful if you want to use the other hand to control something else).  
In both modes a change in the Fret Buttons while playing a note will send a Note On with the new note followed by a Note Off for the old note (the program is written with monophonic synths in mind).
