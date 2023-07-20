//go:build darwin
// +build darwin

package metal

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

// Buffer contains a slice that wraps a block of memory that can be sent to the GPU and used to run
// metal code.
type Buffer[T any] struct {
	// Block of memory accessible to both the CPU and GPU. Only the contents of the underlying array
	// should be modified. The length/capacity of the slice and which block of memory it points to
	// should not be altered,
	Data []T

	// Id of the buffer, as assigned by the underlying code that creates and manages it. This is
	// used to send the buffer as an argument to the metal code.
	id int
}

// BufferId returns the Id of the buffer. This Id is assigned by the process that creates and
// manages the buffer's memory.
func (b Buffer[T]) BufferId() int { return b.id }

// A Pipeline executes computational processes on the default GPU.
type Pipeline struct {
	// Id of the pipeline, as assigned by the underlying code that creates and manages it. This is
	// be used to run the pipeline and execute its computational process on the GPU.
	id int
}

// NewPipeline sets up a new pipeline that will run on the default GPU. The pipeline is built with
// the specified function in the provided metal code.
func NewPipeline(metalSource, funcName string) Pipeline {
	src := C.CString(metalSource)
	defer C.free(unsafe.Pointer(src))

	name := C.CString(funcName)
	defer C.free(unsafe.Pointer(name))

	return Pipeline{
		id: int(C.metal_newFunction(src, name)),
	}
}

// A BufferIder can advertise the Id that references its buffer in metal code.
type BufferIder interface {
	// BufferId returns the Id of the buffer that was created for sending data to the GPU. This Id
	// is used when running the computation pipeline.
	BufferId() int
}

// NewBuffer allocates a block of memory that is accessible to both the CPU and GPU. The returned
// Buffer contains a slice that wraps the new memory and has a length equal to numElems, and the
// underlying memory of the slice has (numElems * sizeof(T)) bytes.
func NewBuffer[T any](numElems int) Buffer[T] {
	if numElems <= 0 {
		return Buffer[T]{}
	}

	size := sizeof[T]()
	bufferSize := numElems * size

	// Allocate memory for the new buffer, and then retrieve a pointer to the beginning of the new
	// memory using the buffer's Id.
	bufferId := C.metal_newBuffer(C.int(bufferSize))
	newBuffer := C.metal_retrieveBuffer(bufferId)

	return Buffer[T]{
		Data: toSlice[T](newBuffer, numElems),
		id:   int(bufferId),
	}
}

// Run executes the computation function stored in pipeline on the GPU. buffers is a list of buffers
// that have a buffer Id, which is used to retrieve the correct block of memory for the buffer. Each
// buffer is supplied as an argument to the pipeline's metal function in the order given here.
func Run(pipeline Pipeline, buffers ...BufferIder) {

	// Make a list of buffer Ids.
	var bufferIds []C.int
	for _, buffer := range buffers {
		bufferIds = append(bufferIds, C.int(buffer.BufferId()))
	}

	// Get a pointer to the beginning of the list of buffer Ids (if we have any).
	var bufferPtr *C.int
	if len(bufferIds) > 0 {
		bufferPtr = &bufferIds[0]
	}

	// Run the computation on the GPU.
	C.metal_runFunction(C.int(pipeline.id), bufferPtr, C.int(len(bufferIds)))
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
