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
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"

	"github.com/vntchain/go-vnt/common"
	inter "github.com/vntchain/go-vnt/core/vm/interface"
	"github.com/vntchain/go-vnt/log"
	"github.com/vntchain/go-vnt/rlp"
)

const (
	VOTERPREFIX        = byte(0)
	CANDIDATEPREFIX    = byte(1)
	STAKEPREFIX        = byte(2)
	BOUNTYPREFIX       = byte(3)
	MAINNETVOTESPREFIX = byte(4)
	PREFIXLENGTH       = 4 // key的结构为，4位表前缀，20位address，8位的value在struct中的位置
)

type getFuncType func(key common.Hash) common.Hash
type setFuncType func(key common.Hash, value common.Hash)

func (ec electionContext) getVoter(addr common.Address) Voter {
	return getVoterFrom(addr, ec.getFromDB)
}

func (ec electionContext) getCandidate(key common.Address) Candidate {
	// var candidate Candidate
	candidate := newCandidate()
	var err error
	if err := convertToStruct(CANDIDATEPREFIX, key, &candidate, ec.getFromDB); err == nil {
		return candidate
	}

	log.Debug("Get Candidate From DB ", "addr", key.String(), "err", err)
	return newCandidate()
}

func (ec electionContext) getStake(addr common.Address) Stake {
	return getStakeFrom(addr, ec.getFromDB)
}

func (ec electionContext) setVoter(voter Voter) error {
	err := convertToKV(VOTERPREFIX, voter, ec.setToDB)
	if err != nil {
		log.Error("setVoter error", "err", err, "voter", voter)
	}
	return err
}

func (ec electionContext) setCandidate(candidate Candidate) error {
	err := convertToKV(CANDIDATEPREFIX, candidate, ec.setToDB)
	if err != nil {
		log.Error("setCandidate error", "err", err, "candidate", candidate)
	}
	return err
}

func (ec electionContext) setStake(stake Stake) error {
	err := convertToKV(STAKEPREFIX, stake, ec.setToDB)
	if err != nil {
		log.Error("setCandidate error", "err", err, "stake", stake)
	}
	return err
}

func (ec electionContext) setToDB(key common.Hash, value common.Hash) {
	ec.context.GetStateDb().SetState(contractAddr, key, value)
}

func (ec electionContext) getFromDB(key common.Hash) common.Hash {
	return ec.context.GetStateDb().GetState(contractAddr, key)
}

// getVoterFrom get a voter's information from a specific stateDB
func getVoterFrom(addr common.Address, getFromDB getFuncType) Voter {
	var voter Voter
	var err error
	if err := convertToStruct(VOTERPREFIX, addr, &voter, getFromDB); err == nil {
		return voter
	}

	log.Debug("Get Voter From DB ", "addr", addr.String(), "err", err)
	return newVoter()
}

// getStakeFrom get a user's information from a specific stateDB
func getStakeFrom(addr common.Address, getFromDB getFuncType) Stake {
	var stake Stake
	var err error
	if err := convertToStruct(STAKEPREFIX, addr, &stake, getFromDB); err == nil {
		return stake
	}

	log.Debug("Get Stake From DB ", "addr", addr.String(), "err", err)
	return Stake{}
}

func convertToKV(prefix byte, v interface{}, setToDB setFuncType) error {
	var key common.Hash
	key[0] = prefix

	value := reflect.ValueOf(v)
	if value.Kind() == reflect.Ptr {
		value = reflect.ValueOf(v).Elem()
	}
	if value.Kind() != reflect.Struct {
		return fmt.Errorf("error : v %v must be struct", v)
	}
	if !value.IsValid() {
		return fmt.Errorf("error: value %v is not valid", v)
	}

	// 结构体的owner作为key
	owner := value.FieldByName("Owner")
	if owner.IsValid() && owner.CanInterface() {
		if k, ok := owner.Interface().(common.Address); ok {
			copy(key[PREFIXLENGTH:], k.Bytes())
		} else {
			return fmt.Errorf("error: owner %v is not address", owner)
		}
	} else {
		copy(key[PREFIXLENGTH:], contractAddr.Bytes())
	}

	// 结构体中的每个元素都要分别存储
	for i := 0; i < value.NumField(); i++ {
		// 根据字段在结构体中的位置，对key进行相应的操作
		binary.BigEndian.PutUint64(key[PREFIXLENGTH+common.AddressLength:], uint64(i))
		fv := value.Field(i)
		isArray := false

		// 若元素为数组，数组中的每个元素也需要分别存储
		if fv.Kind() == reflect.Array || fv.Kind() == reflect.Slice {
			isArray = true
			for j := 0; j < fv.Len(); j++ {
				var subKey common.Hash
				copy(subKey[:], key[:])
				subv := fv.Index(j)
				binary.BigEndian.PutUint32(subKey[PREFIXLENGTH+common.AddressLength:], uint32(j+1))
				if !subv.IsValid() || (subv.Kind() != reflect.Struct && subv.Kind() != reflect.Ptr && subv.Kind() != reflect.Array) {
					isArray = false
					break
				}
				elem, err := rlp.EncodeToBytes(subv.Interface())
				if err != nil {
					return err
				}
				setToDB(subKey, common.BytesToHash(elem))
			}
		}
		// 如果是数组，则数组开始的key，存储数组的长度
		if isArray {
			elem, err := rlp.EncodeToBytes(uint32(fv.Len()))
			if err != nil {
				return err
			}
			setToDB(key, common.BytesToHash(elem))
			continue
		}

		if !fv.IsValid() || !fv.CanInterface() {
			return fmt.Errorf("error: %v is not valid", fv)
		}
		// 普通元素存储rlp
		elem, err := rlp.EncodeToBytes(fv.Interface())
		if err != nil {
			return err
		}

		// 如果要存储的字节过长，就拆分了存
		// 0号位置存储切分的长度，后面按右对齐方式存储，若需要补空位，补在第一个元素处
		valLen := len(elem)/32 + 1
		var j int
		for j = valLen - 1; j >= 0; j-- {
			var subKey common.Hash
			copy(subKey[:], key[:])
			binary.BigEndian.PutUint32(subKey[PREFIXLENGTH+common.AddressLength:], uint32(j))
			cutPos := len(elem) - 32
			if cutPos < 0 {
				setToDB(subKey, common.BytesToHash(elem))
				break
			}
			tmpElem := elem[cutPos:]
			elem = elem[:cutPos]
			setToDB(subKey, common.BytesToHash(tmpElem))
		}
	}
	return nil
}

func convertToStruct(prefix byte, addr common.Address, v interface{}, getFn getFuncType) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("error : v %v must be ptr", v)
	}
	value = value.Elem()

	var key common.Hash
	key[0] = prefix
	copy(key[PREFIXLENGTH:], addr.Bytes())
	// 结构体中的每个元素都要分别获取
	for i := 0; i < value.NumField(); i++ {
		// 根据字段在结构体中的位置，对key进行相应的操作
		binary.BigEndian.PutUint64(key[PREFIXLENGTH+common.AddressLength:], uint64(i))
		fv := value.Field(i)

		if !fv.IsValid() || !fv.CanInterface() {
			return fmt.Errorf("error: %v is not valid", fv)
		}

		// 从数据库中得到对应的数据
		valByte := getFn(key)

		// 按照数据类型对数据进行解析后，赋值给struct
		if _, ok := fv.Interface().(common.Address); ok {
			var tmp common.Address
			if err := rlp.DecodeBytes(valByte.Big().Bytes(), &tmp); err == nil {
				value.Field(i).Set(reflect.ValueOf(tmp))
			} else {
				return fmt.Errorf("decode to common.Address error: %v", err)
			}
		} else if _, ok = fv.Interface().(bool); ok {
			var tmp bool
			if err := rlp.DecodeBytes(valByte.Big().Bytes(), &tmp); err == nil {
				value.Field(i).Set(reflect.ValueOf(tmp))
			} else {
				return err
			}
		} else if _, ok = fv.Interface().(uint64); ok {
			var tmp uint64
			if err := rlp.DecodeBytes(valByte.Big().Bytes(), &tmp); err == nil {
				value.Field(i).Set(reflect.ValueOf(tmp))
			} else {
				return err
			}
		} else if _, ok = fv.Interface().(*big.Int); ok {
			var tmp *big.Int
			if err := rlp.DecodeBytes(valByte.Big().Bytes(), &tmp); err == nil {
				value.Field(i).Set(reflect.ValueOf(tmp))
			} else {
				return err
			}
		} else if _, ok = fv.Interface().([]common.Address); ok {
			var tmp []common.Address
			var valLen uint32

			// 如果是数组，先解析出数组长度，然后获取数组中的元素
			if err := rlp.DecodeBytes(valByte.Big().Bytes(), &valLen); err == nil {
				for j := uint32(0); j < valLen; j++ {
					var tmpArray common.Address
					binary.BigEndian.PutUint32(key[PREFIXLENGTH+common.AddressLength:], uint32(j+1))
					arrayByte := getFn(key)
					if err = rlp.DecodeBytes(arrayByte.Big().Bytes(), &tmpArray); err == nil {
						tmp = append(tmp, tmpArray)
					} else {
						return err
					}
				}
				value.Field(i).Set(reflect.ValueOf(tmp))
			} else {
				return err
			}
		} else if _, ok := fv.Interface().([]byte); ok {
			// 部分byte数组过长，是拆分了之后存储的
			var val []byte
			err := rlp.DecodeBytes(valByte.Big().Bytes(), &val)
			if err == nil {
				value.Field(i).Set(reflect.ValueOf(val))
			} else {
				val = valByte.Big().Bytes()
				var tmp []byte
				for j := 1; ; j++ {
					binary.BigEndian.PutUint32(key[PREFIXLENGTH+common.AddressLength:], uint32(j))
					arrayByte := getFn(key)
					if arrayByte.Big().Sign() == 0 {
						break
					}
					val = append(val, arrayByte.Bytes()...)
					if err = rlp.DecodeBytes(val, &tmp); err == nil {
						value.Field(i).Set(reflect.ValueOf(tmp))
						break
					}
				}
			}
		}

	}
	return nil
}

func getAllCandidate(db inter.StateDB) CandidateList {
	var result CandidateList
	addrs := make(map[common.Address]struct{})
	// 从数据库的value中找到所有的address
	db.ForEachStorage(contractAddr, func(key common.Hash, value common.Hash) bool {
		_, content, _, err := rlp.Split(value.Big().Bytes())
		if err != nil {
			// 这个地方长的bytes做过处理这里split会出错，所以这个错改成debug打印日志
			log.Debug("rlp split error", "err", err)
			return true
		}
		var addr common.Address
		if len(content) == common.AddressLength {
			addr = common.BytesToAddress(content)
		} else if len(content) == common.AddressLength+1 {
			if err := rlp.DecodeBytes(content, &addr); err != nil {
				return true
			}
		}
		if !bytes.Equal(addr.Bytes(), emptyAddress.Bytes()) {
			addrs[addr] = struct{}{}
		}
		return true
	})

	// 用这些address尝试去数据库中找候选者，当没有这个地址的候选者时会报错
	// 有可能并不是见证人所以报错
	for addr := range addrs {
		// var candidate Candidate
		candidate := newCandidate()
		err := convertToStruct(CANDIDATEPREFIX, addr, &candidate, genGetFunc(db))
		if err != nil {
			log.Debug("getAllCandidate maybe error", "address", addr, "err", err)
			continue
		}
		result = append(result, candidate)
	}

	return result
}

func getAllProxy(db inter.StateDB) []*Voter {
	var result []*Voter
	addrs := make(map[common.Address]struct{})

	db.ForEachStorage(contractAddr, func(key common.Hash, value common.Hash) bool {
		if key[0] == VOTERPREFIX {
			var addr common.Address
			copy(addr[:], key[PREFIXLENGTH:PREFIXLENGTH+common.AddressLength])
			addrs[addr] = struct{}{}
		}
		return true
	})

	for addr := range addrs {
		var voter Voter
		err := convertToStruct(VOTERPREFIX, addr, &voter, genGetFunc(db))
		if err != nil {
			log.Error("getAllProxy error", "address", addr, "err", err)
		}

		if voter.IsProxy {
			result = append(result, &voter)
		}
	}
	return result
}
func addCandidateBounty(stateDB inter.StateDB, addr common.Address, bouns *big.Int) error {
	candidate := newCandidate()
	err := convertToStruct(CANDIDATEPREFIX, addr, &candidate, genGetFunc(stateDB))
	if err != nil {
		return err
	}

	candidate.TotalBounty = new(big.Int).Add(candidate.TotalBounty, bouns)
	err = convertToKV(CANDIDATEPREFIX, &candidate, genSetFunc(stateDB))
	if err != nil {
		return err
	}
	return nil
}

func getRestBounty(stateDB inter.StateDB) Bounty {
	var bounty Bounty
	err := convertToStruct(BOUNTYPREFIX, contractAddr, &bounty, genGetFunc(stateDB))
	if err != nil {
		return Bounty{big.NewInt(0)}
	}
	return bounty
}

func setRestBounty(stateDB inter.StateDB, restBounty Bounty) error {
	return convertToKV(BOUNTYPREFIX, restBounty, genSetFunc(stateDB))
}

// getMainNetVotes return a initialized main net votes information,
// when no main net votes information in state db.
func getMainNetVotes(stateDb inter.StateDB) MainNetVotes {
	var mv MainNetVotes
	err := convertToStruct(MAINNETVOTESPREFIX, contractAddr, &mv, genGetFunc(stateDb))
	if err != nil {
		return MainNetVotes{big.NewInt(0), false}
	}

	return mv
}

// setMainNetVotes store the main net votes to db.
func setMainNetVotes(stateDB inter.StateDB, mv MainNetVotes) error {
	return convertToKV(MAINNETVOTESPREFIX, mv, genSetFunc(stateDB))
}

// genGetFunc generate universal get function for read from state db.
func genGetFunc(stateDb inter.StateDB) getFuncType {
	return func(key common.Hash) common.Hash {
		return stateDb.GetState(contractAddr, key)
	}
}

// genSetFunc generate universal get function for write state to state db.
func genSetFunc(stateDb inter.StateDB) setFuncType {
	return func(key common.Hash, value common.Hash) {
		stateDb.SetState(contractAddr, key, value)
	}
}
