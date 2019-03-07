// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"crypto/rand"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vntchain/go-vnt/params"
)

const (
	ipcAPIs  = "admin:1.0 bp:1.0 core:1.0 debug:1.0 dpos:1.0 net:1.0 personal:1.0 rpc:1.0 shh:1.0 txpool:1.0 vnt:1.0"
	httpAPIs = "bp:1.0 core:1.0 net:1.0 personal:1.0 rpc:1.0 vnt:1.0"
)

// Tests that a node embedded within a console can be started up properly and
// then terminated by closing the input stream.
func TestConsoleWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"

	// Start a gvnt console, make sure it's cleaned up and terminate the console
	gvnt := runGvnt(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--coinbase", coinbase, "--shh",
		"console")

	// Gather all the infos the welcome message needs to contain
	gvnt.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	gvnt.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	gvnt.SetTemplateFunc("gover", runtime.Version)
	gvnt.SetTemplateFunc("gvntver", func() string { return params.Version })
	gvnt.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	gvnt.SetTemplateFunc("apis", func() string { return ipcAPIs })

	// Verify the actual welcome message to the required template
	gvnt.Expect(`
Welcome to the Gvnt JavaScript console!

instance: Gvnt/v{{gvntver}}/{{goos}}-{{goarch}}/{{gover}}
coinbase: {{.Coinbase}}
at block: 0 ({{niltime}})
 datadir: {{.Datadir}}
 modules: {{apis}}

> {{.InputLine "exit"}}
`)
	gvnt.ExpectExit()
}

// Tests that a console can be attached to a running node via various means.
func TestIPCAttachWelcome(t *testing.T) {
	// Configure the instance for IPC attachement
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	var ipc string
	if runtime.GOOS == "windows" {
		ipc = `\\.\pipe\gvnt` + strconv.Itoa(trulyRandInt(100000, 999999))
	} else {
		ws := tmpdir(t)
		defer os.RemoveAll(ws)
		ipc = filepath.Join(ws, "gvnt.ipc")
	}
	// Note: we need --shh because testAttachWelcome checks for default
	// list of ipc modules and shh is included there.
	gvnt := runGvnt(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--coinbase", coinbase, "--shh", "--ipcpath", ipc)

	time.Sleep(2 * time.Second) // Simple way to wait for the RPC endpoint to open
	testAttachWelcome(t, gvnt, "ipc:"+ipc, ipcAPIs)

	gvnt.Interrupt()
	gvnt.ExpectExit()
}

func TestHTTPAttachWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P
	gvnt := runGvnt(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--coinbase", coinbase, "--rpc", "--rpcport", port, "--rpcapi", "db,core,net,vnt,personal,bp")

	time.Sleep(2 * time.Second) // Simple way to wait for the RPC endpoint to open
	testAttachWelcome(t, gvnt, "http://localhost:"+port, httpAPIs)

	gvnt.Interrupt()
	gvnt.ExpectExit()
}

func TestWSAttachWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P

	gvnt := runGvnt(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--coinbase", coinbase, "--ws", "--wsport", port, "--wsapi", "db,core,net,vnt,personal,bp")

	time.Sleep(2 * time.Second) // Simple way to wait for the RPC endpoint to open
	testAttachWelcome(t, gvnt, "ws://localhost:"+port, httpAPIs)

	gvnt.Interrupt()
	gvnt.ExpectExit()
}

func testAttachWelcome(t *testing.T, gvnt *testgvnt, endpoint, apis string) {
	// Attach to a running gvnt note and terminate immediately
	attach := runGvnt(t, "attach", endpoint)
	defer attach.ExpectExit()
	attach.CloseStdin()

	// Gather all the infos the welcome message needs to contain
	attach.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	attach.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	attach.SetTemplateFunc("gover", runtime.Version)
	attach.SetTemplateFunc("gvntver", func() string { return params.Version })
	attach.SetTemplateFunc("coinbase", func() string { return gvnt.Coinbase })
	attach.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	attach.SetTemplateFunc("ipc", func() bool { return strings.HasPrefix(endpoint, "ipc") })
	attach.SetTemplateFunc("datadir", func() string { return gvnt.Datadir })
	attach.SetTemplateFunc("apis", func() string { return apis })

	// Verify the actual welcome message to the required template
	attach.Expect(`
Welcome to the Gvnt JavaScript console!

instance: Gvnt/v{{gvntver}}/{{goos}}-{{goarch}}/{{gover}}
coinbase: {{coinbase}}
at block: 0 ({{niltime}}){{if ipc}}
 datadir: {{datadir}}{{end}}
 modules: {{apis}}

> {{.InputLine "exit" }}
`)
	attach.ExpectExit()
}

// trulyRandInt generates a crypto random integer used by the console tests to
// not clash network ports with other tests running cocurrently.
func trulyRandInt(lo, hi int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(hi-lo)))
	return int(num.Int64()) + lo
}
