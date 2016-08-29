package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bionicrm/clifx/protocol"
	"github.com/spf13/cobra"
	"gopkg.in/golifx/controlifx.v1"
	"log"
	"math"
	"net"
	"strconv"
)

var (
	getServiceCmd = &cobra.Command{
		Use:   "service",
		Short: "acquires responses from all devices on the network",
		Run: func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetService())
		},
	}
	getHostInfoCmd = &cobra.Command{
		Use:   "hostinfo",
		Short: "gets the host MCU information",
		Run: func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetHostInfo())
		},
	}
	getHostFirmwareCmd = &cobra.Command{
		Use:   "hostfirmware",
		Short: "gets the host MCU firmware information",
		Run: func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetHostFirmware())
		},
	}
	getWifiInfoCmd = &cobra.Command{
		Use:   "wifiinfo",
		Short: "gets the Wi-Fi subsystem information",
		Run: func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetWifiInfo())
		},
	}
	getWifiFirmwareCmd = &cobra.Command{
		Use:   "wififirmware",
		Short: "gets the Wi-Fi subsystem firmware",
		Run: func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetWifiFirmware())
		},
	}
	powerCmd = &cobra.Command{
		Use:       "power",
		Short:     "gets or sets the power level",
		ValidArgs: []string{"[level]"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				handle(true, controlifx.GetPower())
			} else {
				powerOn, err := argToBool(args[0])
				if err != nil {
					log.Fatalln(err)
				}

				var level uint16
				if powerOn {
					level = 0xffff
				}

				handle(false, controlifx.SetPower(controlifx.SetPowerLanMessage{
					Level: level,
				}))
			}
		},
	}
	labelCmd = &cobra.Command{
		Use:       "label",
		Short:     "gets or sets the label",
		ValidArgs: []string{"[label]"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				handle(true, controlifx.GetLabel())
			} else {
				label := args[0]
				if len(label) > 32 {
					log.Fatalln("Label exceeds 32 characters")
				}

				handle(false, controlifx.SetLabel(controlifx.SetLabelLanMessage{
					Label: label,
				}))
			}
		},
	}
	getVersionCmd = &cobra.Command{
		Use:   "version",
		Short: "gets the hardware version",
		Run: func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetVersion())
		},
	}
	getInfoCmd = &cobra.Command{
		Use:   "info",
		Short: "gets the run-time information",
		Run: func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetInfo())
		},
	}
	getLocationCmd = &cobra.Command{
		Use:   "location",
		Short: "gets the location information",
		Run: func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetLocation())
		},
	}
	getGroupCmd = &cobra.Command{
		Use:   "group",
		Short: "gets the group membership information",
		Run: func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetGroup())
		},
	}
	getOwnerCmd = &cobra.Command{
		Use: "owner",
		Short: "*undocumented*",
		Run: func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.GetOwner())
		},
	}
	echoRequestCmd = &cobra.Command{
		Use:       "echo",
		Short:     "requests an arbitrary payload be echoed back",
		ValidArgs: []string{"<payload>"},
		Run: func(cmd *cobra.Command, args []string) {
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
		Use:   "lightget",
		Short: "gets the light state",
		Run: func(cmd *cobra.Command, args []string) {
			handle(true, controlifx.LightGet())
		},
	}
	lightSetColorCmd = &cobra.Command{
		Use:       "lightcolor",
		Short:     "sets the light color",
		Long:      "Sets the light color to a hex value if 1 argument is supplied, HSBK if 3 or 4 are supplied, or RGB if --rgb is set and 3 arguments are supplied",
		ValidArgs: []string{"<hex> | <hue> <saturation> <brightness> [kelvin] | <red> <green> <blue>"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatalln("No color supplied")
			}

			validateKelvin := func() {
				if kelvin < 2500 || kelvin > 9000 {
					log.Fatalln("Color temperature (Kelvin) is out of the range 2500-9000")
				}
			}

			validateKelvin()

			msgPayload := controlifx.LightSetColorLanMessage{
				Duration: uint32(duration),
			}
			setMsgPayloadColor := func(h, s, l uint16) {
				msgPayload.Color = controlifx.HSBK{
					Hue:        h,
					Saturation: s,
					Brightness: l,
					Kelvin:     uint16(kelvin),
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

				setMsgPayloadColor(h, s, l)
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

				setMsgPayloadColor(h, s, l)
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

				if len(args) == 4 {
					kelvin, err = strconv.Atoi(args[3])
					if err != nil {
						log.Fatalln(err)
					}
					validateKelvin()
				}

				toUint16 := func(x int, max float64) uint16 {
					return uint16(float64(x)/max*math.MaxUint16 + 0.5)
				}

				setMsgPayloadColor(toUint16(h, 360), toUint16(s, 100), toUint16(l, 100))
			} else {
				log.Fatalln("Invalid color supplied")
			}

			handle(false, controlifx.LightSetColor(msgPayload))
		},
	}
	lightPowerCmd = &cobra.Command{
		Use:       "lightpower",
		Short:     "gets or sets the light power level",
		ValidArgs: []string{"[on|off]"},
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				handle(true, controlifx.LightGetPower())
			} else {
				powerOn, err := argToBool(args[0])
				if err != nil {
					log.Fatalln(err)
				}

				var level uint16
				if powerOn {
					level = 0xffff
				}

				handle(false, controlifx.LightSetPower(controlifx.LightSetPowerLanMessage{
					Level:    level,
					Duration: uint32(duration),
				}))
			}
		},
	}

	// Flags.

	rgb      bool
	kelvin   int
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
		getOwnerCmd,
		echoRequestCmd,
		lightGetCmd,
		lightSetColorCmd,
		lightPowerCmd,
	)

	lightSetColorCmd.Flags().BoolVar(&rgb, "rgb", false, "the color values are in red, green, and blue form")
	lightSetColorCmd.Flags().IntVarP(&kelvin, "kelvin", "k", 3500, "the color temperature (Kelvin)")
	lightSetColorCmd.Flags().IntVarP(&duration, "duration", "d", 0, "the duration of the color transition in milliseconds")

	lightPowerCmd.Flags().IntVarP(&duration, "duration", "d", 0, "the duration of the power transition in milliseconds")
}

func rgbToHsl(rI, gI, bI uint8) (uint16, uint16, uint16) {
	r := float64(rI) / 255
	g := float64(gI) / 255
	b := float64(bI) / 255

	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))

	h := (max + min) / 2
	s := h
	l := s

	if max == min {
		h = 0
		s = 0
	} else {
		d := max - min

		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case r:
			if g < b {
				h = (g-b)/d + 6
			} else {
				h = (g - b) / d
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}

		h /= 6
	}

	toUint16 := func(x float64) uint16 {
		return uint16(x*math.MaxUint16 + 0.5)
	}

	return toUint16(h), toUint16(s), toUint16(l)
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

func handle(requireResByDef bool, msg controlifx.SendableLanMessage) {
	var conn controlifx.Connection
	var err error
	var bcastAddr *net.UDPAddr
	bcastAddr, err = net.ResolveUDPAddr("udp", broadcast)
	if err == nil {
		conn, err = controlifx.ManualConnect(bcastAddr)
	}
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	discoveryNeeded := len(labels) > 0 || len(groups) > 0 || len(macs) > 0 || len(ips) > 0 || count > 0
	resCode := getResponseCode(msg)

	if requireAck {
		msg.Header.FrameAddress.AckRequired = true
		resCode = controlifx.AcknowledgementType
	}

	waitForResponses := requireAck || requireResByDef || requireRes

	if !requireAck && waitForResponses {
		msg.Header.FrameAddress.ResRequired = true
	}

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
		return controlifx.StatePowerType
	case controlifx.GetLabelType:
		return controlifx.StateLabelType
	case controlifx.SetLabelType:
		return controlifx.StateLabelType
	case controlifx.GetVersionType:
		return controlifx.StateVersionType
	case controlifx.GetInfoType:
		return controlifx.StateInfoType
	case controlifx.GetLocationType:
		return controlifx.StateLocationType
	case controlifx.GetGroupType:
		return controlifx.StateGroupType
	case controlifx.GetOwnerType:
		return controlifx.StateOwnerType
	case controlifx.EchoRequestType:
		return controlifx.EchoResponseType
	case controlifx.LightGetType:
		return controlifx.LightStateType
	case controlifx.LightSetColorType:
		return controlifx.LightStateType
	case controlifx.LightGetPowerType:
		return controlifx.LightStatePowerType
	case controlifx.LightSetPowerType:
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

		responses = append(responses, struct {
			Device   controlifx.Device
			Response interface{} `json:",omitempty"`
		}{
			Device:   device,
			Response: createFriendlyPayload(msg.Payload),
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
		return struct {
			Payload string
		}{
			Payload: string(bytes.TrimRight(payload.(*controlifx.EchoResponseLanMessage).Payload[:], "\x00")),
		}
	}

	return payload
}
