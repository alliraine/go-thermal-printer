package escpos

import "io"

// Transport is the minimal interface the print service needs.
type Transport interface {
	io.WriteCloser
}

// SerialLike is any serial handle that has Write and Close.
type SerialLike interface {
	Write([]byte) (int, error)
	Close() error
}

// WrapSerial adapts a SerialLike to Transport.
type WrapSerial struct{ S SerialLike }

func (w WrapSerial) Write(p []byte) (int, error) { return w.S.Write(p) }
func (w WrapSerial) Close() error                { return w.S.Close() }
