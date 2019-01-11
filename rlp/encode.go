// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package rlp

import (
	"fmt"
	"io"
	"math/big"
	"reflect"
	"sync"
)

var (
	// Common encoded values.
	// These are useful when implementing EncodeRLP.
	EmptyString = []byte{0x80}
	EmptyList   = []byte{0xC0}
)

// Encoder is implemented by types that require custom
// encoding rules or want to encode private fields.
type Encoder interface {
	// EncodeRLP should write the RLP encoding of its receiver to w.
	// If the implementation is a pointer method, it may also be
	// called for nil pointers.
	//
	// Implementations should generate valid RLP. The data written is
	// not verified at the moment, but a future version might. It is
	// recommended to write only a single value but writing multiple
	// values or no value at all is also permitted.
	EncodeRLP(io.Writer) error
}

// Encode writes the RLP encoding of val to w. Note that Encode may
// perform many small writes in some cases. Consider making w
// buffered.
//
// Encode uses the following type-dependent encoding rules:
//
// If the type implements the Encoder interface, Encode calls
// EncodeRLP. This is true even for nil pointers, please see the
// documentation for Encoder.
//
// To encode a pointer, the value being pointed to is encoded. For nil
// pointers, Encode will encode the zero value of the type. A nil
// pointer to a struct type always encodes as an empty RLP list.
// A nil pointer to an array encodes as an empty list (or empty string
// if the array has element type byte).
//
// Struct values are encoded as an RLP list of all their encoded
// public fields. Recursive struct types are supported.
//
// To encode slices and arrays, the elements are encoded as an RLP
// list of the value's elements. Note that arrays and slices with
// element type uint8 or byte are always encoded as an RLP string.
//
// A Go string is encoded as an RLP string.
//
// An unsigned integer value is encoded as an RLP string. Zero always
// encodes as an empty RLP string. Encode also supports *big.Int.
//
// An interface value encodes as the value contained in the interface.
//
// Boolean values are not supported, nor are signed integers, floating
// point numbers, maps, channels and functions.
func Encode(w io.Writer, val interface{}) error {
	if outer, ok := w.(*Encbuf); ok {
		// Encode was called by some type's EncodeRLP.
		// Avoid copying by writing to the outer encbuf directly.
		return outer.encode(val)
	}
	eb := encbufPool.Get().(*Encbuf)
	defer encbufPool.Put(eb)
	eb.reset()
	if err := eb.encode(val); err != nil {
		return err
	}
	return eb.toWriter(w)
}

// EncodeToBytes returns the RLP encoding of val.
// Please see the documentation of Encode for the encoding rules.
func EncodeToBytes(val interface{}) ([]byte, error) {
	eb := encbufPool.Get().(*Encbuf)
	defer encbufPool.Put(eb)
	eb.reset()
	if err := eb.encode(val); err != nil {
		return nil, err
	}
	return eb.toBytes(), nil
}

// EncodeToReader returns a reader from which the RLP encoding of val
// can be read. The returned size is the total size of the encoded
// data.
//
// Please see the documentation of Encode for the encoding rules.
func EncodeToReader(val interface{}) (size int, r io.Reader, err error) {
	eb := encbufPool.Get().(*Encbuf)
	eb.reset()
	if err := eb.encode(val); err != nil {
		return 0, nil, err
	}
	return eb.size(), &EncReader{Buf: eb}, nil
}

type Encbuf struct {
	Str     []byte      // string data, contains everything except list headers
	Lheads  []*Listhead // all list headers
	Lhsize  int         // sum of sizes of all encoded list headers
	Sizebuf []byte      // 9-byte auxiliary buffer for uint encoding
}

type Listhead struct {
	Offset int // index of this header in string data
	Size   int // total size of encoded data (including list headers)
}

// encode writes head to the given buffer, which must be at least
// 9 bytes long. It returns the encoded bytes.
func (head *Listhead) encode(buf []byte) []byte {
	return buf[:puthead(buf, 0xC0, 0xF7, uint64(head.Size))]
}

// headsize returns the size of a list or string header
// for a value of the given size.
func headsize(size uint64) int {
	if size < 56 {
		return 1
	}
	return 1 + intsize(size)
}

// puthead writes a list or string header to buf.
// buf must be at least 9 bytes long.
func puthead(buf []byte, smalltag, largetag byte, size uint64) int {
	if size < 56 {
		buf[0] = smalltag + byte(size)
		return 1
	}
	sizesize := putint(buf[1:], size)
	buf[0] = largetag + byte(sizesize)
	return sizesize + 1
}

// encbufs are pooled.
var encbufPool = sync.Pool{
	New: func() interface{} { return &Encbuf{Sizebuf: make([]byte, 9)} },
}

func (w *Encbuf) reset() {
	w.Lhsize = 0
	if w.Str != nil {
		w.Str = w.Str[:0]
	}
	if w.Lheads != nil {
		w.Lheads = w.Lheads[:0]
	}
}

// encbuf implements io.Writer so it can be passed it into EncodeRLP.
func (w *Encbuf) Write(b []byte) (int, error) {
	w.Str = append(w.Str, b...)
	return len(b), nil
}

func (w *Encbuf) encode(val interface{}) error {
	rval := reflect.ValueOf(val)
	ti, err := cachedTypeInfo(rval.Type(), tags{})
	if err != nil {
		return err
	}
	return ti.writer(rval, w)
}

func (w *Encbuf) encodeStringHeader(size int) {
	if size < 56 {
		w.Str = append(w.Str, 0x80+byte(size))
	} else {
		// TODO: encode to w.str directly
		sizesize := putint(w.Sizebuf[1:], uint64(size))
		w.Sizebuf[0] = 0xB7 + byte(sizesize)
		w.Str = append(w.Str, w.Sizebuf[:sizesize+1]...)
	}
}

func (w *Encbuf) encodeString(b []byte) {
	if len(b) == 1 && b[0] <= 0x7F {
		// fits single byte, no string header
		w.Str = append(w.Str, b[0])
	} else {
		w.encodeStringHeader(len(b))
		w.Str = append(w.Str, b...)
	}
}

func (w *Encbuf) list() *Listhead {
	lh := &Listhead{Offset: len(w.Str), Size: w.Lhsize}
	w.Lheads = append(w.Lheads, lh)
	return lh
}

func (w *Encbuf) listEnd(lh *Listhead) {
	lh.Size = w.size() - lh.Offset - lh.Size
	if lh.Size < 56 {
		w.Lhsize++ // length encoded into kind tag
	} else {
		w.Lhsize += 1 + intsize(uint64(lh.Size))
	}
}

func (w *Encbuf) size() int {
	return len(w.Str) + w.Lhsize
}

func (w *Encbuf) toBytes() []byte {
	out := make([]byte, w.size())
	strpos := 0
	pos := 0
	for _, head := range w.Lheads {
		// write string data before header
		n := copy(out[pos:], w.Str[strpos:head.Offset])
		pos += n
		strpos += n
		// write the header
		enc := head.encode(out[pos:])
		pos += len(enc)
	}
	// copy string data after the last list header
	copy(out[pos:], w.Str[strpos:])
	return out
}

func (w *Encbuf) toWriter(out io.Writer) (err error) {
	strpos := 0
	for _, head := range w.Lheads {
		// write string data before header
		if head.Offset-strpos > 0 {
			n, err := out.Write(w.Str[strpos:head.Offset])
			strpos += n
			if err != nil {
				return err
			}
		}
		// write the header
		enc := head.encode(w.Sizebuf)
		if _, err = out.Write(enc); err != nil {
			return err
		}
	}
	if strpos < len(w.Str) {
		// write string data after the last list header
		_, err = out.Write(w.Str[strpos:])
	}
	return err
}

// encReader is the io.Reader returned by EncodeToReader.
// It releases its encbuf at EOF.
type EncReader struct {
	Buf    *Encbuf // the buffer we're reading from. this is nil when we're at EOF.
	Lhpos  int     // index of list header that we're reading
	Strpos int     // current position in string buffer
	Piece  []byte  // next piece to be read
}

func (r *EncReader) Read(b []byte) (n int, err error) {
	for {
		if r.Piece = r.next(); r.Piece == nil {
			// Put the encode buffer back into the pool at EOF when it
			// is first encountered. Subsequent calls still return EOF
			// as the error but the buffer is no longer valid.
			if r.Buf != nil {
				encbufPool.Put(r.Buf)
				r.Buf = nil
			}
			return n, io.EOF
		}
		nn := copy(b[n:], r.Piece)
		n += nn
		if nn < len(r.Piece) {
			// piece didn't fit, see you next time.
			r.Piece = r.Piece[nn:]
			return n, nil
		}
		r.Piece = nil
	}
}

// next returns the next piece of data to be read.
// it returns nil at EOF.
func (r *EncReader) next() []byte {
	switch {
	case r.Buf == nil:
		return nil

	case r.Piece != nil:
		// There is still data available for reading.
		return r.Piece

	case r.Lhpos < len(r.Buf.Lheads):
		// We're before the last list header.
		head := r.Buf.Lheads[r.Lhpos]
		sizebefore := head.Offset - r.Strpos
		if sizebefore > 0 {
			// String data before header.
			p := r.Buf.Str[r.Strpos:head.Offset]
			r.Strpos += sizebefore
			return p
		}
		r.Lhpos++
		return head.encode(r.Buf.Sizebuf)

	case r.Strpos < len(r.Buf.Str):
		// String data at the end, after all list headers.
		p := r.Buf.Str[r.Strpos:]
		r.Strpos = len(r.Buf.Str)
		return p

	default:
		return nil
	}
}

var (
	encoderInterface = reflect.TypeOf(new(Encoder)).Elem()
	big0             = big.NewInt(0)
)

// makeWriter creates a writer function for the given type.
func makeWriter(typ reflect.Type, ts tags) (writer, error) {
	kind := typ.Kind()
	switch {
	case typ == rawValueType:
		return writeRawValue, nil
	case typ.Implements(encoderInterface):
		return writeEncoder, nil
	case kind != reflect.Ptr && reflect.PtrTo(typ).Implements(encoderInterface):
		return writeEncoderNoPtr, nil
	case kind == reflect.Interface:
		return writeInterface, nil
	case typ.AssignableTo(reflect.PtrTo(bigInt)):
		return writeBigIntPtr, nil
	case typ.AssignableTo(bigInt):
		return writeBigIntNoPtr, nil
	case isUint(kind):
		return writeUint, nil
	case kind == reflect.Bool:
		return writeBool, nil
	case kind == reflect.String:
		return writeString, nil
	case kind == reflect.Slice && isByte(typ.Elem()):
		return writeBytes, nil
	case kind == reflect.Array && isByte(typ.Elem()):
		return writeByteArray, nil
	case kind == reflect.Slice || kind == reflect.Array:
		return makeSliceWriter(typ, ts)
	case kind == reflect.Struct:
		return makeStructWriter(typ)
	case kind == reflect.Ptr:
		return makePtrWriter(typ)
	default:
		return nil, fmt.Errorf("rlp: type %v is not RLP-serializable", typ)
	}
}

func isByte(typ reflect.Type) bool {
	return typ.Kind() == reflect.Uint8 && !typ.Implements(encoderInterface)
}

func writeRawValue(val reflect.Value, w *Encbuf) error {
	w.Str = append(w.Str, val.Bytes()...)
	return nil
}

func writeUint(val reflect.Value, w *Encbuf) error {
	i := val.Uint()
	if i == 0 {
		w.Str = append(w.Str, 0x80)
	} else if i < 128 {
		// fits single byte
		w.Str = append(w.Str, byte(i))
	} else {
		// TODO: encode int to w.str directly
		s := putint(w.Sizebuf[1:], i)
		w.Sizebuf[0] = 0x80 + byte(s)
		w.Str = append(w.Str, w.Sizebuf[:s+1]...)
	}
	return nil
}

func writeBool(val reflect.Value, w *Encbuf) error {
	if val.Bool() {
		w.Str = append(w.Str, 0x01)
	} else {
		w.Str = append(w.Str, 0x80)
	}
	return nil
}

func writeBigIntPtr(val reflect.Value, w *Encbuf) error {
	ptr := val.Interface().(*big.Int)
	if ptr == nil {
		w.Str = append(w.Str, 0x80)
		return nil
	}
	return writeBigInt(ptr, w)
}

func writeBigIntNoPtr(val reflect.Value, w *Encbuf) error {
	i := val.Interface().(big.Int)
	return writeBigInt(&i, w)
}

func writeBigInt(i *big.Int, w *Encbuf) error {
	if cmp := i.Cmp(big0); cmp == -1 {
		return fmt.Errorf("rlp: cannot encode negative *big.Int")
	} else if cmp == 0 {
		w.Str = append(w.Str, 0x80)
	} else {
		w.encodeString(i.Bytes())
	}
	return nil
}

func writeBytes(val reflect.Value, w *Encbuf) error {
	w.encodeString(val.Bytes())
	return nil
}

func writeByteArray(val reflect.Value, w *Encbuf) error {
	if !val.CanAddr() {
		// Slice requires the value to be addressable.
		// Make it addressable by copying.
		copy := reflect.New(val.Type()).Elem()
		copy.Set(val)
		val = copy
	}
	size := val.Len()
	slice := val.Slice(0, size).Bytes()
	w.encodeString(slice)
	return nil
}

func writeString(val reflect.Value, w *Encbuf) error {
	s := val.String()
	if len(s) == 1 && s[0] <= 0x7f {
		// fits single byte, no string header
		w.Str = append(w.Str, s[0])
	} else {
		w.encodeStringHeader(len(s))
		w.Str = append(w.Str, s...)
	}
	return nil
}

func writeEncoder(val reflect.Value, w *Encbuf) error {
	return val.Interface().(Encoder).EncodeRLP(w)
}

// writeEncoderNoPtr handles non-pointer values that implement Encoder
// with a pointer receiver.
func writeEncoderNoPtr(val reflect.Value, w *Encbuf) error {
	if !val.CanAddr() {
		// We can't get the address. It would be possible to make the
		// value addressable by creating a shallow copy, but this
		// creates other problems so we're not doing it (yet).
		//
		// package json simply doesn't call MarshalJSON for cases like
		// this, but encodes the value as if it didn't implement the
		// interface. We don't want to handle it that way.
		return fmt.Errorf("rlp: game over: unadressable value of type %v, EncodeRLP is pointer method", val.Type())
	}
	return val.Addr().Interface().(Encoder).EncodeRLP(w)
}

func writeInterface(val reflect.Value, w *Encbuf) error {
	if val.IsNil() {
		// Write empty list. This is consistent with the previous RLP
		// encoder that we had and should therefore avoid any
		// problems.
		w.Str = append(w.Str, 0xC0)
		return nil
	}
	eval := val.Elem()
	ti, err := cachedTypeInfo(eval.Type(), tags{})
	if err != nil {
		return err
	}
	return ti.writer(eval, w)
}

func makeSliceWriter(typ reflect.Type, ts tags) (writer, error) {
	etypeinfo, err := cachedTypeInfo1(typ.Elem(), tags{})
	if err != nil {
		return nil, err
	}
	writer := func(val reflect.Value, w *Encbuf) error {
		if !ts.tail {
			defer w.listEnd(w.list())
		}
		vlen := val.Len()
		for i := 0; i < vlen; i++ {
			if err := etypeinfo.writer(val.Index(i), w); err != nil {
				return err
			}
		}
		return nil
	}
	return writer, nil
}

func makeStructWriter(typ reflect.Type) (writer, error) {
	fields, err := structFields(typ)
	if err != nil {
		return nil, err
	}
	writer := func(val reflect.Value, w *Encbuf) error {
		lh := w.list()
		for _, f := range fields {
			if err := f.info.writer(val.Field(f.index), w); err != nil {
				return err
			}
		}
		w.listEnd(lh)
		return nil
	}
	return writer, nil
}

func makePtrWriter(typ reflect.Type) (writer, error) {
	etypeinfo, err := cachedTypeInfo1(typ.Elem(), tags{})
	if err != nil {
		return nil, err
	}

	// determine nil pointer handler
	var nilfunc func(*Encbuf) error
	kind := typ.Elem().Kind()
	switch {
	case kind == reflect.Array && isByte(typ.Elem().Elem()):
		nilfunc = func(w *Encbuf) error {
			w.Str = append(w.Str, 0x80)
			return nil
		}
	case kind == reflect.Struct || kind == reflect.Array:
		nilfunc = func(w *Encbuf) error {
			// encoding the zero value of a struct/array could trigger
			// infinite recursion, avoid that.
			w.listEnd(w.list())
			return nil
		}
	default:
		zero := reflect.Zero(typ.Elem())
		nilfunc = func(w *Encbuf) error {
			return etypeinfo.writer(zero, w)
		}
	}

	writer := func(val reflect.Value, w *Encbuf) error {
		if val.IsNil() {
			return nilfunc(w)
		}
		return etypeinfo.writer(val.Elem(), w)
	}
	return writer, err
}

// putint writes i to the beginning of b in big endian byte
// order, using the least number of bytes needed to represent i.
func putint(b []byte, i uint64) (size int) {
	switch {
	case i < (1 << 8):
		b[0] = byte(i)
		return 1
	case i < (1 << 16):
		b[0] = byte(i >> 8)
		b[1] = byte(i)
		return 2
	case i < (1 << 24):
		b[0] = byte(i >> 16)
		b[1] = byte(i >> 8)
		b[2] = byte(i)
		return 3
	case i < (1 << 32):
		b[0] = byte(i >> 24)
		b[1] = byte(i >> 16)
		b[2] = byte(i >> 8)
		b[3] = byte(i)
		return 4
	case i < (1 << 40):
		b[0] = byte(i >> 32)
		b[1] = byte(i >> 24)
		b[2] = byte(i >> 16)
		b[3] = byte(i >> 8)
		b[4] = byte(i)
		return 5
	case i < (1 << 48):
		b[0] = byte(i >> 40)
		b[1] = byte(i >> 32)
		b[2] = byte(i >> 24)
		b[3] = byte(i >> 16)
		b[4] = byte(i >> 8)
		b[5] = byte(i)
		return 6
	case i < (1 << 56):
		b[0] = byte(i >> 48)
		b[1] = byte(i >> 40)
		b[2] = byte(i >> 32)
		b[3] = byte(i >> 24)
		b[4] = byte(i >> 16)
		b[5] = byte(i >> 8)
		b[6] = byte(i)
		return 7
	default:
		b[0] = byte(i >> 56)
		b[1] = byte(i >> 48)
		b[2] = byte(i >> 40)
		b[3] = byte(i >> 32)
		b[4] = byte(i >> 24)
		b[5] = byte(i >> 16)
		b[6] = byte(i >> 8)
		b[7] = byte(i)
		return 8
	}
}

// intsize computes the minimum number of bytes required to store i.
func intsize(i uint64) (size int) {
	for size = 1; ; size++ {
		if i >>= 8; i == 0 {
			return size
		}
	}
}
