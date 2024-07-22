package json

import (
	jsonv1 "encoding/json"

	jsonv2 "github.com/go-json-experiment/json"
)

type jsonVersion interface {
	Marshal(v any) ([]byte, error)
	Unmarshal([]byte, any) error
}

var (
	marshal   func(v any) ([]byte, error)
	unmarshal func([]byte, any) error
)

func SetVersion(version int) {
	switch version {
	case 1:
		marshal = jsonv1.Marshal
		unmarshal = jsonv1.Unmarshal
	case 2:
		marshal = func(v any) ([]byte, error) {
			return jsonv2.Marshal(v)
		} 
		unmarshal =  func(b []byte, a any) error {
			return jsonv2.Unmarshal(b, a)
		}  
	}
}

func Marshal(v any) ([]byte, error) {
	return marshal(v)
}

func Unmarshal(b []byte, a any) error {
	return unmarshal(b,a)
}

func init() {
	SetVersion(1)
}