package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

func cutUSR(t string) string {
	pt := t
	idx := strings.LastIndex(t, "@FI@")
	if idx != -1 {
		pt = t[0:idx]
	}
	return pt
}

func readfile(filepath string) []string {
	fi, err := os.Open(filepath)
	if err != nil {
		panic(err.Error())
	}
	defer fi.Close()
	contents := []string{}
	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		contents = append(contents, string(a))
	}
	return contents
}

//KEY _complex s3;
var astKeyReg = `([ ]*)(KEY)([ ]+)(int(|32|64)|uint(|32|64|256)|address|string|bool|mapping|array|struct)([\s\S]*)`
var astKeyRegFmt = `([ ]*)(KEY)([ ]+)(int(|32|64)|uint(|32|64|256)|address|string|bool|mapping|array|struct|%s)([\s\S]*)`

const keyTypeReg = `(int(|32|64)|uint(|32|64|256)|address|string|bool|mapping|array|struct)`

func isKey(input string, structnames string) bool {
	keyReg := ""
	if structnames == "" {
		keyReg = astKeyReg
	} else {
		keyReg = fmt.Sprintf(astKeyRegFmt, structnames)
	}

	// fmt.Printf("isKey %s astKeyReg %s\n", input, keyReg)
	reg, err := regexp.Compile(keyReg)
	flag := false
	if err != nil {
		return flag
	}
	idx := reg.FindAllStringIndex(input, -1)
	if len(idx) == 0 {
		return flag
	}
	flag = true
	return flag
}

func isSupportKeyType() bool {
	return false
}

func DeepCopy(dst, src interface{}) error {
	var buffer bytes.Buffer
	if err := gob.NewEncoder(&buffer).Encode(src); err != nil {
		return err
	}

	return gob.NewDecoder(&buffer).Decode(dst)
}

func GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func writeFile(file string, content []byte) error {
	const dirPerm = 0700
	if err := os.MkdirAll(filepath.Dir(file), dirPerm); err != nil {
		return err
	}
	// Atomic write: create a temporary hidden file first
	// then move it into place. TempFile assigns mode 0600.
	f, err := ioutil.TempFile(filepath.Dir(file), "."+filepath.Base(file)+".tmp")
	if err != nil {
		return err
	}
	if _, err := f.Write(content); err != nil {
		f.Close()
		os.Remove(f.Name())
		return err
	}
	f.Close()
	return os.Rename(f.Name(), file)
}
