// Copyright 2016 The go-ethereum Authors
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

package vntclient

import vnt "github.com/vntchain/go-vnt"

// Verify that Client implements the vnt interfaces.
var (
	_ = vnt.ChainReader(&Client{})
	_ = vnt.TransactionReader(&Client{})
	_ = vnt.ChainStateReader(&Client{})
	_ = vnt.ChainSyncReader(&Client{})
	_ = vnt.ContractCaller(&Client{})
	_ = vnt.GasEstimator(&Client{})
	_ = vnt.GasPricer(&Client{})
	_ = vnt.LogFilterer(&Client{})
	_ = vnt.PendingStateReader(&Client{})
	// _ = vnt.PendingStateEventer(&Client{})
	_ = vnt.PendingContractCaller(&Client{})
)
