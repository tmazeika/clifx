package protocol

import (
	"encoding/binary"
	"fmt"
	"gopkg.in/golifx/controlifx.v1"
	"net"
	"os"
	"strings"
)

func Discover(conn controlifx.Connection, labelWhitelist, groupWhitelist, macWhitelist, ipWhitelist []string, timeout, count int) (devices []controlifx.Device, err error) {
	// Used for count enforcement.
	leftToDiscover := count

	if devices, err = conn.DiscoverDevices(timeout, func(msg controlifx.ReceivableLanMessage, device controlifx.Device) (register, cont bool) {
		// Register and continue until set otherwise.
		register = true
		cont = true

		var err error

		// Enforce MAC whitelist.
		register, err = macIsWhitelisted(macWhitelist, device.Mac)
		if err != nil {
			// Error out.
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		if register {
			// Enforce IP whitelist.
			register = ipIsWhitelisted(ipWhitelist, device.Addr.IP.String())
		}

		// Enforce count.
		if count > 0 {
			if register {
				leftToDiscover--
			}

			cont = leftToDiscover > 0
		}

		return
	}); err != nil {
		return
	}

	// Enforce label whitelist.
	if len(labelWhitelist) > 0 {
		msg := controlifx.GetLabel()

		var recMsgs map[controlifx.Device]controlifx.ReceivableLanMessage
		if recMsgs, err = conn.SendToAndGet(msg, devices, controlifx.TypeFilter(controlifx.StateLabelType)); err != nil {
			return
		}

		devices = nil

		for device, recMsg := range recMsgs {
			if labelIsWhitelisted(labelWhitelist, string(recMsg.Payload.(*controlifx.StateLabelLanMessage).Label)) {
				devices = append(devices, device)
			}
		}
	}

	// Enforce group whitelist.
	if len(groupWhitelist) > 0 {
		msg := controlifx.GetGroup()

		var recMsgs map[controlifx.Device]controlifx.ReceivableLanMessage
		if recMsgs, err = conn.SendToAndGet(msg, devices, controlifx.TypeFilter(controlifx.StateGroupType)); err != nil {
			return
		}

		devices = nil

		for device, recMsg := range recMsgs {
			if groupIsWhitelisted(groupWhitelist, string(recMsg.Payload.(*controlifx.StateGroupLanMessage).Label)) {
				devices = append(devices, device)
			}
		}
	}

	return
}

func labelIsWhitelisted(whitelist []string, label string) bool {
	label = strings.ToLower(label)

	for _, v := range whitelist {
		if label == strings.ToLower(v) {
			return true
		}
	}

	return false
}

func groupIsWhitelisted(whitelist []string, group string) bool {
	group = strings.ToLower(group)

	for _, v := range whitelist {
		if group == strings.ToLower(v) {
			return true
		}
	}

	return false
}

func macIsWhitelisted(whitelist []string, mac uint64) (bool, error) {
	// All MACs are wanted if there's no whitelist.
	if len(whitelist) == 0 {
		return true, nil
	}
	var err error
	parsedMacs := make([]net.HardwareAddr, len(whitelist))
	for i, v := range whitelist {
		parsedMacs[i], err = net.ParseMAC(v)
		if err != nil {
			goto NotWanted
		}
	}
	for _, wantedMac := range parsedMacs {
		if macEqual(mac, wantedMac) {
			return true, nil
		}
	}

NotWanted:
	return false, err
}

func macEqual(i uint64, b []byte) bool {
	// Pad b if necessary.
	if len(b) == 6 {
		b = append([]byte{0, 0}, b...)
	}

	return i == binary.BigEndian.Uint64(b)
}

func ipIsWhitelisted(whitelist []string, ip string) bool {
	// All IPs are wanted if there's no whitelist.
	if len(whitelist) == 0 {
		return true
	}

	for _, v := range whitelist {
		if ip == v {
			return true
		}
	}

	return false
}
