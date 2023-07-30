//go:build darwin
// +build darwin

package metal

// frameworks not included:
// -framework Cocoa

/*
#cgo LDFLAGS: -framework Metal -framework CoreGraphics -framework Foundation
#include "metal.h"
*/
import "C"

import (
	"reflect"
	"unsafe"
)

func init() {
	// Initialize the device that will be used to run the computations.
	C.metal_init()
}

// A Function executes computational processes on the default GPU.
type Function struct {
	// Id of the metal function, as assigned by the underlying code that creates and manages it.
	// This is used to run the function and execute its computational process on the GPU.
	id int
}

// NewFunction sets up a new function that will run on the default GPU. It is built with the
// specified function in the provided metal code.
func NewFunction(metalSource, funcName string) Function {
	src := C.CString(metalSource)
	defer C.free(unsafe.Pointer(src))

	name := C.CString(funcName)
	defer C.free(unsafe.Pointer(name))

	return Function{
		id: int(C.metal_newFunction(src, name)),
	}
}

// A BufferId references a specific metal buffer created with NewBuffer.
type BufferId int

// NewBuffer allocates a block of memory that is accessible to both the CPU and GPU. It returns a
// unique Id for the buffer and a slice that wraps the new memory and has a len and cap equal to
// numElems.
//
// The Id is used to reference the buffer as an argument for the metal function.
//
// Only the contents of the slice should be modified. Its length and capacity and the block of
// memory that it points to should not be altered. The slice's length and capacity are equal to
// numElems, and its underlying memory has (numElems * sizeof(T)) bytes.
func NewBuffer[T any](numElems int) (BufferId, []T) {
	if numElems <= 0 {
		return 0, nil
	}

	elemSize := sizeof[T]()

	// Allocate memory for the new buffer, and then retrieve a pointer to the beginning of the new
	// memory using the buffer's Id.
	bufferId := C.metal_newBuffer(C.int(elemSize * numElems))
	newBuffer := C.metal_retrieveBuffer(bufferId)

	return BufferId(bufferId), toSlice[T](newBuffer, numElems)
}

// Run executes the computational function on the GPU. buffers is a list of buffers that have a
// buffer Id, which is used to retrieve the correct block of memory for the buffer. Each buffer is
// supplied as an argument to the metal function in the order given here.
func Run(function Function, buffers ...BufferId) {

	// Make a list of buffer Ids.
	var bufferIds []C.int
	for _, buffer := range buffers {
		bufferIds = append(bufferIds, C.int(buffer))
	}

	// Get a pointer to the beginning of the list of buffer Ids (if we have any).
	var bufferPtr *C.int
	if len(bufferIds) > 0 {
		bufferPtr = &bufferIds[0]
	}

	// Run the computation on the GPU.
	C.metal_runFunction(C.int(function.id), bufferPtr, C.int(len(bufferIds)))
}

// sizeof returns the size in bytes of the generic type T.
func sizeof[T any]() int {
	var t T
	return int(unsafe.Sizeof(t))
}

// toSlice transforms a block of memory into a go slice. It wraps the memory inside a slice header
// and sets the len/cap to the number of elements. This is unsafe behavior and can lead to data
// corruption.
func toSlice[T any](data unsafe.Pointer, numElems int) []T {
	// Create a slice header with the generic type for a slice that has no backing array.
	var s []T

	// Cast the slice header into a reflect.SliceHeader so we can actually access the slice's
	// internals and set our own values. In effect, this wraps a go slice around our data so we can
	// access it natively.
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&s))

	// Set our data in the slice internals.
	hdr.Data = uintptr(data)
	hdr.Len = numElems
	hdr.Cap = numElems

	return s
}
