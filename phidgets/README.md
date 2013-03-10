# phidgets

phidgets is a library for using the awesome USB devices from [Phidgets](http://www.phidgets.com/).

## Description

phidgets wraps the Phidgets C library. You must install the [drivers for your operating system](http://www.phidgets.com/docs/Operating_System_Support).

## Installation

    $ go get github.com/samuelkadolph/go/phidgets

## Usage

```go
package main

import (
  "github.com/samuelkadolph/go/phidgets"
  "log"
  "time"
)

func main() {
  ir, err := phidgets.NewIR()
  if err != nil {
    log.Fatalf("%s", err.Error())
  }

  go func() {
    <-ir.Attached
    log.Printf("Attached")

    code := []byte{0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80}

    for {
      if err := ir.Transmit(code, phidgets.IRCodeInfo{BitCount: len(code) * 8}); err != nil {
        log.Printf("%s", err.Error())
      }
      time.Sleep(500 * time.Millisecond)
    }
  }()

  if err := ir.Open(phidgets.Label{"foosfan1"}); err != nil {
    log.Fatalf("%s", err.Error())
  }

  select {}
}
```

## Implemented Phidgets

* InterfaceKit
* IR

## Contributing

Fork, branch & pull request.
