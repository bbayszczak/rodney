# Rodney

![golangci-lint](https://github.com/bbayszczak/rodney/workflows/golangci-lint/badge.svg)
[![Build Status](https://travis-ci.com/bbayszczak/rodney.svg?token=AWkyENePdvxphuA78oxv&branch=main)](https://travis-ci.com/bbayszczak/rodney)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Rodney is a 2WD robot controlled by a bluetooth controller

## Specs

  - 2WD
  - controlled by a bluetooth controller
  - yellow LED on power on
  - blue LED for bluetooth: blink when not paired, still when paired
  - white LED when process started/running
  - red LED when issue
  - when issue, restart process

## Components

  - 2WD robot chassis
  - raspberry pi
  - USB battery pack x1
  - L293D x1
  - 220Ω resistor x4
  - yellow LED x1
  - white LED x1
  - blue LED x1
  - red LED x1

## requirements

  - should work with Raspberry Pi 2, zero, 3 & 4 but tested only on zero
  - go >= 1.14

## License

MIT License
