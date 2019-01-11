// Copyright 2015 The go-ethereum Authors
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

// bootnode runs a bootstrap node for the VNT Discovery Protocol.
package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"net"
	"os"

	"context"

	"github.com/vntchain/go-vnt/cmd/utils"
	"github.com/vntchain/go-vnt/crypto"
	"github.com/vntchain/go-vnt/log"
	p2p "github.com/vntchain/go-vnt/vntp2p"
	// "log"/
)

func main() {
	var (
		dataDir    = flag.String("datadir", "./", "data directory for the database")
		listenAddr = flag.String("addr", "30301", "listen address")
		genKey     = flag.String("genkey", "", "generate a node key")
		// writeAddr   = flag.Bool("writeaddress", false, "write out the node's pubkey hash and quit")
		nodeKeyFile = flag.String("nodekey", "", "private key filename")
		nodeKeyHex  = flag.String("nodekeyhex", "", "private key as hex (for testing)")
		natdesc     = flag.String("nat", "none", "port mapping mechanism (any|none)")
		netrestrict = flag.String("netrestrict", "", "restrict network communication to the given IP networks (CIDR masks)")
		// runv5       = flag.Bool("v5", false, "run a v5 topic discovery bootnode")
		verbosity = flag.Int("verbosity", int(log.LvlInfo), "log verbosity (0-9)")
		vmodule   = flag.String("vmodule", "", "log verbosity pattern")
		nodeKey   *ecdsa.PrivateKey
		err       error
	)
	flag.Parse()

	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(*verbosity))
	glogger.Vmodule(*vmodule)
	log.Root().SetHandler(glogger)

	natm, err := p2p.NATParse(*natdesc)
	if err != nil {
		utils.Fatalf("-nat: %v", err)
	}
	switch {
	case *genKey != "":
		nodeKey, err = crypto.GenerateKey()
		if err != nil {
			utils.Fatalf("could not generate key: %v", err)
		}
		if err = crypto.SaveECDSA(*genKey, nodeKey); err != nil {
			utils.Fatalf("%v", err)
		}
		return
	case *nodeKeyFile == "" && *nodeKeyHex == "":
		// utils.Fatalf("Use -nodekey or -nodekeyhex to specify a private key")
	case *nodeKeyFile != "" && *nodeKeyHex != "":
		utils.Fatalf("Options -nodekey and -nodekeyhex are mutually exclusive")
	case *nodeKeyFile != "":
		if nodeKey, err = crypto.LoadECDSA(*nodeKeyFile); err != nil {
			utils.Fatalf("-nodekey: %v", err)
		}
	case *nodeKeyHex != "":
		if nodeKey, err = crypto.HexToECDSA(*nodeKeyHex); err != nil {
			utils.Fatalf("-nodekeyhex: %v", err)
		}
	}

	// if *writeAddr {
	// 	fmt.Printf("%v\n", discover.PubkeyID(&nodeKey.PublicKey))
	// 	os.Exit(0)
	// }

	var restrictList []*net.IPNet
	if *netrestrict != "" {
		restrictList, err = p2p.ParseNetlist(*netrestrict)
		if err != nil {
			utils.Fatalf("-netrestrict: %v", err)
		}
	}

	// 以太坊原有代码
	// addr, err := net.ResolveUDPAddr("udp", *listenAddr)
	// if err != nil {
	// 	utils.Fatalf("-ResolveUDPAddr: %v", err)
	// }
	// conn, err := net.ListenUDP("udp", addr)
	// if err != nil {
	// 	utils.Fatalf("-ListenUDP: %v", err)
	// }

	// realaddr := conn.LocalAddr().(*net.UDPAddr)
	// if natm != nil {
	// 	if !realaddr.IP.IsLoopback() {
	// 		go nat.Map(natm, nil, "udp", realaddr.Port, realaddr.Port, "ethereum discovery")
	// 	}
	// 	// TODO: react to external IP changes over time.
	// 	if ext, err := natm.ExternalIP(); err == nil {
	// 		realaddr = &net.UDPAddr{IP: ext, Port: realaddr.Port}
	// 	}
	// }

	// if *runv5 {
	// 	if _, err := discv5.ListenUDP(nodeKey, conn, realaddr, "", restrictList); err != nil {
	// 		utils.Fatalf("%v", err)
	// 	}
	// } else {
	// cfg := discover.Config{
	// 		PrivateKey:   nodeKey,
	// 		AnnounceAddr: realaddr,
	// 		NetRestrict:  restrictList,
	// 	}
	// if _, err := discover.ListenUDP(conn, cfg); err != nil {
	// 		utils.Fatalf("%v", err)
	// 	}
	// }

	// chenqiao修改代码
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, host, err := p2p.ConstructDHT(ctx, p2p.MakePort(*listenAddr), nodeKey, *dataDir, restrictList, natm)
	if err != nil {
		log.Error(fmt.Sprintf("constructDHT a: %s", err))
		return
	}

	log.Info("Listen for connect")
	log.Info(fmt.Sprintf("PID: %s, Addr: %s", host.ID(), host.Addrs()))

	// h := &vntprotocol.HostWrapper{
	// 	Host: host,
	// }

	// host.SetStreamHandler(p2p.PID, h.HandleStream)

	select {}
}
