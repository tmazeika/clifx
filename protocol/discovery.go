package protocol

import (
	"github.com/bionicrm/controlifx"
	"net"
	"strings"
	"encoding/binary"
)

func Discover(mac string, labels, ips []string, timeout int, count int) (*controlifx.Connector, error) {
	conn := &controlifx.Connector{
		DiscoverTimeout:timeout,
	}

	if mac == "" && len(labels) == 0 && len(ips) == 0 {
		if count > -1 {
			return conn, conn.DiscoverNDevices(count)
		}
		return conn, conn.DiscoverAllDevices()
	}

	if mac != "" {
		parsedMac, err := net.ParseMAC(mac)
		if err != nil {
			return nil, err
		}

		conn.DiscoverFilteredDevices(func(msg controlifx.ReceivableLanMessage, device controlifx.Device) (bool, bool) {
			if macEqual(device.Mac, parsedMac) {
				// Do register, don't continue.
				return true, false
			}
			// Don't register, do continue.
			return false, true
		})
	} else if count > -1 {
		if err := conn.DiscoverNDevices(count); err != nil {
			return nil, err
		}
	} else if err := conn.DiscoverAllDevices(); err != nil {
		return nil, err
	}

	if len(labels) > 0 {
		msg := controlifx.LanDeviceMessageBuilder{}.GetLabel()

		recMsgs, err := conn.GetResponseFromAll(msg, controlifx.TypeFilter(&controlifx.StateLabelLanMessage{}))
		if err != nil {
			return nil, err
		}

		// Remove devices that don't have the requested label(s).
		for device, msg := range recMsgs {
			var ok bool
			label := strings.ToLower(string(msg.Payload.(*controlifx.StateLabelLanMessage).Label))

			for _, acceptLabel := range labels {
				if strings.ToLower(acceptLabel) == label {
					ok = true
					break
				}
			}
			if !ok {
				conn.RemoveDevice(device)
			}
		}
	}

	return conn, nil
}

func macEqual(i uint64, b []byte) bool {
	// Pad b if necessary.
	if len(b) == 6 {
		b = append([]byte{0, 0}, b...)
	}
	return i == binary.BigEndian.Uint64(b)
}
