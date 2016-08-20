package cmd

import (
	"github.com/spf13/cobra"
	"github.com/bionicrm/clifx/protocol"
	"github.com/bionicrm/controlifx"
	"fmt"
	"os"
)

var (
	RootCmd = &cobra.Command{
		Use: "lifx",
		Short: "Control LIFX devices from the command line",
		Run: func(cmd *cobra.Command, args []string) {
			var conn controlifx.Connector
			var devices []controlifx.Device
			discovered := len(macs) > 0 || len(labels) > 0 || len(ips) > 0 || count > 0

			if discovered {
				var err error
				conn, devices, err = protocol.Discover(macs, labels, ips, timeout, count)
				if err != nil {
					errorOut(err)
				}
			}

			msg, err := protocol.CreateMessage(msgType, payload)
			if err != nil {
				errorOut(err)
			}
			if get {
				if err := protocol.SendAndReceiveMessages(conn, devices, msg, get, pretty, ackOnly); err != nil {
					errorOut(err)
				}
			} else if discovered {
				if err := conn.SendTo(msg, devices); err != nil {
					errorOut(err)
				}
			} else {
				if err := conn.SendToAll(msg); err != nil {
					errorOut(err)
				}
			}
		},
	}

	// Flags.

	labels []string
	macs   []string
	ips    []string

	timeout int
	count   int

	msgType string
	payload []string
	get     bool
	pretty  bool
	ackOnly bool
)

func init() {
	RootCmd.PersistentFlags().StringSliceVar(&labels, "label", []string{},
		"the message will only be sent to devices with one of the given labels")
	RootCmd.PersistentFlags().StringSliceVar(&macs, "mac", []string{},
		"the message will only be sent to devices with one of the given MAC addresses")
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
	RootCmd.PersistentFlags().BoolVar(&get, "get", false,
		"prints all responses in JSON form")
	RootCmd.PersistentFlags().BoolVar(&pretty, "pretty-json", false,
		"pretty prints any JSON output")
	RootCmd.PersistentFlags().BoolVar(&ackOnly, "ack-only", false,
		"only prints an acknowledgement response; does nothing if 'get' is not specified")
}

func errorOut(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}
