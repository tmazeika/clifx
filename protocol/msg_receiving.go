package protocol

import (
	"github.com/bionicrm/controlifx"
	"fmt"
	"strings"
	"reflect"
	"errors"
)

func SendAndReceiveMessages(conn *controlifx.Connector, msg controlifx.SendableLanMessage, get []string) error {
	responseType, code, err := getResponseType(msg)
	if err != nil {
		return err
	}
	recMsgs, err := conn.GetResponseFromAll(msg, func(msg controlifx.ReceivableLanMessage) bool {
		return msg.Header.ProtocolHeader.Type == code
	})
	if err != nil {
		return err
	}

	for device, msg := range recMsgs {
		response := reflect.Indirect(reflect.ValueOf(msg.Payload).Convert(responseType))

		fmt.Printf("[%s] ", device.Addr.String())

		for i, getValue := range get {
			parts := strings.Split(getValue, ":")

			responseField := response
			for _, field := range parts {
				payloadFieldStr := responseField.Type().Name()
				responseField = responseField.FieldByName(field)
				if !responseField.IsValid() {
					return errors.New("unknown field '" + field + "' on " + payloadFieldStr + " struct")
				}
			}
			fmt.Printf("%s: %v", getValue, responseField.Interface())
			if i < len(get)-1 {
				fmt.Print(",")
			}
			fmt.Print(" ")
		}

		fmt.Println()
	}

	return nil
}

func getResponseType(msg controlifx.SendableLanMessage) (t reflect.Type, code uint16, _ error) {
	// TODO: implement all
	switch msg.Header.ProtocolHeader.Type {
	case controlifx.GetServiceType:
		t = reflect.TypeOf(&controlifx.StateServiceLanMessage{})
		code = controlifx.StateServiceType
	default:
		return t, code, errors.New("no response can be expected for that message type")
	}
	return
}
