package clang

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/termie/go-shutil"

	"github.com/go-clang/gen"
)

// cmdFatal returns a command error
func cmdFatal(msg string, err error) error {
	if err == nil {
		return fmt.Errorf("FATAL %s", msg)
	} else {
		return fmt.Errorf("FATAL %s: %s", msg, err)
	}
}

// Cmd executes a generic go-clang-generate command
func Cmd(args []string, api *gen.API) error {
	rawLLVMVersion, _, err := execToBuffer("llvm-config", "--version")
	if err != nil {
		return cmdFatal("Cannot determine LLVM version", err)
	}

	llvmVersion := ParseVersion(rawLLVMVersion)
	if llvmVersion == nil {
		return cmdFatal("Cannot parse LLVM version", nil)
	}

	fmt.Println("Found LLVM version", llvmVersion.String())

	rawLLVMIncludeDir, _, err := execToBuffer("llvm-config", "--includedir")
	if err != nil {
		return cmdFatal("Cannot determine LLVM include directory", err)
	}

	clangCIncludeDir := strings.TrimSpace(string(rawLLVMIncludeDir)) + "/clang-c/"
	if err := dirExists(clangCIncludeDir); err != nil {
		return cmdFatal(fmt.Sprintf("Cannot find Clang-C include directory %q", clangCIncludeDir), err)
	}

	fmt.Println("Clang-C include directory", clangCIncludeDir)

	clangArguments := []string{
		"-I", ".", // Include the current directory
	}

	// Find Clang's include directory
	for _, d := range []string{
		"/usr/local/lib/clang",
		"/usr/include/clang",
	} {
		for _, di := range []string{
			d + fmt.Sprintf("/%s/include", llvmVersion.String()),
			d + fmt.Sprintf("/%s/include", llvmVersion.StringMinor()),
		} {
			if dirExists(di) == nil {
				clangArguments = append(clangArguments, "-I", di)
			}
		}
	}

	api.ClangArguments = append(clangArguments, api.ClangArguments...)

	fmt.Printf("Using clang arguments: %v\n", clangArguments)

	fmt.Printf("Will generate go-clang for LLVM version %s into the current directory\n", llvmVersion.String())

	clangCDirectory := "./clang-c/"

	// Copy the Clang-C include directory into the current directory
	_ = os.RemoveAll(clangCDirectory)
	if err := shutil.CopyTree(clangCIncludeDir, clangCDirectory, nil); err != nil {
		return cmdFatal(fmt.Sprintf("Cannot copy Clang-C include directory %q into current directory", clangCIncludeDir), err)
	}

	// Remove all generated .go files
	if files, err := ioutil.ReadDir("./"); err != nil {
		return cmdFatal("Cannot read current directory", err)
	} else {
		for _, f := range files {
			fn := f.Name()

			if !f.IsDir() && strings.HasSuffix(fn, "_gen.go") {
				if err := os.Remove(fn); err != nil {
					return cmdFatal(fmt.Sprintf("Cannot remove generated file %q", fn), err)
				}
			}
		}
	}

	headerFiles, err := api.HandleDirectory(clangCDirectory)
	if err != nil {
		return err
	}

	generator := gen.NewGeneration(api)
	generator.AddHeaderFiles(headerFiles)
	if err = generator.Generate(); err != nil {
		return err
	}

	return nil
}
