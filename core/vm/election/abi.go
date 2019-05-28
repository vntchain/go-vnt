// Copyright 2019 The go-vnt Authors
// This file is part of the go-vnt library.
//
// The go-vnt library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-vnt library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-vnt library. If not, see <http://www.gnu.org/licenses/>.

package election

import (
	"strings"

	"github.com/vntchain/go-vnt/accounts/abi"
)

const ElectionAbiJSON = `[
{"name":"registerWitness","inputs":[{"name":"nodeUrl","type":"bytes"},{"name":"website","type":"bytes"},{"name":"nodeName","type":"bytes"}],"outputs":[],"type":"function"},
{"name":"unregisterWitness","inputs":[],"outputs":[],"type":"function"},
{"name":"voteWitnesses","inputs":[{"name":"candidate","type":"address[]"}],"outputs":[],"type":"function"},
{"name":"cancelVote","inputs":[],"outputs":[],"type":"function"},
{"name":"startProxy","inputs":[],"outputs":[],"type":"function"},
{"name":"stopProxy","inputs":[],"outputs":[],"type":"function"},
{"name":"cancelProxy","inputs":[],"outputs":[],"type":"function"},
{"name":"setProxy","inputs":[{"name":"proxy","type":"address"}],"outputs":[],"type":"function"},
{"name":"stake","inputs":[{"name":"stakeCount","type":"uint256"}],"outputs":[],"type":"function"},
{"name":"unStake","inputs":[],"outputs":[],"type":"function"},
{"name":"extractOwnBounty","inputs":[],"outputs":[],"type":"function"}
]`

func GetElectionABI() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(ElectionAbiJSON))
}

func PackInput(abiobj abi.ABI, name string, args ...interface{}) ([]byte, error) {
	abires := abiobj
	var res []byte
	var err error
	if len(args) == 0 {
		res, err = abires.Pack(name)
	} else {
		res, err = abires.Pack(name, args...)
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}
