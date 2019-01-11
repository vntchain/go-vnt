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

// Contains all the wrappers from the accounts package to support client side vnode
// management on mobile platforms.

package gvnt

// import (
// 	"errors"
// )

// // Vnode represents a host on the network.
// type Vnode struct {
// 	node *discv5.Node
// }

// // NewEnode parses a node designator.
// //
// // There are two basic forms of node designators
// //   - incomplete nodes, which only have the public key (node ID)
// //   - complete nodes, which contain the public key and IP/Port information
// //
// // For incomplete nodes, the designator must look like one of these
// //
// //    vnode://<hex node id>
// //    <hex node id>
// //
// // For complete nodes, the node ID is encoded in the username portion
// // of the URL, separated from the host by an @ sign. The hostname can
// // only be given as an IP address, DNS domain names are not allowed.
// // The port in the host name section is the TCP listening port. If the
// // TCP and UDP (discovery) ports differ, the UDP port is specified as
// // query parameter "discport".
// //
// // In the following example, the node URL describes
// // a node with IP address 10.3.58.6, TCP listening port 30303
// // and UDP discovery port 30301.
// //
// //    vnode://<hex node id>@10.3.58.6:30303?discport=30301
// func NewEnode(rawurl string) (vnode *Vnode, _ error) {
// 	node, err := discv5.ParseNode(rawurl)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Vnode{node}, nil
// }

// // Enodes represents a slice of accounts.
// type Enodes struct{ nodes []*discv5.Node }

// // NewEnodes creates a slice of uninitialized enodes.
// func NewEnodes(size int) *Enodes {
// 	return &Enodes{
// 		nodes: make([]*discv5.Node, size),
// 	}
// }

// // NewEnodesEmpty creates an empty slice of Vnode values.
// func NewEnodesEmpty() *Enodes {
// 	return NewEnodes(0)
// }

// // Size returns the number of enodes in the slice.
// func (e *Enodes) Size() int {
// 	return len(e.nodes)
// }

// // Get returns the vnode at the given index from the slice.
// func (e *Enodes) Get(index int) (vnode *Vnode, _ error) {
// 	if index < 0 || index >= len(e.nodes) {
// 		return nil, errors.New("index out of bounds")
// 	}
// 	return &Vnode{e.nodes[index]}, nil
// }

// // Set sets the vnode at the given index in the slice.
// func (e *Enodes) Set(index int, vnode *Vnode) error {
// 	if index < 0 || index >= len(e.nodes) {
// 		return errors.New("index out of bounds")
// 	}
// 	e.nodes[index] = vnode.node
// 	return nil
// }

// // Append adds a new vnode element to the end of the slice.
// func (e *Enodes) Append(vnode *Vnode) {
// 	e.nodes = append(e.nodes, vnode.node)
// }
