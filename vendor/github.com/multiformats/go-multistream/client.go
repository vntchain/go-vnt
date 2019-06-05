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
	fmt.Printf("#### client.go: SelectOneOf: protos: %v \n", protos)
	err := handshake(rwc)
	fmt.Printf("#### client.go: SelectOneOf finished: err: %v \n", err)
	if err != nil {
		return "", err
	}

	for _, p := range protos {
		err := trySelect(p, rwc)
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
	errCh := make(chan error, 1)
	go func() {
		fmt.Printf("#### client.go: handshake: will delimWriteBuffered \n")
		errCh <- delimWriteBuffered(rwc, []byte(ProtocolID))
	}()

	tok, readErr := ReadNextToken(rwc)
	fmt.Printf("#### client.go: handshake: ReadNextToken: token: %s, err: %v \n", tok, readErr)
	writeErr := <-errCh
	fmt.Printf("#### client.go: handshake: ReadNextToken: writeErr: %v \n", writeErr)

	if writeErr != nil {
		return writeErr
	}
	if readErr != nil {
		return readErr
	}

	if tok != ProtocolID {
		return errors.New("received mismatch in protocol id")
	}
	return nil
}

func trySelect(proto string, rwc io.ReadWriteCloser) error {
	fmt.Printf("#### client.go: trySelect: proto %s \n", proto)
	fmt.Printf("#### client.go: trySelect: will delimWriteBuffered \n")
	err := delimWriteBuffered(rwc, []byte(proto))
	fmt.Printf("#### client.go: trySelect: done delimWriteBuffered, err %v \n", err)
	if err != nil {
		return err
	}

	tok, err := ReadNextToken(rwc)
	fmt.Printf("#### client.go: trySelect: done ReadNextToken, token %s, err %v \n", tok, err)

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
