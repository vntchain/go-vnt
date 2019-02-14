package abi

import (
	"fmt"
	"strings"
)

type StorageType uint64

const (
	NormalTy StorageType = iota
	MappingKeyTy
	ArrayIndexTy
	MappingValueTy
	ArrayValueTy
	StructValueTy
	LengthTy
)

const (
	TY_INT32 int32 = 1 + iota
	TY_INT64
	TY_UINT32
	TY_UINT64
	TY_UINT256
	TY_STRING
	TY_ADDRESS
	TY_BOOL
	TY_POINTER
)

type Key struct {
	Name      string           `json:"name"`
	Tables    Tables           `json:"tables"`
	Types     map[string]Type  `json:"-"`
	Keys      map[string]Table `json:"keys"`
	Connector map[string][]string
}

type Connector struct {
	Name string
	Type string
}

func errnotexist(name string) error {
	return fmt.Errorf("Can not find name %s", name)
}

type Root struct {
	Root map[string]*Node
}

type Node struct {
	FieldName     string           `json:"name"`
	FieldType     string           `json:"type"`
	FieldLocation string           `json:"-"`
	StorageType   StorageType      `json:"storagetype"`
	Children      map[string]*Node `json:"children"`
	Tables        []*Node          `json:"tables"`
}

func NewNode(fieldname string, fieldtype string, fieldlocation string) *Node {
	return &Node{
		FieldName:     fieldname,
		FieldType:     fieldtype,
		FieldLocation: fieldlocation,
		Children:      make(map[string]*Node),
		Tables:        []*Node{},
	}
}

func (nd *Node) addchild(fieldname, fieldtype, fieldlocation, base string) *Node {
	child, ok := nd.Children[base]
	if !ok {
		child = NewNode(fieldname, fieldtype, fieldlocation)
		nd.Children[fieldname] = child
		if child.FieldName != "mapping1537182776" && child.FieldName != "array1537182776" {
			nd.Tables = append(nd.Tables, child)
		}

	}
	return child
}

func (nd *Node) Add(fieldname string, fieldtype string, fieldlocation string, base string) *Node {
	return nd.addchild(fieldname, fieldtype, fieldlocation, base)
}

func (nd *Node) Get(name string) (*Node, error) {
	if _, ok := nd.Children[name]; !ok {
		return &Node{}, errnotexist(name)
	}
	return nd, nil
}

func (nd *Node) Traversal(r Root) {
	// fmt.Printf("Node %v\n", nd)
	for _, v := range nd.Children {
		if child, ok := r.Root[v.FieldType]; ok {

			if v.FieldType != child.FieldType {
				if child1, ok := r.Root[child.FieldType]; ok {
					child = child1
				}
			}
			v.FieldType = ""
			cld := make(map[string]*Node)
			tables := []*Node{}
			for k, c := range child.Children {
				cld[k] = c

				if k != "mapping1537182776" && k != "array1537182776" {
					tables = append(tables, c)
				}
			}

			v.Children = cld
			v.Tables = tables

			//todo 优化 clang.go:204
			var allkey = ""
			for k, _ := range v.Children {
				allkey = allkey + k
			}
			if strings.Contains(allkey, "mapping1537182776") {
				v.FieldType = "mapping"
			} else if strings.Contains(allkey, "array1537182776") {
				v.FieldType = "array"
			} else {
				v.FieldType = "struct"
			}

			for k, c := range child.Children {
				if v.FieldType == "mapping" {
					if k == "key" {
						c.StorageType = MappingKeyTy
					} else if k == "value" {
						c.StorageType = MappingValueTy
					} else {
						delete(v.Children, k)
					}
				} else if v.FieldType == "array" {
					if k == "index" {
						c.StorageType = ArrayIndexTy
					} else if k == "value" {
						c.StorageType = ArrayValueTy
					} else if k == "length" {
						c.StorageType = LengthTy
					} else {
						delete(v.Children, k)
					}
				} else if v.FieldType == "struct" {
					c.StorageType = StructValueTy
				}
			}

			v.Traversal(r)
		}

	}
}

func (r Root) Fulling() Root {
	for _, v := range r.Root {
		v.Traversal(r)
	}
	return r
}

func (key *Key) KeyTraversal() {
	if key.Types == nil {
		key.Types = make(map[string]Type)
	}
	if key.Keys == nil {
		key.Keys = make(map[string]Table)
	}
	for _, v := range key.Tables {
		v.Traversal(key.Name, key)
	}
}

func KeyType(tp string) int32 {
	switch tp {
	case "int32":
		return TY_INT32
	case "int64":
		return TY_INT64
	case "uint32":
		return TY_UINT32
	case "uint64":
		return TY_UINT64
	case "uint256":
		return TY_UINT256
	case "string":
		return TY_STRING
	case "address":
		return TY_ADDRESS
	case "bool", "_Bool": //support stdbool.h
		return TY_BOOL
	case "pointer":
		return TY_POINTER
	default:
		panic(fmt.Sprintf("Error: Unsupport Type %s", tp))
	}
}
