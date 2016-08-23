# Clifx
Command line LIFX device control. Uses [Controlifx](https://github.com/bionicrm/controlifx).

### Installation
If you have Go installed and `$GOPATH/bin` in your path, just run
`go get github.com/bionicrm/clifx` and the `clifx` binary will be available.
Otherwise, [download the latest release](https://github.com/bionicrm/clifx/releases) for your platform and install it where
it'll be executable.

### Introduction
Clifx is a simple yet powerful command line utility for controlling LIFX devices over the LAN. As such, every possible device message is able to be sent and received, along with their appropriate payloads. Commands are simple! For example, to set the color of the lights in your bedroom to a nice green shade specified in RGB, you can use `clifx -g bedroom lightcolor --rgb 75 221 88`

In addition, responses are able to be received and printed in JSON form. This makes it easy to parse the responses sent back by devices.

### Usage
- [Specify devices](#specify-devices)
  - By label
  - By group
  - By MAC
  - By IP
  - Count
- [Set light color](#set-light-color)
  - Hex
  - RGB
  - HSL(K)

##### [Specify devices](#specify-devices)
You can choose which devices will receive the message. This can be done either by label, group, IP, MAC, or simply a count of how many should receive it. All of the filters can be combined to be as specific or as general as you want.

###### By label
```bash
$ clifx -l nightstand power
```

###### By group
```bash
$ clifx -g bedroom power
```

###### By MAC
```bash
$ clifx -m 4d:bc:12:d5:73:d0 power
```

###### By IP
```bash
$ clifx -i 10.0.0.23 power
```

###### Count
This will send the message to only a certain number of devices. If you have more devices on the network than this number, the devices that the message is sent to is nondeterministic, as it depends on which devices reply to the initial discover request first.

```bash
$ clifx -c 3 power
```

##### [Set light color](#set-light-color)
You can set the color of your lights to a hex, RGB, or HSL value with a color temperature. All of the following examples will set the color of your lights to the same shade of yellow.

###### Hex
```bash
$ clifx lightcolor ffd642
```

###### RGB
```bash
$ clifx lightcolor --rgb 255 214 66
```

###### HSL(K)
The color temperature (Kelvin) is optional. This value becomes more noticeable as the saturation becomes lower.

Without a color temperature:
```bash
$ clifx lightcolor 47 100 63
```

With a color temperature:

```bash
$ clifx lightcolor 47 50 63 2500
$ clifx lightcolor 47 50 63 9000
```
