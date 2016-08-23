package cmd

import (
	"github.com/spf13/cobra"
	"github.com/bionicrm/controlifx"
)

var (
	RootCmd = &cobra.Command{
		Use:"clifx",
		Short:"Control LIFX devices from the command line",
	}

	// Flags.

	labels []string
	groups []string
	macs   []string
	ips    []string

	timeout int
	count   int

	requireAck bool
	requireRes bool
	pretty     bool

	broadcast  string
)

func init() {
	RootCmd.PersistentFlags().StringSliceVarP(&labels, "label", "l", []string{},
		"the message will only be sent to devices with one of the given labels")
	RootCmd.PersistentFlags().StringSliceVarP(&groups, "group", "g", []string{},
		"the message will only be sent to devices in one of the given groups")
	RootCmd.PersistentFlags().StringSliceVarP(&macs, "mac", "m", []string{},
		"the message will only be sent to devices with one of the given MAC addresses")
	RootCmd.PersistentFlags().StringSliceVarP(&ips, "ip", "i", []string{},
		"the message will only be sent to the given IP addresses")

	RootCmd.PersistentFlags().IntVarP(&timeout, "timeout", "t", controlifx.NormalTimeout,
		"devices will be discovered for the duration of the timeout in milliseconds until continuing with sending the message; 0 = no timeout")
	RootCmd.PersistentFlags().IntVarP(&count, "count", "c", 0,
		"only the given number of devices will be discovered before continuing with sending the message")

	RootCmd.PersistentFlags().BoolVarP(&requireAck, "require-ack", "a", false,
		"acknowledgement responses will be printed; mutually exclusive with --require-res")
	RootCmd.PersistentFlags().BoolVarP(&requireRes, "require-res", "r", false,
		"responses will always be printed; mutually exclusive with --require-ack")
	RootCmd.PersistentFlags().BoolVarP(&pretty, "pretty", "p", false,
		"pretty prints any JSON output")

	RootCmd.PersistentFlags().StringVar(&broadcast, "broadcast-addr", "",
		"overrides the broadcast address when sending messages to all devices on the network")
}
