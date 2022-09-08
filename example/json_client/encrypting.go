package main

import (
	"bytes"
	"encoding/json"
	"github.com/dimonrus/tcpless"
)

type (
	// JsonDataEncryptor encode decode via json serialisation
	JsonDataEncryptor struct {
		// encoder
		encoder *json.Encoder
		// decoder
		decoder *json.Decoder
	}
)

// Encode message
func (d *JsonDataEncryptor) Encode(v any) error {
	return d.encoder.Encode(v)
}

// Decode message
func (d *JsonDataEncryptor) Decode(v any) error {
	return d.decoder.Decode(v)
}

// RegisterType register custom type
func (d *JsonDataEncryptor) RegisterType(v any) {
	panic("no need to json register type call")
}

// NewJSONDataEncryptor init json data encryptor
func NewJSONDataEncryptor(buf *bytes.Buffer) tcpless.DataEncryptor {
	return &JsonDataEncryptor{
		encoder: json.NewEncoder(buf),
		decoder: json.NewDecoder(buf),
	}
}
