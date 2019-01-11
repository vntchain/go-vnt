package dpos

import (
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/consensus"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/rpc"
	"math/big"
)

// API is a user facing RPC API to allow controlling the signer and voting
// mechanisms of the proof-of-authority scheme.
type API struct {
	chain consensus.ChainReader
	dpos  *Dpos
}

// GetSigners retrieves the list of authorized signers at the specified block.
func (api *API) GetSigners(number *rpc.BlockNumber) ([]common.Address, error) {
	// Retrieve the requested block number (or current if none requested)
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	// Ensure we have an actually valid block and return the signers from its snapshot
	if header == nil {
		return nil, errUnknownBlock
	}

	return header.Witnesses, nil
}

// GetSignersAtHash retrieves the state snapshot at a given block.
func (api *API) GetSignersAtHash(hash common.Hash) ([]common.Address, error) {
	header := api.chain.GetHeaderByHash(hash)
	if header == nil {
		return nil, errUnknownBlock
	}

	return header.Witnesses, nil
}

func (api *API) GetAllMessage() []types.ConsensusMsg {
	msgs := api.dpos.bft.roundMp.getAllMsgOf(api.dpos.bft.h, api.dpos.bft.r)
	return msgs
}

func (api *API) GetPrePrepareMsg() *types.PreprepareMsg {
	msg, err := api.dpos.bft.roundMp.getPrePrepareMsg(api.dpos.bft.h, api.dpos.bft.r)
	if err != nil {
		return nil
	}
	return msg
}

func (api *API) GetPrepareMsgs() []*types.PrepareMsg {
	msgs := api.dpos.bft.roundMp.getAllMsgOf(api.dpos.bft.h, api.dpos.bft.r)
	var result []*types.PrepareMsg
	for _, msg := range msgs {
		if prepare, ok := msg.(*types.PrepareMsg); ok {
			result = append(result, prepare)
		}
	}
	return result
}

func (api *API) GetCommitMsgs() []*types.CommitMsg {
	msgs := api.dpos.bft.roundMp.getAllMsgOf(api.dpos.bft.h, api.dpos.bft.r)
	var result []*types.CommitMsg
	for _, msg := range msgs {
		if commit, ok := msg.(*types.CommitMsg); ok {
			result = append(result, commit)
		}
	}
	return result
}

func (api *API) GetCurrentStep() uint32 {
	return api.dpos.bft.step
}

func (api *API) GetCurrentHeight() *big.Int {
	return api.dpos.bft.h
}

func (api *API) GetCurrentRound() uint32 {
	return api.dpos.bft.r
}
