package cmd

import (
	"github.com/spf13/cobra"
	"github.com/bionicrm/controlifx"
	"github.com/pkg/errors"
	"math"
	"strconv"
	"bytes"
	"github.com/bionicrm/clifx/protocol"
	"encoding/json"
	"fmt"
	"log"
)

var (
	getServiceCmd = &cobra.Command{
		Use:"service",
		Short:"Acquires responses from all devices on the network",
		Run:func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetService())
		},
	}
	getHostInfoCmd = &cobra.Command{
		Use:"hostinfo",
		Short:"Gets the host MCU information",
		Run:func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetHostInfo())
		},
	}
	getHostFirmwareCmd = &cobra.Command{
		Use:"hostfirmware",
		Short:"Gets the host MCU firmware information",
		Run:func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetHostFirmware())
		},
	}
	getWifiInfoCmd = &cobra.Command{
		Use:"wifiinfo",
		Short:"Gets the Wi-Fi subsystem information",
		Run:func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetWifiInfo())
		},
	}
	getWifiFirmwareCmd = &cobra.Command{
		Use:"wififirmware",
		Short:"Gets the Wi-Fi subsystem firmware",
		Run:func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetWifiFirmware())
		},
	}
	powerCmd = &cobra.Command{
		Use:"power",
		Short:"Gets or sets the power level",
		ValidArgs:[]string{"[level]"},
		Run:func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				handle(true, controlifx.GetPower())
			} else {
				powerOn, err := argToBool(args[0])
				if err != nil {
					log.Fatalln(err)
				}

				level := 0
				if powerOn {
					level = 65535
				}

				handle(false, controlifx.SetPower(controlifx.SetPowerLanMessage{
					Level:controlifx.PowerLevel(level),
				}))
			}
		},
	}
	labelCmd = &cobra.Command{
		Use:"label",
		Short:"Gets or sets the label",
		ValidArgs:[]string{"[label]"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				handle(true, controlifx.GetLabel())
			} else {
				label := args[0]
				if len(label) > 32 {
					log.Fatalln("Label exceeds 32 characters")
				}

				handle(false, controlifx.SetLabel(controlifx.SetLabelLanMessage{
					Label:controlifx.Label(label),
				}))
			}
		},
	}
	getVersionCmd = &cobra.Command{
		Use:"version",
		Short:"Gets the hardware version",
		Run:func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetVersion())
		},
	}
	getInfoCmd = &cobra.Command{
		Use:"info",
		Short:"Gets the run-time information",
		Run:func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetInfo())
		},
	}
	getLocationCmd = &cobra.Command{
		Use:"location",
		Short:"Gets the location information",
		Run:func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetLocation())
		},
	}
	getGroupCmd = &cobra.Command{
		Use:"group",
		Short:"Gets the group membership information",
		Run:func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetGroup())
		},
	}
	echoRequestCmd = &cobra.Command{
		Use:"echo",
		Short:"Requests an arbitrary payload be echoed back",
		ValidArgs:[]string{"<payload>"},
		Run:func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				args = []string{""}
			}

			payload := args[0]
			if len(payload) > 64 {
				log.Fatalln("Payload exceeds 64 characters")
			}

			msgPayload := controlifx.EchoRequestLanMessage{}
			copy(msgPayload.Payload[:], []byte(payload))

			handle(true, controlifx.EchoRequest(msgPayload))
		},
	}
	lightGetCmd = &cobra.Command{
		Use:"lightget",
		Short:"Gets the light state",
		Run:func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.LightGet())
		},
	}
	lightSetColorCmd = &cobra.Command{
		Use:"lightcolor",
		Short:"Sets the light color",
		Long:"Sets the light color to a hex value if 1 argument is supplied, HSBK if 3 or 4 are supplied, or RGB if --rgb is set and 3 arguments are supplied",
		ValidArgs:[]string{"<hex> | <hue> <saturation> <brightness> [kelvin] | <red> <green> <blue>"},
		Run:func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatalln("No color supplied")
			}

			const DefaultKelvin = 2500

			msgPayload := controlifx.LightSetColorLanMessage{
				Duration:uint32(duration),
			}
			setMsgPayloadColor := func(h, s, l, k uint16) {
				msgPayload.Color = controlifx.HSBK{
					Hue:h,
					Saturation:s,
					Brightness:l,
					Kelvin:k,
				}
			}

			if len(args) == 1 {
				// Hex.
				r, err := strconv.ParseUint(string(args[0][:2]), 16, 8)
				if err != nil {
					log.Fatalln(err)
				}

				g, err := strconv.ParseUint(string(args[0][2:4]), 16, 8)
				if err != nil {
					log.Fatalln(err)
				}

				b, err := strconv.ParseUint(string(args[0][4:]), 16, 8)
				if err != nil {
					log.Fatalln(err)
				}

				h, s, l := rgbToHsl(uint8(r), uint8(g), uint8(b))

				setMsgPayloadColor(h, s, l, DefaultKelvin)
			} else if rgb {
				// RGB.
				if len(args) != 3 {
					log.Fatalln("Red, green, and/or blue not supplied")
				}

				r, err := strconv.Atoi(args[0])
				if err != nil {
					log.Fatalln(err)
				}

				g, err := strconv.Atoi(args[1])
				if err != nil {
					log.Fatalln(err)
				}

				b, err := strconv.Atoi(args[2])
				if err != nil {
					log.Fatalln(err)
				}

				h, s, l := rgbToHsl(uint8(r), uint8(g), uint8(b))

				setMsgPayloadColor(h, s, l, DefaultKelvin)
			} else if len(args) == 3 || len(args) == 4 {
				// HSBK.
				h, err := strconv.Atoi(args[0])
				if err != nil {
					log.Fatalln(err)
				}

				s, err := strconv.Atoi(args[1])
				if err != nil {
					log.Fatalln(err)
				}

				l, err := strconv.Atoi(args[2])
				if err != nil {
					log.Fatalln(err)
				}

				k := DefaultKelvin

				if len(args) == 4 {
					k, err = strconv.Atoi(args[3])
					if err != nil {
						log.Fatalln(err)
					}
					if k < 2500 || k > 9000 {
						log.Fatalln("Color temperature (Kelvin) is out of the range 2500-9000")
					}
				}

				setMsgPayloadColor(lerpToUint16(360, h), lerpToUint16(100, s), lerpToUint16(100, l), uint16(k))
			} else {
				log.Fatalln("Invalid color supplied")
			}

			handle(false, controlifx.LightSetColor(msgPayload))
		},
	}
	lightPowerCmd = &cobra.Command{
		Use:"lightpower",
		Short:"Gets or sets the light power level",
		ValidArgs:[]string{"[on|off]"},
		Run:func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				handle(true, controlifx.LightGetPower())
			} else {
				powerOn, err := argToBool(args[0])
				if err != nil {
					log.Fatalln(err)
				}

				level := 0
				if powerOn {
					level = 65535
				}

				handle(false, controlifx.LightSetPower(controlifx.LightSetPowerLanMessage{
					Level:controlifx.PowerLevel(level),
					Duration:uint32(duration),
				}))
			}
		},
	}

	// Flags.

	rgb      bool
	duration int
)

func init() {
	RootCmd.AddCommand(
		getServiceCmd,
		getHostInfoCmd,
		getHostFirmwareCmd,
		getWifiInfoCmd,
		getWifiFirmwareCmd,
		powerCmd,
		labelCmd,
		getVersionCmd,
		getInfoCmd,
		getLocationCmd,
		getGroupCmd,
		echoRequestCmd,
		lightGetCmd,
		lightSetColorCmd,
		lightPowerCmd,
	)

	lightSetColorCmd.Flags().BoolVar(&rgb, "rgb", false, "specifies that the color values are in red, green, and blue form")
	lightSetColorCmd.Flags().IntVarP(&duration, "duration", "d", 0, "specifies the duration of the color transition in milliseconds")

	lightPowerCmd.Flags().IntVarP(&duration, "duration", "d", 0, "specifies the duration of the power transition in milliseconds")
}

func rgbToHsl(rI, gI, bI uint8) (hI, sI, lI uint16) {
	var (
		r = float64(rI)/255
		g = float64(gI)/255
		b = float64(bI)/255

		min   float64
		max   float64
		delta float64
	)

	if r < g && r < b {
		min = r
	} else if g < r && g < b {
		min = g
	} else {
		min = b
	}

	if r > g && r > b {
		max = r
		delta = max-min
		hI = uint16((math.Mod((g-b)/delta, 6))/6*0xffff)
	} else if g > r && g > b {
		max = g
		delta = max-min
		hI = uint16(((b-r)/delta+2)/6*0xffff)
	} else {
		max = b
		delta = max-min
		hI = uint16(((r-g)/delta+4)/6*0xffff)
	}

	l := (max+min)/2

	if delta == 0 {
		sI = 0
	} else {
		sI = uint16(delta/(1-math.Abs(2*l-1))*0xffff)
	}

	lI = uint16(l*0xffff)

	return
}

func argToBool(arg string) (bool, error) {
	switch arg {
	case "1":
		fallthrough
	case "on":
		return true, nil
	case "0":
		fallthrough
	case "off":
		return false, nil
	default:
		return false, errors.New("Must be 'on' or 'off'")
	}
}

func lerpToUint16(rng float64, x int) uint16 {
	return uint16(float64(x)/rng*math.MaxUint16+0.5)
}

func handle(requireResByDef bool, msg controlifx.SendableLanMessage) {
	conn, err := controlifx.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	discoveryNeeded := len(labels) > 0 || len(groups) > 0 || len(macs) > 0 || len(ips) > 0 || count > 0
	resCode := getResponseCode(msg)

	if requireAck || ((requireResByDef || requireRes) && resCode == controlifx.AcknowledgementType) {
		msg.Header.FrameAddress.AckRequired = true
		resCode = controlifx.AcknowledgementType
	}

	waitForResponses := requireAck || requireResByDef || requireRes

	if discoveryNeeded {
		devices, err := protocol.Discover(conn, labels, groups, macs, ips, timeout, count)
		if err != nil {
			log.Fatalln(err)
		}

		if waitForResponses {
			recMsgs, err := conn.SendToAndGet(msg, devices, controlifx.TypeFilter(resCode))
			if err != nil {
				log.Fatalln(err)
			}

			printResponses(recMsgs)
		} else if err := conn.SendTo(msg, devices); err != nil {
			log.Fatalln(err)
		}
	} else if waitForResponses {
		recMsgs, err := conn.SendToAllAndGet(timeout, msg, controlifx.TypeFilter(resCode))
		if err != nil {
			log.Fatalln(err)
		}

		printResponses(recMsgs)
	} else if err := conn.SendToAll(msg); err != nil {
		log.Fatalln(err)
	}
}

func getResponseCode(msg controlifx.SendableLanMessage) uint16 {
	switch msg.Header.ProtocolHeader.Type {
	case controlifx.GetServiceType:
		return controlifx.StateServiceType
	case controlifx.GetHostInfoType:
		return controlifx.StateHostInfoType
	case controlifx.GetHostFirmwareType:
		return controlifx.StateHostFirmwareType
	case controlifx.GetWifiInfoType:
		return controlifx.StateWifiInfoType
	case controlifx.GetWifiFirmwareType:
		return controlifx.StateWifiFirmwareType
	case controlifx.GetPowerType:
		return controlifx.StatePowerType
	case controlifx.SetPowerType:
		// ACK only.
		return controlifx.AcknowledgementType
	case controlifx.GetLabelType:
		return controlifx.StateLabelType
	case controlifx.SetLabelType:
		// ACK only.
		return controlifx.AcknowledgementType
	case controlifx.GetVersionType:
		return controlifx.StateVersionType
	case controlifx.GetInfoType:
		return controlifx.StateInfoType
	case controlifx.GetLocationType:
		return controlifx.StateLocationType
	case controlifx.GetGroupType:
		return controlifx.StateGroupType
	case controlifx.EchoRequestType:
		return controlifx.EchoResponseType
	case controlifx.LightGetType:
		return controlifx.LightStateType
	case controlifx.LightSetColorType:
		// As per the protocol.
		return controlifx.LightStateType
	case controlifx.LightGetPowerType:
		return controlifx.LightStatePowerType
	case controlifx.LightSetPowerType:
		// As per the protocol.
		return controlifx.LightStatePowerType
	}

	return 0
}

func printResponses(recMsgs map[controlifx.Device]controlifx.ReceivableLanMessage) {
	var (
		responses []interface{}
		err       error
		b         []byte
	)

	for device, msg := range recMsgs {
		if requireAck {
			msg.Payload = nil
		}

		responses = append(responses, struct{
			Device   controlifx.Device
			Response interface{} `json:",omitempty"`
		}{
			Device:device,
			Response:createFriendlyPayload(msg.Payload),
		})
	}

	if pretty {
		b, err = json.MarshalIndent(responses, "", "  ")
	} else {
		b, err = json.Marshal(responses)
	}
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(b))
}

func createFriendlyPayload(payload interface{}) interface{} {
	switch payload.(type) {
	case *controlifx.EchoResponseLanMessage:
		return struct{
			Payload string
		}{
			Payload:string(bytes.TrimRight(payload.(*controlifx.EchoResponseLanMessage).Payload[:], "\x00")),
		}
	}

	return payload
}
