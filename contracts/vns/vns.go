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

package vns

import (
	"strings"

	"github.com/vntchain/go-vnt/accounts/abi/bind"
	"github.com/vntchain/go-vnt/common"
	"github.com/vntchain/go-vnt/contracts/vns/contract"
	"github.com/vntchain/go-vnt/core/types"
	"github.com/vntchain/go-vnt/crypto"
)

var (
	MainNetAddress = common.HexToAddress("0x314159265dD8dbb310642f98f50C066173C1259b")
	TestNetAddress = common.HexToAddress("0x112234455c3a32fd11230c42e7bccd4a84e02010")
)

// swarm domain name registry and resolver
type VNS struct {
	*contract.VNSSession
	contractBackend bind.ContractBackend
}

// NewVNS creates a struct exposing convenient high-level operations for interacting with
// the VNT Name Service.
func NewVNS(transactOpts *bind.TransactOpts, contractAddr common.Address, contractBackend bind.ContractBackend) (*VNS, error) {
	vns, err := contract.NewVNS(contractAddr, contractBackend)
	if err != nil {
		return nil, err
	}

	return &VNS{
		&contract.VNSSession{
			Contract:     vns,
			TransactOpts: *transactOpts,
		},
		contractBackend,
	}, nil
}

// DeployVNS deploys an instance of the VNS nameservice, with a 'first-in, first-served' root registrar.
func DeployVNS(transactOpts *bind.TransactOpts, contractBackend bind.ContractBackend) (common.Address, *VNS, error) {
	// Deploy the VNS registry.
	vnsAddr, _, _, err := contract.DeployVNS(transactOpts, contractBackend)
	if err != nil {
		return vnsAddr, nil, err
	}
	vns, err := NewVNS(transactOpts, vnsAddr, contractBackend)
	if err != nil {
		return vnsAddr, nil, err
	}
	// Deploy the registrar.
	regAddr, _, _, err := contract.DeployFIFSRegistrar(transactOpts, contractBackend, vnsAddr, "")
	if err != nil {
		return vnsAddr, nil, err
	}
	// Set the registrar as owner of the VNS root.
	if _, err = vns.SetOwner("", regAddr); err != nil {
		return vnsAddr, nil, err
	}

	return vnsAddr, vns, nil
}

func vnsParentNode(name string) (string, string) {
	parts := strings.SplitN(name, ".", 2)
	label := crypto.Keccak256Hash([]byte(parts[0]))
	if len(parts) == 1 {
		return "", string(label.Bytes())
	} else {
		parentNode, parentLabel := vnsParentNode(parts[1])
		return string(crypto.Keccak256Hash([]byte(parentNode), []byte(parentLabel)).Bytes()), string(label.Bytes())
	}
}

func VnsNode(name string) string {
	parentNode, parentLabel := vnsParentNode(name)
	return string(crypto.Keccak256Hash([]byte(parentNode), []byte(parentLabel)).Bytes())
}

func (self *VNS) getResolver(node string) (*contract.PublicResolverSession, error) {
	resolverAddr, err := self.Resolver(node)
	if err != nil {
		return nil, err
	}

	resolver, err := contract.NewPublicResolver(resolverAddr, self.contractBackend)
	if err != nil {
		return nil, err
	}

	return &contract.PublicResolverSession{
		Contract:     resolver,
		TransactOpts: self.TransactOpts,
	}, nil
}

func (self *VNS) getRegistrar(node string) (*contract.FIFSRegistrarSession, error) {
	registrarAddr, err := self.Owner(node)
	if err != nil {
		return nil, err
	}

	registrar, err := contract.NewFIFSRegistrar(registrarAddr, self.contractBackend)
	if err != nil {
		return nil, err
	}

	return &contract.FIFSRegistrarSession{
		Contract:     registrar,
		TransactOpts: self.TransactOpts,
	}, nil
}

// Resolve is a non-transactional call that returns the content hash associated with a name.
func (self *VNS) Resolve(name string) (string, error) {
	node := VnsNode(name)

	resolver, err := self.getResolver(node)
	if err != nil {
		return "", err
	}

	ret, err := resolver.Content(node)
	if err != nil {
		return "", err
	}

	return ret, nil
}

// Addr is a non-transactional call that returns the address associated with a name.
func (self *VNS) Addr(name string) (common.Address, error) {
	node := VnsNode(name)

	resolver, err := self.getResolver(node)
	if err != nil {
		return common.Address{}, err
	}

	ret, err := resolver.Addr(node)
	if err != nil {
		return common.Address{}, err
	}

	return common.BytesToAddress(ret[:]), nil
}

// SetAddress sets the address associated with a name. Only works if the caller
// owns the name, and the associated resolver implements a `setAddress` function.
func (self *VNS) SetAddr(name string, addr common.Address) (*types.Transaction, error) {
	node := VnsNode(name)

	resolver, err := self.getResolver(node)
	if err != nil {
		return nil, err
	}

	opts := self.TransactOpts
	opts.GasLimit = 200000
	return resolver.Contract.SetAddr(&opts, node, addr)
}

// Register registers a new domain name for the caller, making them the owner of the new name.
// Only works if the registrar for the parent domain implements the FIFS registrar protocol.
func (self *VNS) Register(name string) (*types.Transaction, error) {
	parentNode, label := vnsParentNode(name)
	registrar, err := self.getRegistrar(parentNode)
	if err != nil {
		return nil, err
	}
	return registrar.Contract.Register(&self.TransactOpts, label, self.TransactOpts.From)
}

// SetContentHash sets the content hash associated with a name. Only works if the caller
// owns the name, and the associated resolver implements a `setContent` function.
func (self *VNS) SetContentHash(name string, hash string) (*types.Transaction, error) {
	node := VnsNode(name)

	resolver, err := self.getResolver(node)
	if err != nil {
		return nil, err
	}

	opts := self.TransactOpts
	opts.GasLimit = 200000
	return resolver.Contract.SetContent(&opts, node, hash)
}
