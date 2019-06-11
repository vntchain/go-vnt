// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

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
const PublicResolverABI = "[{\"name\":\"PublicResolver\",\"constant\":false,\"inputs\":[{\"name\":\"ensAddr\",\"type\":\"address\",\"indexed\":false}],\"outputs\":[],\"type\":\"constructor\"},{\"name\":\"ABIRecord\",\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"contentTypes\",\"type\":\"uint256\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"string\",\"indexed\":false}],\"type\":\"function\"},{\"name\":\"ABIContentType\",\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"contentTypes\",\"type\":\"uint256\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"uint256\",\"indexed\":false}],\"type\":\"function\"},{\"name\":\"setABI\",\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"contentType\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"data\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[],\"type\":\"function\"},{\"name\":\"text\",\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"key\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"string\",\"indexed\":false}],\"type\":\"function\"},{\"name\":\"setAddr\",\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"addr\",\"type\":\"address\",\"indexed\":false}],\"outputs\":[],\"type\":\"function\"},{\"name\":\"setText\",\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"key\",\"type\":\"string\",\"indexed\":false},{\"name\":\"value\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[],\"type\":\"function\"},{\"name\":\"supportsInterface\",\"constant\":true,\"inputs\":[{\"name\":\"interfaceID\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"bool\",\"indexed\":false}],\"type\":\"function\"},{\"name\":\"addr\",\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"address\",\"indexed\":false}],\"type\":\"function\"},{\"name\":\"content\",\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"string\",\"indexed\":false}],\"type\":\"function\"},{\"name\":\"setContent\",\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"hash\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[],\"type\":\"function\"},{\"name\":\"name\",\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"string\",\"indexed\":false}],\"type\":\"function\"},{\"name\":\"setName\",\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[],\"type\":\"function\"},{\"name\":\"pubkeyX\",\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"string\",\"indexed\":false}],\"type\":\"function\"},{\"name\":\"pubkeyY\",\"constant\":true,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"string\",\"indexed\":false}],\"type\":\"function\"},{\"name\":\"setPubkey\",\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false},{\"name\":\"x\",\"type\":\"string\",\"indexed\":false},{\"name\":\"y\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[],\"type\":\"function\"},{\"name\":\"ContentChanged\",\"anonymous\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":true},{\"name\":\"hash\",\"type\":\"string\",\"indexed\":false}],\"type\":\"event\"},{\"name\":\"NameChanged\",\"anonymous\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":true},{\"name\":\"name\",\"type\":\"string\",\"indexed\":false}],\"type\":\"event\"},{\"name\":\"ABIChanged\",\"anonymous\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":true},{\"name\":\"contentType\",\"type\":\"uint256\",\"indexed\":true}],\"type\":\"event\"},{\"name\":\"PubkeyChanged\",\"anonymous\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":true},{\"name\":\"x\",\"type\":\"string\",\"indexed\":false},{\"name\":\"y\",\"type\":\"string\",\"indexed\":false}],\"type\":\"event\"},{\"name\":\"TextChanged\",\"anonymous\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":true},{\"name\":\"indexedKey\",\"type\":\"string\",\"indexed\":true},{\"name\":\"key\",\"type\":\"string\",\"indexed\":false}],\"type\":\"event\"},{\"name\":\"AddrChanged\",\"anonymous\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":true},{\"name\":\"a\",\"type\":\"address\",\"indexed\":false}],\"type\":\"event\"},{\"name\":\"owner\",\"constant\":false,\"inputs\":[{\"name\":\"node\",\"type\":\"string\",\"indexed\":false}],\"outputs\":[{\"name\":\"output\",\"type\":\"address\",\"indexed\":false}],\"type\":\"call\"}]"

// PublicResolverBin is the compiled bytecode used for deploying new contracts.
const PublicResolverBin = `0x0161736db912a30100789cd45a0b701cc599feba7b7a76a491e49525db929fe3150ffb62f4966c7306bc7e4014c070c63c7c294e1eed8e568b57bbca3e14898b25e190e0048738319582aac49847aa280a282a49e5b0a1ee5295a48253dc4172397239928a7370099784546213122e505cfd3d333bfb908d0cf61d8f62e6efbfffc7f7fffd774faff8dfb8f8274756ff16766e9401603db5bbd83476f1a929ec6253d8c5a7a7d9f42e39353d353d8d5dc02eb06962629798a6279b66d3fc05de2c9cf4b8bedd1977b2798006e11bb3c9bc7363323f726d2699ce3b593062cfdbeed8f12aae1e8dc7afb073e03430aeefeeeb1fdc3c3a06110ca3f178e9f0bac2903bac8dc6e3573a9303e9e10c3462340da493f9a49d4adeeadc606793f650cac941d28ccc7c3ced645db59a2b9cfc754e3aee64a1abb9ad1f2bd82977ce8cc6e3d9cd23763ae1c411224ec3e64c3aefa4f3654c739b3dea94713ca4e972a423291fe9a68132f1fa6b0b43bb9d499f6728a33b9c897cc0696e9252a2a62644ff0a61d4d418c2d0241b638c49c1c0f505629a457fff19cd9c86ff7ab9c60cbdca9af5516734939de4300707471c7b6c70c8ce3982d50e0ec6edbc3de8a4e38287e34e2c65679df8e585742c9fcca4b1a0e1dac2502a19dbeee432a971278b458db9c2d858269bcf0dd07a0ddb31072d9a1d8f67d11aca39794a161687626e86b0a436e7e4bd7461a996b6471d2c2339ca1696d744370d6c7762996c1c2b1a2821aee48ec93107964ee6360d6065684c65e626443c6a27da6a724ede4d18ced3f2ce441ee79359ca172ea87daeb589c35c8e4dd884f06d3333332eb5d7a5c29fa4b7f945c9b4a9e8e1cfcccce0b108a2a1e83d2ec9a2355184ef20190bd190c5a27ad918c17cf4215f5bb7cad5f40a35bd4cedd177a55665e676d70c8fea162f613f1bb04badf30aeb3caa479f714545859888caa2b599d99d1c989d7ddc6787cad84fcc195259c095d6ef0ef2569a8737f7cdcade3fbbf4d3734e7ef4491f76b9180fc4a20fbfbb0ad8e7d75a999ad98af0978808ef5352b0d0cf554acd7b1863d3ee6494ad647ca385285671b53ae13b895d0f97757b356b5f356b7f35eb4035ebee6a964ab6cb5ac94c8b99470859db8c42f62dec8eb0f314ed615521a89d65315257fb2d827ebece62d1ce5b88f52362f5731258754b14abb92ae0b5821896627c2760b428c6d30123ac18df081886623c1a30b09a3f148cbe855b5c8016ccaf9640b776477829728b478d5b361dfe9c4e6a3c3ab34f816cb638413eee420e5bdc9568b1b8c5578bf05a01bed1e216c29f51963e3b5392bfe8cf54ea940fb324315158acdf4d8c058b85f77bee2d0fa8f9cd33cfb09b5fac7223a77cf7fae915cf57a4573c53915ef1ad8af48a272bd22b9e284daf78f8fd9fde878ae9fddc7b49ef6a71d88dd54befc37e7a3bbdf49695b35559ce2d95e51cae2c67a3b29ccb6bf9fd9ade878be9bdab32bdcf72264ad3cb4f955e1e5d575e7d5c1d1025d5c7bd0a7e2260b815fc70c00896a8c58df69962b434528b1111fdbcce95f53743f4150a6f25a27f566f8d6fe41b5d9ec5c2b7bbc7eee5f5101be9e48d68160b7fdefd3caaac2881adf5b084abef266d6bbd302d2dcac2076654950496aea88769714bebe761dfa169f192d5d5e6f9f3656cf141cbe6ccff753667ce209b9af9d3f22d24cab69028d94222d842a26c0b097f0b094b14b79098c316e21b2d1665e14fcfb849f994bb7ec5d8bdc84fb1dd9ef148deef2e81da795ff0501777ded2eab220ecea3839d5a42a06f36bff7f79a98af5d120d627fc582d1efe6265b417ce16101978d28fede9f75d6c4f06b13d5d8c8d850f56c666bec604a233863959e8eeebefeaeb59dbb5ae7beddafe4e744e74760d0fc7d6db6bd139d133d433d4b736eea07322beae67ddfa78ac0f9d13fdebbb867b7a7bbad039d1ddddd9630ff5f5a37322b6ae7f7d67774f0f3a27fad6c7bbe2bd3d316ccbe42def7770b99f2e00d1370d736958039668c0051ad0a5011b34e00a0dd8ae91c43d35e655389bffbc39cf6c8f3b4385c460323d9cf96d08d0145f63f4ac439b27f765efbd9101fc45006b017021c52734406c208d2b009c00c03a2561d45fa507e36ff8c29a14331ac04de631a414b713c3f219ba14fb88d1ed334252ec27c6669f614871801837fa8c1a29ee26469201df2746448ae304e9391fd22803580f2f30e03724d026c53d24f07b00a10ff336630f0336a81918b731a095426c338caf9386a237423b8f01bfa781155ac75b8c3b18401ef86218fb7dfd259a71b74f2f358c7b19f06d525956673cce8038d1cb5b4261be3cf0b9dcf5e9d25ae073b901e325069011be625528cc5728a5b7d4b844694589d20a03d0fe85010fd16861c8e0f38dc3be60338c077c7a8106ed7506fc9472b346ff2e034286cad62bc4e9aea1d56f04d0c181dab5f4d8c4019812c072007fc301f6bbba3f91f4cdacfe6df5dec681fa7934733379815ec301ce1ac818d91ea2a9d729827a42e5910f14c9794b05c0fe4c64981cf701984705f586eb1dcad2308034597a2bb0f4566049918dca1258d1d422d226536f9744339f6676f84e38f3bc349117f27d2b79d1959166e5c6a31f08e8051ddcdd17cc2853df08e0d3a46e96a89b25ea2eede2ac0f4cb590a906df94793e802f72f73fbe66e1931c6805708886d7d42f17f4a232845ab069001769406d9f062cda488f2be971036de9867fe5c03d001e23e5845aab0bc84042ad55c3a8008e00a0cae69364b47e19494ed2748bd1b6eea0d68683bddda4b287b608ea9a68fe06e6cadd4072f394f64ed29e1721eedf2a723591370d110c6d54032ea7ccfc3a1466ff6574727723b15fc158cfddfa66af684694bb35cdfedbc0c2551af05500ff4476f6b9fef691bf7ae56f1f1968a5c83ba9b8efac881f68784e03be0980b62abf774331b87b4b82eb3da8a9e0bea2825b7887747d7e8fe41e71c51f519aeb854f2a9f1f269f8f55fb6cd75d9f7424f0a3814f45b6d41e34da3a554ad728afffe87afda5ee7afd21491e73158e295d65e658d1eb2ef2fa6cb5d7cf8580ab18f01fa4f962e0d52507c8d58b94e3968683abdac207b5b63a05c250207e4e20165391ccbb81c6c74972f1048d2fa5f14b1bdc626f78d170bdbc4ca65f559b4f797935f0f2ea29bdfc417979bee8e584f2f244d1cb1f3d2f4b56d7023f0570924c8758f332f74dd3cd4a33e46ea82151e42fa0dcbc42b9a96555c969d24d602980ff2143cb5c43cb021a4d6d9e00278be79708b83496444de06b006a49a0c783d4e30ad21ef7e907025a41028dfb6781f49a095c08803e107c4b89478fde13f0b164599debbd99bceff0bcef2811f4e8db025a79a70f2cbfb1da3b106ef0befea3768c2e00c097b47a8f95b5d30927475f8126002b007ca8e4e2b006403b806e009796f03769c0831af0a8063cae05fcfb257de0015307eaf51285106087805b42402a14b0ffce00f619c08306f03503f8b6017cdf08e6bb4ec17fb40678a106f8430d2d13d0584bbb2d989f3e057f732d707f2df0cd5ae0c912fe474de06913f8ae097ccf049e37811f9ac1fcda3ae09e3ae0fe3ae0c13ae0d13ae0f1ba609eb6f57c00f4e56cf3c61ff6de09efddefbda300e8c8ff070ed059ff03eef26700d0f1bd5b00744ed3dd90cec47fd6003ada3e25013aafd6e8001d3b2feb009d247786003a107e6200b4652fac0568534913a0bd1031012af9cb4c802afba40950412ea903a8d69e2b890333c22f0c7b6828eb8cb346767e4393140de1454b1a1ad9f58b00de0bd130d07459ebc5e65f9bbc05107d186802b4851868120dc405e43cc5d3cf8368b8d4344d20d4c44cd39d34ea95051aac33819a76d6c8e6eb5d4d40ad04bdcc7610c73576592b50d7cedcc185ad97b55a26502fe132c84b836bc097206cc0bcde1289f032d6d5443237993b4da0d19b933437bf644e024d25c6a46fac5996282c2085eb17b9e20bcb7d93f0a25e15444b2f784b00a1b5a865028b4d95ca25e5be4819f84dd8bfaca79269e7a374dbc36b543b8cfda5a19e5eccad24d6717dcec9e63a3eee247376f2d6113b9de8d8928915469d743ed791c874e4b2b18e44323f52186a8f65463bc6d3f9d8889d4c772432178da7f31db14c3a9fb563f95c8793ce15471ded3837763385fc58210f8ca7f3a9e450fb088550febfc506c7b24e2c333a964c39ed31a8cb91e414acf81998aca95d39a28bfba6f74aa68b97607106c62039ed39f13bc6a4acdd79df27254bf8134af375c66463ed25b2fe2979a17e42f66c907f7548d6e8e297537b2f256b47a5d4c55f98c54589d25b8cc9fada4b6444372c123e3e6575bac2f5bad8cb2d19d13f221bf53532725876eb7b65eb21593f4c823f9b1a0eacdececbad7e9a93d59572a9fe02c9fe78ca9275e211d62ea57e48fc698ff2f40b10f36146060e950572272f06b2543f2a23963cbf5df61e25a51f4c59528a2758bb2efeb067836f468ac7d9510540e3909c8e1cf1212659ed178cc24963a7d1e43e4e9e349e223a71d2481809459f0c4823e10a196e5695253ab4c47320532b13b22621d963c1a4c2fa02986cae352e95ad47e462fd19d97a4036eb4fc985c64e5993d087f54e5d3c387d94d6f1e7b0789daf4ba79ff895327c89acbb4faed03f22eb0ec8ee6159a777c8151775c83a3d2ffbbca1de3ecb6457d9646ff96477e9a43830bd9756e9d7b08224d1b92bde50088c2b6563422e7b4c365e249775c8c6b7d5f3a5d33e8fd053dec7e7f9f6e8fc16fb98b277896c16d70cbf434a16ebe2df61c946d1db2e178bab0ec816d1d72ee51724bb2fc8137d0dc4bd6c369447e680b50a257d55c42301cac7a6e70cf30510ce47a60fc825e2c7a8424adf2971f4ec21a5ef9d381620bd6bee48ffa890de394d297dbd1a297d41c58b15486f3c1daeb2e7bfd173f57dba383cb5571c677be5d5c606292d7178aa5d2ed4c5cb6c58d6eb27e4c0b0bc6a83ac7f417c796a83b475719c5972b57e485e7d544a4bb6ea97eae2d0d4cbb4317ecd2c5d7c69ea65d9aa8bff6496649f08b2405f79f1ea7bc77ad7d45e7122c07a978bf5f54aac9f75b19ea8c6badfc5fa2661bdddc5fa5a0556ba89883b7871c59e9f7aa7155ba28bfddc92f375b578cf4eed9535e21bac5d1f268f27f61c958dba38c02db9f4986cfcbc5cde59b99274e7115f511e575272561e9352cf072714dd85c4574f3d4f7724f1f500f143ef88b858633f675463f74f1d904bc571451f56f42f985b6f4f0528e90626be13a070a359d1598e85ee67e24701964fce1dcb0c27ff33caff6d8a7e7bcf01b95cece5a5587042aff56e1bb97c3696b2d3096bdcc9e69299b4b5aebdb3bdd35a95cf16d2bbad9edeaef5fd9dabcfe9dda0e374d781f7e61803db766cdd7e7974f3d6c1abb7ee880e0e6c416cc4ce2297cf26d30944b76cd93e18880c6cc1e66bb6edd8ba6d4739735bf4eaade59ce8a68172c6b5d76fba72ebce72de8ead37551872d239d8f178d6c9e590559d3b39ec7626316ea70a8e9a81dffba39a7ddc861d4c60d2bb325de94c42f5eb8cda6363c97422f81b33529974c2528f423a974ca49db8954ce75148a6f3fdbdb08792394577f7f5c36b1adaee7cac90cc3ae4329e540d4b839b32991446730964d2a9c96b548b573a13779073ffac3d6667edd11c06fd1806edd14c219dc760c2ce61b39d4a5deb0a5475425ddfddd77f9d3dec5c5d4815e92dc9f1227d5d61a84847e3714acaf8d0f848d21e5b5b715944750795df3855d22ce5f748051d5215fd515e7794df1be57746057d517e3bd4adb468ca7cd27738b0c55b3c8cd8b9118f56098e051eca0688db791bf881503d5cec98b8105e7f1fab6aede3955d7dc26be8d38abd7cb2d8c6a7173bf84225cd7bc66c7d7b356ecb5e6dd0ad67ba8d7a75a53d7af515ed790da59d79f38a4d79e1623f5e63492bdefcf22ebca6d206bce6c6c1c18fdbb9d1c1989d4a0dc6f2996c6e4155c7dc42b364e1175574cfb55437cfb5aadeb9c57eebdc12bf736e6949e3dc3295f3e57edbdc8aa06bceaa689a5be9f5ccf98d7237b5f91d73e7050d73e7ab7eb90bfc76b923e15f7ef4ef23e4217271a41c6f644d249649e7f2763a1fb978d84ee59c3591647aac90cf452e0e74bcea8aac89e427c788e1edad0809c79d0927ee29efb9794dc43d2f49ffe6a2bcf2912d5046237bd614ed16a32c8391cf166645419b3c80e09e8ed50802eba5851ee879e7cb3b402f1a7159efec368875d8ab938a404bd6f08313ede9d44f13ae57a273aaadf710e71c70067a74bacd691967abdfd982a42d764e5672b733798640cf72b57a27d5b9593ffb3d1d24a780bba3722dce1adc392d4620ae2e4967b5ccaabe2773aab9927bc0d9aaa5a14c2675a695649fe5329a1bd2d315d569c07a27cb59dcd3e778936e9e05f0592b7cba369ed54256af0f4c6eb755a23d6b89f5f270f612eb5d043f20b9f56eab1f10b4c52bf5b9a9858933113ed35bc1e9232bff09155913b1d399f4e468a6907b0ff1d142be9b33c49370c6e9382b3150f2bbeedc209cf3663c15c2e087e5b90138d7ebae523b0dceb21fbde706ea3928e6534553f29bfddcc4e2f1af3cfdedaf4269ae17f753d652f0278e7313957d4677a42a7cea6f33efcbeb5ccc4ea5227b6e9ef95f000000ffff`

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
