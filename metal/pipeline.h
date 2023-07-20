// go:build darwin
//  +build darwin

#ifndef HEADER_PIPELINE
#define HEADER_PIPELINE

#import <Metal/Metal.h>

// Structure of various metal resources needed to execute a computational
// process on the GPU. We have to bundle this in a header that cgo doesn't
// import because of a bug in LLVM that leads to a compilation error of "struct
// size calculation error off=8 bytesize=0". Doesn't seem to be another solution
// to this at the moment.
typedef struct {
  // Metal resources
  id<MTLComputePipelineState> pipeline;
  id<MTLCommandQueue> commandQueue;
} _pipeline;

_pipeline *pipeline_new();

#endif