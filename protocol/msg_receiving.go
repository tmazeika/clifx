package protocol

import (
	"github.com/bionicrm/controlifx"
	"errors"
	"fmt"
	"encoding/json"
	"encoding"
)

type responseEntry struct {
	Device   controlifx.Device
	Response encoding.BinaryUnmarshaler
}

type ackResponseEntry struct {
	Device   controlifx.Device
}

func SendAndReceiveMessages(conn controlifx.Connector, devices []controlifx.Device, msg controlifx.SendableLanMessage, get, pretty, ackOnly bool) error {
	code, err := getResponseCode(msg)
	if err != nil {
		return err
	}
	if ackOnly {
		msg.Header.FrameAddress.AckRequired = true
	}
	recMsgs, err := conn.SendToAndGet(msg, func(msg controlifx.ReceivableLanMessage) bool {
		if ackOnly {
			return msg.Header.ProtocolHeader.Type == controlifx.AcknowledgementType
		} else {
			return msg.Header.ProtocolHeader.Type == code
		}
	}, devices)
	if err != nil {
		return err
	}
	if get {
		var jsonOut []interface{}
		for device, msg := range recMsgs {
			if msg.Header.ProtocolHeader.Type == controlifx.AcknowledgementType {
				jsonOut = append(jsonOut, ackResponseEntry{
					Device: device,
				})
			} else {
				jsonOut = append(jsonOut, responseEntry{
					Device: device,
					Response: msg.Payload,
				})
			}
		}
		var err error
		var jsonB []byte
		if pretty {
			jsonB, err = json.MarshalIndent(jsonOut, "", "  ")
		} else {
			jsonB, err = json.Marshal(jsonOut)
		}
		if err != nil {
			return err
		}
		fmt.Println(string(jsonB))
	}

	return nil
}

func getResponseCode(msg controlifx.SendableLanMessage) (code uint16, _ error) {
	switch msg.Header.ProtocolHeader.Type {
	case controlifx.GetServiceType:
		code = controlifx.StateServiceType
	case controlifx.GetHostInfoType:
		code = controlifx.StateHostInfoType
	case controlifx.GetHostFirmwareType:
		code = controlifx.StateHostFirmwareType
	case controlifx.GetWifiInfoType:
		code = controlifx.StateWifiInfoType
	case controlifx.GetWifiFirmwareType:
		code = controlifx.StateWifiFirmwareType
	case controlifx.GetPowerType:
		code = controlifx.StatePowerType
	case controlifx.SetPowerType:
		code = controlifx.AcknowledgementType // ack only
	case controlifx.GetLabelType:
		code = controlifx.StateLabelType
	case controlifx.SetLabelType:
		code = controlifx.AcknowledgementType // ack only
	case controlifx.GetVersionType:
		code = controlifx.StateVersionType
	case controlifx.GetInfoType:
		code = controlifx.StateInfoType
	case controlifx.GetLocationType:
		code = controlifx.StateLocationType
	case controlifx.GetGroupType:
		code = controlifx.StateGroupType
	case controlifx.EchoRequestType:
		code = controlifx.EchoResponseType
	case controlifx.LightGetType:
		code = controlifx.LightStateType
	case controlifx.LightSetColorType:
		code = controlifx.LightStateType // as per protocol
	case controlifx.LightGetPowerType:
		code = controlifx.LightStatePowerType
	case controlifx.LightSetPowerType:
		code = controlifx.LightStatePowerType // as per protocol
	default:
		return code, errors.New("no response can be expected for that message type")
	}
	return
}
