// Copyright 2019 The go-vnt Authors
// This file is part of the go-vnt library.
//
// The go-vnt library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-vnt library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-vnt library. If not, see <http://www.gnu.org/licenses/>.

package vntp2p

import (
	"context"
	"encoding/json"
	"path/filepath"
	"time"
	// "crypto/rand"
	// "flag"
	"net"
	"strings"

	"crypto/ecdsa"

	ds "github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/whyrusleeping/base32"
	// crypto "github.com/libp2p/go-libp2p-crypto"
	p2phost "github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p-kad-dht/opts"
	crypto2 "github.com/libp2p/go-libp2p-crypto"
	// net "github.com/libp2p/go-libp2p-net"
	"github.com/libp2p/go-libp2p-peer"
	bucket "github.com/libp2p/go-libp2p-kbucket"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/multiformats/go-multiaddr-net"
	"github.com/vntchain/go-vnt/log"

	// "github.com/vntchain/go-vnt/crypto"

	"fmt"

	"github.com/vntchain/go-vnt/crypto"
	"encoding/hex"
)

const (
	// PID vnt protocol basic id
	PID                 = "/p2p/1.0.0"
	persistDataInterval = 10 * time.Second
)

const BootnodeCon = 1

type blankValidator struct{}

func (blankValidator) Validate(_ string, _ []byte) error        { return nil }
func (blankValidator) Select(_ string, _ [][]byte) (int, error) { return 0, nil }


// PersistentData record network info
type PersistentData struct {
	PrivKey   []byte
	PeerInfos []pstore.PeerInfo
	KBuckets  []peer.ID
}

// GetKBuckets return k-bucket data as array -- by hx
func GetKBuckets(rt *bucket.RoutingTable) []peer.ID {
	var m []peer.ID

	for i := range rt.Buckets {
		m = append(m, rt.Buckets[i].Peers()...)
	}

	return m
}


// GetPersistentData to save network info periodly
func GetPersistentData(dht *dht.IpfsDHT) *PersistentData {
	pd := &PersistentData{}

	// try to get privateKey from peerstore
	privKey := dht.Host().Peerstore().PrivKey(dht.PeerID())
	bDump, _ := privKey.Raw()
	k := hex.EncodeToString(bDump)
	pd.PrivKey = []byte(k)

	pd.PeerInfos = pstore.PeerInfos(dht.Host().Peerstore(), dht.Host().Peerstore().Peers())
	pd.KBuckets = GetKBuckets(dht.RoutingTable())

	return pd
}

// SaveData save PersistentData to local database
func SaveData(ctx context.Context, dht *dht.IpfsDHT, vdb *LevelDB, key string, value []byte) {
	//fmt.Printf("saveData -->, key = %s, value = %v\n", key, string(value))
	//dht.datastore.Put(mkDsKey(key), value)
	dsKey := ds.NewKey(base32.RawStdEncoding.EncodeToString([]byte(key)))
	err := vdb.Put(dsKey, value)
	//err := dht.PutValue(ctx, key, value)
	if err != nil {
		log.Error("Failed to save data", "error", err)
	}
}


func recoverPersistentData(vdb *LevelDB) *PersistentData {
	pd := &PersistentData{}
	pdKey := ds.NewKey(base32.RawStdEncoding.EncodeToString([]byte("/PersistentData")))
	pdValue, err := vdb.Get(pdKey)
	if err != nil {
		// don't need to care about err != nil
		return nil
	}
	//fmt.Printf("R- pdValue = %v\n", pdValue.([]byte))
	err = json.Unmarshal(pdValue, pd)
	if err != nil {
		log.Error("recoverPersistentData", "unmarshal pd error", err)
		return nil
	}
	return pd
}

// ConstructDHT create Kademlia DHT
func ConstructDHT(ctx context.Context, listenstring string, nodekey *ecdsa.PrivateKey, datadir string, restrictList []*net.IPNet, natm libp2p.Option) (*dht.IpfsDHT, p2phost.Host, *LevelDB, error) {

	var pd *PersistentData
	var vntp2pDB *LevelDB
	var err error
	// if datadir is empty, it means don't need persistentation
	if datadir != "" {
		dbpath := filepath.Join(datadir, "vntdb")
		vntp2pDB, err = GetDatastore(dbpath)
		if err != nil {
			log.Error("ConstructDHT", "getDatastore error", err, "dbpath", dbpath)
			return nil, nil, nil, err
		}
		pd = recoverPersistentData(vntp2pDB)
	}

	var privKey crypto2.PrivKey = nil
	if nodekey == nil && pd != nil {
		k := string(pd.PrivKey)
		bDump, err := hex.DecodeString(k)
		if err != nil {
			log.Error("ConstructDHT", "decode key error", err)
			return nil, nil, nil, err
		}
		nodekey, err = crypto.ToECDSA(bDump)
		if err != nil {
			log.Error("ConstructDHT", "toECDSA error", err)
			return nil, nil, nil, err
		}
	}

	if nodekey != nil {
		privKey, _, err = crypto2.ECDSAKeyPairFromKey(nodekey)
		if err != nil {
			log.Error("Bad private key:", "err", err)
			return nil, nil, nil, err
		}
	}

	//if nodekey == nil && pd != nil {
	//	// try to recover nodekey from database
	//	k := string(pd.PrivKey)
	//	bDump, err := hex.DecodeString(k)
	//	if err != nil {
	//		log.Error("ConstructDHT", "decode key error", err)
	//		return nil, nil, err
	//	}
	//	privKey, err := crypto.ToECDSA(bDump)
	//	if err != nil {
	//		log.Error("ConstructDHT", "toECDSA error", err)
	//		return nil, nil, err
	//	}
	//	nodekey = privKey
	//} // host private key recover finished

	host, err := constructPeerHost(ctx, listenstring, privKey, restrictList, natm)
	if err != nil {
		log.Error("ConstructDHT", "constructPeerHost error", err)
		return nil, nil, nil, err
	}

	var vdht *dht.IpfsDHT

	if vntp2pDB != nil {
		vdht, err = dht.New(
			ctx, host,
			dhtopts.NamespacedValidator("v", blankValidator{}),
			dhtopts.Datastore(vntp2pDB),
		)
	} else {
		vdht, err = dht.New(
			ctx, host,
			dhtopts.NamespacedValidator("v", blankValidator{}),
		)
	}

	hostAddr, err := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", ToString(host.ID())))

	addr := host.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)
	log.Info("ConstructDHT", "addr", fullAddr)

	// fmt.Println("TEST: ", ip)

	if pd != nil {
		recoverPeerInfosAndBuckets(ctx, pd, host, vdht)
	}

	if vntp2pDB != nil {
		go loop(ctx, vdht, vntp2pDB)
	}

	return vdht, host, vntp2pDB, err
}

// some loop handler for p2p itself can be put here
func loop(ctx context.Context, vdht *dht.IpfsDHT, vdb *LevelDB) {
	var persistData = time.NewTicker(persistDataInterval)
	for {
		<-persistData.C
		go persistDataPeriodly(ctx, vdht, vdb)
	}
}

// persist data unified entrance, both for bootnode and membernode
func persistDataPeriodly(ctx context.Context, vdht *dht.IpfsDHT, vdb *LevelDB) {
	pd := GetPersistentData(vdht)
	/* fmt.Printf("host privKey is: %v \n", string(pd.PrivKey))
	fmt.Printf("peerInfos is: \n")
	for i := range pd.PeerInfos {
		fmt.Printf("-->peerID = %v, peerAddr = %v\n", pd.PeerInfos[i].ID, pd.PeerInfos[i].Addrs)
	}
	fmt.Printf("buckets is: \n")
	for i := range pd.KBuckets {
		fmt.Println("----> real e:", pd.KBuckets[i])
	} */
	pdByte, err := json.Marshal(pd)
	if err != nil {
		log.Error("persistDataPeriodly", "marshal error", err)
		return
	}
	//log.Info("persistDataPeriodly TIME TO PERSIST DATA")
	SaveData(ctx, vdht, vdb, "/PersistentData", pdByte)
}

func recoverPeerInfosAndBuckets(ctx context.Context, pd *PersistentData, host p2phost.Host, vdht *dht.IpfsDHT) {

	/* fmt.Printf("R- host privKey is: %v \n", string(pd.PrivKey))
	fmt.Printf("R- peerInfos is: \n")
	for i := range pd.PeerInfos {
		fmt.Printf("R- -->peerID = %v, peerAddr = %v\n", pd.PeerInfos[i].ID, pd.PeerInfos[i].Addrs)
	}
	fmt.Printf("R- buckets is: \n")
	for i := range pd.KBuckets {
		fmt.Println("R- ----> real e:", pd.KBuckets[i])
	} */

	for i := range pd.PeerInfos {
		host.Peerstore().AddAddrs(pd.PeerInfos[i].ID, pd.PeerInfos[i].Addrs, pstore.PermanentAddrTTL)
	}
	for i := range pd.KBuckets {
		vdht.Update(ctx, pd.KBuckets[i])
	}
}

func constructPeerHost(ctx context.Context, listenstring string, nodekey crypto2.PrivKey, restrictList []*net.IPNet, natm libp2p.Option) (p2phost.Host, error) {
	var options []libp2p.Option
	if nodekey != nil {
		//priv, _, err := crypto2.ECDSAKeyPairFromKey(nodekey)
		//if err != nil {
		//	return nil, err
		//}
		options = append(options, libp2p.ListenAddrStrings(listenstring), libp2p.Identity(nodekey))
	} else {
		options = append(options, libp2p.ListenAddrStrings(listenstring))
	}

	options = append(options, libp2p.FilterAddresses(restrictList...))
	if natm != nil {
		options = append(options, natm)
	}

	return libp2p.New(ctx, options...)
}

func MakePort(port string) string {

	return "/ip4/0.0.0.0/tcp/" + port
}

func GetIPfromAddr(a ma.Multiaddr) string {
	// TODO:
	// 将ma地址转换为可读的IP

	_, ip, _ := manet.DialArgs(a)

	return strings.Split(ip, ":")[0]
}

func GetAddr(target string) (ma.Multiaddr, peer.ID, error) {
	ipfsaddr, err := ma.NewMultiaddr(target)
	if err != nil {
		log.Error("GetAddr", "NewMultiaddr Err", err)
		return nil, "", err
	}

	pid, err := ipfsaddr.ValueForProtocol(ma.P_IPFS)
	if err != nil {
		log.Error("GetAddr", "ValueForProtocol Err", err)
		return nil, "", err
	}

	peerid, err := peer.IDB58Decode(pid)
	if err != nil {
		log.Error("GetAddr", "IDB58Decode Err", err)
		return nil, "", err
	}

	// Decapsulate the /ipfs/<peerID> part from the target
	// /ip4/<a.b.c.d>/ipfs/<peer> becomes /ip4/<a.b.c.d>
	targetPeerAddr, _ := ma.NewMultiaddr(
		fmt.Sprintf("/ipfs/%s", peer.IDB58Encode(peerid)))
	targetAddr := ipfsaddr.Decapsulate(targetPeerAddr)

	return targetAddr, peerid, nil
}

func ParseNetlist(s string) ([]*net.IPNet, error) {
	ws := strings.NewReplacer(" ", "", "\n", "", "\t", "")
	masks := strings.Split(ws.Replace(s), ",")
	l := []*net.IPNet{}
	for _, mask := range masks {
		if mask == "" {
			continue
		}
		// n := net.ParseIP(mask)
		_, n, err := net.ParseCIDR(mask)
		if err != nil {
			return nil, err
		}
		l = append(l, n)
	}

	return l, nil
}

func NATParse(spec string) (libp2p.Option, error) {
	switch spec {
	case "", "none", "off":
		return nil, nil
	case "any":
		return libp2p.NATPortMap(), nil
	default:
		return nil, fmt.Errorf("unknow mechanism")
	}
}
