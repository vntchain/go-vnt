package multistream

import (
	"errors"
	"io"
	"fmt"
)

// ErrNotSupported is the error returned when the muxer does not support
// the protocol specified for the handshake.
var ErrNotSupported = errors.New("protocol not supported")

// SelectProtoOrFail performs the initial multistream handshake
// to inform the muxer of the protocol that will be used to communicate
// on this ReadWriteCloser. It returns an error if, for example,
// the muxer does not know how to handle this protocol.
func SelectProtoOrFail(proto string, rwc io.ReadWriteCloser) error {
	err := handshake(rwc)
	if err != nil {
		return err
	}

	return trySelect(proto, rwc)
}

// SelectOneOf will perform handshakes with the protocols on the given slice
// until it finds one which is supported by the muxer.
func SelectOneOf(protos []string, rwc io.ReadWriteCloser) (string, error) {
	fmt.Printf("#### SelectOnOf starting... \n")
	err := handshake(rwc)
	fmt.Printf("#### SelectOnOf starting..., err: %v \n", err)
	if err != nil {
		return "", err
	}

	fmt.Printf("#### SelectOnOf protocals: %v \n", protos)
	for _, p := range protos {
		err := trySelect(p, rwc)
		fmt.Printf("#### SelectOnOf.trySelect err: %v \n", err)
		switch err {
		case nil:
			return p, nil
		case ErrNotSupported:
		default:
			return "", err
		}
	}
	return "", ErrNotSupported
}

func handshake(rwc io.ReadWriteCloser) error {
	tok, err := ReadNextToken(rwc)
	if err != nil {
		return err
	}

	if tok != ProtocolID {
		return errors.New("received mismatch in protocol id")
	}

	err = delimWrite(rwc, []byte(ProtocolID))
	if err != nil {
		return err
	}

	return nil
}

func trySelect(proto string, rwc io.ReadWriteCloser) error {
	err := delimWrite(rwc, []byte(proto))
	if err != nil {
		return err
	}

	tok, err := ReadNextToken(rwc)
	if err != nil {
		return err
	}

	switch tok {
	case proto:
		return nil
	case "na":
		return ErrNotSupported
	default:
		return errors.New("unrecognized response: " + tok)
	}
}
