# Rodney

![golangci-lint](https://github.com/bbayszczak/rodney/workflows/golangci-lint/badge.svg)
[![Build Status](https://travis-ci.com/bbayszczak/rodney.svg?token=AWkyENePdvxphuA78oxv&branch=main)](https://travis-ci.com/bbayszczak/rodney)
[![Go Report Card](https://goreportcard.com/badge/github.com/bbayszczak/rodney)](https://goreportcard.com/report/github.com/bbayszczak/rodney)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Rodney is a 2WD robot controlled by a Nintendo Switch Pro bluetooth controller

![Rodney](media/rodney.gif)

## Specs

  - 2WD
  - controlled by a Nintendo Switch Pro bluetooth controller
  - yellow LED on power on
  - blue LED for bluetooth: blink when not paired, still when paired
  - white LED when process started/running
  - red LED obstacle too close

## Components

  - 2WD robot chassis
  - Raspberry Pi with Bluetooth
  - USB battery pack x1
  - L293D x1
  - HCSR04 range sensor x1
  - 220Ω resistor x4
  - 1kΩ resistor x1
  - 2kΩ resistor x1
  - yellow LED x1
  - white LED x1
  - blue LED x1
  - red LED x1

## requirements

  - should work with any Raspberry Pi with Bluetooth but tested only on zero WH
  - go >= 1.14

## License

MIT License
