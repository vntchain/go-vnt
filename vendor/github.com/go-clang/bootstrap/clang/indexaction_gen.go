package clang

// #include "./clang-c/Index.h"
// #include "go-clang.h"
import "C"
import (
	"reflect"
	"unsafe"
)

// An indexing action/session, to be applied to one or multiple translation units.
type IndexAction struct {
	c C.CXIndexAction
}

/*
	Destroy the given index action.

	The index action must not be destroyed until all of the translation units
	created within that index action have been destroyed.
*/
func (ia IndexAction) Dispose() {
	C.clang_IndexAction_dispose(ia.c)
}

/*
	Index the given source file and the translation unit corresponding
	to that file via callbacks implemented through #IndexerCallbacks.

	Parameter client_data pointer data supplied by the client, which will
	be passed to the invoked callbacks.

	Parameter index_callbacks Pointer to indexing callbacks that the client
	implements.

	Parameter index_callbacks_size Size of #IndexerCallbacks structure that gets
	passed in index_callbacks.

	Parameter index_options A bitmask of options that affects how indexing is
	performed. This should be a bitwise OR of the CXIndexOpt_XXX flags.

	Parameter out_TU [out] pointer to store a CXTranslationUnit that can be reused
	after indexing is finished. Set to NULL if you do not require it.

	Returns If there is a failure from which the there is no recovery, returns
	non-zero, otherwise returns 0.

	The rest of the parameters are the same as #clang_parseTranslationUnit.
*/
func (ia IndexAction) IndexSourceFile(clientData ClientData, indexCallbacks *IndexerCallbacks, indexCallbacksSize uint32, indexOptions uint32, sourceFilename string, commandLineArgs []string, unsavedFiles []UnsavedFile, outTU *TranslationUnit, tUOptions uint32) int32 {
	ca_commandLineArgs := make([]*C.char, len(commandLineArgs))
	var cp_commandLineArgs **C.char
	if len(commandLineArgs) > 0 {
		cp_commandLineArgs = &ca_commandLineArgs[0]
	}
	for i := range commandLineArgs {
		ci_str := C.CString(commandLineArgs[i])
		defer C.free(unsafe.Pointer(ci_str))
		ca_commandLineArgs[i] = ci_str
	}
	gos_unsavedFiles := (*reflect.SliceHeader)(unsafe.Pointer(&unsavedFiles))
	cp_unsavedFiles := (*C.struct_CXUnsavedFile)(unsafe.Pointer(gos_unsavedFiles.Data))

	c_sourceFilename := C.CString(sourceFilename)
	defer C.free(unsafe.Pointer(c_sourceFilename))

	return int32(C.clang_indexSourceFile(ia.c, clientData.c, &indexCallbacks.c, C.uint(indexCallbacksSize), C.uint(indexOptions), c_sourceFilename, cp_commandLineArgs, C.int(len(commandLineArgs)), cp_unsavedFiles, C.uint(len(unsavedFiles)), &outTU.c, C.uint(tUOptions)))
}

/*
	Index the given translation unit via callbacks implemented through
	#IndexerCallbacks.

	The order of callback invocations is not guaranteed to be the same as
	when indexing a source file. The high level order will be:

	-Preprocessor callbacks invocations
	-Declaration/reference callbacks invocations
	-Diagnostic callback invocations

	The parameters are the same as #clang_indexSourceFile.

	Returns If there is a failure from which the there is no recovery, returns
	non-zero, otherwise returns 0.
*/
func (ia IndexAction) IndexTranslationUnit(clientData ClientData, indexCallbacks *IndexerCallbacks, indexCallbacksSize uint32, indexOptions uint32, tu TranslationUnit) int32 {
	return int32(C.clang_indexTranslationUnit(ia.c, clientData.c, &indexCallbacks.c, C.uint(indexCallbacksSize), C.uint(indexOptions), tu.c))
}
