package ioutils

import "io"

type nopWriteCloser struct {
	io.Writer
}

func (nopWriteCloser) Close() error {
	return nil
}

func NopWriteCloser(w io.Writer) nopWriteCloser {
	return nopWriteCloser{w}
}
