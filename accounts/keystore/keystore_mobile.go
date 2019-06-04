package keystore

// import (
// 	"encoding/json"
// 	"io/ioutil"
// )

type PackKeyStore struct {
	fileName    string
	VNTKeystore []encryptedKeyJSONV3 `json:"vnt_keystore"`
	keyJson     []byte
}

// func NewPackKeyStore(filename string) (*PackKeyStore, error) {
// 	keyjson, err := ioutil.ReadFile(fileName)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var ks *PackKeyStore
// 	err = json.Unmarshal(keyjson, ks)
// 	return nil, &PackKeyStore{
// 		fileName: filename,
// 	}
// }

// func (pk *PackKeyStore) Read() []byte {

// 	return keyjson
// }

// func (pk *PackKeyStore) ImportEthKeyStore(keyjson string) {

// }
