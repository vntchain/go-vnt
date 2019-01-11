package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/wavm"
	"github.com/vntchain/go-vnt/core/wavm/utils"
)

var (
	erc20Code = filepath.Join(basepath, "erc20/erc20.wasm")
	// erc20Code = filepath.Join(basepath, "qlang/PERC20.wasm")
	erc20Abi = filepath.Join(basepath, "erc20/abi.json")
)

// func newErc20VM(iscreated bool) *wavm.Wavm {
// 	return newWavm(erc20Code, erc20Abi, iscreated)
// }
// func TestERC20(t *testing.T) {
// 	initialSupply := new(big.Int)
// 	initialSupply.SetString("1000000000000", 10)
// 	tokenName := "bitcoin"
// 	tokenSymbol := "BTC"
// 	amount1 := new(big.Int)
// 	amount1.SetString("10000000", 10)
// 	//to1 := common.HexToAddress("0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
// 	to1 := common.HexToAddress("0x02")
// 	vm := newErc20VM(true)
// 	mutable := wavm.MutableFunction(vm.ChainContext.Abi, vm.Module)
// 	testCreate := func(initialSupply *big.Int, tokenName, tokenSymbol string) {
// 		funcName := ""
// 		res, err := vm.Apply(pack(vm, funcName, initialSupply, tokenName, tokenSymbol), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		t.Logf("create contract %s", res)
// 	}
// 	testTransfer := func(to common.Address, value *big.Int) {
// 		vm.ChainContext.IsCreated = false
// 		funcName := "transfer"
// 		res, err := vm.Apply(pack(vm, funcName, to, value), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		var flag bool
// 		unPack(vm, &flag, funcName, res)
// 		if flag != true {
// 			t.Errorf("unexpected value : want %t, got %t", true, flag)
// 		} else {
// 			t.Logf("%s success get %t", funcName, flag)
// 		}
// 		type TransferEvent struct {
// 			From  common.Address
// 			To    common.Address
// 			Value *big.Int
// 		}
// 		var ev TransferEvent
// 		logs := vm.ChainContext.StateDB.GetLogs(common.HexToHash("0x1111"))
// 		res, _ = logs[0].MarshalJSON()
// 		t.Logf("success get event json %s", res)
// 		unPack(vm, &ev, "Transfer", logs[0].Data)
// 		t.Logf("success get event %v", ev)
// 	}

// 	testGetTokenName := func() {
// 		vm.ChainContext.IsCreated = false
// 		funcName := "GetTokenName"
// 		res, err := vm.Apply(pack(vm, funcName), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		var name string
// 		unPack(vm, &name, funcName, res)
// 		if name != tokenName {
// 			t.Errorf("unexpected value : want %s, got %s", tokenName, name)
// 		} else {
// 			t.Logf("%s success get %s", funcName, name)
// 		}
// 	}

// 	testGetTotalSupply := func() {
// 		vm.ChainContext.IsCreated = false
// 		funcName := "GetTotalSupply"
// 		res, err := vm.Apply(pack(vm, funcName), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		var supply *big.Int
// 		unPack(vm, &supply, funcName, res)
// 		a := new(big.Int)
// 		if supply.Cmp(a.Mul(initialSupply, new(big.Int).SetUint64(100000000))) != 0 {
// 			t.Errorf("unexpected value : want %d, got %d", a, supply)
// 		} else {
// 			t.Logf("%s success get %d", funcName, supply)
// 		}
// 	}

// 	testGetSymbol := func() {
// 		vm.ChainContext.IsCreated = false
// 		funcName := "GetSymbol"
// 		res, err := vm.Apply(pack(vm, funcName), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		var sym string
// 		unPack(vm, &sym, funcName, res)
// 		if sym != tokenSymbol {
// 			t.Errorf("unexpected value : want %s, got %s", tokenSymbol, sym)
// 		} else {
// 			t.Logf("%s success get %s", funcName, sym)
// 		}
// 	}

// 	testGetDecimals := func() {
// 		vm.ChainContext.IsCreated = false
// 		funcName := "GetDecimals"
// 		res, err := vm.Apply(pack(vm, funcName), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		var deci *big.Int
// 		unPack(vm, &deci, funcName, res)
// 		if deci.Uint64() != 8 {
// 			t.Errorf("unexpected value : want %d, got %d", 8, deci)
// 		} else {
// 			t.Logf("%s success get %d", funcName, deci)
// 		}
// 	}

// 	testGetAmount := func(addr common.Address) *big.Int {
// 		vm.ChainContext.IsCreated = false
// 		funcName := "GetAmount"
// 		res, err := vm.Apply(pack(vm, funcName, addr), nil, mutable)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		var amount *big.Int
// 		unPack(vm, &amount, funcName, res)
// 		return amount
// 		// if amount != 4 {
// 		// 	t.Errorf("unexpected value : want %d, got %d", 4, amount)
// 		// } else {
// 		// 	t.Logf("%s success get %d", funcName, amount)
// 		// }
// 	}

// 	testCreate(initialSupply, tokenName, tokenSymbol)
// 	amount := testGetAmount(caller)
// 	a := new(big.Int)
// 	if amount.Cmp(a.Mul(initialSupply, new(big.Int).SetUint64(100000000))) != 0 {
// 		t.Errorf("unexpected value : want %d, got %d", a, amount)
// 	} else {
// 		t.Logf("%s success get %d", "GetAmount", amount)
// 	}
// 	now := T.Now()
// 	testTransfer(to1, amount1)
// 	after := T.Since(now)
// 	t.Logf("time %s", after.String())

// 	testGetTokenName()
// 	testGetTotalSupply()
// 	testGetSymbol()
// 	testGetDecimals()
// 	amount = testGetAmount(to1)
// 	if amount.Cmp(amount1) != 0 {
// 		t.Errorf("unexpected value : want %d, got %d", amount1, amount)
// 	} else {
// 		t.Logf("%s success get %d", "GetAmount", amount)
// 	}
// 	amount = testGetAmount(caller)
// 	b := new(big.Int)
// 	if new(big.Int).Add(amount, amount1).Cmp(b.Mul(initialSupply, new(big.Int).SetUint64(100000000))) != 0 {
// 		t.Errorf("unexpected value : want %d, got %d", b.Sub(b, amount1), amount)
// 	} else {
// 		t.Logf("%s success get %d", "GetAmount", amount)
// 	}
// }

var ercJsonPath = filepath.Join("", "erc20.json")

func TestERC(t *testing.T) {
	run(t, ercJsonPath)
}

func TestRlp(t *testing.T) {
	str := "0061736d0100789cec5b5b93aab816fe2ffdbaa7ea00eadee354cdc3020997567a82122e6f8674030aea8c17d453e7bf9fe2a2626bbbbbb46dfbece3836508c9ca5adfba26c0bf1fa4317b7ef8e30114b4a658060010679a040140db05fddb4b2006e0d4c1050017f47fbda4620080beb900e28b00e24b9ab5d137709a0138bf677dd24bedf7e25e80bec1df422bd0001430623a72156c9b4353406b00b9e5d90d01592c66b21e5205cd5d81287ed29c6563fc84acfd141410c89c054bd5734c892acd419f0b435f08404c737a4fa640389c2c17ae307b62361f798e0662b1563bbf27342774a887be608188f3fe6ed93f67aa1eba3553f41473e539c61ac0ca796ae53c99132ff1266e8d4c993adccc55b2b905ff381f2b71d95816fadc7242ed98dbf6e7349aabbecd26340054c8116772f43cbb31ccf92c6876329a4489679e6370aecda75b1a328b19f242aa9278c39bccc5f3be634e5c7b39794e88ed3afac4b5f5a9e77400829c9e66a9fa848ecc1e53358002a7bfac84841e6213a684980af19c29640500394d11191aee75402a70d3484fe38db8297a0a0609eff5b5a803affbe4bee3beee33989dbeeec3b4e6bfeeeb315bdbf46152d3274c214edf31c796d2cc6488fd0000bef721b3211d6786193c4a30ce1a9a08e34729c87e1612b1a588d89745b03411b02ee161f693c3c296453f336cf729a703184471ac146d786e64634072156a9319ade90d08007b0a195005ad3cdb5498dd1830255ed0005c4998c5cf3d63fe649bb5e7c41c8238549942d64ccef40b7dc93163bf66ea7484d79d0068e55a30007c69945f93bedde03a6b0d543ce132bd133b9eb396cb81586fe7f41099d0115e193d17d434ccc7e0a4b9a0004149a3c7323d273e687856dcdff94d24d6f2314f19e65a1ae7f76d9b84fea8033a14d7442071bfd6011d17f44d0535fc149292678dd6c8ca4f612415b48cbedde033de400a4a1e8dd8159a734fedc023e60b1a02e12c81ccbdccbea4b15160b3e70bffb4d49c5ee60bb66b2f790fc35476b67d72dff692be6df054351b2071de8ec69b71628e46dbf92d2a34383f410326a3a1a7eab15feb407b296ef92be385e63ac6185ad0add09f52c1f8a737223337212b68c9787b0f99136a933553d0cc0758b56ad5f5e2394ed0dac3b0960ad97a5d595e430b3fee68b3d04f81afcc535dc7e8f46d7ec200841213d2195882c3a1a10750dbebe30dde0da0be3f8ee7294063bfcf9cb000beefcf15172c851ffbe38c4c8edfa5899dc767dd7b2cfc41d60007801157fcb7380d64807630cefd08a4fadfc538110c6c8aa226c29354ef167d991b5db32d427bbbbe0480c7ad5108cd6f523da9f6f1ba5a6f476ec1af06a2143c76d04b30d782cccf0102901fb1aad6db617d37e6afec5adb5cc3b62dc3b8d21e56da69a5edb6c3f1762e063c97cbb5b088e7ca765d11a48a2cad7db9b66d59aaf7b7f2c8629517bfd20edad18edf004c30a57ab49da7818c515566094459c4afe7f40ee6e0d773b4b6f4bde42d104519a69dcce7fd9ab962b6c1d19a0620cb819496fe6f37044dc972148687df1e80460f7f3cd8b5e9631eb125886882665e57fbae49fa6e42946eef97c4433ac2d1532426aebd5c7bddddfdbedd489963aeb5c138b0796e2a45797fa0254d9ea922cfd44e368f3a95394c8dd372cd01151a6ba6eabc5b33177ea43dbe74d3e0391d1734247deeda7c9cd1d6f6836cd486ed98828682e64cd2be6b4a96a08db8727f42473b3e76b4377c36426a5bd1530451bfabedf83fe4757db8f686469c05d8ba675730c2e5189567155e16cc317778edf47028ef063fe9ea3c71cf8e58aee9f13431b8bedd9c6bc1a4d996c41f3b7dee702a82a8193e3b9dabe881462765e6b3646d0cf057d2c347f0749e1e7649ffae8b9bebc25a193def252b84ce8ea1a7e41f0587328dcc86af943ce445abc61b4774d6cf373961ec55ec6133663fe61ecaf67c3a2ee67867454c45e69d9d2566fcac9061c6bf97a0a92f58e518b159cd09b446385f259cbfd2be6b68f6e3f57a953c55dacb11fd7e0416a492c32af4bc1199bb357342857a248db8a9a6be9d43bbc566e216396a926f742e91f955de3e191b4a9d1dd17bc57f323d745697d84645be64bbf9fb79ad400ee3c0a7e9e11e17f76a856cc366c499befdaa6d26c6828e0cceb51b5cc5062af7e37925265c139b7073e874033b3dd8645fa3febed88f47db0390dbd84f7918d055d0ba2f9d6d43a7f24a941f680a0d9edafabe9dfe6ff8f7edf60265dd63d54874af7bfe7feb9ef240f3be0fb975be452cf6ecd9fa2abe98d4df617feeb15c72639fb88ca73373bb78a4aef8405d0c4fcaed6e0ee4575f693fb07d0870833a225c3047bffbc5edfde2ad074367e78e6be2f4217e74664e7df5d0fc8ecf117c4882a6ccfe9ad8dcb21edbbc1cf185f7938271a486fa8cf3b91b9d0b5d2ef355cee76e92039e68edfce778bf685edcbed471af516eae8bfca591ebe8c10c99229f94dbaf99abbedd187d255d7c044fe79ee56f5fb0b9ced95ffd3a679f7b79e8886ca7f76f1f523bfc4ae7c09be7ab029fdecf196f1b0bfa07cf383fc906e4fd177cefb5ed7e6d5bbcb8baf7bcf5d3ce9bd8e1da37b7d38fe0e9cc9c2551a139bdc9b95fd25cd02fb5c7d01754486f142bca179defcf9d8f9dc149aebd0c6962c4e6485f50ebf273b853b6929f139d7e8ef32178fc427b90ca8736479e4b7f1d3bbef1bb25958f01cedf1ffc9af8a063f5fae7be07f5d56ad75b9ec9563ed2ba8a3e9edf6187e7e863cfc60e656b7cc29ef12a3e79a96d5e749ec18b0bb6ba9f67fccae7190efef3cf87df1ea4713289e267567c7d1361a139cc344547644aa534220a8a688652d2e0a9924eb5211f3ecb66ecabe6581b8cd3b6a45b34db39c8cb852ba02949d02ab3828eca4d9f577a8b0a66ac0d2673662fa7d9d8fc13d40afd9e8deaa64252a604d1d300a65a6c2e9882a63da53970eda5ed3a7ae655e94b77f9431b1a0bafca4b6c866eb2acd2375c27943d47e4fa92f6bd93f3dfe45c7bd9a1821152e4e56f013d0dc4665bbdcb7a97f52eeb5dd6bbac77593f5f562d46a33e87446c6d7ea4fcd75bbbbe63d7c834472cfb47ddfdfed7d7122662deb66c56ed17cd6d5b8baaf8d83cf71aebd5fbb01e4796cdd0996bbc539ff91ae4e81a04c9568121b24a8cac125bcb3247c578aba94515db8a67ec4077d1fb75b78ff5467768b4e5adc7fd4ceee59b723b0772eb47e5dee915fdd5d5aab6499dee6b3b97dfb2f3e59e9d47bae971f19e5d6ddaded65edf8165f72d2cd11dcb3b961f8aa56519b5cbb13cf471f7687cde6289ba1c321e392295d79a215f104f0fb11d9cc696b4b1d5ace49080c3d1cfb0d6de8fb5609cb45bd39a8984ecfabb2d533a3fbe1edaf14fb12748def9d0a763af3fc955fe3e18fbaf11334c2cb0cb63c6610e3f1e33086a5997c78cf4fd31c3fc002cdf5d179838b90a96c6693bd5453c2ce5ed719f6333f107e49973e4e43e594ee7737d83719fe91bfa557ca377ba369130d9c82b373fa76ed0af5237bc6933dc557cc33ced1bc62606203c2c79393f57ad0eb0e5d1e95cc5a1b23e22925deaf9fc3a413fc0da3a6d532d5cdab269e9254e17d409a34bb0e76f8a3db92df6fae5d8f3ece13fff050000ffff"
	res := common.Hex2Bytes(str)
	fmt.Printf("res %+v\n", res)
	code := wavm.WasmCode{}
	decompress, err := utils.DeCompress(res)
	if err != nil {
		panic(err)
	}

	sep := []byte{0x7d} // 分割符'}',是{Code: "0x2da32be...", Abi: "0x23290da98acb032..."}的最后一位
	sepIdx := bytes.Index(decompress, sep)

	err = json.Unmarshal(decompress[:sepIdx+1], &code)
	if err != nil {
		panic(err)
	}
	fmt.Printf("code %+v\n", code)
}

func FromUInt64(n uint64) (out []byte) {
	more := true
	for more {
		b := byte(n & 0x7F)
		n >>= 7
		if n == 0 {
			more = false
		} else {
			b = b | 0x80
		}
		out = append(out, b)
	}
	return
}
