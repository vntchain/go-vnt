package utils

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"encoding/binary"
	"io"
	"io/ioutil"
	"time"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/wavm/contract"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/rlp"
)

const (
	HEADERLEN   int    = 6
	MAGIC       uint32 = 0x6d736101
	MagicBase64 uint32 = 0x7a464741
	ZLIBUTIL    uint16 = 0x01
	GZIPUTIL    uint16 = 0x02
)

func Compress(src []byte) []byte {
	log.Debug("before compress", "code length", len(src))
	defer func(start time.Time) { log.Debug("Compress finished", "runtime", time.Since(start)) }(time.Now())
	// now only support zlib compression
	dst := compressZlib(src)
	log.Debug("after compress", "compressed len", len(dst), "src length", len(src))

	// have a check if the compress is worth it
	isCompressed := len(dst)+HEADERLEN < len(src)

	// if compress is ok, add the header in front of the code
	if isCompressed {
		compressType := make([]byte, 2)
		binary.LittleEndian.PutUint16(compressType, ZLIBUTIL)
		dst = append(compressType, dst...)
		return dst
	}

	log.Debug("after compress", "code length", len(src))
	// if not compressed, return the raw bytes
	return src
}

func DeCompress(src []byte) ([]byte, error) {
	defer func(start time.Time) { log.Debug("DeCompress finished", "runtime", time.Since(start)) }(time.Now())
	if len(src) == 0 {
		return src, nil
	}

	log.Debug("before decompress", "code length", len(src))

	compressType, err := readCompressType(src)
	if err != nil {
		return nil, err
	}

	log.Debug("during decompress", "compressType", compressType)
	// get rid of the compress type bytes
	src = src[2:]
	var dst []byte
	switch compressType {
	case ZLIBUTIL:
		dst, err = deZlib(src)
	case GZIPUTIL:
		dst, err = deGzip(src)
	default:
		dst, err = src, nil
	}

	// log.Debug("after decompress", "decompressed code", dst, "err", err)
	log.Debug("after decompress", "decompressed length", len(dst), "src length", len(src))
	return dst, err
}

func compressZlib(src []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	defer w.Close()

	w.Write(src)
	w.Flush()
	dst := b.Bytes()
	return dst
}

func CompressZlib(src []byte) []byte {
	return compressZlib(src)
}

func deZlib(src []byte) (dst []byte, err error) {
	var out bytes.Buffer
	b := bytes.NewReader(src)
	r, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	io.Copy(&out, r)
	return out.Bytes(), nil
}

func deGzip(src []byte) (dst []byte, err error) {
	b := bytes.NewReader(src)
	r, _ := gzip.NewReader(b)
	defer r.Close()

	dst, err = ioutil.ReadAll(r)
	return
}

func ReadMagic(src []byte) (uint32, error) {
	r := bytes.NewBuffer(src)
	var buf [4]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}
	log.Debug("read magic", "magic", common.ToHex(buf[:]))
	return binary.LittleEndian.Uint32(buf[:]), nil
}

func readCompressType(src []byte) (uint16, error) {
	r := bytes.NewBuffer(src)
	var buf [2]byte
	_, err := io.ReadFull(r, buf[:])
	if err != nil {
		return 0, err
	}
	log.Debug("read CompressType", "CompressType", common.ToHex(buf[:]))
	return binary.LittleEndian.Uint16(buf[:]), nil
}

//将abi和wasm压缩后进行rlp编码
func CompressWasmAndAbi(abijson, wasm, compiled []byte) []byte {
	wasmcode := contract.WasmCode{
		Code:     wasm,
		Abi:      abijson,
		Compiled: compiled,
	}
	res, err := rlp.EncodeToBytes(wasmcode)
	if err != nil {
		panic(err)
	}
	rlpcps := Compress(res)
	cpsres, err := rlp.EncodeToBytes(rlpcps)
	if err != nil {
		panic(err)
	}
	magic := make([]byte, 4)
	binary.LittleEndian.PutUint32(magic, MAGIC)
	return append(magic, cpsres...)
}
