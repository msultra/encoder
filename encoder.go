package encoder

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"strings"
)

type EncoderTags struct {
	Padding32 bool
	Padding64 bool
}

func extractTags(sf reflect.StructField) *EncoderTags {
	et := &EncoderTags{}
	tag := sf.Tag.Get("encoder")
	if tag == "" {
		return et
	}

	for _, t := range strings.Split(tag, ",") {
		switch t {
		case "padding32":
			et.Padding32 = true
		case "padding64":
			et.Padding64 = true
		}
	}
	return et
}

func marshal(v any, tags *EncoderTags) ([]byte, error) {
	tf := reflect.TypeOf(v)
	vf := reflect.ValueOf(v)

	if tf.Kind() == reflect.Ptr {
		vf = reflect.Indirect(vf)
		tf = vf.Type()
	}

	var b bytes.Buffer
	switch tf.Kind() {
	case reflect.Struct:
		for j := 0; j < vf.NumField(); j++ {
			tags := extractTags(tf.Field(j))
			buf, err := marshal(vf.Field(j).Interface(), tags)
			if err != nil {
				return nil, err
			}
			if err := binary.Write(&b, binary.LittleEndian, buf); err != nil {
				return nil, err
			}
		}
	case reflect.Slice, reflect.Array:
		if tf.Elem().Kind() != reflect.Uint8 {
			return nil, errors.New("marshal not implemented for slice element kind: " + tf.Elem().Kind().String())
		}

		var value []byte
		if tf.Kind() == reflect.Array {
			value = make([]byte, tf.Len())
			reflect.Copy(reflect.ValueOf(value), vf)
		} else {
			value = vf.Interface().([]byte)
		}

		if tags.Padding32 {
			value = append(value, make([]byte, len(value)%4)...)
		} else if tags.Padding64 {
			value = append(value, make([]byte, len(value)%8)...)
		}

		if err := binary.Write(&b, binary.LittleEndian, value); err != nil {
			return nil, err
		}
	case reflect.Uint8:
		if err := binary.Write(&b, binary.LittleEndian, vf.Interface().(uint8)); err != nil {
			return nil, err
		}
	case reflect.Uint16:
		if err := binary.Write(&b, binary.LittleEndian, vf.Interface().(uint16)); err != nil {
			return nil, err
		}
	case reflect.Uint32:
		if err := binary.Write(&b, binary.LittleEndian, vf.Interface().(uint32)); err != nil {
			return nil, err
		}
	case reflect.Uint64:
		if err := binary.Write(&b, binary.LittleEndian, vf.Interface().(uint64)); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("marshal not implemented for kind: " + tf.Kind().String())
	}
	return b.Bytes(), nil
}

func Marshal(v any) ([]byte, error) {
	return marshal(v, &EncoderTags{})
}

func unmarshal(b []byte, v any, tags *EncoderTags) (r any, off int, err error) {
	tf := reflect.TypeOf(v)
	vf := reflect.ValueOf(v)

	if tf.Kind() == reflect.Ptr {
		vf = reflect.Indirect(vf)
		tf = vf.Type()
	}

	if tags != nil && (tags.Padding32 || tags.Padding64) {
		return nil, 0, errors.New("padding not implemented for unmarshal")
	}

	buffer := bytes.NewBuffer(b)
	switch tf.Kind() {
	case reflect.Struct:
		var offset int
		for j := 0; j < vf.NumField(); j++ {
			tags := extractTags(tf.Field(j))

			var data any
			var off int
			switch tf.Field(j).Type.Kind() {
			case reflect.Struct:
				data, off, err = unmarshal(b[offset:], vf.Field(j).Addr().Interface(), tags)
			default:
				data, off, err = unmarshal(b[offset:], vf.Field(j).Interface(), tags)
			}
			if err != nil {
				return nil, 0, err
			}

			if tf.Field(j).Type.Kind() == reflect.Array {
				// convert slice to array
				arr := reflect.New(tf.Field(j).Type).Elem()
				reflect.Copy(arr, reflect.ValueOf(data))
				vf.Field(j).Set(arr)
			} else {
				vf.Field(j).Set(reflect.ValueOf(data))
			}

			offset += off
		}
		v = reflect.Indirect(reflect.ValueOf(v)).Interface()
		return v, offset, nil
	case reflect.Slice, reflect.Array:
		if tf.Elem().Kind() != reflect.Uint8 {
			return nil, 0, errors.New("unmarshal not implemented for slice element kind: " + tf.Elem().Kind().String())
		}

		var length int
		if tf.Kind() == reflect.Array {
			// Use the array's fixed length
			length = tf.Len()
		} else {
			// Use remaining buffer length for slices
			length = len(b)
		}

		return b[:length], length, nil
	case reflect.Uint8:
		var ret uint8
		if err := binary.Read(buffer, binary.LittleEndian, &ret); err != nil {
			return nil, 0, err
		}
		return ret, 1, nil
	case reflect.Uint16:
		var ret uint16
		if err := binary.Read(buffer, binary.LittleEndian, &ret); err != nil {
			return nil, 0, err
		}
		return ret, 2, nil
	case reflect.Uint32:
		var ret uint32
		if err := binary.Read(buffer, binary.LittleEndian, &ret); err != nil {
			return nil, 0, err
		}
		return ret, 4, nil
	case reflect.Uint64:
		var ret uint64
		if err := binary.Read(buffer, binary.LittleEndian, &ret); err != nil {
			return nil, 0, err
		}
		return ret, 8, nil
	default:
		return nil, 0, errors.New("unmarshal not implemented for kind: " + tf.Kind().String())
	}
}

func Unmarshal(b []byte, v any) error {
	_, _, err := unmarshal(b, v, &EncoderTags{})
	return err
}
