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

// VNSABI is the input ABI used to generate the binding from.
const VNSABI = `[{"name":"VNS","constant":false,"inputs":[],"outputs":[],"type":"constructor"},{"name":"setOwner","constant":false,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"owner","type":"address","indexed":false}],"outputs":[],"type":"function"},{"name":"setSubnodeOwner","constant":false,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"label","type":"string","indexed":false},{"name":"owner","type":"address","indexed":false}],"outputs":[],"type":"function"},{"name":"setResolver","constant":false,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"resolver","type":"address","indexed":false}],"outputs":[],"type":"function"},{"name":"setTTL","constant":false,"inputs":[{"name":"node","type":"string","indexed":false},{"name":"ttl","type":"uint64","indexed":false}],"outputs":[],"type":"function"},{"name":"owner","constant":true,"inputs":[{"name":"node","type":"string","indexed":false}],"outputs":[{"name":"output","type":"address","indexed":false}],"type":"function"},{"name":"resolver","constant":true,"inputs":[{"name":"node","type":"string","indexed":false}],"outputs":[{"name":"output","type":"address","indexed":false}],"type":"function"},{"name":"ttl","constant":true,"inputs":[{"name":"node","type":"string","indexed":false}],"outputs":[{"name":"output","type":"uint64","indexed":false}],"type":"function"},{"name":"NewResolver","anonymous":false,"inputs":[{"name":"node","type":"string","indexed":true},{"name":"resolver","type":"address","indexed":false}],"type":"event"},{"name":"NewTTL","anonymous":false,"inputs":[{"name":"node","type":"string","indexed":true},{"name":"ttl","type":"uint64","indexed":false}],"type":"event"},{"name":"NewOwner","anonymous":false,"inputs":[{"name":"node","type":"string","indexed":true},{"name":"label","type":"string","indexed":true},{"name":"owner","type":"address","indexed":false}],"type":"event"},{"name":"Transfer","anonymous":false,"inputs":[{"name":"node","type":"string","indexed":true},{"name":"owner","type":"address","indexed":false}],"type":"event"}]`

// VNSBin is the compiled bytecode used for deploying new contracts.
const VNSBin = `0x0161736db90a5c0100789ccc586b6c1cd515feeebd7bc7eb1dafbd9b4decb80930d98008e0aef3b493901036d086501a909d84a614d9e3ddd9f58af58c999975086dbc6ec22329491b54a915a27543ff94d282a2b4954aa25654a895f8c10f2a41a44a6d512bf851448b78a8140a54e7ceecc376121ea969f363e6dcf3f8ce77ef9c7b9c3dff5af2cd27535598de1803c036b70db32a86f9e42486d92486e56475b25ac530300c561de6d56af0c430234954437172985727c1cfb0b8b0ec096dc09ab05c1fa045e256b7e45bb796fcd15b9c92ed5b2e18a93b062c333f47ab65f3f96da6074e8b58369fff82b56fbb5d70204891da6e97fc92592edd6ded36dd923952b63c44c8d2bacdf2072d3b6fb990b4969fbbb36296a1911cdde99ab657b05cb4a814d73976cef4035b64f086ec1a4495db0e6befcd7b6dcb452b2df51dd6de01cb73ca132ab04d8f4422916834d6d2daa247241b678c49c1c0b545a2cab2270e45f42a6aaf9f46f59607589b36668d39ee3e0e7d6868d432c787464ccf122c363494377d73c8b2f38227f256ae6cba56fef3153be7971c1bed62f78e4124a4a3b824a36e8dc402e1fb65a4a29ee5073c177678963f5819b19dbc156816e99ee5d759776a9ee5efdc7913ba62cf4a9d43bf185bb11589af4f4d4d05d281404a1ca4b77e82313e993d72686a0a8fa7916dc94e0522cbb66691b8877c0c645b0c96d566acd1b0678fd5a223c6ccb0c8acb0c88cb043f5a4c6c7c8a6f720f12d1212f72a14648f134e1f57c4b348dc37a5960a5d5f1a3a6761345c569cdf78e4dcc62b84daabbe67ae915f5b070e2824eea7671cd947946e4a913760b0c4a139b806ab117618ab5e08f66165fe063d9709d2183cf14008271af9782ddf85ede4c85977121ce0cd9f1479165e7f78e6fa068eec545437f061ffb2c7a37a1b6ccbca1bc185c233513d93b7462ac5a1925d7036448088f28c5003441b9687910f87ef4701f03f00e827419782fcc5268a30b82eff02c05216c8d7017c4dc951b99e01a795bc18da388057a8bfae6ca90088be0c4063dafd00563180b746123c22df0670824234488d85724b4476d4a0a2d100ea7d825aa3fd0ec0950c603dd114035aa2adb48524803701c4fe4d8f1652e912c0c5448ebc5f6bbb915e432c3ea0de141b1f51226d04d1af1025d64e91db001c25d3dbab22405cf0bac8ea628724f11d121394781d800e62f86e901d0aa900e01142faa081f4410349894985c419c90b08aa8ba2098ab1a6ed28d3ce5a96080bd3d03e54f2c7294d54a12c5479429935e4456f865f97c566845f0be0490a8f3785c79bc2033920dad1805a4c50891a947e59f8257f4b277991fe02801e00cfd17285fa4a6b00bcca81d85b1ce884003a7541bb150c580a808a815f4720f1710abb4e7d9b762334537df081867920306f0ecd54327cf874dd3caccc6d0592fb59a0eb275d87cab391803a2aa4bd9a44747d99017b003c4f3a2f08f02820ae023cf25a4c1b5849cb89d9db40d7cf19e000f813451e0e000e2b80165e1715d661953c4edaa3560df616b21c990b7b290f78fd9df0a603d86985a54e64bace8bee097f642ec024076e06f00eb99f0c004e2a007566274fd7008e11c02fe600a80bf81095c86a20d11eb692313347dd04f8358f872ad7b48b9647b73105e012005735f5252a870c80d500ae69d29b0c789001df63c034037ec4801fb386fd570c7895016d1c68e7401707ba79c37e15076ee3408103a31cb893035e93fd5e0efc9203bfe1c0d34d7a75af0050a7581eae6f08dfc5f0dd17beb300a8c29f0740a5cd59a05f0b80caf71206509d6e620015e41e06503dfd8c015416cb3940df713f07e8733cdbc40367ea07688e8cb8d6044bb2cbda5352b427ba96b427d9ae2e80af8568df9edad2bd51bf5ae78b01b10edb534024c5749d74808c2b175aacd701ad13db53a23db0b57428efe8a510edd7e8ba0eb46658922dd056a5809804bdf40c4813846ce906da322c585cdebda5dbd081b844a020acf600a0e641dc808eb54d1e898bd8aa14f97c49dfa303c9d026c9b6a0c92681541398ac812d944d018b28605757e0deb95631ee9a49816216d7dd7475b2f53f7ce5926d6d537ff8de0b9aebbbed717ab1e003b3de5d9ee57abd7bad926796ee1e35ed62eff54eae3266d9bed75b747a3d37d75b2cf9a395914cce19eb9db0fddca859b27b8bce67276cbf37e7d8be6be67caf77c2f6eaab267506f39dc1a9f8e3151f98b0fd726924334adbdabd637068dcb572ced878a96c6572506527a97a21fe0c265b63cb4635f1fdea01c934f1120ccec01824a73b215e634ccad89e1f1c94ac5833a8c8b71993c9d866193f2d2fd75e976b36c92ba765ab265e9e3c700da19d925213ef31838ba6a00f1893f1d86699d6a20639bf3869ac0c9ce39ab8871b32addd28935a8f4c1f97abb503b27b5ac60be4f8c7c94203f53e3e13f51027d46572a9f602f99e9934649b788c65a4d4a6c55bfb55a6bf82948f3202989eb191a3bcbe91a5da299936e46519b9f614053d376948294eb08c26feb17f530d468a27d8294520c22139b504b19449167b305a7923ba279a0a1e6fa82cca85ba8558412ecb8a3259944ba229d9ad1524d32a0d2745a69fb82c93cb5f929d15d95ad40ada4a4d6c3d455fe6a9f0cb2867ea37e269284842fccc33526a7ec34e7d48fcfedc76ea4fe2c573dba96f89376a76717d612e25a9897fc29049b1755a6cc988ec31d92d9e4246b2071b30d4f6c46146309b65523c513d0bce124d1c6586ecd4a6a5d40ecaa4f849f594780119f1589520cfcc82a40e2aa6598dd9436783949af821236adfad4e8bbf2123be533d26978a57664151131627eb50f79e0d2aa9890d86780f1971b07a4c7689f71b18381e89859dc5f3dd5cd9b48bc684e57a25c736d6675666561a2b7cb762df61ac59bb6a43dfca2b3ea5dbdf3be7c25f585e10949bf77087b50fb951d385e7bb25bb8809b35cb110fcbc30f379d7f23cd47fb5d36ff6b263170df5a8d85ea9685b79a364fba8946cbf6f2d06142cc6ccf1f1925d5cb56e4dffaaf5abfbfbfb3060dd5929b916728e9d2fa9d9c0d056c72963cc2bc2b1cbfb82dffcf4eb1f5e30fb98334cd8b57a5ddfa059b0be5829ab74abd7f5d575d79726eaf26065a42e67f379da62de1cebeb2fac5f477d13f599c3ec8943f3bc219c36ecc35db81b6573c42ac30b7c818779c436c72cf66dbe18e17888cd990cf1d9432111ce83224da32079b62990d61800b504b39f687decd31a4e7c626ad8a3d7e73c6dcd239e78726868afe98d0de5cc727928e73baed73e6736d3a1379d4b42ecde31980cc6340bea539a94f0fdf2c2fa8c66d1ec114d67f384a62b1cd03c198ddef6d5349d4f7a637af78ec1744f3ae7d89e6fda7e7a63c12c7b564fba648f577c2fbdf1b6db7bd2417d870b7fdf3885a900b742ccd3fb7bea683522e785ac7b13cd741d3228ef3439e6adbbac7c18d884ee84d0614058fc7323cec1b9109eeb2cc2cde7353fbc5571feef375aab83f9d9a4db40ffefd2deb9f3a6f961ecfb4d1f25e88e17c4d5995341be5bb900a6335237b228d5473be5f39075cf560cffc77c838ff5e9523d4f4d9c8769539b4ff7a44ddbb1f78d3915ef020a97f6fa896f5ae8624d58b63f8b6670b3e681e147bd59e7e1566bc7f3c0ee43dbf12cff8fd18dcfb5a1da7f0ee667439f98e0ed53ff010000ffff`

// DeployVNS deploys a new VNT contract, binding an instance of VNS to it.
func DeployVNS(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VNS, error) {
	parsed, err := abi.JSON(strings.NewReader(VNSABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VNSBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VNS{VNSCaller: VNSCaller{contract: contract}, VNSTransactor: VNSTransactor{contract: contract}, VNSFilterer: VNSFilterer{contract: contract}}, nil
}

// VNS is an auto generated Go binding around an VNT contract.
type VNS struct {
	VNSCaller     // Read-only binding to the contract
	VNSTransactor // Write-only binding to the contract
	VNSFilterer   // Log filterer for contract events
}

// VNSCaller is an auto generated read-only Go binding around an VNT contract.
type VNSCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VNSTransactor is an auto generated write-only Go binding around an VNT contract.
type VNSTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VNSFilterer is an auto generated log filtering Go binding around an VNT contract events.
type VNSFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VNSSession is an auto generated Go binding around an VNT contract,
// with pre-set call and transact options.
type VNSSession struct {
	Contract     *VNS              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VNSCallerSession is an auto generated read-only Go binding around an VNT contract,
// with pre-set call options.
type VNSCallerSession struct {
	Contract *VNSCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VNSTransactorSession is an auto generated write-only Go binding around an VNT contract,
// with pre-set transact options.
type VNSTransactorSession struct {
	Contract     *VNSTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VNSRaw is an auto generated low-level Go binding around an VNT contract.
type VNSRaw struct {
	Contract *VNS // Generic contract binding to access the raw methods on
}

// VNSCallerRaw is an auto generated low-level read-only Go binding around an VNT contract.
type VNSCallerRaw struct {
	Contract *VNSCaller // Generic read-only contract binding to access the raw methods on
}

// VNSTransactorRaw is an auto generated low-level write-only Go binding around an VNT contract.
type VNSTransactorRaw struct {
	Contract *VNSTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVNS creates a new instance of VNS, bound to a specific deployed contract.
func NewVNS(address common.Address, backend bind.ContractBackend) (*VNS, error) {
	contract, err := bindVNS(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VNS{VNSCaller: VNSCaller{contract: contract}, VNSTransactor: VNSTransactor{contract: contract}, VNSFilterer: VNSFilterer{contract: contract}}, nil
}

// NewVNSCaller creates a new read-only instance of VNS, bound to a specific deployed contract.
func NewVNSCaller(address common.Address, caller bind.ContractCaller) (*VNSCaller, error) {
	contract, err := bindVNS(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VNSCaller{contract: contract}, nil
}

// NewVNSTransactor creates a new write-only instance of VNS, bound to a specific deployed contract.
func NewVNSTransactor(address common.Address, transactor bind.ContractTransactor) (*VNSTransactor, error) {
	contract, err := bindVNS(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VNSTransactor{contract: contract}, nil
}

// NewVNSFilterer creates a new log filterer instance of VNS, bound to a specific deployed contract.
func NewVNSFilterer(address common.Address, filterer bind.ContractFilterer) (*VNSFilterer, error) {
	contract, err := bindVNS(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VNSFilterer{contract: contract}, nil
}

// bindVNS binds a generic wrapper to an already deployed contract.
func bindVNS(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VNSABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VNS *VNSRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VNS.Contract.VNSCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VNS *VNSRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VNS.Contract.VNSTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VNS *VNSRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VNS.Contract.VNSTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VNS *VNSCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _VNS.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VNS *VNSTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VNS.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VNS *VNSTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VNS.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x02571be3.
//
// function owner(node bytes32) constant returns(address)
func (_VNS *VNSCaller) Owner(opts *bind.CallOpts, node string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _VNS.contract.Call(opts, out, "owner", node)
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x02571be3.
//
// function owner(node bytes32) constant returns(address)
func (_VNS *VNSSession) Owner(node string) (common.Address, error) {
	return _VNS.Contract.Owner(&_VNS.CallOpts, node)
}

// Owner is a free data retrieval call binding the contract method 0x02571be3.
//
// function owner(node bytes32) constant returns(address)
func (_VNS *VNSCallerSession) Owner(node string) (common.Address, error) {
	return _VNS.Contract.Owner(&_VNS.CallOpts, node)
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// function resolver(node bytes32) constant returns(address)
func (_VNS *VNSCaller) Resolver(opts *bind.CallOpts, node string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _VNS.contract.Call(opts, out, "resolver", node)
	return *ret0, err
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// function resolver(node bytes32) constant returns(address)
func (_VNS *VNSSession) Resolver(node string) (common.Address, error) {
	return _VNS.Contract.Resolver(&_VNS.CallOpts, node)
}

// Resolver is a free data retrieval call binding the contract method 0x0178b8bf.
//
// function resolver(node bytes32) constant returns(address)
func (_VNS *VNSCallerSession) Resolver(node string) (common.Address, error) {
	return _VNS.Contract.Resolver(&_VNS.CallOpts, node)
}

// Ttl is a free data retrieval call binding the contract method 0x16a25cbd.
//
// function ttl(node bytes32) constant returns(uint64)
func (_VNS *VNSCaller) Ttl(opts *bind.CallOpts, node string) (uint64, error) {
	var (
		ret0 = new(uint64)
	)
	out := ret0
	err := _VNS.contract.Call(opts, out, "ttl", node)
	return *ret0, err
}

// Ttl is a free data retrieval call binding the contract method 0x16a25cbd.
//
// function ttl(node bytes32) constant returns(uint64)
func (_VNS *VNSSession) Ttl(node string) (uint64, error) {
	return _VNS.Contract.Ttl(&_VNS.CallOpts, node)
}

// Ttl is a free data retrieval call binding the contract method 0x16a25cbd.
//
// function ttl(node bytes32) constant returns(uint64)
func (_VNS *VNSCallerSession) Ttl(node string) (uint64, error) {
	return _VNS.Contract.Ttl(&_VNS.CallOpts, node)
}

// SetOwner is a paid mutator transaction binding the contract method 0x5b0fc9c3.
//
// function setOwner(node bytes32, owner address) returns()
func (_VNS *VNSTransactor) SetOwner(opts *bind.TransactOpts, node string, owner common.Address) (*types.Transaction, error) {
	return _VNS.contract.Transact(opts, "setOwner", node, owner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x5b0fc9c3.
//
// function setOwner(node bytes32, owner address) returns()
func (_VNS *VNSSession) SetOwner(node string, owner common.Address) (*types.Transaction, error) {
	return _VNS.Contract.SetOwner(&_VNS.TransactOpts, node, owner)
}

// SetOwner is a paid mutator transaction binding the contract method 0x5b0fc9c3.
//
// function setOwner(node bytes32, owner address) returns()
func (_VNS *VNSTransactorSession) SetOwner(node string, owner common.Address) (*types.Transaction, error) {
	return _VNS.Contract.SetOwner(&_VNS.TransactOpts, node, owner)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// function setResolver(node bytes32, resolver address) returns()
func (_VNS *VNSTransactor) SetResolver(opts *bind.TransactOpts, node string, resolver common.Address) (*types.Transaction, error) {
	return _VNS.contract.Transact(opts, "setResolver", node, resolver)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// function setResolver(node bytes32, resolver address) returns()
func (_VNS *VNSSession) SetResolver(node string, resolver common.Address) (*types.Transaction, error) {
	return _VNS.Contract.SetResolver(&_VNS.TransactOpts, node, resolver)
}

// SetResolver is a paid mutator transaction binding the contract method 0x1896f70a.
//
// function setResolver(node bytes32, resolver address) returns()
func (_VNS *VNSTransactorSession) SetResolver(node string, resolver common.Address) (*types.Transaction, error) {
	return _VNS.Contract.SetResolver(&_VNS.TransactOpts, node, resolver)
}

// SetSubnodeOwner is a paid mutator transaction binding the contract method 0x06ab5923.
//
// function setSubnodeOwner(node bytes32, label bytes32, owner address) returns()
func (_VNS *VNSTransactor) SetSubnodeOwner(opts *bind.TransactOpts, node string, label string, owner common.Address) (*types.Transaction, error) {
	return _VNS.contract.Transact(opts, "setSubnodeOwner", node, label, owner)
}

// SetSubnodeOwner is a paid mutator transaction binding the contract method 0x06ab5923.
//
// function setSubnodeOwner(node bytes32, label bytes32, owner address) returns()
func (_VNS *VNSSession) SetSubnodeOwner(node string, label string, owner common.Address) (*types.Transaction, error) {
	return _VNS.Contract.SetSubnodeOwner(&_VNS.TransactOpts, node, label, owner)
}

// SetSubnodeOwner is a paid mutator transaction binding the contract method 0x06ab5923.
//
// function setSubnodeOwner(node bytes32, label bytes32, owner address) returns()
func (_VNS *VNSTransactorSession) SetSubnodeOwner(node string, label string, owner common.Address) (*types.Transaction, error) {
	return _VNS.Contract.SetSubnodeOwner(&_VNS.TransactOpts, node, label, owner)
}

// SetTTL is a paid mutator transaction binding the contract method 0x14ab9038.
//
// function setTTL(node bytes32, ttl uint64) returns()
func (_VNS *VNSTransactor) SetTTL(opts *bind.TransactOpts, node string, ttl uint64) (*types.Transaction, error) {
	return _VNS.contract.Transact(opts, "setTTL", node, ttl)
}

// SetTTL is a paid mutator transaction binding the contract method 0x14ab9038.
//
// function setTTL(node bytes32, ttl uint64) returns()
func (_VNS *VNSSession) SetTTL(node string, ttl uint64) (*types.Transaction, error) {
	return _VNS.Contract.SetTTL(&_VNS.TransactOpts, node, ttl)
}

// SetTTL is a paid mutator transaction binding the contract method 0x14ab9038.
//
// function setTTL(node bytes32, ttl uint64) returns()
func (_VNS *VNSTransactorSession) SetTTL(node string, ttl uint64) (*types.Transaction, error) {
	return _VNS.Contract.SetTTL(&_VNS.TransactOpts, node, ttl)
}

// VNSNewOwnerIterator is returned from FilterNewOwner and is used to iterate over the raw logs and unpacked data for NewOwner events raised by the VNS contract.
type VNSNewOwnerIterator struct {
	Event *VNSNewOwner // Event containing the contract specifics and raw log

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
func (it *VNSNewOwnerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VNSNewOwner)
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
		it.Event = new(VNSNewOwner)
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
func (it *VNSNewOwnerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VNSNewOwnerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VNSNewOwner represents a NewOwner event raised by the VNS contract.
type VNSNewOwner struct {
	Node  string
	Label string
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterNewOwner is a free log retrieval operation binding the contract event 0xce0457fe73731f824cc272376169235128c118b49d344817417c6d108d155e82.
//
// e NewOwner(node indexed bytes32, label indexed bytes32, owner address)
func (_VNS *VNSFilterer) FilterNewOwner(opts *bind.FilterOpts, node []string, label []string) (*VNSNewOwnerIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}
	var labelRule []interface{}
	for _, labelItem := range label {
		labelRule = append(labelRule, labelItem)
	}

	logs, sub, err := _VNS.contract.FilterLogs(opts, "NewOwner", nodeRule, labelRule)
	if err != nil {
		return nil, err
	}
	return &VNSNewOwnerIterator{contract: _VNS.contract, event: "NewOwner", logs: logs, sub: sub}, nil
}

// WatchNewOwner is a free log subscription operation binding the contract event 0xce0457fe73731f824cc272376169235128c118b49d344817417c6d108d155e82.
//
// e NewOwner(node indexed bytes32, label indexed bytes32, owner address)
func (_VNS *VNSFilterer) WatchNewOwner(opts *bind.WatchOpts, sink chan<- *VNSNewOwner, node []string, label []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}
	var labelRule []interface{}
	for _, labelItem := range label {
		labelRule = append(labelRule, labelItem)
	}

	logs, sub, err := _VNS.contract.WatchLogs(opts, "NewOwner", nodeRule, labelRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VNSNewOwner)
				if err := _VNS.contract.UnpackLog(event, "NewOwner", log); err != nil {
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

// VNSNewResolverIterator is returned from FilterNewResolver and is used to iterate over the raw logs and unpacked data for NewResolver events raised by the VNS contract.
type VNSNewResolverIterator struct {
	Event *VNSNewResolver // Event containing the contract specifics and raw log

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
func (it *VNSNewResolverIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VNSNewResolver)
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
		it.Event = new(VNSNewResolver)
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
func (it *VNSNewResolverIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VNSNewResolverIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VNSNewResolver represents a NewResolver event raised by the VNS contract.
type VNSNewResolver struct {
	Node     string
	Resolver common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewResolver is a free log retrieval operation binding the contract event 0x335721b01866dc23fbee8b6b2c7b1e14d6f05c28cd35a2c934239f94095602a0.
//
// e NewResolver(node indexed bytes32, resolver address)
func (_VNS *VNSFilterer) FilterNewResolver(opts *bind.FilterOpts, node []string) (*VNSNewResolverIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _VNS.contract.FilterLogs(opts, "NewResolver", nodeRule)
	if err != nil {
		return nil, err
	}
	return &VNSNewResolverIterator{contract: _VNS.contract, event: "NewResolver", logs: logs, sub: sub}, nil
}

// WatchNewResolver is a free log subscription operation binding the contract event 0x335721b01866dc23fbee8b6b2c7b1e14d6f05c28cd35a2c934239f94095602a0.
//
// e NewResolver(node indexed bytes32, resolver address)
func (_VNS *VNSFilterer) WatchNewResolver(opts *bind.WatchOpts, sink chan<- *VNSNewResolver, node []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _VNS.contract.WatchLogs(opts, "NewResolver", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VNSNewResolver)
				if err := _VNS.contract.UnpackLog(event, "NewResolver", log); err != nil {
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

// VNSNewTTLIterator is returned from FilterNewTTL and is used to iterate over the raw logs and unpacked data for NewTTL events raised by the VNS contract.
type VNSNewTTLIterator struct {
	Event *VNSNewTTL // Event containing the contract specifics and raw log

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
func (it *VNSNewTTLIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VNSNewTTL)
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
		it.Event = new(VNSNewTTL)
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
func (it *VNSNewTTLIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VNSNewTTLIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VNSNewTTL represents a NewTTL event raised by the VNS contract.
type VNSNewTTL struct {
	Node string
	Ttl  uint64
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNewTTL is a free log retrieval operation binding the contract event 0x1d4f9bbfc9cab89d66e1a1562f2233ccbf1308cb4f63de2ead5787adddb8fa68.
//
// e NewTTL(node indexed bytes32, ttl uint64)
func (_VNS *VNSFilterer) FilterNewTTL(opts *bind.FilterOpts, node []string) (*VNSNewTTLIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _VNS.contract.FilterLogs(opts, "NewTTL", nodeRule)
	if err != nil {
		return nil, err
	}
	return &VNSNewTTLIterator{contract: _VNS.contract, event: "NewTTL", logs: logs, sub: sub}, nil
}

// WatchNewTTL is a free log subscription operation binding the contract event 0x1d4f9bbfc9cab89d66e1a1562f2233ccbf1308cb4f63de2ead5787adddb8fa68.
//
// e NewTTL(node indexed bytes32, ttl uint64)
func (_VNS *VNSFilterer) WatchNewTTL(opts *bind.WatchOpts, sink chan<- *VNSNewTTL, node []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _VNS.contract.WatchLogs(opts, "NewTTL", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VNSNewTTL)
				if err := _VNS.contract.UnpackLog(event, "NewTTL", log); err != nil {
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

// VNSTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the VNS contract.
type VNSTransferIterator struct {
	Event *VNSTransfer // Event containing the contract specifics and raw log

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
func (it *VNSTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VNSTransfer)
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
		it.Event = new(VNSTransfer)
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
func (it *VNSTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VNSTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VNSTransfer represents a Transfer event raised by the VNS contract.
type VNSTransfer struct {
	Node  string
	Owner common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xd4735d920b0f87494915f556dd9b54c8f309026070caea5c737245152564d266.
//
// e Transfer(node indexed bytes32, owner address)
func (_VNS *VNSFilterer) FilterTransfer(opts *bind.FilterOpts, node []string) (*VNSTransferIterator, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _VNS.contract.FilterLogs(opts, "Transfer", nodeRule)
	if err != nil {
		return nil, err
	}
	return &VNSTransferIterator{contract: _VNS.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xd4735d920b0f87494915f556dd9b54c8f309026070caea5c737245152564d266.
//
// e Transfer(node indexed bytes32, owner address)
func (_VNS *VNSFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *VNSTransfer, node []string) (event.Subscription, error) {

	var nodeRule []interface{}
	for _, nodeItem := range node {
		nodeRule = append(nodeRule, nodeItem)
	}

	logs, sub, err := _VNS.contract.WatchLogs(opts, "Transfer", nodeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VNSTransfer)
				if err := _VNS.contract.UnpackLog(event, "Transfer", log); err != nil {
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
