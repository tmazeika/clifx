package protocol

import (
	"github.com/bionicrm/controlifx"
	"strconv"
	"strings"
	"github.com/pkg/errors"
	"reflect"
)

func CreateMessage(msgType string, payloadStrs []string) (msg controlifx.SendableLanMessage, _ error) {
	method := reflect.ValueOf(&controlifx.LanDeviceMessageBuilder{}).MethodByName(msgType)

	if !method.IsValid() || method.Type().NumOut() != 1 || method.Type().Out(0) != reflect.TypeOf(controlifx.SendableLanMessage{}) {
		return msg, errors.New("unknown message type '" + msgType + "'")
	}

	if method.Type().NumIn() == 1 {
		payload := reflect.Indirect(reflect.New(method.Type().In(0)))

		for _, payloadValue := range payloadStrs {
			parts := strings.Split(payloadValue, ":")

			payloadField := payload
			for i, field := range parts {
				// Check if value, rather than field.
				if i == len(parts)-1 {
					value, err := sToValFor(field, payloadField)
					if err != nil {
						return msg, err
					}
					payloadField.Set(value)
				} else {
					payloadField = payloadField.FieldByName(field)
					if !payloadField.IsValid() {
						return msg, errors.New("unknown field '" + field + "' on " + payloadField.String())
					}
				}
			}
		}

		return method.Call([]reflect.Value{payload})[0].Interface().(controlifx.SendableLanMessage), nil
	}

	return method.Call([]reflect.Value{})[0].Interface().(controlifx.SendableLanMessage), nil
}

func sToValFor(s string, field reflect.Value) (v reflect.Value, err error) {
	var (
		// If the field is an array, we should switch on the slice.
		safeField = field
		i int
		f float64
	)

	if field.Kind() == reflect.Array {
		safeField = field.Slice(0, field.Len())
	}

	switch safeField.Interface().(type) {
	case controlifx.Label:
		v = reflect.ValueOf(controlifx.Label(s))
	case string:
		v = reflect.ValueOf(s)
	case []byte:
		b := []byte(s)
		v = field
		reflect.Copy(v, reflect.ValueOf(b))
	case uint8:
		if i, err = strconv.Atoi(s); err == nil {
			v = reflect.ValueOf(uint8(i))
		}
	case controlifx.PowerLevel:
		if i, err = strconv.Atoi(s); err == nil {
			v = reflect.ValueOf(controlifx.PowerLevel(uint16(i)))
		}
	case uint16:
		if i, err = strconv.Atoi(s); err == nil {
			v = reflect.ValueOf(uint16(i))
		}
	case controlifx.Port:
		if i, err = strconv.Atoi(s); err == nil {
			v = reflect.ValueOf(controlifx.Port(uint32(i)))
		}
	case uint32:
		if i, err = strconv.Atoi(s); err == nil {
			v = reflect.ValueOf(uint32(i))
		}
	case controlifx.Time:
		if i, err = strconv.Atoi(s); err == nil {
			v = reflect.ValueOf(controlifx.Time(uint64(i)))
		}
	case uint64:
		if i, err = strconv.Atoi(s); err == nil {
			v = reflect.ValueOf(uint64(i))
		}
	case float32:
		if f, err = strconv.ParseFloat(s, 32); err == nil {
			v = reflect.ValueOf(float32(f))
		}
	default:
		return v, errors.New("could not convert '" + s + "' to " + field.Type().String())
	}
	return
}
