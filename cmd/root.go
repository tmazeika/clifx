package cmd

import (
	"github.com/spf13/cobra"
	"github.com/bionicrm/clifx/protocol"
	"github.com/bionicrm/controlifx"
	"log"
	"math/rand"
	"time"
)

var (
	RootCmd = &cobra.Command{
		Use: "lifx",
		Short: "Control LIFX devices from the command line",
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := protocol.Discover(mac, labels, ips, timeout, count)
			if err != nil {
				log.Fatalln(err)
			}

			// TODO: remove testing
			builder := controlifx.LanDeviceMessageBuilder{}
			payload := controlifx.LightSetColorLanMessage{
				Color:controlifx.HSBK{
					Hue:uint16(rand.Intn(0xffff)),
					Saturation:0xffff,
					Brightness:0xffff,
					Kelvin:3500,
				},
				Duration:1024,
			}
			msg := builder.LightSetColor(payload)

			if err := conn.SendToAll(msg); err != nil {
				log.Fatalln(err)
			}
		},
	}

	// Flags.

	labels []string
	mac    string
	ips    []string

	timeout int
	count   int
)

func init() {
	// TODO: remove testing
	rand.Seed(time.Now().UTC().UnixNano())

	RootCmd.PersistentFlags().StringSliceVar(&labels, "label", []string{},
		"the message will only be sent to devices with one of the given labels")
	RootCmd.PersistentFlags().StringVar(&mac, "mac", "",
		"the message will only be sent to the device with the given MAC address")
	RootCmd.PersistentFlags().StringSliceVar(&ips, "ip", []string{},
		"the message will only be sent to the given IPv4/6 addresses")

	RootCmd.PersistentFlags().IntVar(&timeout, "timeout", controlifx.NormalDiscoverTimeout,
		"devices will be discovered for the duration of the timeout in milliseconds until continuing with sending the message; 0 = no timeout")
	RootCmd.PersistentFlags().IntVar(&count, "count", -1,
		"only the given number of devices will be discovered before continuing with sending the message")
}
