# Clifx

Command line LIFX device control.

### Installation

If you have Go installed and `$GOPATH/bin` in your path, just run
`go get github.com/bionicrm/clifx` and the `clifx` binary will be available.
Otherwise, download the latest release for your platform and install it where
it'll be executable.

### Introduction

Clifx is a very low level yet powerful command line utility for controlling LIFX
devices over the LAN. As such, every possible device message is able to be sent
and received, along with their appropriate payloads. This means you won't be
seeing commands like `clifx power on floor` or anything like that, but rather
`clifx --label nightstand --type SetPower --payload Level:65535` (powers on the
device with the 'nightstand' label).

In addition, responses are able to be received and outputted in JSON form. This
makes it easy to parse the responses sent back by devices.

While the commands may seem complicated, they all make sense and are able to be
easily constructed with a quick look at the documentation. This also makes it
easy for other programs (those not written in Go and therefore ones not able to
use the API) to execute the command and receive responses.

### Usage

Both the official [Devices Messages](https://lan.developer.lifx.com/docs/device-messages)
and [Light Messages](https://lan.developer.lifx.com/docs/light-messages) docs
contain all possible message types (specified with `--type`) and payload fields
(given by `--payload`), along with their descriptions. However, a bunch of
examples are listed below.

- [Specify devices](#specify-devices)
- [Set light color](#set-light-color)

##### [Specify devices](#specify-devices)

You can choose which devices will receive the message. This can be done either
by IP, MAC, label, or simply a count of how many should receive it.

###### By IP

```bash
$ clifx --type LightSetColor --payload Color:Hue:28399,Color:Saturation:43253,Color:Brightness:55000,Color:Kelvin:2500,Duration:2500
```

##### [Set light color](#set-light-color)

This sets the color of all lights to a hue of 156, saturation of 66%, brightness
of 84%, and a duration of 2.5 seconds.

```bash
$ clifx --type LightSetColor --payload Color:Hue:28399,Color:Saturation:43253,Color:Brightness:55000,Color:Kelvin:2500,Duration:2500
```
