// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"strings"

	hubble "github.com/vntchain/go-vnt"
	"github.com/vntchain/go-vnt/common"

	"github.com/vntchain/go-vnt/accounts/abi"
	"github.com/vntchain/go-vnt/accounts/abi/bind"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/event"
)

// ENSABI is the input ABI used to generate the binding from.
const ENSABI = "[{\"name\":\"ENS\",\"constant\":false,\"inputs\":[],\"outputs\":[],\"type\":\"constructor\"},{\"name\":\"owner\",\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"address\",\"indexed\":false}],\"type\":\"function\"},{\"name\":\"resolver\",\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"address\",\"indexed\":false}],\"type\":\"function\"},{\"name\":\"ttl\",\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"uint64\",\"indexed\":false}],\"type\":\"function\"},{\"name\":\"setOwner\",\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":false}],\"outputs\":[],\"type\":\"function\"},{\"name\":\"setSubnodeOwner\",\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"label\",\"type\":\"string\",\"indexed\":false},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":false}],\"outputs\":[],\"type\":\"function\"},{\"name\":\"setResolver\",\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"resolver\",\"type\":\"address\",\"indexed\":false}],\"outputs\":[],\"type\":\"function\"},{\"name\":\"setTTL\",\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"ttl\",\"type\":\"uint64\",\"indexed\":false}],\"outputs\":[],\"type\":\"function\"},{\"name\":\"NewResolver\",\"anonymous\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":true},{\"name\":\"resolver\",\"type\":\"address\",\"indexed\":false}],\"type\":\"event\"},{\"name\":\"NewTTL\",\"anonymous\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":true},{\"name\":\"ttl\",\"type\":\"uint64\",\"indexed\":false}],\"type\":\"event\"},{\"name\":\"NewOwner\",\"anonymous\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":true},{\"name\":\"label\",\"type\":\"string\",\"indexed\":true},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":false}],\"type\":\"event\"},{\"name\":\"Transfer\",\"anonymous\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":true},{\"name\":\"owner\",\"type\":\"address\",\"indexed\":false}],\"type\":\"event\"}]"

// ENSBin is the compiled bytecode used for deploying new contracts.
const ENSBin = `0x0161736db90a5f0100789ccc586d6c1cd5d57eeebd7bc7eb1def7a379bc4f19bbc30d980022f7ed7ce9713f226840de485b434409c405384d6e3ddf17acb7ac6cccc3a9812af9bf0e136911a54a955d5c604a955296d51a2f2a3246ad552a92a3ff8c10f12f507fd54914a5baa024a8b68a03a77663f6c878414dcd63f66ce3d1fcf7dce9d738fb5e7ede50f3d971e85e98d32006c6bc720ab61904f4e62904d62504ed6266b350c028360b5415eab054f0c3292442d142707796d12fc2c8b0bcb1ed7765be396eb03b448deed967debeeb23f728753b67dcb052375e76ecb2cced36ab962f116d303a7452c572c7edc9ad8690f3b10a448efb4cb7ed9ac941fb4ee32ddb23954b13c44c8d27e8be50f5876d17221692d77dc5f352bd0488eee714ddb1bb65cb4a92d6e72ec82e907b6c8c0adb975882ab75dd6fedbf7db968b765aeabbacfdbb2dcfa98cabc00e3d128944a2d1585b7b9b1e916c8c31260503d796881acb9d988ee835d45fdf89ea6d875987366a8d3aee04879ecf8f58e6587ec8f42cc162f97cd1f4cdbc6517054f16ad42c574ade2ff57ed825f766c24c48e5d03484a47714945dd3a8945c2f72b48473dcb0f782eeef42c7fa03a643b452bd02cd13dcb6fb05eaa7996bf67cf6de88abd28750efd0a6cc776243f3b353515480703297988defa09c6f8646e7a7a6a0adfcd20d7969b0a44966bcf21f930f918c8b5192ca7cd5aa369cf1d69441b971376b41e16991316991316690dd37b90fc0209c947140a72c709a79f2be239241f9d524b9592be2274cec168ba5c7371e391f7375e2b14697ddf7c23bfb1011c50483e46cf38724f2add94226fc060c9e979b806ab137618ab7d18eccf29f3e7e9b95290c6e0c9c3219c68eec7ebfb7db84c8e5c3093e0006fff6791e7e06d0ccf5cbf9e233715d50d5cea2f773caa77c0b6aca2115c28fc2caa678bd650b5942fdbc3cea60810519e116a80e8c0aa30f2f1f0fd1400fe0b001b49884941fe620b45183c265f0560290be439000f29392ab730e0b49297411b03f067eaaf7d6d5500d1d700684c7b0cc00606f06824c985fc3b8013142221632c94b5884cd7a1daa2011467005ba7fd144096c49e681703daa2ed94420ac05b006204176b23952e015c01608abc5fefd845af7b59fc2ef526f8f8b0122911440789124b50e42d7410643ad71701e26dbc21ca86d8d94ee2df484cd2c61b007452b26f07bb43210d03f806219d6f229d6f222931a590c01a505d144d50efb564b3882c7bea9b7016ee4269a8bd4fd22e9a0259acb60965d99497bc157e5c169d157e23801f50b8de12aeb7840772c033de845a4650893a947e75f8215fa0835ca19f01d003e0655aae561f691d803f7220f6160796be478f98a06405035600a074f9760289fb14b65d7d9a84119aa93cf89d4df39d81796b68a66478fe74c39c57e68e4f93dccf025d3fe93ad53ed71350e70469379388ae4f31601f809f93ce0d025c0a88ab0097bc9651027db4acce4d035dcf32c001f01b8a9c0e00a61580ce1ba2c29a569ba7487bc4aac3de4196c3f361afe201af3708ef58007b4c61a91339d6e0a5aec9f1f900931cb81dc0bbe47e320038a900d4999d3c5d07384a00cfce0350f7ef1895c85a2099083bc9a859a06602fc90c743956bda25cba3cb9806702580eb5ada12954316c05a0037b4e84d16dcb7af31608601df64c0b758d34ec5f92706747020c1812e0e74f3a6fd3a0edcc381610e8c70e07e0e782df64738f07d0efc98033f69d1ab7b05801ac5aa707d6bf82e85effef09d034015fe328033f56e04603d002adf2b1940754a1d900a721f03a89ebec7002a8b551ca0ef788003f4395e6ce181b38d033487865c6b9ca5d8d589b4148964d7f2448aeded02f87a88c4cef4b6eecdfaffe97c19203660671a88a499ae930e9071e5428b4d3aa02dc5ceb44804b6b64ee51dbd0a227183aeeb407b96a5d8226d4d1a8849d04bcf823441c8b66ea023cb82c5eaee6ddd864e5d0b8182b0120140dd83b8019deb5b3c92ffcdd6a4c9e793fa3e1d48853649b6452d3609a45bc0641d6cb16c095842017bbb02f7a5eb15e3aed914286659c34da7937dbef17faf52b6adb5eaffdeefe99332f64e224e2f167c60d6bbd7b35caf77bf55f6ccf28323a65deabdd92954472ddbf77a4b4eafe7167a4b657fa43a942d38a3bde3b65f1831cb766fc9f9df71dbef2d38b6ef9a05dfebb56cafb1eacd6261709daa3f56f58171dbaf9487b22394c28e5d03f931d72a38a363e58a952d409598a44a8578054cb6c7568e68e289da41c934f15b189c8131484ef52f5e674ccad8be270e4956aa1b54e439c6642ab655c64fcbd5da1b72dd16f93f33b25d13af4e1ebc81d04e49a9897798c1454bd079c6643cb65566b4a841cebf9a34fa02e7b8260e724366b48fc994d62333c7e55aeda0ec9e91f161727c6572b889fa309f8dfa2827d495728576867ccf4e1ab2433ccdb2526a33e2af07d44ebf06299f620430332b91c3bc91c80aed94cc18f2eaac5c7f8a825e9a34a414275856137f39b0a50e23c533ec942210e1909caebf58ce248b3d1eadbe19dd174d078f37d52eca853a83584d2e2b4b325592cba369d9ad0d4ba6559b4e8a4c3f71592957fd4e2eadcaf69236acf569e2a653f465b6194d5f6a2de279284402fcaf17a4d4fca69d5a8e78e9fdedd48ac42fdfdf4e2d4abc51b78b9b87e733929a380743a6c4f619b12d2b724765b7f811b2923dde84a10e27a619c16c9529f14ced0238cb3571841972a93623a57648a6c4b76ba7c41964c5d335823c3b07929aa538c6eaccbe722148a989271951fb726d46bc86acf852eda85c21fe30078afaad38d9807ae44250294d6c32c47964c5a1da51d925de6d62e0eb9158d8443cdf2d544cbb648c5bae57766c6353b62fdb675ce3bb55fb3e63ddfa35d7f7f75dbba057be77de2dff28770301bb450ff75913288c982e3cdf2ddb258c9b95aa85e0c784592cba9647bee14f6efa855e71ec92a11e55db2b976cab68946d1fd5b2edf7afc76e058b51736cac6c97d66c58b771cda6b51b37f663b7757fb5ec5a283876b1ac2601f9ed8e53c1a85782635726825ff8f45b1f5e30e998373ad8bb7643ff80396c7da25a51dbadddd0dfd0dd5c1e6fc803d5a1869c2b1629c562d1ab0c8d5747a975a23161983b5f689d2e84b385093c8007513187ac0abcc017f82a8fd8e6a8c5bec897211c06b17973203e770424c2e94fa465f0232f34f3d19ae39eb660d2136d0c79dac3f94e4c8d76f4c654a7a375a0134fe5f3fb4d6f345f302b957cc1775c2f316f12d3a9b79c4b52ecd835900a86328b1a3399b4f0fdcae2c64466c9dc81ccd2d6794c57388e792e1abde733193a9fcce6cc8e5d03999e4cc1b13ddfb4fdcce661b3e2593d99b23d56f5bdcce67beeedc904d51e2efc89310a53016e9598670ef434d07cbf320bcd77abad600d3f2298698005859d21c7a2f580550c591c98b57723365035a383babe6074e8311c1e682bd3fa915d34f9cbe4db447742e83020bca697c8f092845bbfecc2f056d7e8df9f68bd62172649b789fed1d2deb3e7b685611c5cab4bd7fb07e5eaccaba07fc545bdd8295f84ac7ba162f80fe65b6ff8999e8c693bf6c4a853f53e442d50a6977343e7f85fc6050d5dac71cbf65b13aaff675b98843e02822dff581786e3e5758c8bd00c3ac40230fca01d6236b77ba7fe010000ffff`

// DeployENS deploys a new VNT contract, binding an instance of ENS to it.
func DeployENS(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ENS, error) {
	parsed, err := abi.JSON(strings.NewReader(ENSABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ENSBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ENS{ENSCaller: ENSCaller{contract: contract}, ENSTransactor: ENSTransactor{contract: contract}, ENSFilterer: ENSFilterer{contract: contract}}, nil
}

// ENS is an auto generated Go binding around an VNT contract.
type ENS struct {
	ENSCaller     // Read-only binding to the contract
	ENSTransactor // Write-only binding to the contract
	ENSFilterer   // Log filterer for contract events
}

// ENSCaller is an auto generated read-only Go binding around an VNT contract.
type ENSCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSTransactor is an auto generated write-only Go binding around an VNT contract.
type ENSTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSFilterer is an auto generated log filtering Go binding around an VNT contract events.
type ENSFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ENSSession is an auto generated Go binding around an VNT contract,
// with pre-set call and transact options.
type ENSSession struct {
	Contract     *ENS              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ENSCallerSession is an auto generated read-only Go binding around an VNT contract,
// with pre-set call options.
type ENSCallerSession struct {
	Contract *ENSCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ENSTransactorSession is an auto generated write-only Go binding around an VNT contract,
// with pre-set transact options.
type ENSTransactorSession struct {
	Contract     *ENSTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ENSRaw is an auto generated low-level Go binding around an VNT contract.
type ENSRaw struct {
	Contract *ENS // Generic contract binding to access the raw methods on
}

// ENSCallerRaw is an auto generated low-level read-only Go binding around an VNT contract.
type ENSCallerRaw struct {
	Contract *ENSCaller // Generic read-only contract binding to access the raw methods on
}

// ENSTransactorRaw is an auto generated low-level write-only Go binding around an VNT contract.
type ENSTransactorRaw struct {
	Contract *ENSTransactor // Generic write-only contract binding to access the raw methods on
}

// NewENS creates a new instance of ENS, bound to a specific deployed contract.
func NewENS(address common.Address, backend bind.ContractBackend) (*ENS, error) {
	contract, err := bindENS(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ENS{ENSCaller: ENSCaller{contract: contract}, ENSTransactor: ENSTransactor{contract: contract}, ENSFilterer: ENSFilterer{contract: contract}}, nil
}

// NewENSCaller creates a new read-only instance of ENS, bound to a specific deployed contract.
func NewENSCaller(address common.Address, caller bind.ContractCaller) (*ENSCaller, error) {
	contract, err := bindENS(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ENSCaller{contract: contract}, nil
}

// NewENSTransactor creates a new write-only instance of ENS, bound to a specific deployed contract.
func NewENSTransactor(address common.Address, transactor bind.ContractTransactor) (*ENSTransactor, error) {
	contract, err := bindENS(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ENSTransactor{contract: contract}, nil
}

// NewENSFilterer creates a new log filterer instance of ENS, bound to a specific deployed contract.
func NewENSFilterer(address common.Address, filterer bind.ContractFilterer) (*ENSFilterer, error) {
	contract, err := bindENS(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ENSFilterer{contract: contract}, nil
}

// bindENS binds a generic wrapper to an already deployed contract.
func bindENS(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ENSABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ENS *ENSRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ENS.Contract.ENSCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ENS *ENSRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ENS.Contract.ENSTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ENS *ENSRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ENS.Contract.ENSTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ENS *ENSCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ENS.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ENS *ENSTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ENS.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ENS *ENSTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ENS.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x02571be3.
//
// Solidity: function owner(node bytes32) constant returns(address)
func (_ENS *ENSCaller) Owner(opts *bind.CallOpts, node string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ENS.contract.Call(opts, out, "owner", node)
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x02571be3.
//
// Solidity: function owner(node bytes32) constant returns(address)
func (_ENS *ENSSession) Owner(node string) (common.Address, error) {
	return _ENS.Contract.Owner(&_ENS.CallOpts, node)
}

// Owner is a free data retrieval call binding the contract method 0x02571be3.
//
// Solidity: function owner(node bytes32) constant returns(address)
func (_ENS *ENSCallerSession) Owner(node string) (common.Address, error) {
	return _ENS.Contract.Owner(&_ENS.CallOpts, node)
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// Solidity: function resolver(node bytes32) constant returns(address)
func (_ENS *ENSCaller) Resolver(opts *bind.CallOpts, node string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ENS.contract.Call(opts, out, "resolver", node)
	return *ret0, err
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// Solidity: function resolver(node bytes32) constant returns(address)
func (_ENS *ENSSession) Resolver(node string) (common.Address, error) {
	return _ENS.Contract.Resolver(&_ENS.CallOpts, node)
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// Solidity: function resolver(node bytes32) constant returns(address)
func (_ENS *ENSCallerSession) Resolver(node string) (common.Address, error) {
	return _ENS.Contract.Resolver(&_ENS.CallOpts, node)
}

// Ttl is a free data retrieval call binding the contract method 0x16a25cbd.
//
// Solidity: function ttl(node bytes32) constant returns(uint64)
func (_ENS *ENSCaller) Ttl(opts *bind.CallOpts, node string) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _ENS.contract.Call(opts, out, "ttl", node)
	return *ret0, err
}

// Ttl is a free data retrieval call binding the contract method 0x16a25cbd.
//
// Solidity: function ttl(node bytes32) constant returns(uint64)
func (_ENS *ENSSession) Ttl(node string) (uint64, error) {
	return _ENS.Contract.Ttl(&_ENS.CallOpts, node)
}

// Ttl is a free data retrieval call binding the contract method 0x16a25cbd.
//
// Solidity: function ttl(node bytes32) constant returns(uint64)
func (_ENS *ENSCallerSession) Ttl(node string) (uint64, error) {
	return _ENS.Contract.Ttl(&_ENS.CallOpts, node)
}

// SetOwner is a paid mutator transaction binding the contract method 0x5b0fc9c3.
//
// Solidity: function setOwner(node bytes32, owner address) returns()
func (_ENS *ENSTransactor) SetOwner(opts *bind.TransactOpts, node string, owner common.Address) (*types.Transaction, error) {
	return _ENS.contract.Transact(opts, "setOwner", node, owner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x5b0fc9c3.
//
// Solidity: function setOwner(node bytes32, owner address) returns()
func (_ENS *ENSSession) SetOwner(node string, owner common.Address) (*types.Transaction, error) {
	return _ENS.Contract.SetOwner(&_ENS.TransactOpts, node, owner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x5b0fc9c3.
//
// Solidity: function setOwner(node bytes32, owner address) returns()
func (_ENS *ENSTransactorSession) SetOwner(node string, owner common.Address) (*types.Transaction, error) {
	return _ENS.Contract.SetOwner(&_ENS.TransactOpts, node, owner)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// Solidity: function setResolver(node bytes32, resolver address) returns()
func (_ENS *ENSTransactor) SetResolver(opts *bind.TransactOpts, node string, resolver common.Address) (*types.Transaction, error) {
	return _ENS.contract.Transact(opts, "setResolver", node, resolver)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// Solidity: function setResolver(node bytes32, resolver address) returns()
func (_ENS *ENSSession) SetResolver(node string, resolver common.Address) (*types.Transaction, error) {
	return _ENS.Contract.SetResolver(&_ENS.TransactOpts, node, resolver)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// Solidity: function setResolver(node bytes32, resolver address) returns()
func (_ENS *ENSTransactorSession) SetResolver(node string, resolver common.Address) (*types.Transaction, error) {
	return _ENS.Contract.SetResolver(&_ENS.TransactOpts, node, resolver)
}

// SetSubnodeOwner is a paid mutator transaction binding the contract method 0x06ab5923.
//
// Solidity: function setSubnodeOwner(node bytes32, label bytes32, owner address) returns()
func (_ENS *ENSTransactor) SetSubnodeOwner(opts *bind.TransactOpts, node string, label string, owner common.Address) (*types.Transaction, error) {
	return _ENS.contract.Transact(opts, "setSubnodeOwner", node, label, owner)
}

// SetSubnodeOwner is a paid mutator transaction binding the contract method 0x06ab5923.
//
// Solidity: function setSubnodeOwner(node bytes32, label bytes32, owner address) returns()
func (_ENS *ENSSession) SetSubnodeOwner(node string, label string, owner common.Address) (*types.Transaction, error) {
	return _ENS.Contract.SetSubnodeOwner(&_ENS.TransactOpts, node, label, owner)
}

// SetSubnodeOwner is a paid mutator transaction binding the contract method 0x06ab5923.
//
// Solidity: function setSubnodeOwner(node bytes32, label bytes32, owner address) returns()
func (_ENS *ENSTransactorSession) SetSubnodeOwner(node string, label string, owner common.Address) (*types.Transaction, error) {
	return _ENS.Contract.SetSubnodeOwner(&_ENS.TransactOpts, node, label, owner)
}

// SetTTL is a paid mutator transaction binding the contract method 0x14ab9038.
//
// Solidity: function setTTL(node bytes32, ttl uint64) returns()
func (_ENS *ENSTransactor) SetTTL(opts *bind.TransactOpts, node string, ttl uint64) (*types.Transaction, error) {
	return _ENS.contract.Transact(opts, "setTTL", node, ttl)
}

// SetTTL is a paid mutator transaction binding the contract method 0x14ab9038.
//
// Solidity: function setTTL(node bytes32, ttl uint64) returns()
func (_ENS *ENSSession) SetTTL(node string, ttl uint64) (*types.Transaction, error) {
	return _ENS.Contract.SetTTL(&_ENS.TransactOpts, node, ttl)
}

// SetTTL is a paid mutator transaction binding the contract method 0x14ab9038.
//
// Solidity: function setTTL(node bytes32, ttl uint64) returns()
func (_ENS *ENSTransactorSession) SetTTL(node string, ttl uint64) (*types.Transaction, error) {
	return _ENS.Contract.SetTTL(&_ENS.TransactOpts, node, ttl)
}

// ENSNewOwnerIterator is returned from FilterNewOwner and is used to iterate over the raw logs and unpacked data for NewOwner events raised by the ENS contract.
type ENSNewOwnerIterator struct {
	Event *ENSNewOwner // Event containing the contract specifics and raw log

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
func (it *ENSNewOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ENSNewOwner)
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
		it.Event = new(ENSNewOwner)
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
func (it *ENSNewOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ENSNewOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ENSNewOwner represents a NewOwner event raised by the ENS contract.
type ENSNewOwner struct {
	Node  string
	Label string
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterNewOwner is a free log retrieval operation binding the contract event 0xce0457fe73731f824cc272376169235128c118b49d344817417c6d108d155e82.
//
// Solidity: e NewOwner(node indexed bytes32, label indexed bytes32, owner address)
func (_ENS *ENSFilterer) FilterNewOwner(opts *bind.FilterOpts, node []string, label []string) (*ENSNewOwnerIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}
	var labelRule []interface{}
	for _, labelItem := range label {
		labelRule = append(labelRule, labelItem)
	}

	logs, sub, err := _ENS.contract.FilterLogs(opts, "NewOwner", nodeRule, labelRule)
	if err != nil {
		return nil, err
	}
	return &ENSNewOwnerIterator{contract: _ENS.contract, event: "NewOwner", logs: logs, sub: sub}, nil
}

// WatchNewOwner is a free log subscription operation binding the contract event 0xce0457fe73731f824cc272376169235128c118b49d344817417c6d108d155e82.
//
// Solidity: e NewOwner(node indexed bytes32, label indexed bytes32, owner address)
func (_ENS *ENSFilterer) WatchNewOwner(opts *bind.WatchOpts, sink chan<- *ENSNewOwner, node []string, label []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}
	var labelRule []interface{}
	for _, labelItem := range label {
		labelRule = append(labelRule, labelItem)
	}

	logs, sub, err := _ENS.contract.WatchLogs(opts, "NewOwner", nodeRule, labelRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ENSNewOwner)
				if err := _ENS.contract.UnpackLog(event, "NewOwner", log); err != nil {
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

// ENSNewResolverIterator is returned from FilterNewResolver and is used to iterate over the raw logs and unpacked data for NewResolver events raised by the ENS contract.
type ENSNewResolverIterator struct {
	Event *ENSNewResolver // Event containing the contract specifics and raw log

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
func (it *ENSNewResolverIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ENSNewResolver)
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
		it.Event = new(ENSNewResolver)
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
func (it *ENSNewResolverIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ENSNewResolverIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ENSNewResolver represents a NewResolver event raised by the ENS contract.
type ENSNewResolver struct {
	Node     string
	Resolver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewResolver is a free log retrieval operation binding the contract event 0x335721b01866dc23fbee8b6b2c7b1e14d6f05c28cd35a2c934239f94095602a0.
//
// Solidity: e NewResolver(node indexed bytes32, resolver address)
func (_ENS *ENSFilterer) FilterNewResolver(opts *bind.FilterOpts, node []string) (*ENSNewResolverIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ENS.contract.FilterLogs(opts, "NewResolver", nodeRule)
	if err != nil {
		return nil, err
	}
	return &ENSNewResolverIterator{contract: _ENS.contract, event: "NewResolver", logs: logs, sub: sub}, nil
}

// WatchNewResolver is a free log subscription operation binding the contract event 0x335721b01866dc23fbee8b6b2c7b1e14d6f05c28cd35a2c934239f94095602a0.
//
// Solidity: e NewResolver(node indexed bytes32, resolver address)
func (_ENS *ENSFilterer) WatchNewResolver(opts *bind.WatchOpts, sink chan<- *ENSNewResolver, node []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ENS.contract.WatchLogs(opts, "NewResolver", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ENSNewResolver)
				if err := _ENS.contract.UnpackLog(event, "NewResolver", log); err != nil {
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

// ENSNewTTLIterator is returned from FilterNewTTL and is used to iterate over the raw logs and unpacked data for NewTTL events raised by the ENS contract.
type ENSNewTTLIterator struct {
	Event *ENSNewTTL // Event containing the contract specifics and raw log

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
func (it *ENSNewTTLIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ENSNewTTL)
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
		it.Event = new(ENSNewTTL)
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
func (it *ENSNewTTLIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ENSNewTTLIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ENSNewTTL represents a NewTTL event raised by the ENS contract.
type ENSNewTTL struct {
	Node string
	Ttl  uint64
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNewTTL is a free log retrieval operation binding the contract event 0x1d4f9bbfc9cab89d66e1a1562f2233ccbf1308cb4f63de2ead5787adddb8fa68.
//
// Solidity: e NewTTL(node indexed bytes32, ttl uint64)
func (_ENS *ENSFilterer) FilterNewTTL(opts *bind.FilterOpts, node []string) (*ENSNewTTLIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ENS.contract.FilterLogs(opts, "NewTTL", nodeRule)
	if err != nil {
		return nil, err
	}
	return &ENSNewTTLIterator{contract: _ENS.contract, event: "NewTTL", logs: logs, sub: sub}, nil
}

// WatchNewTTL is a free log subscription operation binding the contract event 0x1d4f9bbfc9cab89d66e1a1562f2233ccbf1308cb4f63de2ead5787adddb8fa68.
//
// Solidity: e NewTTL(node indexed bytes32, ttl uint64)
func (_ENS *ENSFilterer) WatchNewTTL(opts *bind.WatchOpts, sink chan<- *ENSNewTTL, node []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ENS.contract.WatchLogs(opts, "NewTTL", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ENSNewTTL)
				if err := _ENS.contract.UnpackLog(event, "NewTTL", log); err != nil {
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

// ENSTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the ENS contract.
type ENSTransferIterator struct {
	Event *ENSTransfer // Event containing the contract specifics and raw log

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
func (it *ENSTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ENSTransfer)
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
		it.Event = new(ENSTransfer)
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
func (it *ENSTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ENSTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ENSTransfer represents a Transfer event raised by the ENS contract.
type ENSTransfer struct {
	Node  string
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xd4735d920b0f87494915f556dd9b54c8f309026070caea5c737245152564d266.
//
// Solidity: e Transfer(node indexed bytes32, owner address)
func (_ENS *ENSFilterer) FilterTransfer(opts *bind.FilterOpts, node []string) (*ENSTransferIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ENS.contract.FilterLogs(opts, "Transfer", nodeRule)
	if err != nil {
		return nil, err
	}
	return &ENSTransferIterator{contract: _ENS.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xd4735d920b0f87494915f556dd9b54c8f309026070caea5c737245152564d266.
//
// Solidity: e Transfer(node indexed bytes32, owner address)
func (_ENS *ENSFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *ENSTransfer, node []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _ENS.contract.WatchLogs(opts, "Transfer", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ENSTransfer)
				if err := _ENS.contract.UnpackLog(event, "Transfer", log); err != nil {
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
