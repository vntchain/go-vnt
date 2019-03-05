// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// package vntext contains gvnt specific vnt.js extensions.
package vntjsext

var Modules = map[string]string{
	"admin":      Admin_JS,
	"chequebook": Chequebook_JS,
	"debug":      Debug_JS,
	"core":       Core_JS,
	"bp":         Bp_JS,
	"net":        Net_JS,
	"personal":   Personal_JS,
	"rpc":        RPC_JS,
	"shh":        Shh_JS,
	"swarmfs":    SWARMFS_JS,
	"txpool":     TxPool_JS,
	"dpos":       Dpos_JS,
}

const Chequebook_JS = `
vnt._extend({
	property: 'chequebook',
	methods: [
		new vnt._extend.Method({
			name: 'deposit',
			call: 'chequebook_deposit',
			params: 1,
			inputFormatter: [null]
		}),
		new vnt._extend.Property({
			name: 'balance',
			getter: 'chequebook_balance',
			outputFormatter: vnt._extend.utils.toDecimal
		}),
		new vnt._extend.Method({
			name: 'cash',
			call: 'chequebook_cash',
			params: 1,
			inputFormatter: [null]
		}),
		new vnt._extend.Method({
			name: 'issue',
			call: 'chequebook_issue',
			params: 2,
			inputFormatter: [null, null]
		}),
	]
});
`

const Admin_JS = `
vnt._extend({
	property: 'admin',
	methods: [
		new vnt._extend.Method({
			name: 'addPeer',
			call: 'admin_addPeer',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'removePeer',
			call: 'admin_removePeer',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'exportChain',
			call: 'admin_exportChain',
			params: 1,
			inputFormatter: [null]
		}),
		new vnt._extend.Method({
			name: 'importChain',
			call: 'admin_importChain',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'sleepBlocks',
			call: 'admin_sleepBlocks',
			params: 2
		}),
		new vnt._extend.Method({
			name: 'startRPC',
			call: 'admin_startRPC',
			params: 4,
			inputFormatter: [null, null, null, null]
		}),
		new vnt._extend.Method({
			name: 'stopRPC',
			call: 'admin_stopRPC'
		}),
		new vnt._extend.Method({
			name: 'startWS',
			call: 'admin_startWS',
			params: 4,
			inputFormatter: [null, null, null, null]
		}),
		new vnt._extend.Method({
			name: 'stopWS',
			call: 'admin_stopWS'
		}),
	],
	properties: [
		new vnt._extend.Property({
			name: 'nodeInfo',
			getter: 'admin_nodeInfo'
		}),
		new vnt._extend.Property({
			name: 'peers',
			getter: 'admin_peers'
		}),
		new vnt._extend.Property({
			name: 'datadir',
			getter: 'admin_datadir'
		}),
	]
});
`

const Debug_JS = `
vnt._extend({
	property: 'debug',
	methods: [
		new vnt._extend.Method({
			name: 'printBlock',
			call: 'debug_printBlock',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'getBlockRlp',
			call: 'debug_getBlockRlp',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'setHead',
			call: 'debug_setHead',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'seedHash',
			call: 'debug_seedHash',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'dumpBlock',
			call: 'debug_dumpBlock',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'chaindbProperty',
			call: 'debug_chaindbProperty',
			params: 1,
			outputFormatter: console.log
		}),
		new vnt._extend.Method({
			name: 'chaindbCompact',
			call: 'debug_chaindbCompact',
		}),
		new vnt._extend.Method({
			name: 'metrics',
			call: 'debug_metrics',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'verbosity',
			call: 'debug_verbosity',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'vmodule',
			call: 'debug_vmodule',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'backtraceAt',
			call: 'debug_backtraceAt',
			params: 1,
		}),
		new vnt._extend.Method({
			name: 'stacks',
			call: 'debug_stacks',
			params: 0,
			outputFormatter: console.log
		}),
		new vnt._extend.Method({
			name: 'freeOSMemory',
			call: 'debug_freeOSMemory',
			params: 0,
		}),
		new vnt._extend.Method({
			name: 'setGCPercent',
			call: 'debug_setGCPercent',
			params: 1,
		}),
		new vnt._extend.Method({
			name: 'memStats',
			call: 'debug_memStats',
			params: 0,
		}),
		new vnt._extend.Method({
			name: 'gcStats',
			call: 'debug_gcStats',
			params: 0,
		}),
		new vnt._extend.Method({
			name: 'cpuProfile',
			call: 'debug_cpuProfile',
			params: 2
		}),
		new vnt._extend.Method({
			name: 'startCPUProfile',
			call: 'debug_startCPUProfile',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'stopCPUProfile',
			call: 'debug_stopCPUProfile',
			params: 0
		}),
		new vnt._extend.Method({
			name: 'goTrace',
			call: 'debug_goTrace',
			params: 2
		}),
		new vnt._extend.Method({
			name: 'startGoTrace',
			call: 'debug_startGoTrace',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'stopGoTrace',
			call: 'debug_stopGoTrace',
			params: 0
		}),
		new vnt._extend.Method({
			name: 'blockProfile',
			call: 'debug_blockProfile',
			params: 2
		}),
		new vnt._extend.Method({
			name: 'setBlockProfileRate',
			call: 'debug_setBlockProfileRate',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'writeBlockProfile',
			call: 'debug_writeBlockProfile',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'mutexProfile',
			call: 'debug_mutexProfile',
			params: 2
		}),
		new vnt._extend.Method({
			name: 'setMutexProfileFraction',
			call: 'debug_setMutexProfileFraction',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'writeMutexProfile',
			call: 'debug_writeMutexProfile',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'writeMemProfile',
			call: 'debug_writeMemProfile',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'traceBlock',
			call: 'debug_traceBlock',
			params: 2,
			inputFormatter: [null, null]
		}),
		new vnt._extend.Method({
			name: 'traceBlockFromFile',
			call: 'debug_traceBlockFromFile',
			params: 2,
			inputFormatter: [null, null]
		}),
		new vnt._extend.Method({
			name: 'traceBlockByNumber',
			call: 'debug_traceBlockByNumber',
			params: 2,
			inputFormatter: [null, null]
		}),
		new vnt._extend.Method({
			name: 'traceBlockByHash',
			call: 'debug_traceBlockByHash',
			params: 2,
			inputFormatter: [null, null]
		}),
		new vnt._extend.Method({
			name: 'traceTransaction',
			call: 'debug_traceTransaction',
			params: 2,
			inputFormatter: [null, null]
		}),
		new vnt._extend.Method({
			name: 'preimage',
			call: 'debug_preimage',
			params: 1,
			inputFormatter: [null]
		}),
		new vnt._extend.Method({
			name: 'getBadBlocks',
			call: 'debug_getBadBlocks',
			params: 0,
		}),
		new vnt._extend.Method({
			name: 'storageRangeAt',
			call: 'debug_storageRangeAt',
			params: 5,
		}),
		new vnt._extend.Method({
			name: 'getModifiedAccountsByNumber',
			call: 'debug_getModifiedAccountsByNumber',
			params: 2,
			inputFormatter: [null, null],
		}),
		new vnt._extend.Method({
			name: 'getModifiedAccountsByHash',
			call: 'debug_getModifiedAccountsByHash',
			params: 2,
			inputFormatter:[null, null],
		}),
	],
	properties: []
});
`

const Core_JS = `
vnt._extend({
	property: 'core',
	methods: [
		new vnt._extend.Method({
			name: 'resend',
			call: 'core_resend',
			params: 3,
			inputFormatter: [vnt._extend.formatters.inputTransactionFormatter, vnt._extend.utils.fromDecimal, vnt._extend.utils.fromDecimal]
		}),
		new vnt._extend.Method({
			name: 'getRawTransaction',
			call: 'core_getRawTransactionByHash',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'getRawTransactionFromBlock',
			call: function(args) {
				return (vnt._extend.utils.isString(args[0]) && args[0].indexOf('0x') === 0) ? 'core_getRawTransactionByBlockHashAndIndex' : 'core_getRawTransactionByBlockNumberAndIndex';
			},
			params: 2,
			inputFormatter: [vnt._extend.formatters.inputBlockNumberFormatter, vnt._extend.utils.toHex]
		}),
	],
	properties: [
		new vnt._extend.Property({
			name: 'pendingTransactions',
			getter: 'core_pendingTransactions',
			outputFormatter: function(txs) {
				var formatted = [];
				for (var i = 0; i < txs.length; i++) {
					formatted.push(vnt._extend.formatters.outputTransactionFormatter(txs[i]));
					formatted[i].blockHash = null;
				}
				return formatted;
			}
		}),
	]
});
`

const Bp_JS = `
vnt._extend({
	property: 'bp',
	methods: [
		new vnt._extend.Method({
			name: 'start',
			call: 'bp_start',
			params: 1,
			inputFormatter: [null]
		}),
		new vnt._extend.Method({
			name: 'stop',
			call: 'bp_stop'
		}),
		new vnt._extend.Method({
			name: 'setCoinbase',
			call: 'bp_setCoinbase',
			params: 1,
			inputFormatter: [vnt._extend.formatters.inputAddressFormatter]
		}),
		new vnt._extend.Method({
			name: 'setExtra',
			call: 'bp_setExtra',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'setGasPrice',
			call: 'bp_setGasPrice',
			params: 1,
			inputFormatter: [vnt._extend.utils.fromDecimal]
		}),
	],
	properties: []
});
`

const Net_JS = `
vnt._extend({
	property: 'net',
	methods: [],
	properties: [
		new vnt._extend.Property({
			name: 'version',
			getter: 'net_version'
		}),
	]
});
`

const Personal_JS = `
vnt._extend({
	property: 'personal',
	methods: [
		new vnt._extend.Method({
			name: 'importRawKey',
			call: 'personal_importRawKey',
			params: 2
		}),
		new vnt._extend.Method({
			name: 'sign',
			call: 'personal_sign',
			params: 3,
			inputFormatter: [null, vnt._extend.formatters.inputAddressFormatter, null]
		}),
		new vnt._extend.Method({
			name: 'ecRecover',
			call: 'personal_ecRecover',
			params: 2
		}),
		new vnt._extend.Method({
			name: 'deriveAccount',
			call: 'personal_deriveAccount',
			params: 3
		}),
		new vnt._extend.Method({
			name: 'signTransaction',
			call: 'personal_signTransaction',
			params: 2,
			inputFormatter: [vnt._extend.formatters.inputTransactionFormatter, null]
		}),
	],
	properties: [
		new vnt._extend.Property({
			name: 'listWallets',
			getter: 'personal_listWallets'
		}),
	]
})
`

const RPC_JS = `
vnt._extend({
	property: 'rpc',
	methods: [],
	properties: [
		new vnt._extend.Property({
			name: 'modules',
			getter: 'rpc_modules'
		}),
	]
});
`

const Shh_JS = `
vnt._extend({
	property: 'shh',
	methods: [
	],
	properties:
	[
		new vnt._extend.Property({
			name: 'version',
			getter: 'shh_version',
			outputFormatter: vnt._extend.utils.toDecimal
		}),
		new vnt._extend.Property({
			name: 'info',
			getter: 'shh_info'
		}),
	]
});
`

const SWARMFS_JS = `
vnt._extend({
	property: 'swarmfs',
	methods:
	[
		new vnt._extend.Method({
			name: 'mount',
			call: 'swarmfs_mount',
			params: 2
		}),
		new vnt._extend.Method({
			name: 'unmount',
			call: 'swarmfs_unmount',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'listmounts',
			call: 'swarmfs_listmounts',
			params: 0
		}),
	]
});
`

const TxPool_JS = `
vnt._extend({
	property: 'txpool',
	methods: [],
	properties:
	[
		new vnt._extend.Property({
			name: 'content',
			getter: 'txpool_content'
		}),
		new vnt._extend.Property({
			name: 'inspect',
			getter: 'txpool_inspect'
		}),
		new vnt._extend.Property({
			name: 'status',
			getter: 'txpool_status',
			outputFormatter: function(status) {
				status.pending = vnt._extend.utils.toDecimal(status.pending);
				status.queued = vnt._extend.utils.toDecimal(status.queued);
				return status;
			}
		}),
	]
});
`
const Dpos_JS = `
vnt._extend({
	property: 'dpos',
	methods: [
		new vnt._extend.Method({
			name: 'getSigners',
			call: 'dpos_getSigners',
			params: 1,
			inputFormatter: [null]
		}),
		new vnt._extend.Method({
			name: 'getSignersAtHash',
			call: 'dpos_getSignersAtHash',
			params: 1
		}),
		new vnt._extend.Method({
			name: 'getPrePrepareMsg',
			call: 'dpos_getPrePrepareMsg',
		}),
		new vnt._extend.Method({
			name: 'getPrepareMsgs',
			call: 'dpos_getPrepareMsgs',
		}),
		new vnt._extend.Method({
			name: 'getCommitMsgs',
			call: 'dpos_getCommitMsgs',
		}),
		new vnt._extend.Method({
			name: 'getAllMessage',
			call: 'dpos_getAllMessage',
		}),
		new vnt._extend.Property({
			name: 'step',
			getter: 'dpos_getCurrentStep',
		}),
		new vnt._extend.Property({
			name: 'height',
			getter: 'dpos_getCurrentHeight',
		}),
		new vnt._extend.Property({
			name: 'round',
			getter: 'dpos_getCurrentRound',
		}),
		
	]
});
`
