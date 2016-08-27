# Clifx
Command-line interface for LIFX device control.

Clifx is a simple yet powerful command line utility for controlling LIFX devices over the LAN, or even your local machine (see [Implifx](https://github.com/golifx/implifx)). As such, every possible device message is able to be sent and received, along with their appropriate payloads. Commands are simple. For example, to set the color of the lights in your bedroom to a nice green shade specified in RGB, you can use `clifx -g bedroom lightcolor --rgb 75 221 88`

In addition, responses are able to be received and printed in JSON form. This makes it easy to parse the responses sent back by devices.

**Contents:**
- [Installation](#installation)
- [Usage](#usage)
  - [Choosing devices](#choosing-devices)
    - [By label](#by-label)
    - [By group](#by-group)
    - [By MAC](#by-mac)
    - [By IP](#by-ip)
    - [Count](#count)
  - [Light color](#light-color)
    - [Color temperature](#color-temperature)
    - [Hex](#hex)
    - [RGB](#rgb)
    - [HSL(K)](#hslk)
  - [Getting responses](#getting-responses)
    - [Acknowledgements](#acknowledgements)
    - [Pretty printing](#pretty-printing)
    - [Forcing responses](#forcing-responses)
- [Additional Help](#additional-help)

## Installation
If you have Go installed and `$GOPATH/bin` in your path, just run `go get -u github.com/golifx/clifx` and the `clifx` binary will be available. Otherwise, [download the latest release](https://github.com/golifx/clifx/releases) for your platform, unarchive it, and move the binary to some location in your path (try `/usr/bin/`).

## Usage
#### Choosing devices
When these flags are not used, the message you send will be emitted to every device on the LAN. However, if you only want to target certain devices, you can filter them by their label, group, MAC address, IP address, or simply by the number of devices to receive it. All of these options can be combined to be as specific or as general as you want.

###### By label
This gets the power state of devices with a label of "nightstand", case insensitive:

```bash
$ clifx -l nightstand power
```

###### By group
This gets the power state of devices belonging to the group "bedroom", case insensitive:

```bash
$ clifx -g bedroom power
```

###### By MAC
This gets the power state the device with the given MAC address:

```bash
$ clifx -m 4d:bc:12:d5:73:d0 power
```

###### By IP
This gets the power state of the device with the given IP address:

```bash
$ clifx -i 10.0.0.23 power
```

###### Count
Only a certain number of devices will receive this message. If you have more devices on the network than this number, the devices that the message is sent to is nondeterministic, as it depends on which devices reply to the initial discover request first:

```bash
$ clifx -c 3 power
```

#### Light color
You can set the color of your lights to a hex, RGB, or HSL value with a color temperature. All of the following examples will set the color of your lights to the same shade of yellow.

###### Color temperature
This example uses a hex triplet for specifying the color, along with the `-k` option to specify the color temperature. This value must be within the range 2500..9000, inclusive. The lower the saturation, the more the color temperature comes into play:

```bash
$ clifx lightcolor d0bb71 -k 2500
$ clifx lightcolor d0bb71 -k 9000
```

###### Hex
```bash
$ clifx lightcolor ffd642
```

###### RGB
When the `--rgb` flag is used, instead of parsing the 3 arguments as HSL (hue, saturation, and lightness), they'll be recognized as red, green, and blue, in that order. All 3 values must be within the range 0..255, inclusive:

```bash
$ clifx lightcolor --rgb 255 214 66
```

###### HSL(K)
The first value is the hue, which is within the range 0..360, inclusive. The saturation and lightness come next, in order, which are in the range 0..100, inclusive.

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

#### Getting responses
Sometimes you may only care about getting an acknowledgement that a message was received. Or, you may wish to have the JSON output from either an acknowledgement or response to be pretty printed. Also, some commands, like `power <on|off>`, won't return a response unless you ask for one. Here's how to do all of that...

###### Acknowledgements
An acknowledgement is simply a response sent from the device to *acknowledge* that it received and processed a message. The `-a` flag requires such a response, and only prints out the device address and its MAC address. This flag *cannot* be used in addition to the `-r` flag, as mentioned below:

```bash
$ clifx -a power on
```

Example output:
```
[{"Device":{"Addr":{"IP":"10.0.0.23","Port":56700,"Zone":""},"Mac":85470165169104}}]
```

###### Pretty printing
Normally, `clifx power` would print something like this:
```
[{"Device":{"Addr":{"IP":"10.0.0.23","Port":56700,"Zone":""},"Mac":85470165169104},"Response":{"Level":65535}}]
```

However, and especially when you have many devices on the network, it would be easier to read if the JSON was formatted. Use the `-p` flag to pretty print:

```bash
$ clifx -p power
```

Example output:
```
[
  {
    "Device": {
      "Addr": {
        "IP": "10.0.0.23",
        "Port": 56700,
        "Zone": ""
      },
      "Mac": 85470165169104
    },
    "Response": {
      "Level": 65535
    }
  }
]

```

###### Forcing responses
You'll notice that `clifx power on` won't give you a response. To force the device to send one back, use the `-r` flag. This flag *cannot* be used in addition to the `-a` flag, as mentioned above:

```bash
$ clifx -r power on
```

Example output:
```
[{"Device":{"Addr":{"IP":"10.0.0.23","Port":56700,"Zone":""},"Mac":85470165169104},"Response":{"Level":0}}]
```

## Additional Help
Run `clifx` or `clifx <command> --help` to learn more about all or one of the commands, flags, and options. Visit [#golifx](http://webchat.freenode.net?randomnick=1&channels=%23golifx&prompt=1) on chat.freenode.net to get help, ask questions, or discuss ideas.
