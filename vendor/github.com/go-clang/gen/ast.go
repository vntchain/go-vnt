package gen

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"strings"
)

type ASTFunc struct {
	*ast.FuncDecl

	f   *Function
	ret *ast.ReturnStmt
}

func newASTFunc(f *Function) *ASTFunc {
	return &ASTFunc{
		FuncDecl: &ast.FuncDecl{
			Name: &ast.Ident{
				Name: f.Name,
			},
			Type: &ast.FuncType{
				Results: &ast.FieldList{
					List: []*ast.Field{},
				},
			},
			Body: &ast.BlockStmt{},
		},

		f: f,
		ret: &ast.ReturnStmt{
			Results: []ast.Expr{},
		},
	}
}

func generateFunctionString(fa *ASTFunc) string {
	var b bytes.Buffer
	err := format.Node(&b, token.NewFileSet(), []ast.Decl{fa.FuncDecl})
	if err != nil {
		panic(err)
	}

	fStr := b.String()
	fStr = strings.Replace(fStr, "REMOVE()", "", -1)

	return fStr
}

func (fa *ASTFunc) generate() {
	// TODO maybe name the return arguments ... because of clang_getDiagnosticOption -> the normal return can be always just "o"? https://github.com/go-clang/gen/issues/57

	// TODO reenable this, see the comment at the bottom of the generate function for details https://github.com/go-clang/gen/issues/54
	// Add function comment
	/*if f.Comment != "" {
		fa.Doc = &ast.CommentGroup{
			List: []*ast.Comment{
				&ast.Comment{
					Text: f.Comment,
				},
			},
		}
	}*/

	fa.generateReceiver()

	if fa.f.Member != nil {
		if fa.f.ReturnType.IsSlice {
			fa.addStatement(doDeclare("s", doGoType(fa.f.ReturnType)))
			fa.addCToGoSliceConversion("s", fa.f.Receiver.Name+".c."+fa.f.Member.Name, fa.f.Receiver.Name+".c."+fa.f.ReturnType.LengthOfSlice)

			fa.addReturnItem(&ast.Ident{
				Name: "s",
			})
			fa.addReturnType("", fa.f.ReturnType)
			fa.addEmptyLine()

			// Add the return statement
			fa.addStatement(fa.ret)
		} else {
			fa.generateReturn(&ast.SelectorExpr{
				X: accessMember(fa.f.Receiver.Name, "c"),
				Sel: &ast.Ident{
					Name: fa.f.Member.Name,
				},
			})
		}
	} else {
		// Basic call to the C function
		call := doCCast(fa.f.CName)

		if callArguments := fa.generateParameters(); len(callArguments) > 0 {
			call.Args = callArguments
		}

		fa.generateReturn(call)
	}
}

func (fa *ASTFunc) generateReceiver() {
	// Add receiver to make function a method
	if fa.f.Receiver.Name == "" {
		return
	}

	if len(fa.f.Parameters) > 0 { // TODO maybe to not set the receiver at all? -> do this outside of the generate function? https://github.com/go-clang/gen/issues/52
		fa.Recv = &ast.FieldList{
			List: []*ast.Field{
				doField(fa.f.Receiver.Name, fa.f.Receiver.Type),
			},
		}
	}
}

func (fa *ASTFunc) generateParameters() []ast.Expr {
	if len(fa.f.Parameters) == 0 {
		return nil
	}

	var callArguments []ast.Expr

	if fa.f.Receiver.Name != "" {
		fa.f.Parameters[0].Name = fa.f.Receiver.Name
	}

	fa.Type.Params = &ast.FieldList{
		List: []*ast.Field{},
	}

	hasDeclaration := false

	// Add parameters to the function
	for i, p := range fa.f.Parameters {
		if i == 0 && fa.f.Receiver.Name != "" {
			continue
		}

		// Ingore length parameters since they will be filled by the slice itself
		if p.Type.LengthOfSlice != "" && !p.Type.IsReturnArgument {
			continue
		}

		if p.Type.IsSlice && !p.Type.IsReturnArgument {
			hasDeclaration = true

			if p.Type.CGoName == CSChar && p.Type.PointerLevel >= 1 { // one pointer level from being a string, one from being an array
				fa.addGoToCSliceConversion(p.Name, p.Type)
			} else {
				fa.addCArrayFromGoSlice(p.Name, p.Type)
			}
		} else if p.Type.IsReturnArgument {
			if p.Type.LengthOfSlice == "" {
				// Add the return type to the function return arguments
				retType := p.Type
				if p.Type.GoName == "cxstring" {
					retType.GoName = "string"
				}

				fa.addReturnType("", retType)
			}

			if p.Type.IsSlice && p.Type.IsReturnArgument {
				fa.addStatement(doDeclare(
					"cp_"+p.Name,
					getSliceType(p.Type),
				))
			}

			// Declare the return argument's variable
			var varType ast.Expr
			if p.Type.PointerLevel > 0 && p.Type.CGoName == CSChar {
				varType = doPointer(doCType("char"))
			} else if p.Type.IsPrimitive {
				varType = doCType(p.Type.CGoName)
			} else {
				varType = &ast.Ident{
					Name: p.Type.GoName,
				}
			}
			if p.Type.IsSlice {
				varType = &ast.ArrayType{
					Elt: varType,
				}
			}

			fa.addStatement(doDeclare(
				p.Name,
				varType,
			))

			if p.Type.GoName == "cxstring" {
				fa.addDefer(doCall(p.Name, "Dispose"))
			} else if p.Type.PointerLevel > 0 && p.Type.CGoName == CSChar {
				fa.addDefer(doCCast(
					"free",
					doCall(
						"unsafe",
						"Pointer",
						&ast.Ident{
							Name: p.Name,
						},
					),
				))
			}

			if p.Type.LengthOfSlice == "" {
				// Add the return argument to the return statement
				if p.Type.PointerLevel > 0 && p.Type.CGoName == CSChar {
					fa.addReturnItem(doCCast(
						"GoString",
						&ast.Ident{
							Name: p.Name,
						},
					))
				} else if p.Type.GoName == "cxstring" {
					fa.addReturnItem(doCall(p.Name, "String"))
				} else if p.Type.IsPrimitive {
					fa.addReturnItem(doCast(
						p.Type.GoName,
						&ast.Ident{
							Name: p.Name,
						},
					))
				} else {
					fa.addReturnItem(&ast.Ident{
						Name: p.Name,
					})
				}
			}

			continue
		} else if p.Type.PointerLevel > 0 && p.Type.IsPrimitive && p.Type.CGoName != CSChar {
			hasDeclaration = true

			var varType = doCType(p.Type.CGoName)

			fa.addStatement(doDeclare(
				"cp_"+p.Name,
				varType,
			))
			fa.addStatement(&ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X: &ast.Ident{
						Name: p.Name,
					},
					Op: token.NEQ,
					Y: &ast.Ident{
						Name: "nil",
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.AssignStmt{
							Lhs: []ast.Expr{
								&ast.Ident{
									Name: "cp_" + p.Name,
								},
							},
							Tok: token.ASSIGN,
							Rhs: []ast.Expr{
								&ast.CallExpr{
									Fun: varType,
									Args: []ast.Expr{
										doPointer(&ast.Ident{
											Name: p.Name,
										}),
									},
								},
							},
						},
					},
				},
			})
		}

		fa.Type.Params.List = append(fa.Type.Params.List, doField(p.Name, p.Type))
	}

	if len(fa.ret.Results) > 0 || hasDeclaration {
		fa.addEmptyLine()
	}

	goToCTypeConversions := false

	// Add arguments to the C function call
	for _, p := range fa.f.Parameters {
		var pf ast.Expr

		if p.Type.IsSlice {
			pf = &ast.Ident{
				Name: "cp_" + p.Name,
			}
		} else if p.Type.IsPrimitive {
			// Handle Go type to C type conversions
			if p.Type.PointerLevel == 1 && p.Type.CGoName == CSChar {
				goToCTypeConversions = true

				fa.addAssignment(
					"c_"+p.Name,
					doCCast(
						"CString",
						&ast.Ident{
							Name: p.Name,
						},
					),
				)
				fa.addDefer(doCCast(
					"free",
					doCall(
						"unsafe",
						"Pointer",
						&ast.Ident{
							Name: "c_" + p.Name,
						},
					),
				))

				pf = &ast.Ident{
					Name: "c_" + p.Name,
				}
			} else if p.Type.CGoName == "cxstring" { // TODO try to get cxstring and "String" completely out of this function since it is just a struct which can be handled by the struct code https://github.com/go-clang/gen/issues/25
				pf = accessMember(p.Name, "c")
			} else {
				if p.Type.IsReturnArgument {
					// Return arguments already have a cast
					pf = &ast.Ident{
						Name: p.Name,
					}
				} else if p.Type.LengthOfSlice != "" {
					pf = doCCast(
						p.Type.CGoName,
						doCast(
							"len",
							&ast.Ident{
								Name: p.Type.LengthOfSlice,
							},
						),
					)
				} else if p.Type.PointerLevel > 0 {
					pf = doReference(&ast.Ident{
						Name: "cp_" + p.Name,
					})
				} else {
					pf = doCCast(
						p.Type.CGoName,
						&ast.Ident{
							Name: p.Name,
						},
					)
				}
			}
		} else {
			pf = accessMember(p.Name, "c")

			if p.Type.PointerLevel > 0 && !p.Type.IsReturnArgument && !p.Type.IsPointerComposition {
				pf = doReference(pf)
			}
		}

		if p.Type.IsReturnArgument && !p.Type.IsPointerComposition {
			pf = doReference(pf)
		}

		callArguments = append(callArguments, pf)
	}

	if goToCTypeConversions {
		fa.addEmptyLine()
	}

	return callArguments
}

func (fa *ASTFunc) generateReturn(call ast.Expr) {
	returnType := fa.f.ReturnType

	// Check if we need to add a return
	if returnType.GoName != "void" || len(fa.ret.Results) > 0 {
		if returnType.GoName == "cxstring" {
			// Do the C function call and save the result into the new variable "o" while transforming it into a cxstring
			fa.addAssignment("o", doCompose("cxstring", call))
			fa.addDefer(doCall("o", "Dispose"))
			fa.addEmptyLine()

			// Call the String method on the cxstring instance
			fa.addReturnItem(doCall("o", "String"))

			// Change the return type to "string"
			fa.addReturnType("", Type{
				GoName: "string",
			})
		} else {
			if returnType.GoName != "void" {
				// Add the function return type
				fa.addReturnType("", returnType)
			}

			// Do we need to convert the return of the C function into a boolean?
			if returnType.GoName == "bool" {
				// Do the C function call and save the result into the new variable "o"
				fa.addAssignment("o", call)
				fa.addEmptyLine()

				// Check if o is not equal to zero and return the result
				fa.addReturnItem(&ast.BinaryExpr{
					X: &ast.Ident{
						Name: "o",
					},
					Op: token.NEQ,
					Y: doCCast(
						returnType.CGoName,
						doZero(),
					),
				})
			} else if returnType.CGoName == CSChar && returnType.PointerLevel == 1 { // TODO refactor the const char * check so that one function is used everywhere to check for that C type https://github.com/go-clang/gen/issues/56
				// If this is a normal const char * C type there is not so much to do
				fa.addReturnItem(doCCast(
					"GoString",
					call,
				))
			} else if returnType.GoName == "time.Time" {
				fa.addReturnItem(doCall(
					"time",
					"Unix",
					doCast("int64", call),
					doZero(),
				))
			} else if returnType.GoName == "void" {
				// Handle the case where the C function has no return argument but parameters that are return arguments

				// Do the C function call
				fa.addStatement(&ast.ExprStmt{
					X: call,
				})
				fa.addEmptyLine()
			} else if returnType.PointerLevel > 0 {
				// Do the C function call and save the result into the new variable "o"
				fa.addAssignment("o", call)
				fa.addEmptyLine()

				fa.addStatement(doDeclare(
					"gop_o",
					doGoType(returnType),
				))
				var compositionValue ast.Expr = &ast.Ident{
					Name: "o",
				}
				if returnType.IsPointerComposition && returnType.PointerLevel == 0 {
					compositionValue = doReference(compositionValue)
				} else if !returnType.IsPointerComposition && returnType.PointerLevel > 0 {
					compositionValue = doUnreference(compositionValue)
				}
				fa.addStatement(&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "o",
						},
						Op: token.NEQ,
						Y: &ast.Ident{
							Name: "nil",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.AssignStmt{
								Lhs: []ast.Expr{
									&ast.Ident{
										Name: "gop_o",
									},
								},
								Tok: token.ASSIGN,
								Rhs: []ast.Expr{
									doReference(doCompose(
										returnType.GoName,
										compositionValue,
									)),
								},
							},
						},
					},
				})
				fa.addEmptyLine()

				fa.addReturnItem(&ast.Ident{
					Name: "gop_o",
				})
			} else {
				var convCall ast.Expr

				// Structs are literals, everything else is a cast
				if !returnType.IsPrimitive {
					convCall = doCompose(returnType.GoName, call)
				} else {
					convCall = doCast(returnType.GoName, call)
				}

				if len(fa.ret.Results) > 0 {
					// Do the C function call and save the result into the new variable "o"
					fa.addAssignment("o", convCall)
					fa.addEmptyLine()

					// Add the C function call result to the return statement
					fa.addReturnItem(&ast.Ident{
						Name: "o",
					})
				} else {
					fa.addReturnItem(convCall)
				}
			}
		}

		fa.addCToGoConversions()

		// Add the return statement
		fa.addStatement(fa.ret)
	} else {
		fa.addCToGoConversions()

		// No return needed, just add the C function call
		fa.addStatement(&ast.ExprStmt{
			X: call,
		})
	}
}

func (fa *ASTFunc) addStatement(stmt ast.Stmt) {
	fa.Body.List = append(fa.Body.List, stmt)
}

func (fa *ASTFunc) addAssignment(variable string, e ast.Expr) {
	fa.addStatement(&ast.AssignStmt{
		Lhs: []ast.Expr{
			&ast.Ident{
				Name: variable,
			},
		},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{
			e,
		},
	})
}

func (fa *ASTFunc) addDefer(call *ast.CallExpr) {
	fa.addStatement(&ast.DeferStmt{
		Call: call,
	})
}

func (fa *ASTFunc) addEmptyLine() {
	// TODO this should be done using something else than a fake statement. https://github.com/go-clang/gen/issues/53
	fa.addStatement(&ast.ExprStmt{
		X: &ast.CallExpr{
			Fun: &ast.Ident{
				Name: "REMOVE",
			},
		},
	})
}

func (fa *ASTFunc) addReturnItem(item ast.Expr) {
	fa.ret.Results = append(fa.ret.Results, item)
}

func (fa *ASTFunc) addReturnType(name string, typ Type) {
	fa.Type.Results.List = append(fa.Type.Results.List, doField(name, typ))
}

func (fa *ASTFunc) addCToGoConversions() {
	cToGoTypeConversions := false

	for _, p := range fa.f.Parameters {
		if p.Type.IsSlice && p.Type.IsReturnArgument {
			cToGoTypeConversions = true

			var lengthOfSlice string
			for _, pl := range fa.f.Parameters {
				if pl.Type.LengthOfSlice == p.Name {
					lengthOfSlice = pl.Name

					break
				}
			}

			fa.addCToGoSliceConversion(p.Name, "cp_"+p.Name, lengthOfSlice)
		}
	}

	if cToGoTypeConversions {
		fa.addEmptyLine()
	}
}

func (fa *ASTFunc) addCToGoSliceConversion(name string, cname string, lengthOfSlice string) {
	fa.addAssignment("gos_"+name, &ast.CallExpr{
		Fun: &ast.ParenExpr{
			X: doPointer(accessMember("reflect", "SliceHeader")),
		},
		Args: []ast.Expr{
			doCall(
				"unsafe",
				"Pointer",
				doReference(&ast.Ident{
					Name: name,
				}),
			),
		},
	})
	fa.addStatement(&ast.AssignStmt{
		Lhs: []ast.Expr{
			accessMember("gos_"+name, "Cap"),
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			doCast(
				"int",
				&ast.Ident{
					Name: lengthOfSlice,
				},
			),
		},
	})
	fa.addStatement(&ast.AssignStmt{
		Lhs: []ast.Expr{
			accessMember("gos_"+name, "Len"),
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			doCast(
				"int",
				&ast.Ident{
					Name: lengthOfSlice,
				},
			),
		},
	})
	fa.addStatement(&ast.AssignStmt{
		Lhs: []ast.Expr{
			accessMember("gos_"+name, "Data"),
		},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{
			doCast(
				"uintptr",
				doCall(
					"unsafe",
					"Pointer",
					&ast.Ident{
						Name: cname,
					},
				),
			),
		},
	})
}

func (fa *ASTFunc) addCArrayFromGoSlice(name string, typ Type) {
	sliceType := getSliceType(typ)

	fa.addAssignment("gos_"+name, &ast.CallExpr{
		Fun: &ast.ParenExpr{
			X: doPointer(accessMember("reflect", "SliceHeader")),
		},
		Args: []ast.Expr{
			doCall(
				"unsafe",
				"Pointer",
				doReference(&ast.Ident{
					Name: name,
				}),
			),
		},
	})
	fa.addAssignment("cp_"+name, &ast.CallExpr{
		Fun: &ast.ParenExpr{
			X: doPointer(sliceType),
		},
		Args: []ast.Expr{
			doCall(
				"unsafe",
				"Pointer",
				accessMember("gos_"+name, "Data"),
			),
		},
	})
}

func (fa *ASTFunc) addGoToCSliceConversion(name string, typ Type) {
	// Declare the slice
	sliceType := getSliceType(typ)

	fa.addAssignment(
		"ca_"+name,
		doCast(
			"make",
			&ast.ArrayType{
				Elt: sliceType,
			},
			doCast(
				"len",
				&ast.Ident{
					Name: name,
				},
			),
		),
	)
	fa.addStatement(doDeclare(
		"cp_"+name,
		doPointer(sliceType),
	))
	fa.addStatement(&ast.IfStmt{
		Cond: &ast.BinaryExpr{
			X: doCast(
				"len",
				&ast.Ident{
					Name: name,
				},
			),
			Op: token.GTR,
			Y:  doZero(),
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.Ident{
							Name: "cp_" + name,
						},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						doReference(&ast.IndexExpr{
							X: &ast.Ident{
								Name: "ca_" + name,
							},
							Index: doZero(),
						}),
					},
				},
			},
		},
	})

	// Assign elements
	var loopStatements []ast.Stmt

	// Handle our good old friend the const char * differently...
	if typ.CGoName == CSChar {
		loopStatements = append(loopStatements, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.Ident{
					Name: "ci_str",
				},
			},
			Tok: token.DEFINE,
			Rhs: []ast.Expr{
				doCCast(
					"CString",
					&ast.IndexExpr{
						X: &ast.Ident{
							Name: name,
						},
						Index: &ast.Ident{
							Name: "i",
						},
					},
				),
			},
		})
		loopStatements = append(loopStatements, &ast.DeferStmt{
			Call: doCCast(
				"free",
				doCall(
					"unsafe",
					"Pointer",
					&ast.Ident{
						Name: "ci_str",
					},
				),
			),
		})
		loopStatements = append(loopStatements, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.IndexExpr{
					X: &ast.Ident{
						Name: "ca_" + name,
					},
					Index: &ast.Ident{
						Name: "i",
					},
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.Ident{
					Name: "ci_str",
				},
			},
		})
	} else {
		loopStatements = append(loopStatements, &ast.AssignStmt{
			Lhs: []ast.Expr{
				&ast.IndexExpr{
					X: &ast.Ident{
						Name: "ca_" + name,
					},
					Index: &ast.Ident{
						Name: "i",
					},
				},
			},
			Tok: token.ASSIGN,
			Rhs: []ast.Expr{
				&ast.SelectorExpr{
					X: &ast.IndexExpr{
						X: &ast.Ident{
							Name: name,
						},
						Index: &ast.Ident{
							Name: "i",
						},
					},
					Sel: &ast.Ident{
						Name: "c",
					},
				},
			},
		})
	}

	fa.addStatement(&ast.RangeStmt{
		Key: &ast.Ident{
			Name: "i",
		},
		Tok: token.DEFINE,
		X: &ast.Ident{
			Name: name,
		},
		Body: &ast.BlockStmt{
			List: loopStatements,
		},
	})
}
