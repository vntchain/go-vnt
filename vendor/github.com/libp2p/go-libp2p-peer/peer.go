// package peer implements an object used to represent peers in the ipfs network.
package peer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"crypto/ecdsa"

	"github.com/vntchain/go-vnt/crypto"

	logging "github.com/ipfs/go-log" // ID represents the identity of a peer.
	// ic "github.com/libp2p/go-libp2p-crypto"
	b58 "github.com/mr-tron/base58/base58"
	mh "github.com/multiformats/go-multihash"
)

// MaxInlineKeyLength is the maximum length a key can be for it to be inlined in
// the peer ID.
//
// * When `len(pubKey.Bytes()) <= MaxInlineKeyLength`, the peer ID is the
//   identity multihash hash of the public key.
// * When `len(pubKey.Bytes()) > MaxInlineKeyLength`, the peer ID is the
//   sha2-256 multihash of the public key.
const MaxInlineKeyLength = 42

var log = logging.Logger("peer")

// ID is a libp2p peer identity.
type ID string

// MarshalJSON implement marshal for ID
func (id *ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.Pretty())
}

// UnmarshalJSON implement unmarshal for ID
func (id *ID) UnmarshalJSON(b []byte) error {
	var data string
	err := json.Unmarshal(b, &data)
	if err != nil {
		return err
	}
	*id, err = IDB58Decode(data)
	if err != nil {
		return err
	}
	return nil
}

// Pretty returns a b58-encoded string of the ID
func (id ID) Pretty() string {
	return IDB58Encode(id)
}

func (id ID) Loggable() map[string]interface{} {
	return map[string]interface{}{
		"peerID": id.Pretty(),
	}
}

// String prints out the peer.
//
// TODO(brian): ensure correctness at ID generation and
// enforce this by only exposing functions that generate
// IDs safely. Then any peer.ID type found in the
// codebase is known to be correct.
func (id ID) String() string {
	pid := id.Pretty()

	//All sha256 nodes start with Qm
	//We can skip the Qm to make the peer.ID more useful
	if strings.HasPrefix(pid, "Qm") {
		pid = pid[2:]
	}

	maxRunes := 64
	if len(pid) < maxRunes {
		maxRunes = len(pid)
	}
	return fmt.Sprintf("<peer.ID %s>", pid[:maxRunes])
}

func (id ID) ToString() string {
	pid := id.Pretty()

	//All sha256 nodes start with Qm
	//We can skip the Qm to make the peer.ID more useful
	if strings.HasPrefix(pid, "Qm") {
		pid = pid[2:]
	}

	maxRunes := 64
	if len(pid) < maxRunes {
		maxRunes = len(pid)
	}
	return fmt.Sprintf("%s", pid[:maxRunes])
}

// MatchesPrivateKey tests whether this ID was derived from sk
func (id ID) MatchesPrivateKey(pk *ecdsa.PrivateKey) bool {
	return id.MatchesPublicKey(&pk.PublicKey)
}

// MatchesPublicKey tests whether this ID was derived from pk
func (id ID) MatchesPublicKey(pk *ecdsa.PublicKey) bool {
	oid, err := IDFromPublicKey(pk)
	if err != nil {
		return false
	}
	return oid == id
}

// ExtractPublicKey attempts to extract the public key from an ID
//
// This method returns nil, nil if the peer ID looks valid but it can't extract
// the public key.
// 从ID中解压出公钥
func (id ID) ExtractPublicKey() (*ecdsa.PublicKey, error) {

	decoded, err := mh.Decode([]byte(id))
	if err != nil {
		return nil, err
	}
	if decoded.Code != mh.ID {
		return nil, nil
	}

	return crypto.DecompressPubkey(decoded.Digest)
}

// IDFromString cast a string to ID type, and validate
// the id to make sure it is a multihash.
func IDFromString(s string) (ID, error) {
	if _, err := mh.Cast([]byte(s)); err != nil {
		return ID(""), err
	}
	return ID(s), nil
}

// IDFromBytes cast a string to ID type, and validate
// the id to make sure it is a multihash.
func IDFromBytes(b []byte) (ID, error) {
	if _, err := mh.Cast(b); err != nil {
		return ID(""), err
	}
	return ID(b), nil
}

// IDB58Decode returns a b58-decoded Peer
func IDB58Decode(s string) (ID, error) {
	m, err := mh.FromB58String(s)
	if err != nil {
		return "", err
	}
	return ID(m), err
}

// IDB58Encode returns b58-encoded string
func IDB58Encode(id ID) string {
	return b58.Encode([]byte(id))
}

// IDHexDecode returns a hex-decoded Peer
func IDHexDecode(s string) (ID, error) {
	m, err := mh.FromHexString(s)
	if err != nil {
		return "", err
	}
	return ID(m), err
}

// IDHexEncode returns hex-encoded string
func IDHexEncode(id ID) string {
	return hex.EncodeToString([]byte(id))
}

// IDFromPublicKey returns the Peer ID corresponding to pk
// 从公钥中得到对应的ID
func IDFromPublicKey(pk *ecdsa.PublicKey) (ID, error) {
	b := crypto.CompressPubkey(pk)

	var alg uint64 = mh.SHA2_256
	if len(b) <= MaxInlineKeyLength {
		alg = mh.ID
	}
	hash, _ := mh.Sum(b, alg, -1)
	return ID(hash), nil
}

// IDFromPrivateKey returns the Peer ID corresponding to sk
func IDFromPrivateKey(sk *ecdsa.PrivateKey) (ID, error) {
	return IDFromPublicKey(&sk.PublicKey)
}

// IDSlice for sorting peers
type IDSlice []ID

func (es IDSlice) Len() int           { return len(es) }
func (es IDSlice) Swap(i, j int)      { es[i], es[j] = es[j], es[i] }
func (es IDSlice) Less(i, j int) bool { return string(es[i]) < string(es[j]) }
