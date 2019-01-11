package main

import (
	"strings"

	"github.com/vntchain/go-vnt/accounts/abi"
)

type KVTree struct {
	Root map[string]*KVNode
}

type KVNode struct {
	Name        string
	StorageType abi.StorageType
	Type        string
	Children    map[string]*KVNode
}

func NewKVTree() *KVTree {
	return &KVTree{
		Root: make(map[string]*KVNode),
	}
}

func NewKVNode(name string, styp abi.StorageType, typ string) *KVNode {
	return &KVNode{
		Name:        name,
		StorageType: styp,
		Type:        typ,
		Children:    make(map[string]*KVNode),
	}
}

func (node *KVNode) AddNode(name string, styp abi.StorageType, typ string) {
	if _, ok := node.Children[name]; !ok {
		n := NewKVNode(name, styp, typ)
		node.Children[name] = n
	}
}

func (tree *KVTree) AddNode(name string, styp abi.StorageType, typ string, path string) {
	keys := strings.Split(path, ".")
	if len(keys) <= 1 {
		root := NewKVNode(name, styp, typ)
		tree.Root[name] = root
	} else {
		node := tree.Root[keys[0]]
		for i := 1; i < len(keys)-1; i++ {
			node = node.Children[keys[i]]
		}

		node.AddNode(name, styp, typ)
	}
}
