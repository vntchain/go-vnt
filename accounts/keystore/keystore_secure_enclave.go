package keystore

import (
	"crypto/rand"
	"fmt"

	"encoding/json"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/crypto"
	"github.com/vntchain/go-vnt/crypto/ecies"
)

func Encrypt(msg []byte, pub []byte) ([]byte, error) {
	ecdsapub, err := crypto.UnmarshalPubkey(pub)
	if err != nil {
		return nil, err
	}
	eciespub := ecies.ImportECDSAPublic(ecdsapub)
	return ecies.Encrypt(rand.Reader, eciespub, msg, nil, nil)
}

func Decrypt(msg []byte, pri []byte) ([]byte, error) {
	ecdaspri, err := crypto.ToECDSA(pri)
	if err != nil {
		return nil, err
	}
	eciespri := ecies.ImportECDSA(ecdaspri)
	return eciespri.Decrypt(msg, nil, nil)
}

type keyStoreSecureEnclave struct {
	eth_keystore map[common.Address]encryptedKeyJSONV3 //keystore的集合
}

func (ks keyStoreSecureEnclave) GetKey(addr common.Address, filename, auth string) (*Key, error) {
	// Load the key from the keystore and decrypt its contents
	ksjson, err := json.Marshal(ks.eth_keystore[addr])
	if err != nil {
		return nil, err
	}
	key, err := DecryptKey(ksjson, auth)
	if err != nil {
		return nil, err
	}
	// Make sure we're really operating on the requested key (no swap attacks)
	if key.Address != addr {
		return nil, fmt.Errorf("key content mismatch: have account %x, want %x", key.Address, addr)
	}
	return key, nil
}

func (ks keyStoreSecureEnclave) StoreKey(filename string, key *Key, auth string) error {
	// pub, err := crypto.UnmarshalPubkey(ks.publicKey)
	// if err != nil {
	// 	return err
	// }
	// eciespub := ecies.ImportECDSAPublic(pub)
	// file, err := ecies.Encrypt(rand.Reader, eciespub, []byte(ks.keyJson), nil, nil)
	// if err != nil {
	// 	return err
	// }
	// return writeKeyFile(filename, file)
	return nil
}

func (ks keyStoreSecureEnclave) JoinPath(filename string) string {
	// if filepath.IsAbs(filename) {
	// 	return filename
	// }
	// return filepath.Join(ks.keysDirPath, filename)
	return ""
}
