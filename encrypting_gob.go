package tcpless

import (
	"bytes"
	"encoding/gob"
	"io"
)

type (
	// Encoder interface
	Encoder interface {
		// Encode message
		Encode(v any) error
	}

	// Decoder interface
	Decoder interface {
		// Decode message
		Decode(v any) error
	}

	// DataEncryptor provide encoding and decoding and register type
	DataEncryptor interface {
		// Encode message
		Encode(v any) error
		// Decode message
		Decode(v any) error
		// RegisterType register custom type
		RegisterType(v any)
	}

	// EmptyEncoder encoder do nothing
	EmptyEncoder struct {
		w io.Writer
	}

	// EmptyDecoder decoder do nothing
	EmptyDecoder struct {
		r io.Reader
	}

	// EmptyDataEncryptor empty data encryptor
	EmptyDataEncryptor struct {
		// encoder
		encoder Encoder
		// decoder
		decoder Decoder
	}

	// GobDataEncryptor encode decode via gob serialisation
	GobDataEncryptor struct {
		// encoder
		encoder *gob.Encoder
		// decoder
		decoder *gob.Decoder
	}

	// DataEncryptorConstructor apply constructor in set stream method
	DataEncryptorConstructor func(buf *bytes.Buffer) DataEncryptor
)

// Encode custom variable
func (e *EmptyEncoder) Encode(v any) error {
	// nothing to do
	return nil
}

// NewEmptyEncoder init empty encoder
func NewEmptyEncoder(w io.Writer) Encoder {
	return &EmptyEncoder{w: w}
}

// Decode custom variable
func (e *EmptyDecoder) Decode(v any) error {
	// nothing to do
	return nil
}

// NewEmptyDecoder init empty decoder
func NewEmptyDecoder(r io.Reader) Decoder {
	return &EmptyDecoder{r: r}
}

// Encode message
func (d *EmptyDataEncryptor) Encode(v any) error {
	return d.encoder.Encode(v)
}

// Decode message
func (d *EmptyDataEncryptor) Decode(v any) error {
	return d.decoder.Decode(v)
}

// RegisterType register custom type
func (d *EmptyDataEncryptor) RegisterType(v any) {
	// nothing register
}

// NewEmptyDataEncryptor init empty data encryptor
func NewEmptyDataEncryptor(buf *bytes.Buffer) DataEncryptor {
	return &EmptyDataEncryptor{
		encoder: NewEmptyEncoder(buf),
		decoder: NewEmptyDecoder(buf),
	}
}

// Encode message
func (d *GobDataEncryptor) Encode(v any) error {
	return d.encoder.Encode(v)
}

// Decode message
func (d *GobDataEncryptor) Decode(v any) error {
	return d.decoder.Decode(v)
}

// RegisterType register custom type
func (d *GobDataEncryptor) RegisterType(v any) {
	gob.Register(v)
}

// NewGobDataEncryptor init gob data encryptor
func NewGobDataEncryptor(buf *bytes.Buffer) DataEncryptor {
	return &GobDataEncryptor{
		encoder: gob.NewEncoder(buf),
		decoder: gob.NewDecoder(buf),
	}
}
