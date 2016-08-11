package cmd

import (
	"github.com/spf13/cobra"
	"github.com/bionicrm/clifx/protocol"
	"github.com/bionicrm/controlifx"
	"math/rand"
	"time"
	"fmt"
	"os"
)

var (
	RootCmd = &cobra.Command{
		Use: "lifx",
		Short: "Control LIFX devices from the command line",
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := protocol.Discover(mac, labels, ips, timeout, count)
			if err != nil {
				errorOut(err)
			}

			msg, err := protocol.CreateMessage(msgType, payload)
			if err != nil {
				errorOut(err)
			}
			if len(get) > 0 {
				if err := protocol.SendAndReceiveMessages(conn, msg, get); err != nil {
					errorOut(err)
				}
			} else if err := conn.SendToAll(msg); err != nil {
				errorOut(err)
			}
		},
	}

	// Flags.

	labels []string
	mac    string
	ips    []string

	timeout int
	count   int

	msgType string
	payload []string
	get     []string
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

	RootCmd.PersistentFlags().StringVar(&msgType, "type", "GetService",
		"the name of the type of message to be sent")
	RootCmd.PersistentFlags().StringSliceVar(&payload, "payload", []string{},
		"the payload values (if applicable) in the form 'FieldName:value,FieldName:SubFieldName:value,...'")
	RootCmd.PersistentFlags().StringSliceVar(&get, "get", []string{},
		"the payload values to print out of all responses in the form 'FieldName,FieldName:SubFieldName,...'")
}

func errorOut(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
