package main

import (
	"fmt"
	"regexp"
)

//key 类型检查，如果key写在其他方法之间，报错并提示
//没有constructor
//call,第一个参数检查
//mutable
//unmutable，key的写入检查，sendfromcontract,transferfromcontract检查
//payable,input的检查，sendfromcontract,transferfromcontract
//导出方法的类型和参数检查
//没有一个导出方法
//uint256及address类型

type Hint struct {
	Code           []byte
	ConstructorPos int
}

func newHint(code []byte) *Hint {

	return &Hint{
		Code: code,
	}
}

func (h *Hint) contructorCheck() {
	reg := regexp.MustCompile(constructorReg)
	idx := reg.FindAllStringIndex(string(h.Code), -1)
	if len(idx) == 0 {
		panic("Can't find Contructor function")
	}
	if len(idx) > 1 {
		panic("Can only have one constructor")
	}
	h.ConstructorPos = idx[0][0]
	fmt.Printf("pos %d\n", h.ConstructorPos)
	fmt.Printf("code %s\n", h.Code[idx[0][0]:idx[0][1]])
}

func (h *Hint) keyCheck() {

}
