package main

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/go-clang/bootstrap/clang"
	"github.com/vntchain/go-vnt/accounts/abi"
)

var KeyPos [][]int

var index = 0

func cmd(args []string) int {
	idx := clang.NewIndex(0, 1)
	defer idx.Dispose()
	tu := idx.ParseTranslationUnit(args[0], []string{"-I", includeDir}, nil, 0)
	defer tu.Dispose()

	diagnostics := tu.Diagnostics()
	for _, d := range diagnostics {
		// fmt.Printf("d %+v\n", d)
		fmt.Println("PROBLEM:", d.Spelling())
	}

	cursor := tu.TranslationUnitCursor()

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.IsNull() {
			return clang.ChildVisit_Continue
		}
		createStructList(cursor, parent)
		switch cursor.Kind() {
		case clang.Cursor_ClassDecl, clang.Cursor_EnumDecl, clang.Cursor_StructDecl, clang.Cursor_Namespace:
			return clang.ChildVisit_Recurse
		}
		return clang.ChildVisit_Continue
	})
	structLists.Fulling()

	cursor.Visit(func(cursor, parent clang.Cursor) clang.ChildVisitResult {
		if cursor.IsNull() {
			return clang.ChildVisit_Continue
		}
		getVarDecl(cursor, parent)
		switch cursor.Kind() {
		case clang.Cursor_ClassDecl, clang.Cursor_EnumDecl, clang.Cursor_StructDecl, clang.Cursor_Namespace:
			return clang.ChildVisit_Recurse
		}
		return clang.ChildVisit_Continue
	})
	// structLists.Fulling()
	// jsonres, _ := json.Marshal(structLists)
	// fmt.Printf("structLists %s\n", jsonres)

	if len(diagnostics) > 0 {
		fmt.Println("NOTE: There were problems while analyzing the given file")
	}

	return 0
}

func createStructList(cursor, parent clang.Cursor) {

	decl := cursor.Kind()
	cursorname := cursor.Spelling()
	cursortype := cursor.Type().Spelling()
	usr := cursor.USR()

	pdecl := parent.Kind()
	pcursorname := parent.Spelling()
	// //pusr := parent.USR()
	pcursortype := parent.Type().Spelling()
	// fmt.Printf("\n******          %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
	// fmt.Printf("******parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())
	if decl == clang.Cursor_StructDecl && pdecl == clang.Cursor_TranslationUnit { //声明结构体
		if cursorname == "" {
			if strings.Contains(cursortype, "struct (anonymous") { //匿名结构
				// fmt.Printf("匿名结构体\n")
				index = index + 1
				fieldtype := fmt.Sprintf("%s@@@%d", usr, index)
				node := abi.NewNode(cursorname, fieldtype, "")
				structStack = append(structStack, node)
				structLists.Root[fieldtype] = node
				// fmt.Printf("cursorname %s fieldtype %s cursor", cursorname, fieldtype)
			} else {
				// fmt.Printf("typedef 匿名结构体\n")
				// fmt.Printf("cursortype %s\n", cursortype)
				node := abi.NewNode(cursorname, cursortype, "")
				structLists.Root[cursortype] = node
				// fmt.Printf("typedef 匿名结构体 end\n")
			}
		} else {
			node := abi.NewNode(cursorname, cursortype, "")
			structLists.Root[cursortype] = node
			if strings.Contains(cursortype, "struct") {
				structLists.Root[cursortype[7:]] = node
			} else {
				structLists.Root[fmt.Sprintf("struct %s", cursortype)] = node
			}

			//structStack = append(structStack, node)
		}
	} else if decl == clang.Cursor_TypedefDecl && pdecl == clang.Cursor_TranslationUnit { //使用typedef定义的结构体解析
		// fmt.Printf("cursor.TypedefDeclUnderlyingType().Spelling() %s\n", cursor.TypedefDeclUnderlyingType().Spelling())
		// fmt.Printf("cursor.Type() %s\n", cursor.Type().Spelling())
		if strings.Contains(cursor.TypedefDeclUnderlyingType().Spelling(), "struct") {
			fieldname := cursor.TypedefDeclUnderlyingType().Spelling()[7:]
			// fmt.Printf("cursorname %s fieldname %s\n", cursorname, fieldname)
			if fieldname == cursorname { //匿名结构体
			} else {
				node := abi.NewNode(cursorname, fieldname, "")
				structLists.Root[cursorname] = node
			}
		}
	} else if decl == clang.Cursor_StructDecl && pdecl == clang.Cursor_StructDecl { //结构体内部定义的结构体
		if cursorname == "" {
			index = index + 1
			fieldtype := fmt.Sprintf("%s@@@%d", usr, index)
			node := abi.NewNode(cursorname, fieldtype, "")
			structStack = append(structStack, node)
			structLists.Root[fieldtype] = node
		} else {
			node := abi.NewNode(cursorname, cursortype, "")
			structLists.Root[cursorname] = node
			structLists.Root[fmt.Sprintf("struct %s", cursorname)] = node
		}

	} else if decl == clang.Cursor_FieldDecl && pdecl == clang.Cursor_StructDecl {
		// fmt.Printf(" decl == clang.Cursor_FieldDecl && pdecl == clang.Cursor_StructDecl \n")
		if len(structStack) != 0 {
			node := structStack[len(structStack)-1]
			if cutUSR(usr) == strings.Split(node.FieldType, "@@@")[0] {
				// fmt.Printf("cutUSR(usr) == strings.Split(node.FieldType, `@@@`)[0]\n")
				//fmt.Printf("struct element %s\n", strings.Split(node.FieldType, "@@@")[1])
				node.Add(cursorname, cursortype, "", node.FieldType)
			} else {
				// fmt.Printf("node.FieldType %s\n", node.FieldType)
				// fmt.Printf("struct element %v\n", strings.Split(node.FieldType, "@@@"))
				if pcursorname == "" {
					childnode := structStack[len(structStack)-1]
					structStack = structStack[0 : len(structStack)-1]
					node := structStack[len(structStack)-1]
					node.Add(cursorname, childnode.FieldType, "", node.FieldType)
				} else if strings.Contains(cursortype, "anonymous struct") {
					childnode := structStack[len(structStack)-1]
					node = structLists.Root[pcursorname]
					node.Add(cursorname, childnode.FieldType, "", pcursorname)
				} else {
					node = structLists.Root[pcursorname]
					node.Add(cursorname, cursortype, "", pcursorname)
				}
			}
		} else {
			if pcursorname != "" {
				node := structLists.Root[pcursorname]
				node.Add(cursorname, cursortype, "", pcursorname)
			} else {
				node := structLists.Root[pcursortype]
				// fmt.Printf("node %+v\n", node)
				node.Add(cursorname, cursortype, "", pcursortype)
				// fmt.Printf("cursorname %s cursortype %s  pcursortype %s\n", cursorname, cursortype, pcursortype)
			}

		}

	} else if decl == clang.Cursor_VarDecl && strings.Contains(cursortype, "anonymous struct") {
		// fmt.Printf(` decl == clang.Cursor_VarDecl && strings.Contains(cursortype, "anonymous struct")\n`)
		node := structStack[len(structStack)-1]
		structStack = structStack[0 : len(structStack)-1]
		structLists.Root[cursortype] = node
	}

}

//c:main6.cpp@S@main6.cpp@8255
func getVarDecl(cursor, parent clang.Cursor) {
	decl := cursor.Kind()
	cursortype := cursor.Type().Spelling()
	cursorname := cursor.Spelling()
	allstruct := []string{}
	for k, _ := range structLists.Root {
		allstruct = append(allstruct, k)
	}
	structnames := strings.Join(allstruct, "|")
	if decl == clang.Cursor_VarDecl {
		// fmt.Printf("\n******          %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
		// fmt.Printf("******parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())
		// if strings.Contains(cursortype, "struct") || strings.Contains(cursortype, "volatile _S") {
		// fmt.Printf("!!!cursortype %s\n", cursortype)
		if !strings.Contains(cursortype, "volatile") {
			return
		}
		if strings.Contains(strings.Join(allstruct, ""), cursortype[9:]) {
			// fmt.Printf("======cursortype=======%s\n", cursortype)
			// fmt.Printf("\n******          %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
			// fmt.Printf("******parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())
			if strings.Contains(cursortype, "struct (anonymous") {
				contents := strings.Split(cursortype, ":")
				num, err := strconv.Atoi(contents[len(contents)-2])
				if err != nil {
					panic(err)
				}
				if !isKey(fileContent[num-1], structnames) {
					return
				}
			} else {
				_, x1, _, _ := cursor.Location().FileLocation()
				if !isKey(fileContent[x1-1], structnames) {
					return
				}
			}
			for k, v := range structLists.Root {
				//volatile struct (anonymous xxxx
				if k == cursortype || cursortype == "volatile "+k {
					varLists.Root[cursorname] = v
					v.FieldName = cursorname
					v.FieldType = ""
					v.FieldLocation = ""

					//todo 优化 key.go:84
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
					for k, c := range v.Children {
						if v.FieldType == "mapping" {
							if k == "key" {
								c.StorageType = abi.MappingKeyTy
							} else if k == "value" {
								c.StorageType = abi.MappingValueTy
							} else {
								delete(v.Children, k)
							}
						} else if v.FieldType == "array" {
							if k == "index" {
								c.StorageType = abi.ArrayIndexTy
							} else if k == "value" {
								c.StorageType = abi.ArrayValueTy
							} else if k == "length" {
								c.StorageType = abi.LengthTy
							} else {
								delete(v.Children, k)
							}
						} else if v.FieldType == "struct" {
							c.StorageType = abi.StructValueTy
						}
					}
				}
			}
		} else {
			sourceFile, x1, _, _ := cursor.Location().FileLocation()
			ext := path.Ext(sourceFile.Name())
			if strings.Compare(ext, ".h") == 0 {
				return
			}
			if !isKey(fileContent[x1-1], structnames) {
				return
			}
			if strings.Contains(cursortype, "volatile") {
				//volatile int64
				node := abi.NewNode(cursorname, cursortype[9:], "")
				node.StorageType = abi.NormalTy
				varLists.Root[cursorname] = node
			} else {
				node := abi.NewNode(cursorname, cursortype, "")
				node.StorageType = abi.NormalTy
				varLists.Root[cursorname] = node
			}

		}
	}
}

func getFunc(cursor, parent clang.Cursor) {

	if cursor.Kind() == clang.Cursor_FunctionDecl {
		_, x1, _, _ := cursor.Location().FileLocation()
		fmt.Println("func =======================")
		fmt.Printf("cursor %v \n", fileContent[x1-1])
		fmt.Printf("func          %s: %s (%s) (%s)\n", cursor.Kind().Spelling(), cursor.Spelling(), cursor.USR(), cursor.Type().Spelling())
		fmt.Printf("func parent    %s: %s (%s) (%s)\n", parent.Kind().Spelling(), parent.Spelling(), parent.USR(), parent.Type().Spelling())
	}

}
