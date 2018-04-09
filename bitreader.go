//===- bitreader.go - Bindings for bitreader ------------------------------===//
//
//                     The LLVM Compiler Infrastructure
//
// This file is distributed under the University of Illinois Open Source
// License. See LICENSE.TXT for details.
//
//===----------------------------------------------------------------------===//
//
// This file defines bindings for the bitreader component.
//
//===----------------------------------------------------------------------===//

package llvm

/*
#include "llvm-c/BitReader.h"
#include "llvm-c/Core.h"
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

// ParseBitcodeFile parses the LLVM IR (bitcode) in the file with the
// specified name, and returns a new LLVM module.
func ParseBitcodeFile(name string) (Module, error) {
	var buf C.LLVMMemoryBufferRef
	var errmsg *C.char
	var cfilename *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(cfilename))
	result := C.LLVMCreateMemoryBufferWithContentsOfFile(cfilename, &buf, &errmsg)
	if result != 0 {
		err := errors.New(C.GoString(errmsg))
		C.free(unsafe.Pointer(errmsg))
		return Module{}, err
	}
	defer C.LLVMDisposeMemoryBuffer(buf)

	var m Module
	if C.LLVMParseBitcode2(buf, &m.C) == 0 {
		return m, nil
	}

	err := errors.New(C.GoString(errmsg))
	C.free(unsafe.Pointer(errmsg))
	return Module{}, err
}

// ParseBitcodeFileInContext parses the LLVM IR (bitcode) in the file with the
// specified name, and returns a new LLVM module in the given context
func ParseBitcodeFileInContext(context Context, name string) (Module, error) {
	var buf C.LLVMMemoryBufferRef
	var errmsg *C.char
	var cfilename *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(cfilename))
	result := C.LLVMCreateMemoryBufferWithContentsOfFile(cfilename, &buf, &errmsg)
	if result != 0 {
		err := errors.New(C.GoString(errmsg))
		C.free(unsafe.Pointer(errmsg))
		return Module{}, err
	}
	defer C.LLVMDisposeMemoryBuffer(buf)

	var m Module
	if C.LLVMParseBitcodeInContext2(context.C, buf, &m.C) == 0 {
		return m, nil
	}

	err := errors.New(C.GoString(errmsg))
	C.free(unsafe.Pointer(errmsg))
	return Module{}, err
}

// ParseBitcodeInContext parses the LLVM IR (bitcode) from the given byteslice
// and returns a new LLVM module in the given context.
func ParseBitcodeInContext(context Context, name string, bitcode []byte) (Module, error) {
	var buf C.LLVMMemoryBufferRef
	var csize = C.size_t(uint(len(bitcode)))

	var cname *C.char = C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cbitcode := C.CString(string(bitcode))
	defer C.free(unsafe.Pointer(cbitcode))

	buf = C.LLVMCreateMemoryBufferWithMemoryRangeCopy(cbitcode, csize, cname)
	defer C.LLVMDisposeMemoryBuffer(buf)

	var m Module
	if C.LLVMParseBitcodeInContext2(context.C, buf, &m.C) == 0 {
		return m, nil
	}

	err := errors.New("failed to parse bitcode")
	return Module{}, err
}
