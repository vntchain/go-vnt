package keystore_test

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/vntchain/go-vnt/accounts"
	"github.com/vntchain/go-vnt/accounts/keystore"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/crypto"
)

const (
	pubkey = "0455000cfef889b4750f4ae53cda20697aa372394c7b3f41be287eade4f5be703a93f646b40c82690f078cf24729df45ec37ad4d57b67caba7be2e7a6c7a1666bd"
	prikey = "c2c3546e88ca9e05dfe454306dbd4c31175e13b9ce46d1420f305a3c767698f6"
)

func TestGenerateKey(t *testing.T) {
	privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("private key %s\n", hex.EncodeToString(crypto.FromECDSA(privateKeyECDSA)))
	fmt.Printf("public key %s\n", hex.EncodeToString(crypto.FromECDSAPub(&privateKeyECDSA.PublicKey)))
	//prikeybytes := crypto.FromECDSA(privateKeyECDSA)
	pubkeybytes := crypto.FromECDSAPub(&privateKeyECDSA.PublicKey)
	res, err := keystore.Encrypt([]byte("111111"), pubkeybytes)
	if err != nil {
		panic(err)
	}
	fmt.Printf("res %s\n", res)
}

func TestEncryptAndDecrypt(t *testing.T) {
	pub, err := hex.DecodeString(pubkey)
	if err != nil {
		panic(err)
	}
	msg := []byte("111111")
	encRes, err := keystore.Encrypt(msg, pub)
	if err != nil {
		panic(err)
	}
	fmt.Printf("res %s\n", encRes)

	pri, err := hex.DecodeString(prikey)
	if err != nil {
		panic(err)
	}
	decRes, err := keystore.Decrypt(encRes, pri)
	if err != nil {
		panic(err)
	}
	if !bytes.Equal(decRes, msg) {
		t.Errorf("need msg %s,get msg %s", msg, decRes)
	}
}

func TestSecureEnclave(t *testing.T) {
	ksjson := `{
		"vnt_keystore":[
			{"address":"45dea0fb0bba44f4fcf290bba71fd57d7117cbb8","crypto":{"cipher":"aes-128-ctr","ciphertext":"b87781948a1befd247bff51ef4063f716cf6c2d3481163e9a8f42e1f9bb74145","cipherparams":{"iv":"dc4926b48a105133d2f16b96833abf1e"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":2,"p":1,"r":8,"salt":"004244bbdc51cadda545b1cfa43cff9ed2ae88e08c61f1479dbb45410722f8f0"},"mac":"39990c1684557447940d4c69e06b1b82b2aceacb43f284df65c956daf3046b85"},"id":"ce541d8d-c79b-40f8-9f8c-20f59616faba","version":3},
			{"address":"7ef5a6135f1fd6a02593eedc869c6d41d934aef8","crypto":{"cipher":"aes-128-ctr","ciphertext":"1d0839166e7a15b9c1333fc865d69858b22df26815ccf601b28219b6192974e1","cipherparams":{"iv":"8df6caa7ff1b00c4e871f002cb7921ed"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":8,"p":16,"r":8,"salt":"e5e6ef3f4ea695f496b643ebd3f75c0aa58ef4070e90c80c5d3fb0241bf1595c"},"mac":"6d16dfde774845e4585357f24bce530528bc69f4f84e1e22880d34fa45c273e5"},"id":"950077c7-71e3-4c44-a4a1-143919141ed4","version":3},
			{"address":"89057fb3c003605b1d7f9e94dfdd191c1265c9ff","crypto":{"cipher":"aes-128-ctr","ciphertext":"896ba766b2c962e7445feb0d68d104b0bfcca6df1548f09ae45bfc22f713d109","cipherparams":{"iv":"d8ccfce66a90c3336e7e762d4a59be4f"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"a5673385c9ebbd84eefc8c0b7216f4d03b204c8c049b9038b6cb98f62991aae0"},"mac":"578a78e35279fa3000bcbd365c6bdce5dd985713036a55990d06fb77250081f0"},"id":"6c644f99-aa3a-43e9-a837-81c96e5ea095","version":3}
		]
	}`
	ks := keystore.NewSecureEnclaveKeyStore([]byte(ksjson))
	wanted := []accounts.Account{
		accounts.Account{
			Address: common.HexToAddress("45dea0fb0bba44f4fcf290bba71fd57d7117cbb8"),
		},
		accounts.Account{
			Address: common.HexToAddress("7ef5a6135f1fd6a02593eedc869c6d41d934aef8"),
		},
		accounts.Account{
			Address: common.HexToAddress("89057fb3c003605b1d7f9e94dfdd191c1265c9ff"),
		},
	}
	wantedJson, _ := json.Marshal(wanted)
	got, _ := json.Marshal(ks.Accounts())
	if string(got) != string(wantedJson) {
		t.Errorf("Expected accounts %s, get: %s", wantedJson, got)
	}
}
