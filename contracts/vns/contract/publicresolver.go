package contract

import (
	"math/big"
	"strings"

	hubble "github.com/vntchain/go-vnt"
	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/accounts/abi/bind"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/event"
)

// PublicResolverABI is the input ABI used to generate the binding from.
const PublicResolverABI = `[{"name":"PublicResolver","constant":false,"inputs":[{"name":"vnsAddr","type":"address","indexed":false}],"outputs":[],"type":"constructor"},{"name":"content","constant":true,"inputs":[{"name":"node","type":"string","indexed":false}],"outputs":[{"name":"output","type":"string","indexed":false}],"type":"function"},{"name":"ABIRecord","constant":true,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"contentTypes","type":"uint256","indexed":false}],"outputs":[{"name":"output","type":"string","indexed":false}],"type":"function"},{"name":"ABIContentType","constant":true,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"contentTypes","type":"uint256","indexed":false}],"outputs":[{"name":"output","type":"uint256","indexed":false}],"type":"function"},{"name":"setABI","constant":false,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"contentType","type":"uint256","indexed":false},{"name":"data","type":"string","indexed":false}],"outputs":[],"type":"function"},{"name":"pubkeyX","constant":true,"inputs":[{"name":"node","type":"string","indexed":false}],"outputs":[{"name":"output","type":"string","indexed":false}],"type":"function"},{"name":"setContent","constant":false,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"hash","type":"string","indexed":false}],"outputs":[],"type":"function"},{"name":"name","constant":true,"inputs":[{"name":"node","type":"string","indexed":false}],"outputs":[{"name":"output","type":"string","indexed":false}],"type":"function"},{"name":"setName","constant":false,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"name","type":"string","indexed":false}],"outputs":[],"type":"function"},{"name":"pubkeyY","constant":true,"inputs":[{"name":"node","type":"string","indexed":false}],"outputs":[{"name":"output","type":"string","indexed":false}],"type":"function"},{"name":"setPubkey","constant":false,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"x","type":"string","indexed":false},{"name":"y","type":"string","indexed":false}],"outputs":[],"type":"function"},{"name":"setText","constant":false,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"key","type":"string","indexed":false},{"name":"value","type":"string","indexed":false}],"outputs":[],"type":"function"},{"name":"supportsInterface","constant":true,"inputs":[{"name":"interfaceID","type":"string","indexed":false}],"outputs":[{"name":"output","type":"bool","indexed":false}],"type":"function"},{"name":"addr","constant":false,"inputs":[{"name":"node","type":"string","indexed":false}],"outputs":[{"name":"output","type":"address","indexed":false}],"type":"function"},{"name":"setAddr","constant":false,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"addr","type":"address","indexed":false}],"outputs":[],"type":"function"},{"name":"text","constant":true,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"key","type":"string","indexed":false}],"outputs":[{"name":"output","type":"string","indexed":false}],"type":"function"},{"name":"AddrChanged","anonymous":false,"inputs":[{"name":"node","type":"string","indexed":true},{"name":"a","type":"address","indexed":false}],"type":"event"},{"name":"ContentChanged","anonymous":false,"inputs":[{"name":"node","type":"string","indexed":true},{"name":"hash","type":"string","indexed":false}],"type":"event"},{"name":"NameChanged","anonymous":false,"inputs":[{"name":"node","type":"string","indexed":true},{"name":"name","type":"string","indexed":false}],"type":"event"},{"name":"ABIChanged","anonymous":false,"inputs":[{"name":"node","type":"string","indexed":true},{"name":"contentType","type":"uint256","indexed":true}],"type":"event"},{"name":"PubkeyChanged","anonymous":false,"inputs":[{"name":"node","type":"string","indexed":true},{"name":"x","type":"string","indexed":false},{"name":"y","type":"string","indexed":false}],"type":"event"},{"name":"TextChanged","anonymous":false,"inputs":[{"name":"node","type":"string","indexed":true},{"name":"indexedKey","type":"string","indexed":true},{"name":"key","type":"string","indexed":false}],"type":"event"},{"name":"owner","constant":false,"inputs":[{"name":"node","type":"string","indexed":false}],"outputs":[{"name":"output","type":"address","indexed":false}],"type":"call"}]`

// PublicResolverBin is the compiled bytecode used for deploying new contracts.
const PublicResolverBin = `0x0161736db912990100789cd45a0d705cd575feeebbefbe7ddaa79556c83f927f9fd706dba99056bf365480d73f0105703cc62eb84c2a3fed3e498b57bbcafe0889c6926248b00b4e4c9c6668dada5048c6931886c96452dbd06966d24c81490a691b3299928999109a9686129c1f9ac0d039f7bdb7efedca32729153c033ef9d7beef9f9ceb9e7de77579cffb9facdd31f326015461800d619decba6b057999cc45e3689bdcad4149bda2b26a726a7a6b017d80b36454cece553f464536c4a794159c0edec98b6d31eb3f3458006d15bf3e9a27d6bba38bc2397ce16ed3c18b1eb77da566a06574ba452d75b052834d0777774f7f46f191905f78789542a38bca534e00cc38954ea467ba22f3b98834a8cc6be6cba98b632e9bbec3fb2f2696b206317206846e4eeccda7947ade67abb788b9d4dd97968726edbc74b56c6993312a9547ecbb0951db2530811a76e4b2e5bb4b3c50aa6b1dd1ab12b382ed26c25d2e18c8774735f8578644769609f3de1f1746974973d5ef4390b1a8510a8a909d13fcef59a1a9debaa60a38c31c119146d219f6289d70fa9c614bcd7cb3546e8e76c8136628fe4f2130a8cfefe61db1aed1fb00a3667e1fefe9455b4faed6c8a2bd1949dcc58793bf5e15236594ce7b25858b7a33490492777da855c66ccce637143a1343a9acb170b7db45e8356d246936aa5527934870a7691928525a1a493212c0d17eca29b2e2c53b3d6888de52447d9c28a9ac4e6be9d7632974f61651d25c491dc35316ac3d4c8dce63eac0a8dcaccdc86984bedc1ea9a825d741286356ad11e2fe272324bf9c215e1e79a1b15182bb0199b11fde4f4f4b4431d70a8e8ddf4363e27983a993879687a1a8fc590d0120f3a244bd42410bd97644c24349325b48a31aae603e3c4771c138a54f3c5942a334a424b3ced88f22a319e10656b6f1d2c5b53024e1e9db3930ab5fb0f9dd7dae7cfcf3eebb14315ece9f34b3fe44b57420a5541aab4f6d4dcd3756ab67495c5124fcc3d2fa89a0f82baa75c11154b7bc2675f4485240e7a4555a19638725eb6d18ce81788881e945230d1a3c895301e648c4d399309b68a299b4c24b04e91eb11bd8fd81138ac7b66b20ece64dd3f93756426ebf33359b2281dd62a6698cc384dc8564f4b64dfc4be185b236917ab0c41ee2c9391baac94187a948d264bc4ef20d6bf12ab472181757724b05e91bb680327862919ffe0339a24e3299f11958caffb0c5d324efa0cac571ef547dfc41d0e4013c69702d0cd7d312588dc5412fa1d9b1f3aac919a92983e28412e3015827cd6811c351547a2c9544c653d8f6ee05036998a89e82169e9cfa603f94bfc48a64efa30028949c0643d4e624c982c7abfebde74811adfb8f80c3bf9c53a2772ca7797975efe7c557af9d355e9e5dfac4a2f3f55955efe4430bdfcc4fb3fbd8f96d37bf8bda4773d7fc889d54def092fbd7137bd15e56c56977353753947abcb59af2ee7ca5a7ebfa6f74439bd9fa94eef7714c683e955664baf92d858597d8a3c2002d5a7b815fc84cf702af884cff097a8c989f6e972b434928b11e33d4aad23eb6d86c4cf28bc5548bc29dfaab249d9e4f04c16bdc739763f1c01df44276f4c3559f4b3ce075166450a6c8bc0e48ebe93b46d116e986a82458f4ccb2af12d5d1f81612aa6daa3443d8786a9045657adf7e62bd8fc8396cde9df7736a72f229baaf162e516e2155b8807b610f7b710afd842dcdb42dce4e52dc4e7b085944d264bb0e8a7a79da47cca59bf72ec6ee4b36cb7a75d52e9719640eebc075cd4e59db76c66591076799ccc36298bc1f8daff5f5e66c47ad28ff5092f5653897eae3adab5e70b880c9cf2627bea7d17db293fb6a7cab1b1e8d1ead88c5f328ec4b46e4c943aba7bdabb3b37b46fecd8b0a1278ef878bc7d70307995b501f1f1ce81ce81ee0d291bf1f1d4c6ce8d57a592dd888ff75cd53ed8d9d5d98ef8784747bcd31ae8ee417c3cb9b1e7aa78476727e2e3dd57a5da535d9d496ccf154df77770a59f760089b7746359540596aac0152ad0ae02bd2a70bd0aec5449e2c11ae326cce77f6fd71bad297ba034d49fce0ee67e1e0254c957193d6bb1da95fbb2fbdec400e5a70036005054c13fa102bc9734ae07f0160016178451fb0d3d98c2992b2c049f560165a1c7d004bf87181ff21821c10f12a3d763e882df4f8c9b3d468de0478861798cb0e09f27468101cf1263b5e06709d2731ea43b19c03a952906bc4a026b047f90045e0710ba4159a3dfcb805e3903fd3e063453886b74fdef4843d29ba0ae61c06f69b02ab45169d61f600079509642ff82a7bf4cd58f79f4725d7f8401df229515b5fa2906a4885ed9148a2a2b7d9f2b1d9f0eadfa3e57ead07fce80bfa081b92e14554ca9f4b61c0794cc8092a903ea3f31e0711a2d0ee94aa3fe154f7021f4931ebd4885fa6b06bc42b969d19e63404897d97a83381d35b4fa0d00da1420bc811e9b15008600b002c0ed0ac05eaf65f4ea67911af9deae009166499217688d0aa0b03a3246b63334f526451021542e79b24cd65fc101f65b22a3e4b81b403d15d4ef1cef909606018c93a5777c4beff89624d9202d298ce8cbc8d462d226538c05c29153bb3c2f2a73dd34921b727e0fb9d1a59505d28f4b9ff4e9856d8ab33158b8427d138023a41e09a84702ea0eed00adf74d3591a9a867cab81cc0e714e0af2893572e3aa500cd004ed07047641da717d521e48a4d01685581708f0a2c4ed0e3267adc4a7bbaee5f14e041007f4bcac372b1e26460582e56dd08074e03a0d256ee22a391b52479174d37e9ab371e5557e368572fa94cd21e41ed729abf953972b7925cbdd4fe63d2ae6f21eeed92ec24724f9260a85915d84199f9af5094fdbb7eb5e2ec24f61fd0b7284e81b3ff54f51b15a7a8d9ab3a16ad53812f017886ec1c72fc1d227f11e9ef101968a6c8e354ddf757c50fd43da702df00407b55f9626f39b82f0682eb3aaacae08ec9e016dd2b1c9fdf23b9af3ae25f959a5bb8474a9f3790cfc767fa6cd51c9f7426284ffa3e25d9143eaaaf8ecb94b648af7fef787d4573bcbe4892cf3a0acf4a5d69e6d9b2d7bde4f5bb33bd1e0e013731e065d27cd1f7ea90b790ab1729c74d7547d7ad8e1e5557d74a10ba047196402ca122a9df4be3974872c9388d6fa0f1cbbd4eb1d7fd9bee78798d4cbf26779ff4f29aefe5b559bdbc21bd3c5ff6724e7a79a2ece5d7ae97a5ebc3c08bf49122d33a5bb0d679d3f402a9e9eec70c2ff317526e7e46b931d88ce4346a06b08c8a9014563886a4518746e36a57a09604ae0808383496260ce06bc425812e17529723487bdca34ffab484041a6f380fa45f1ac05a00f48550b6053cbaf4bd3e1f4b97d73ade5790f7ddaef7dd014197becfa7a57793c6b7cdf40e44ebdccfff8895a41b00f00535e2b2f25676c82ed067a011c04a007f10b839b4d06103a003c0b501fe66157844054eaac0e3aacf7f58d0271f303420a2051442801502ee08019990cffe131d38a8038fe8c0d774e05b3af0aceecfb7cfc23f5903bc5003fca20640186808d36ef3e7a766e16f09030f87816f84815301feed06f094017cdb00fed1009e37807f36fcf90db5c083b5c0c3b5c023b5c0c95ae0f15a7f9eb6f56500e8d3b9da1ddfe0be87dc778ffb4e00b8dc3da2e9aca7b387f8d300e8f8dec7013aa7e9724867e27755808eb64f0980ceab160da063e7650da093e4be104007c20f7580b6ecda30409b4a1800ed85980150c95f670054d9e70c800a72692d40b5f65c200e4c73af30ac8181bc3dc61ad8e5758d82d745172fad6b60bb1753c583d7f5355ed77cb5f18786d204f06ef43502ea22f435f23ae202a25ef2b435e075d71a8601841a996138937a445aa0c14603a869650dec32adbd11080bd0cb6805711c63d73503b5adcc19ac6dbeaed934808880c3202f758e014f82b001f55d0189e872d6de4832b7197b0ca0c19d13347759604e008d0163c233b64004141692c2eec58ef8a24adf24bcb84b06d1d405a5c987d05cd6328025864ce5d24a5fa40cbc13f56eeb9974d6ced1750f3a93d79adfd545e8c59c4a626dbb0b76bed076a79d2e58e9bb86adec50dbd65cb23462678b85b6a15c5b219f6c1b4a17874b03adc9dc48db58b6981cb6d2d9b6a1dc9563d9625b32972de6ad64b1d036962d944701762b2eb5875ca9385a2a0263d962263dd03a4c6155febfb2fed1bc9dcc8d8ca633766b12f2c224144a00ff3198a809af1ad6f85f4f1d104ce33f85a9303006a1d03ee4af33264478cff1bb051bf226a4e69b8c8986f03522f2a458abbd213a7bc5878e891a8dbf3279e05ab27646088dbfcd4c850794de614c44c2d78898a69b247c76d28c3bc2118ddfa39822a67d4434682d22f690e8d00e88e663223248823f9a1cf4ad7e5aa9b47a5021ababc432ed0592fdc1a4296af95758ab10da31feabfdd2d34f40cc138c0c1cab08e4b0520e64997646c44c7179abe83a434adf9b3485e04fb0568dfff7fe5ecf8ce08fb3331280aa4028740cf1162658f801bd74ee9cbe476f3ca70f9dd39fd41bcf9dd387e43fa2dd09f93a17989660a4253ac8f8f32053ab8644cd90608ff99312eb0fc0c482b07ead683e2d96684f8be6236281f6a458a4ef113543daa016d7f8df4c9da1753c0b53a9f574e944e43f9386af11b5c7c54aed23a2f688e81814b55a9b5879659ba8d58aa2db1d6aade7996caf98ecaa9cec084ef2cf4e1da0557a15a69f243a8bf96f2502fd46d13024963f261aae14cbdb44c33bf2f9930b3e4fd3531c57ea3d7b74a6f3434cdabb462ce01f1d7c97942cd1f80f618a06ded52a96f09b8e8826dedd2ac403821df7f3445f08fe45763e94a7e78075064afad2f0affa281f9b9a33cc174038bf3275442ce53fc00ca4f4ede24fce1f52fa06f2677da49f993bd25f49a4f74d514a7f3d13297d55f98b55486fbd10ae8ae7f7e9b9feb8c68f4f1ee02fb103e266bd5708931f9f6c158b34fe53362822da1ba26f50dcd42b222ff0bf9cec1596c65f62a658af1d13379f11c214cd1a1d1e932fd3c67895991afff3c99745b3c67fc24cc13ee16781befcfcb5f78ef5f0e4017ecec77ad8c1fa9b6aac871cace76662bdcfc1fa3661bddbc1faab2aac743b91279fbb62cf4fbedb8a2dd5f861c514976972f1be337940d4f0afb3566d903cfe62ff19d1a0f10714532c7b46347c56ac8857af24dd83f831e971152567d533426845ff84a2fb11fff2ecf3746fe25ff7113ffaae88cb35f6634635f6f0e411b18c9f95f443927e8939f5f6a48f926e65fcdb3e0a279a95f14a2c7467e3dff7b1dc3d772cd30af99f96fe3f29e977f61f112bf8012588056f6961f7065228e693192b3b648ed9f9423a973537b6c65be3e6ba62be94dd677676b55fd5135fff7bba25b45de862f0de20a06ffbae6d3b3f9cd8b2adffe66dbb12fd7d5b911cb6f22814f3e9ec10125bb7eeecf745fab662cb47b7efdab67d5725737be2e66d959cc4e6be4ac68edd9b6fdcb6a792b76bdb6d5586c6b20558a954de2e1490978d3d05ecb3273066654ab69c81d71a247b819c7e1e8c63c2bd3cdd684f40b6f38c58a3a3e9ec90ff27686472d921533e4ad9427a286ba7cc74b688523a5bece98235902e48baa3bb076e4fd14efbe3a574de2697a9b4ec67eadf9ccb65305218422e9b99f8a8ec00cbe652360ace5fbd47adbc355240bf1743bf35922b658be81fb20ad86265323b1c81198d52bb3bba7b6eb106ed9b4b9932bd353d56a66f290d94e9442a4549d9d7355eb48b77de59756dc4cc062bafaf2ad04be5b550f90d5455ed536ef394d73ae5354ef96d535eb7d45db468d27cda73d8b7d55d3c0c5b85619796094efa1e2a064859450bf81e972d5eec19be166efb1f9bd1f9a75437fd71b7df4f2db7fa8972979f566ef00b057afbf4f3b5f5d5381d7d61bf99cf70faf86a832d7c91aaeebdba60e35e7db9672f5a6ed76b0874ea5d56d9a4d718eccf5bd0d0df7fa75518e94f5a994c7fb298cb1716ce68a85b6404167e7155735dd3ccdeba66d95ab7c4ebac5bea35d62d0bf4d52d97395fe175d5adf49beaccaa9eba556e4b9dd74777db6aafa16e8ddf4f77b96ca7bbc2eba63b1d7de5f63f8d9187d8d5b14abcb1965832972d14ad6c3176f5a09529d82db17476b4542cc4aef675dcea8ab5c48a13a3c470f7568c8453f6b89d7295f77fac25e69c97a4ffb1b2bcf4912f514663fb5bca76dd4c548028e64be7c5405bdc07e09c8defe2bfacebb0e6a4ed4a0cba8b1d445b5e9379c43b2317725bfa7aee69f87b0f3450711f9c682fa47e8170dd0d35a79df01ee29c034e5f8fcee28bacf60b07e99e161f90ade69f8c976659e89b38afe995af0f4c6eb757a39db7c4ba7998efbaddf3c1c9adf305be34d91dbf18e189795d07f72a7169e27212366771f94b647ea3abbeb4cda9de0297edf92abb815c2e73b145673957b3795b97b921bdd015f05d3eb6f38c774626fe8f97d4f3c12d56d7fc7c5d80e654f297f492e7ff9e8ab5c4ac6c2e3b31922b15dec36a506a828b7151b5628fd1c73ea05ef913efd2409cf33560369081df9d9706e19cbfa7b321f47ff85e1a8073bde04ab50be0acf8517e69a05e82afe76cd104fea670696271f9375ef814a9529aeb99335b54f26f33efcb2f4dd2ca6462fb3f36fdbf000000ffff`

// DeployPublicResolver deploys a new VNT contract, binding an instance of PublicResolver to it.
func DeployPublicResolver(auth *bind.TransactOpts, backend bind.ContractBackend, ensAddr common.Address) (common.Address, *types.Transaction, *PublicResolver, error) {
	parsed, err := abi.JSON(strings.NewReader(PublicResolverABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(PublicResolverBin), backend, ensAddr)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PublicResolver{PublicResolverCaller: PublicResolverCaller{contract: contract}, PublicResolverTransactor: PublicResolverTransactor{contract: contract}, PublicResolverFilterer: PublicResolverFilterer{contract: contract}}, nil
}

// PublicResolver is an auto generated Go binding around an VNT contract.
type PublicResolver struct {
	PublicResolverCaller     // Read-only binding to the contract
	PublicResolverTransactor // Write-only binding to the contract
	PublicResolverFilterer   // Log filterer for contract events
}

// PublicResolverCaller is an auto generated read-only Go binding around an VNT contract.
type PublicResolverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicResolverTransactor is an auto generated write-only Go binding around an VNT contract.
type PublicResolverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicResolverFilterer is an auto generated log filtering Go binding around an VNT contract events.
type PublicResolverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PublicResolverSession is an auto generated Go binding around an VNT contract,
// with pre-set call and transact options.
type PublicResolverSession struct {
	Contract     *PublicResolver   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PublicResolverCallerSession is an auto generated read-only Go binding around an VNT contract,
// with pre-set call options.
type PublicResolverCallerSession struct {
	Contract *PublicResolverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// PublicResolverTransactorSession is an auto generated write-only Go binding around an VNT contract,
// with pre-set transact options.
type PublicResolverTransactorSession struct {
	Contract     *PublicResolverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// PublicResolverRaw is an auto generated low-level Go binding around an VNT contract.
type PublicResolverRaw struct {
	Contract *PublicResolver // Generic contract binding to access the raw methods on
}

// PublicResolverCallerRaw is an auto generated low-level read-only Go binding around an VNT contract.
type PublicResolverCallerRaw struct {
	Contract *PublicResolverCaller // Generic read-only contract binding to access the raw methods on
}

// PublicResolverTransactorRaw is an auto generated low-level write-only Go binding around an VNT contract.
type PublicResolverTransactorRaw struct {
	Contract *PublicResolverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPublicResolver creates a new instance of PublicResolver, bound to a specific deployed contract.
func NewPublicResolver(address common.Address, backend bind.ContractBackend) (*PublicResolver, error) {
	contract, err := bindPublicResolver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PublicResolver{PublicResolverCaller: PublicResolverCaller{contract: contract}, PublicResolverTransactor: PublicResolverTransactor{contract: contract}, PublicResolverFilterer: PublicResolverFilterer{contract: contract}}, nil
}

// NewPublicResolverCaller creates a new read-only instance of PublicResolver, bound to a specific deployed contract.
func NewPublicResolverCaller(address common.Address, caller bind.ContractCaller) (*PublicResolverCaller, error) {
	contract, err := bindPublicResolver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PublicResolverCaller{contract: contract}, nil
}

// NewPublicResolverTransactor creates a new write-only instance of PublicResolver, bound to a specific deployed contract.
func NewPublicResolverTransactor(address common.Address, transactor bind.ContractTransactor) (*PublicResolverTransactor, error) {
	contract, err := bindPublicResolver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PublicResolverTransactor{contract: contract}, nil
}

// NewPublicResolverFilterer creates a new log filterer instance of PublicResolver, bound to a specific deployed contract.
func NewPublicResolverFilterer(address common.Address, filterer bind.ContractFilterer) (*PublicResolverFilterer, error) {
	contract, err := bindPublicResolver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PublicResolverFilterer{contract: contract}, nil
}

// bindPublicResolver binds a generic wrapper to an already deployed contract.
func bindPublicResolver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PublicResolverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicResolver *PublicResolverRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PublicResolver.Contract.PublicResolverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicResolver *PublicResolverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicResolver.Contract.PublicResolverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicResolver *PublicResolverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicResolver.Contract.PublicResolverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PublicResolver *PublicResolverCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PublicResolver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PublicResolver *PublicResolverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PublicResolver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PublicResolver *PublicResolverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PublicResolver.Contract.contract.Transact(opts, method, params...)
}

// ABI is a free data retrieval call binding the contract method 0x2203ab56.
//
// Solidity: function ABI(node bytes32, contentTypes uint256) constant returns(contentType uint256, data bytes)
func (_PublicResolver *PublicResolverCaller) ABIRecord(opts *bind.CallOpts, node string, contentTypes *big.Int) (string, error) {
	var (
		ret = new(string)
	)
	out := ret
	err := _PublicResolver.contract.Call(opts, out, "ABIRecord", node, contentTypes)
	return *ret, err
}

func (_PublicResolver *PublicResolverCaller) ABIContentType(opts *bind.CallOpts, node string, contentTypes *big.Int) (*big.Int, error) {
	var (
		ret = new(*big.Int)
	)
	out := ret
	err := _PublicResolver.contract.Call(opts, out, "ABIContentType", node, contentTypes)
	return *ret, err
}

// ABI is a free data retrieval call binding the contract method 0x2203ab56.
//
// Solidity: function ABI(node bytes32, contentTypes uint256) constant returns(contentType uint256, data bytes)
func (_PublicResolver *PublicResolverSession) ABIRecord(node string, contentTypes *big.Int) (string, error) {
	return _PublicResolver.Contract.ABIRecord(&_PublicResolver.CallOpts, node, contentTypes)
}

func (_PublicResolver *PublicResolverSession) ABIContentType(node string, contentTypes *big.Int) (*big.Int, error) {
	return _PublicResolver.Contract.ABIContentType(&_PublicResolver.CallOpts, node, contentTypes)
}

// ABI is a free data retrieval call binding the contract method 0x2203ab56.
//
// Solidity: function ABI(node bytes32, contentTypes uint256) constant returns(contentType uint256, data bytes)
func (_PublicResolver *PublicResolverCallerSession) ABIRecord(node string, contentTypes *big.Int) (string, error) {
	return _PublicResolver.Contract.ABIRecord(&_PublicResolver.CallOpts, node, contentTypes)
}

func (_PublicResolver *PublicResolverCallerSession) ABIContentType(node string, contentTypes *big.Int) (*big.Int, error) {
	return _PublicResolver.Contract.ABIContentType(&_PublicResolver.CallOpts, node, contentTypes)
}

// Addr is a free data retrieval call binding the contract method 0x3b3b57de.
//
// Solidity: function addr(node bytes32) constant returns(ret address)
func (_PublicResolver *PublicResolverCaller) Addr(opts *bind.CallOpts, node string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _PublicResolver.contract.Call(opts, out, "addr", node)
	return *ret0, err
}

// Addr is a free data retrieval call binding the contract method 0x3b3b57de.
//
// Solidity: function addr(node bytes32) constant returns(ret address)
func (_PublicResolver *PublicResolverSession) Addr(node string) (common.Address, error) {
	return _PublicResolver.Contract.Addr(&_PublicResolver.CallOpts, node)
}

// Addr is a free data retrieval call binding the contract method 0x3b3b57de.
//
// Solidity: function addr(node bytes32) constant returns(ret address)
func (_PublicResolver *PublicResolverCallerSession) Addr(node string) (common.Address, error) {
	return _PublicResolver.Contract.Addr(&_PublicResolver.CallOpts, node)
}

// Content is a free data retrieval call binding the contract method 0x2dff6941.
//
// Solidity: function content(node bytes32) constant returns(ret bytes32)
func (_PublicResolver *PublicResolverCaller) Content(opts *bind.CallOpts, node string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _PublicResolver.contract.Call(opts, out, "content", node)
	return *ret0, err
}

// Content is a free data retrieval call binding the contract method 0x2dff6941.
//
// Solidity: function content(node bytes32) constant returns(ret bytes32)
func (_PublicResolver *PublicResolverSession) Content(node string) (string, error) {
	return _PublicResolver.Contract.Content(&_PublicResolver.CallOpts, node)
}

// Content is a free data retrieval call binding the contract method 0x2dff6941.
//
// Solidity: function content(node bytes32) constant returns(ret bytes32)
func (_PublicResolver *PublicResolverCallerSession) Content(node string) (string, error) {
	return _PublicResolver.Contract.Content(&_PublicResolver.CallOpts, node)
}

// Name is a free data retrieval call binding the contract method 0x691f3431.
//
// Solidity: function name(node bytes32) constant returns(ret string)
func (_PublicResolver *PublicResolverCaller) Name(opts *bind.CallOpts, node string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _PublicResolver.contract.Call(opts, out, "name", node)
	return *ret0, err
}

// Name is a free data retrieval call binding the contract method 0x691f3431.
//
// Solidity: function name(node bytes32) constant returns(ret string)
func (_PublicResolver *PublicResolverSession) Name(node string) (string, error) {
	return _PublicResolver.Contract.Name(&_PublicResolver.CallOpts, node)
}

// Name is a free data retrieval call binding the contract method 0x691f3431.
//
// Solidity: function name(node bytes32) constant returns(ret string)
func (_PublicResolver *PublicResolverCallerSession) Name(node string) (string, error) {
	return _PublicResolver.Contract.Name(&_PublicResolver.CallOpts, node)
}

// Pubkey is a free data retrieval call binding the contract method 0xc8690233.
//
// Solidity: function pubkey(node bytes32) constant returns(x bytes32, y bytes32)
func (_PublicResolver *PublicResolverCaller) PubkeyX(opts *bind.CallOpts, node string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _PublicResolver.contract.Call(opts, &out, "pubkeyX", node)
	return *ret0, err
}

func (_PublicResolver *PublicResolverCaller) PubkeyY(opts *bind.CallOpts, node string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _PublicResolver.contract.Call(opts, &out, "pubkeyY", node)
	return *ret0, err
}

// Pubkey is a free data retrieval call binding the contract method 0xc8690233.
//
// Solidity: function pubkey(node bytes32) constant returns(x bytes32, y bytes32)
func (_PublicResolver *PublicResolverSession) PubkeyX(node string) (string, error) {
	return _PublicResolver.Contract.PubkeyX(&_PublicResolver.CallOpts, node)
}

func (_PublicResolver *PublicResolverSession) PubkeyY(node string) (string, error) {
	return _PublicResolver.Contract.PubkeyY(&_PublicResolver.CallOpts, node)
}

// Pubkey is a free data retrieval call binding the contract method 0xc8690233.
//
// Solidity: function pubkey(node bytes32) constant returns(x bytes32, y bytes32)
func (_PublicResolver *PublicResolverCallerSession) PubkeyX(node string) (string, error) {
	return _PublicResolver.Contract.PubkeyX(&_PublicResolver.CallOpts, node)
}

func (_PublicResolver *PublicResolverCallerSession) Pubkey(node string) (string, error) {
	return _PublicResolver.Contract.PubkeyY(&_PublicResolver.CallOpts, node)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceID bytes4) constant returns(bool)
func (_PublicResolver *PublicResolverCaller) SupportsInterface(opts *bind.CallOpts, interfaceID [4]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _PublicResolver.contract.Call(opts, out, "supportsInterface", interfaceID)
	return *ret0, err
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceID bytes4) constant returns(bool)
func (_PublicResolver *PublicResolverSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _PublicResolver.Contract.SupportsInterface(&_PublicResolver.CallOpts, interfaceID)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(interfaceID bytes4) constant returns(bool)
func (_PublicResolver *PublicResolverCallerSession) SupportsInterface(interfaceID [4]byte) (bool, error) {
	return _PublicResolver.Contract.SupportsInterface(&_PublicResolver.CallOpts, interfaceID)
}

// Text is a free data retrieval call binding the contract method 0x59d1d43c.
//
// Solidity: function text(node bytes32, key string) constant returns(ret string)
func (_PublicResolver *PublicResolverCaller) Text(opts *bind.CallOpts, node string, key string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _PublicResolver.contract.Call(opts, out, "text", node, key)
	return *ret0, err
}

// Text is a free data retrieval call binding the contract method 0x59d1d43c.
//
// Solidity: function text(node bytes32, key string) constant returns(ret string)
func (_PublicResolver *PublicResolverSession) Text(node string, key string) (string, error) {
	return _PublicResolver.Contract.Text(&_PublicResolver.CallOpts, node, key)
}

// Text is a free data retrieval call binding the contract method 0x59d1d43c.
//
// Solidity: function text(node bytes32, key string) constant returns(ret string)
func (_PublicResolver *PublicResolverCallerSession) Text(node string, key string) (string, error) {
	return _PublicResolver.Contract.Text(&_PublicResolver.CallOpts, node, key)
}

// SetABI is a paid mutator transaction binding the contract method 0x623195b0.
//
// Solidity: function setABI(node bytes32, contentType uint256, data bytes) returns()
func (_PublicResolver *PublicResolverTransactor) SetABI(opts *bind.TransactOpts, node string, contentType *big.Int, data []byte) (*types.Transaction, error) {
	return _PublicResolver.contract.Transact(opts, "setABI", node, contentType, data)
}

// SetABI is a paid mutator transaction binding the contract method 0x623195b0.
//
// Solidity: function setABI(node bytes32, contentType uint256, data bytes) returns()
func (_PublicResolver *PublicResolverSession) SetABI(node string, contentType *big.Int, data []byte) (*types.Transaction, error) {
	return _PublicResolver.Contract.SetABI(&_PublicResolver.TransactOpts, node, contentType, data)
}

// SetABI is a paid mutator transaction binding the contract method 0x623195b0.
//
// Solidity: function setABI(node bytes32, contentType uint256, data bytes) returns()
func (_PublicResolver *PublicResolverTransactorSession) SetABI(node string, contentType *big.Int, data []byte) (*types.Transaction, error) {
	return _PublicResolver.Contract.SetABI(&_PublicResolver.TransactOpts, node, contentType, data)
}

// SetAddr is a paid mutator transaction binding the contract method 0xd5fa2b00.
//
// Solidity: function setAddr(node bytes32, addr address) returns()
func (_PublicResolver *PublicResolverTransactor) SetAddr(opts *bind.TransactOpts, node string, addr common.Address) (*types.Transaction, error) {
	return _PublicResolver.contract.Transact(opts, "setAddr", node, addr)
}

// SetAddr is a paid mutator transaction binding the contract method 0xd5fa2b00.
//
// Solidity: function setAddr(node bytes32, addr address) returns()
func (_PublicResolver *PublicResolverSession) SetAddr(node string, addr common.Address) (*types.Transaction, error) {
	return _PublicResolver.Contract.SetAddr(&_PublicResolver.TransactOpts, node, addr)
}

// SetAddr is a paid mutator transaction binding the contract method 0xd5fa2b00.
//
// Solidity: function setAddr(node bytes32, addr address) returns()
func (_PublicResolver *PublicResolverTransactorSession) SetAddr(node string, addr common.Address) (*types.Transaction, error) {
	return _PublicResolver.Contract.SetAddr(&_PublicResolver.TransactOpts, node, addr)
}

// SetContent is a paid mutator transaction binding the contract method 0xc3d014d6.
//
// Solidity: function setContent(node bytes32, hash bytes32) returns()
func (_PublicResolver *PublicResolverTransactor) SetContent(opts *bind.TransactOpts, node string, hash string) (*types.Transaction, error) {
	return _PublicResolver.contract.Transact(opts, "setContent", node, hash)
}

// SetContent is a paid mutator transaction binding the contract method 0xc3d014d6.
//
// Solidity: function setContent(node bytes32, hash bytes32) returns()
func (_PublicResolver *PublicResolverSession) SetContent(node string, hash string) (*types.Transaction, error) {
	return _PublicResolver.Contract.SetContent(&_PublicResolver.TransactOpts, node, hash)
}

// SetContent is a paid mutator transaction binding the contract method 0xc3d014d6.
//
// Solidity: function setContent(node bytes32, hash bytes32) returns()
func (_PublicResolver *PublicResolverTransactorSession) SetContent(node string, hash string) (*types.Transaction, error) {
	return _PublicResolver.Contract.SetContent(&_PublicResolver.TransactOpts, node, hash)
}

// SetName is a paid mutator transaction binding the contract method 0x77372213.
//
// Solidity: function setName(node bytes32, name string) returns()
func (_PublicResolver *PublicResolverTransactor) SetName(opts *bind.TransactOpts, node string, name string) (*types.Transaction, error) {
	return _PublicResolver.contract.Transact(opts, "setName", node, name)
}

// SetName is a paid mutator transaction binding the contract method 0x77372213.
//
// Solidity: function setName(node bytes32, name string) returns()
func (_PublicResolver *PublicResolverSession) SetName(node string, name string) (*types.Transaction, error) {
	return _PublicResolver.Contract.SetName(&_PublicResolver.TransactOpts, node, name)
}

// SetName is a paid mutator transaction binding the contract method 0x77372213.
//
// Solidity: function setName(node bytes32, name string) returns()
func (_PublicResolver *PublicResolverTransactorSession) SetName(node string, name string) (*types.Transaction, error) {
	return _PublicResolver.Contract.SetName(&_PublicResolver.TransactOpts, node, name)
}

// SetPubkey is a paid mutator transaction binding the contract method 0x29cd62ea.
//
// Solidity: function setPubkey(node bytes32, x bytes32, y bytes32) returns()
func (_PublicResolver *PublicResolverTransactor) SetPubkey(opts *bind.TransactOpts, node string, x string, y string) (*types.Transaction, error) {
	return _PublicResolver.contract.Transact(opts, "setPubkey", node, x, y)
}

// SetPubkey is a paid mutator transaction binding the contract method 0x29cd62ea.
//
// Solidity: function setPubkey(node bytes32, x bytes32, y bytes32) returns()
func (_PublicResolver *PublicResolverSession) SetPubkey(node string, x string, y string) (*types.Transaction, error) {
	return _PublicResolver.Contract.SetPubkey(&_PublicResolver.TransactOpts, node, x, y)
}

// SetPubkey is a paid mutator transaction binding the contract method 0x29cd62ea.
//
// Solidity: function setPubkey(node bytes32, x bytes32, y bytes32) returns()
func (_PublicResolver *PublicResolverTransactorSession) SetPubkey(node string, x string, y string) (*types.Transaction, error) {
	return _PublicResolver.Contract.SetPubkey(&_PublicResolver.TransactOpts, node, x, y)
}

// SetText is a paid mutator transaction binding the contract method 0x10f13a8c.
//
// Solidity: function setText(node bytes32, key string, value string) returns()
func (_PublicResolver *PublicResolverTransactor) SetText(opts *bind.TransactOpts, node string, key string, value string) (*types.Transaction, error) {
	return _PublicResolver.contract.Transact(opts, "setText", node, key, value)
}

// SetText is a paid mutator transaction binding the contract method 0x10f13a8c.
//
// Solidity: function setText(node bytes32, key string, value string) returns()
func (_PublicResolver *PublicResolverSession) SetText(node string, key string, value string) (*types.Transaction, error) {
	return _PublicResolver.Contract.SetText(&_PublicResolver.TransactOpts, node, key, value)
}

// SetText is a paid mutator transaction binding the contract method 0x10f13a8c.
//
// Solidity: function setText(node bytes32, key string, value string) returns()
func (_PublicResolver *PublicResolverTransactorSession) SetText(node string, key string, value string) (*types.Transaction, error) {
	return _PublicResolver.Contract.SetText(&_PublicResolver.TransactOpts, node, key, value)
}

// PublicResolverABIChangedIterator is returned from FilterABIChanged and is used to iterate over the raw logs and unpacked data for ABIChanged events raised by the PublicResolver contract.
type PublicResolverABIChangedIterator struct {
	Event *PublicResolverABIChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  hubble.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PublicResolverABIChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicResolverABIChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PublicResolverABIChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PublicResolverABIChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicResolverABIChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicResolverABIChanged represents a ABIChanged event raised by the PublicResolver contract.
type PublicResolverABIChanged struct {
	Node        string
	ContentType *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterABIChanged is a free log retrieval operation binding the contract event 0xaa121bbeef5f32f5961a2a28966e769023910fc9479059ee3495d4c1a696efe3.
//
// Solidity: e ABIChanged(node indexed bytes32, contentType indexed uint256)
func (_PublicResolver *PublicResolverFilterer) FilterABIChanged(opts *bind.FilterOpts, node []string, contentType []*big.Int) (*PublicResolverABIChangedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}
	var contentTypeRule []interface{}
	for _, contentTypeItem := range contentType {
		contentTypeRule = append(contentTypeRule, contentTypeItem)
	}

	logs, sub, err := _PublicResolver.contract.FilterLogs(opts, "ABIChanged", nodeRule, contentTypeRule)
	if err != nil {
		return nil, err
	}
	return &PublicResolverABIChangedIterator{contract: _PublicResolver.contract, event: "ABIChanged", logs: logs, sub: sub}, nil
}

// WatchABIChanged is a free log subscription operation binding the contract event 0xaa121bbeef5f32f5961a2a28966e769023910fc9479059ee3495d4c1a696efe3.
//
// Solidity: e ABIChanged(node indexed bytes32, contentType indexed uint256)
func (_PublicResolver *PublicResolverFilterer) WatchABIChanged(opts *bind.WatchOpts, sink chan<- *PublicResolverABIChanged, node []string, contentType []*big.Int) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}
	var contentTypeRule []interface{}
	for _, contentTypeItem := range contentType {
		contentTypeRule = append(contentTypeRule, contentTypeItem)
	}

	logs, sub, err := _PublicResolver.contract.WatchLogs(opts, "ABIChanged", nodeRule, contentTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicResolverABIChanged)
				if err := _PublicResolver.contract.UnpackLog(event, "ABIChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PublicResolverAddrChangedIterator is returned from FilterAddrChanged and is used to iterate over the raw logs and unpacked data for AddrChanged events raised by the PublicResolver contract.
type PublicResolverAddrChangedIterator struct {
	Event *PublicResolverAddrChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  hubble.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PublicResolverAddrChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicResolverAddrChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PublicResolverAddrChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PublicResolverAddrChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicResolverAddrChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicResolverAddrChanged represents a AddrChanged event raised by the PublicResolver contract.
type PublicResolverAddrChanged struct {
	Node string
	A    common.Address
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterAddrChanged is a free log retrieval operation binding the contract event 0x52d7d861f09ab3d26239d492e8968629f95e9e318cf0b73bfddc441522a15fd2.
//
// Solidity: e AddrChanged(node indexed bytes32, a address)
func (_PublicResolver *PublicResolverFilterer) FilterAddrChanged(opts *bind.FilterOpts, node []string) (*PublicResolverAddrChangedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _PublicResolver.contract.FilterLogs(opts, "AddrChanged", nodeRule)
	if err != nil {
		return nil, err
	}
	return &PublicResolverAddrChangedIterator{contract: _PublicResolver.contract, event: "AddrChanged", logs: logs, sub: sub}, nil
}

// WatchAddrChanged is a free log subscription operation binding the contract event 0x52d7d861f09ab3d26239d492e8968629f95e9e318cf0b73bfddc441522a15fd2.
//
// Solidity: e AddrChanged(node indexed bytes32, a address)
func (_PublicResolver *PublicResolverFilterer) WatchAddrChanged(opts *bind.WatchOpts, sink chan<- *PublicResolverAddrChanged, node []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _PublicResolver.contract.WatchLogs(opts, "AddrChanged", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicResolverAddrChanged)
				if err := _PublicResolver.contract.UnpackLog(event, "AddrChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PublicResolverContentChangedIterator is returned from FilterContentChanged and is used to iterate over the raw logs and unpacked data for ContentChanged events raised by the PublicResolver contract.
type PublicResolverContentChangedIterator struct {
	Event *PublicResolverContentChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  hubble.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PublicResolverContentChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicResolverContentChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PublicResolverContentChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PublicResolverContentChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicResolverContentChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicResolverContentChanged represents a ContentChanged event raised by the PublicResolver contract.
type PublicResolverContentChanged struct {
	Node string
	Hash string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterContentChanged is a free log retrieval operation binding the contract event 0x0424b6fe0d9c3bdbece0e7879dc241bb0c22e900be8b6c168b4ee08bd9bf83bc.
//
// Solidity: e ContentChanged(node indexed bytes32, hash bytes32)
func (_PublicResolver *PublicResolverFilterer) FilterContentChanged(opts *bind.FilterOpts, node []string) (*PublicResolverContentChangedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _PublicResolver.contract.FilterLogs(opts, "ContentChanged", nodeRule)
	if err != nil {
		return nil, err
	}
	return &PublicResolverContentChangedIterator{contract: _PublicResolver.contract, event: "ContentChanged", logs: logs, sub: sub}, nil
}

// WatchContentChanged is a free log subscription operation binding the contract event 0x0424b6fe0d9c3bdbece0e7879dc241bb0c22e900be8b6c168b4ee08bd9bf83bc.
//
// Solidity: e ContentChanged(node indexed bytes32, hash bytes32)
func (_PublicResolver *PublicResolverFilterer) WatchContentChanged(opts *bind.WatchOpts, sink chan<- *PublicResolverContentChanged, node []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _PublicResolver.contract.WatchLogs(opts, "ContentChanged", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicResolverContentChanged)
				if err := _PublicResolver.contract.UnpackLog(event, "ContentChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PublicResolverNameChangedIterator is returned from FilterNameChanged and is used to iterate over the raw logs and unpacked data for NameChanged events raised by the PublicResolver contract.
type PublicResolverNameChangedIterator struct {
	Event *PublicResolverNameChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  hubble.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PublicResolverNameChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicResolverNameChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PublicResolverNameChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PublicResolverNameChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicResolverNameChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicResolverNameChanged represents a NameChanged event raised by the PublicResolver contract.
type PublicResolverNameChanged struct {
	Node string
	Name string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNameChanged is a free log retrieval operation binding the contract event 0xb7d29e911041e8d9b843369e890bcb72c9388692ba48b65ac54e7214c4c348f7.
//
// Solidity: e NameChanged(node indexed bytes32, name string)
func (_PublicResolver *PublicResolverFilterer) FilterNameChanged(opts *bind.FilterOpts, node []string) (*PublicResolverNameChangedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _PublicResolver.contract.FilterLogs(opts, "NameChanged", nodeRule)
	if err != nil {
		return nil, err
	}
	return &PublicResolverNameChangedIterator{contract: _PublicResolver.contract, event: "NameChanged", logs: logs, sub: sub}, nil
}

// WatchNameChanged is a free log subscription operation binding the contract event 0xb7d29e911041e8d9b843369e890bcb72c9388692ba48b65ac54e7214c4c348f7.
//
// Solidity: e NameChanged(node indexed bytes32, name string)
func (_PublicResolver *PublicResolverFilterer) WatchNameChanged(opts *bind.WatchOpts, sink chan<- *PublicResolverNameChanged, node []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _PublicResolver.contract.WatchLogs(opts, "NameChanged", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicResolverNameChanged)
				if err := _PublicResolver.contract.UnpackLog(event, "NameChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PublicResolverPubkeyChangedIterator is returned from FilterPubkeyChanged and is used to iterate over the raw logs and unpacked data for PubkeyChanged events raised by the PublicResolver contract.
type PublicResolverPubkeyChangedIterator struct {
	Event *PublicResolverPubkeyChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  hubble.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PublicResolverPubkeyChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicResolverPubkeyChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PublicResolverPubkeyChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PublicResolverPubkeyChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicResolverPubkeyChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicResolverPubkeyChanged represents a PubkeyChanged event raised by the PublicResolver contract.
type PublicResolverPubkeyChanged struct {
	Node string
	X    string
	Y    string
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterPubkeyChanged is a free log retrieval operation binding the contract event 0x1d6f5e03d3f63eb58751986629a5439baee5079ff04f345becb66e23eb154e46.
//
// Solidity: e PubkeyChanged(node indexed bytes32, x bytes32, y bytes32)
func (_PublicResolver *PublicResolverFilterer) FilterPubkeyChanged(opts *bind.FilterOpts, node []string) (*PublicResolverPubkeyChangedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _PublicResolver.contract.FilterLogs(opts, "PubkeyChanged", nodeRule)
	if err != nil {
		return nil, err
	}
	return &PublicResolverPubkeyChangedIterator{contract: _PublicResolver.contract, event: "PubkeyChanged", logs: logs, sub: sub}, nil
}

// WatchPubkeyChanged is a free log subscription operation binding the contract event 0x1d6f5e03d3f63eb58751986629a5439baee5079ff04f345becb66e23eb154e46.
//
// Solidity: e PubkeyChanged(node indexed bytes32, x bytes32, y bytes32)
func (_PublicResolver *PublicResolverFilterer) WatchPubkeyChanged(opts *bind.WatchOpts, sink chan<- *PublicResolverPubkeyChanged, node []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _PublicResolver.contract.WatchLogs(opts, "PubkeyChanged", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicResolverPubkeyChanged)
				if err := _PublicResolver.contract.UnpackLog(event, "PubkeyChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PublicResolverTextChangedIterator is returned from FilterTextChanged and is used to iterate over the raw logs and unpacked data for TextChanged events raised by the PublicResolver contract.
type PublicResolverTextChangedIterator struct {
	Event *PublicResolverTextChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log      // Log channel receiving the found contract events
	sub  hubble.Subscription // Subscription for errors, completion and termination
	done bool                // Whether the subscription completed delivering logs
	fail error               // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PublicResolverTextChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PublicResolverTextChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PublicResolverTextChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PublicResolverTextChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PublicResolverTextChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PublicResolverTextChanged represents a TextChanged event raised by the PublicResolver contract.
type PublicResolverTextChanged struct {
	Node       string
	IndexedKey common.Hash
	Key        string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTextChanged is a free log retrieval operation binding the contract event 0xd8c9334b1a9c2f9da342a0a2b32629c1a229b6445dad78947f674b44444a7550.
//
// Solidity: e TextChanged(node indexed bytes32, indexedKey indexed string, key string)
func (_PublicResolver *PublicResolverFilterer) FilterTextChanged(opts *bind.FilterOpts, node []string, indexedKey []string) (*PublicResolverTextChangedIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}
	var indexedKeyRule []interface{}
	for _, indexedKeyItem := range indexedKey {
		indexedKeyRule = append(indexedKeyRule, indexedKeyItem)
	}

	logs, sub, err := _PublicResolver.contract.FilterLogs(opts, "TextChanged", nodeRule, indexedKeyRule)
	if err != nil {
		return nil, err
	}
	return &PublicResolverTextChangedIterator{contract: _PublicResolver.contract, event: "TextChanged", logs: logs, sub: sub}, nil
}

// WatchTextChanged is a free log subscription operation binding the contract event 0xd8c9334b1a9c2f9da342a0a2b32629c1a229b6445dad78947f674b44444a7550.
//
// Solidity: e TextChanged(node indexed bytes32, indexedKey indexed string, key string)
func (_PublicResolver *PublicResolverFilterer) WatchTextChanged(opts *bind.WatchOpts, sink chan<- *PublicResolverTextChanged, node []string, indexedKey []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}
	var indexedKeyRule []interface{}
	for _, indexedKeyItem := range indexedKey {
		indexedKeyRule = append(indexedKeyRule, indexedKeyItem)
	}

	logs, sub, err := _PublicResolver.contract.WatchLogs(opts, "TextChanged", nodeRule, indexedKeyRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PublicResolverTextChanged)
				if err := _PublicResolver.contract.UnpackLog(event, "TextChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
