package protocol

import (
	"github.com/bionicrm/controlifx"
	"net"
	"strings"
	"encoding/binary"
	"fmt"
	"os"
)

func Discover(macWhitelist, labelWhitelist, ipWhitelist []string, timeout int, count int) (conn controlifx.Connector, devices []controlifx.Device, err error) {
	conn.DiscoverTimeout = timeout

	// Used for count enforcement.
	leftToDiscover := count

	if devices, err = conn.DiscoverFilteredDevices(func(msg controlifx.ReceivableLanMessage, device controlifx.Device) (register bool, cont bool) {
		// Register and continue until set otherwise.
		register = true
		cont = true

		var err error

		// Enforce MAC whitelist.
		register, err = macIsWhitelisted(macWhitelist, device.Mac)
		if err != nil {
			// Error out.
			fmt.Println(err)
			os.Exit(-1)
		}

		// Enforce IP whitelist.
		register = ipIsWhitelisted(ipWhitelist, device.Addr.IP.String())

		// Enforce count.
		if count > -1 {
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
		msg := controlifx.LanDeviceMessageBuilder{}.GetLabel()
		var recMsgs map[controlifx.Device]controlifx.ReceivableLanMessage
		if recMsgs, err = conn.SendToAndGet(msg, controlifx.TypeFilter(&controlifx.StateLabelLanMessage{}), devices); err != nil {
			return
		}
		devices = nil
		for device, recMsg := range recMsgs {
			if labelIsWhitelisted(labelWhitelist, recMsg.Payload.(*controlifx.StateLabelLanMessage).Label) {
				devices = append(devices, device)
			}
		}
	}
	return
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

func labelIsWhitelisted(whitelist []string, label controlifx.Label) bool {
	labelStr := strings.ToLower(string(label))
	for _, v := range whitelist {
		if labelStr == strings.ToLower(v) {
			return true
		}
	}
	return false
}
