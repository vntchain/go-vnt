package contract

import (
	"strings"

	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/accounts/abi/bind"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/core/types"
)

// FIFSRegistrarABI is the input ABI used to generate the binding from.
const FIFSRegistrarABI = `[{"name":"FIFSRegistrar","constant":false,"inputs":[{"name":"vnsAddr","type":"address","indexed":false},{"name":"node","type":"string","indexed":false}],"outputs":[],"type":"constructor"},{"name":"registernode","constant":false,"inputs":[{"name":"subnode","type":"string","indexed":false},{"name":"owner","type":"address","indexed":false}],"outputs":[],"type":"function"},{"name":"owner","constant":false,"inputs":[{"name":"node","type":"string","indexed":false}],"outputs":[{"name":"output","type":"address","indexed":false}],"type":"call"},{"name":"setSubnodeOwner","constant":false,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"label","type":"string","indexed":false},{"name":"owner","type":"address","indexed":false}],"outputs":[],"type":"call"}]`

// FIFSRegistrarBin is the compiled bytecode used for deploying new contracts.
const FIFSRegistrarBin = `0x0161736db9091c0100789cbc575b6c1cd5f9ff9d73e68c7777bceb7516e2d849f064317fc2bf6677e338769a92840d10705b2e4d9a02ade87abc3bbb1eb29e7166669d38155ed7a434b44d45a5be504440a24808b534aa78684c1ffad6f2d056952ae84b5524e80d7a419590682348f5cdcc7a77636823dab20fe77ce7bbfcbedb39a3fdfe9ee93d9fe630bc3906808d27a65913d37c6909d36c09d372a9b9d46c621a98e6cd266b4e335ac09ad34ab34902de6c82bfc87a85692fa887cc05d3f5013aa4ef762ddfbcdbf267ef722cdb375d3062f71d328dca3aae5aac546e353c703a248a95ca27ccc529bbea40102333655bbe65d4ad93e6670cd73266eaa60725b0bbc9b1cb860f4907e5f06dc59d508994ce71db7443b6bce558c3a88774fc56d33f6cda15d3454f108d67fa871b33b65331ef0c2c62422a8a128f2b92cd33c6a460e0ea95a2c98a3f3aad684db4b61fc4b49ebb5575ce9c73dc450ead549a358df9d28ce19982254aa58ae11b25d3ae089eae98e5bae19a95830dbbec5b8e8d64f2e0d4c1c387cc9ae5f9aee122d5eb06b4e9521ce84b3c25140eed2a1cc001a4bfb8bcbc1c522b21957e9076ed0ec6968ae74e2f2fe3bb595eecd179315e44fa14c98acfb7d86a273bfda540081d133c302c4267133c50d6bec7056bb225d6bc7a39d02a1ccdf291908c9c88628f2ed63b1145b58b8ded21a28ef44381d72fd3ba4dd179317d7f3176ff8127cfa893023a2f2e93da041fd139199d0b4fbace438ded3ad7f975429f14697e6360ac2be9d38495558a8f9072fa613a25a12be9af2c779c8b4f07d220784de71fc827428528139da5bf4a683a2f16ee0faba269b7701497635aa131b66b62c7ae9d933b768f4d4e4e1460542aaee9791dac1305d8a659d1c34b89e2b9981643ebf73355cb55cc9946ad64d955679f009480afd07b442fae8ef49e8ef66700f0d7004c12a148515000710359dc0ae0cff48a77ca4f0250ff400be36f01a890b294629c94ad96f2bba45ce8214ffd001e03107b8a96e7c8535c02b80ac0130c606f26246d25a625837d1580b63520090fea000338eb2530c23e47a2b77fca01ed9d3679618d4c5e24f21f44a6c8f12e007d14cf85d03b02a42a801f92e2c536d2c5365240f605489c119d26a80100298262ac239d40f4e9961785456efa59e4fc2784120b5036bcd3415f68d399c7a226b04497f98d007e496ac90ef3648779488781f6b5a1361154ba057505418d51e13970e5362aa7c6368e00b803c02b74eed318a78ddaaded0f3854fc01b228508bfbd1c3013c4b2de740acce814dc76959a1e50c0913d793e100d326839d009201d66040de4ce466f290fc1491438f9291f26d0e7c93e2ffd3609afd7ee81e06900efb23860c06504eec7565e818037e45f41bb1f08e596433a6bcca8133448eaa0b0ce8896df401e419f0177232120633421168940d1f6947b09d220872bc8b72fcbfcd94da3c91d75e9260d0e754f49ee68c323d29e0eb3c19b15cc3ae991eddf50c8061001f41fb370a2017f5605f079f5a3001603780db01dc09e04887fc6bd1db791cc0f7013c0f60b543fe5eb22906500d3fcf8032034c069c62c043ac6d47e40600f40eaf8eceb7457b2dda27a2bd18ed3ba3785d2a0f80ef00a04a3fd5818b77d70a62ccccb8e602eb67d7a43252a4d2039b53fdecc800c0c721525399fd837bb48f697c132076612a03281b31951129e202b22fe0a92310a97d9aa6013d39d6cf36a83b32404c82b6780ec4094df60f02891c0b0fd70eee1fd4354093081984d51b02b43428022039dea191daca766448e71eed5e0de88b649264e90e9904fa3bc0640b6c83ec30c890c1918150fd8a30548afaca1c5a1153881bbba3229881ad882c3560d37890eb6086695a589aa164503f3aecd680cd6b7e346ac0a33dad8f7eddb2cd87838f3e7d861963175249da58d84f963fe299ae973f6e5a9e619d9c35ec5afe66a7dc98336ddfcbd79cbce796f335cb9f6dcce4cace5c7ec1f6cbb38665e76bcef50bb69f2f3bb6ef1a65dfcb2fd8dedaa9839dc3ffda83d3f0e71b3eb060fb756b26374b6975fd012acdbb66d9999bb7ea66ae1c5c5a484ef98bdf80c97862dbac2a1e6fae48a68adf42e70c8c41727a14e24dc6a44cdcfbc48392d55a82c0f26dc6647f62af4cbe20af55ff2677de20ffffac8cabe2774b2bfb086d554a55bcc3742e3a8c2e32269389bd32abc674527e65492f84ca49559ce2baccaa1f97fdeaa8cc3e29c7d41539785626aba4f8eba56a1bf521de8d7a9a13ea36b9457d89745f5ed265af7896e5a454cf8ab71e083cbd0a623ec308e06c572267f85a225bd45599d5e5353939be4a46bf58d2a514e7584e157f7de086168c14cfb1d5200085b740124cb2c4de582d16162910d04742f49160b78cbf219958ccd564fc75c9c4c9dc5119afc9d48b927da3ad4f7d11034c0e27627be526f5119959915bd45372f0bc1cfeb11ccacb41d5975bd4176422f65999535f926339793027c7aa325e550baa3056a981db74dedb42a4cf9218a108442de6cb61917c4d32e1e56a7258f412e9e788aa76797b759db753efe34daae23a5d0e9e176e4ea4aa72581cbb470e89be5c602ad517e4b0985f955b459a3ae14ba63e1184866f2989e86d7abe5bae1b764d5f305dcf726c7d77ae902be8db7db7611fd5778eeff8e844e1ba0fe9fde4ffc593f9cf22c082eda13c6bb8ad7fae701dc7bf834611cf772dbb8643e6b186e59a283b76c50a4698d201c7a963ceabe1a8b978c23dee9b93f35538767d311ca3bc70a642b0941bae6bda7e2899375c63ce43a9e5ab64cc390ddb47c3b2fdb15d1328d50c0f75c7aee9c1d2b03dab669b15dd8a7426c6719351afdf15c2ac9bad8e8ced9a386c54cddb1bf535fa666b618d3edc9859a38b954af787085d73d8224ee02415a758a9b888febbbfcc14db9833d9cf591ad1b8cbd64dbafcd2215744f3add231dacaf79a6ad568a0ed0966d95838c6c6c30936d11e5eb54be7d6defe52e9b8e1cd95ca46bd5e2afb8eeb25d70d9ea9eeb9b3af6bec3c2fc4e7be90a5dcb27bb25d7ad9d16cd9b13ddfb0fdec9eaa51f7ccd1ac65cf377c2fbba76d12d5293b9af517e7891135384bca15f38459898c1f185db321bf6d83f0b2add7bf6f341b3e01f277df9a7a1093dba05cb31d989d395d5ee4d15dfdf781b49d046db98c54df27f46ad490ec7b405e46c01fa06c6d3701ebf2426f15daa8d73b23bde4eafd37636e3ba91b3366fd436c4998e47dcbff040000ffff`

// DeployFIFSRegistrar deploys a new VNT contract, binding an instance of FIFSRegistrar to it.
func DeployFIFSRegistrar(auth *bind.TransactOpts, backend bind.ContractBackend, ensAddr common.Address, node string) (common.Address, *types.Transaction, *FIFSRegistrar, error) {
	parsed, err := abi.JSON(strings.NewReader(FIFSRegistrarABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(FIFSRegistrarBin), backend, ensAddr, node)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &FIFSRegistrar{FIFSRegistrarCaller: FIFSRegistrarCaller{contract: contract}, FIFSRegistrarTransactor: FIFSRegistrarTransactor{contract: contract}, FIFSRegistrarFilterer: FIFSRegistrarFilterer{contract: contract}}, nil
}

// FIFSRegistrar is an auto generated Go binding around an VNT contract.
type FIFSRegistrar struct {
	FIFSRegistrarCaller     // Read-only binding to the contract
	FIFSRegistrarTransactor // Write-only binding to the contract
	FIFSRegistrarFilterer   // Log filterer for contract events
}

// FIFSRegistrarCaller is an auto generated read-only Go binding around an VNT contract.
type FIFSRegistrarCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FIFSRegistrarTransactor is an auto generated write-only Go binding around an VNT contract.
type FIFSRegistrarTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FIFSRegistrarFilterer is an auto generated log filtering Go binding around an VNT contract events.
type FIFSRegistrarFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FIFSRegistrarSession is an auto generated Go binding around an VNT contract,
// with pre-set call and transact options.
type FIFSRegistrarSession struct {
	Contract     *FIFSRegistrar    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FIFSRegistrarCallerSession is an auto generated read-only Go binding around an VNT contract,
// with pre-set call options.
type FIFSRegistrarCallerSession struct {
	Contract *FIFSRegistrarCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// FIFSRegistrarTransactorSession is an auto generated write-only Go binding around an VNT contract,
// with pre-set transact options.
type FIFSRegistrarTransactorSession struct {
	Contract     *FIFSRegistrarTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// FIFSRegistrarRaw is an auto generated low-level Go binding around an VNT contract.
type FIFSRegistrarRaw struct {
	Contract *FIFSRegistrar // Generic contract binding to access the raw methods on
}

// FIFSRegistrarCallerRaw is an auto generated low-level read-only Go binding around an VNT contract.
type FIFSRegistrarCallerRaw struct {
	Contract *FIFSRegistrarCaller // Generic read-only contract binding to access the raw methods on
}

// FIFSRegistrarTransactorRaw is an auto generated low-level write-only Go binding around an VNT contract.
type FIFSRegistrarTransactorRaw struct {
	Contract *FIFSRegistrarTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFIFSRegistrar creates a new instance of FIFSRegistrar, bound to a specific deployed contract.
func NewFIFSRegistrar(address common.Address, backend bind.ContractBackend) (*FIFSRegistrar, error) {
	contract, err := bindFIFSRegistrar(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FIFSRegistrar{FIFSRegistrarCaller: FIFSRegistrarCaller{contract: contract}, FIFSRegistrarTransactor: FIFSRegistrarTransactor{contract: contract}, FIFSRegistrarFilterer: FIFSRegistrarFilterer{contract: contract}}, nil
}

// NewFIFSRegistrarCaller creates a new read-only instance of FIFSRegistrar, bound to a specific deployed contract.
func NewFIFSRegistrarCaller(address common.Address, caller bind.ContractCaller) (*FIFSRegistrarCaller, error) {
	contract, err := bindFIFSRegistrar(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FIFSRegistrarCaller{contract: contract}, nil
}

// NewFIFSRegistrarTransactor creates a new write-only instance of FIFSRegistrar, bound to a specific deployed contract.
func NewFIFSRegistrarTransactor(address common.Address, transactor bind.ContractTransactor) (*FIFSRegistrarTransactor, error) {
	contract, err := bindFIFSRegistrar(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FIFSRegistrarTransactor{contract: contract}, nil
}

// NewFIFSRegistrarFilterer creates a new log filterer instance of FIFSRegistrar, bound to a specific deployed contract.
func NewFIFSRegistrarFilterer(address common.Address, filterer bind.ContractFilterer) (*FIFSRegistrarFilterer, error) {
	contract, err := bindFIFSRegistrar(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FIFSRegistrarFilterer{contract: contract}, nil
}

// bindFIFSRegistrar binds a generic wrapper to an already deployed contract.
func bindFIFSRegistrar(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(FIFSRegistrarABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FIFSRegistrar *FIFSRegistrarRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FIFSRegistrar.Contract.FIFSRegistrarCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FIFSRegistrar *FIFSRegistrarRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FIFSRegistrar.Contract.FIFSRegistrarTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FIFSRegistrar *FIFSRegistrarRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FIFSRegistrar.Contract.FIFSRegistrarTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FIFSRegistrar *FIFSRegistrarCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _FIFSRegistrar.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FIFSRegistrar *FIFSRegistrarTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FIFSRegistrar.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FIFSRegistrar *FIFSRegistrarTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FIFSRegistrar.Contract.contract.Transact(opts, method, params...)
}

// Register is a paid mutator transaction binding the contract method 0xd22057a9.
//
// function register(subnode bytes32, owner address) returns()
func (_FIFSRegistrar *FIFSRegistrarTransactor) Register(opts *bind.TransactOpts, subnode string, owner common.Address) (*types.Transaction, error) {
	return _FIFSRegistrar.contract.Transact(opts, "registernode", subnode, owner)
}

// Register is a paid mutator transaction binding the contract method 0xd22057a9.
//
// function register(subnode bytes32, owner address) returns()
func (_FIFSRegistrar *FIFSRegistrarSession) Register(subnode string, owner common.Address) (*types.Transaction, error) {
	return _FIFSRegistrar.Contract.Register(&_FIFSRegistrar.TransactOpts, subnode, owner)
}

// Register is a paid mutator transaction binding the contract method 0xd22057a9.
//
// function register(subnode bytes32, owner address) returns()
func (_FIFSRegistrar *FIFSRegistrarTransactorSession) Register(subnode string, owner common.Address) (*types.Transaction, error) {
	return _FIFSRegistrar.Contract.Register(&_FIFSRegistrar.TransactOpts, subnode, owner)
}
